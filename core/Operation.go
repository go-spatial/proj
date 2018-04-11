package core

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"

	"github.com/go-spatial/proj4go/merror"
	"github.com/go-spatial/proj4go/support"
)

// DatumType is the enum for the types of datums we support
type DatumType int

// All the DatumType constants
const (
	DatumTypeUnknown   DatumType = 0
	DatumType3Param              = 1
	DatumType7Param              = 2
	DatumTypeGridShift           = 3
	DatumTypeWGS84               = 4 /* WGS84 (or anything considered equivalent) */
)

// IOUnitsType is the enum for the types of input/output units we support
type IOUnitsType int

// All the IOUnitsType constants
const (
	IOUnitsWhatever  IOUnitsType = 0 /* Doesn't matter (or depends on pipeline neighbours) */
	IOUnitsClassic               = 1 /* Scaled meters (right), projected system */
	IOUnitsProjected             = 2 /* Meters, projected system */
	IOUnitsCartesian             = 3 /* Meters, 3D cartesian system */
	IOUnitsAngular               = 4 /* Radians */
)

// DirectionType is the enum for the operation's direction
type DirectionType int

// All the DirectionType constants
const (
	DirectionForward  DirectionType = 1  /* Forward    */
	DirectionIdentity               = 0  /* Do nothing */
	DirectionInverse                = -1 /* Inverse    */
)

const epsLat = 1.0e-12

// Operation contains all the info needed to describe an "operation",
// i.e. a "conversion" (no datum change) or a "transformation".
//
// In PROJ.4, a "projection" is a conversion from "angular" input to "scaled linear" output.
type Operation struct {
	ProjString *ProjString
	Info       *OperationInfo

	//
	// COORDINATE HANDLING
	//
	Over         bool /* Over-range flag */
	Geoc         bool /* Geocentric latitude flag */
	IsLatLong    bool /* proj=latlong ... not really a projection at all */
	IsGeocentric bool /* proj=geocent ... not really a projection at all */
	NeedEllps    bool /* 0 for operations that are purely cartesian */

	Left  IOUnitsType /* Flags for input/output coordinate types */
	Right IOUnitsType

	//
	// ELLIPSOID
	//
	Ellipsoid *Ellipsoid

	//
	// CARTOGRAPHIC OFFSETS
	//
	Lam0, Phi0     float64 /* central meridian, parallel */
	X0, Y0, Z0, T0 float64 /* false easting and northing (and height and time) */

	//
	// SCALING
	//
	K0                   float64 /* General scaling factor - e.g. the 0.9996 of UTM */
	ToMeter, FromMeter   float64 /* Plane coordinate scaling. Internal unit [m] */
	VToMeter, VFromMeter float64 /* Vertical scaling. Internal unit [m] */

	//
	// DATUMS AND HEIGHT SYSTEMS
	//
	DatumType   DatumType  /* PJD_UNKNOWN/3PARAM/7PARAM/GRIDSHIFT/WGS84 */
	DatumParams [7]float64 /* Parameters for 3PARAM and 7PARAM */

	//struct _pj_gi **gridlist;
	//int     gridlist_count;

	HasGeoidVgrids bool
	//struct _pj_gi **vgridlist_geoid;
	//int     vgridlist_geoid_count;

	FromGreenwich  float64 /* prime meridian offset (in radians) */
	LongWrapCenter float64 /* 0.0 for -180 to 180, actually in radians*/
	IsLongWrapSet  bool
	Axis           string /* Axis order, pj_transform/pj_adjust_axis */

	/* New Datum Shift Grid Catalogs */
	CatalogName string
	//struct _PJ_GridCatalog *catalog;
	DatumDate float64

	//struct _pj_gi *last_before_grid;    /* TODO: Description needed */
	//PJ_Region     last_before_region;   /* TODO: Description needed */
	//double        last_before_date;     /* TODO: Description needed */

	//struct _pj_gi *last_after_grid;     /* TODO: Description needed */
	//PJ_Region     last_after_region;    /* TODO: Description needed */
	//double        last_after_date;      /* TODO: Description needed */

	//
	// OPAQUE
	//
	Q interface{} // pointer to the "opaque" object
}

// NewOperation returns a new Operation object
func NewOperation(ps *ProjString) (*Operation, error) {

	op := &Operation{
		ProjString: ps,
		NeedEllps:  true,
		Left:       IOUnitsAngular,
		Right:      IOUnitsClassic,
		Axis:       "enu",
	}

	err := op.initialize()
	if err != nil {
		return nil, err
	}

	return op, nil
}

// Forward executes an operation
func (op *Operation) Forward(input interface{}) (interface{}, error) {

	input2, err := op.forwardPrepare(input.(*CoordLP))
	if err != nil {
		return nil, err
	}

	f := op.Info.forward
	output, err := f(op, input2)

	output2, err := op.forwardFinalize(output.(*CoordXY))
	if err != nil {
		return nil, err
	}

	return output2, nil
}

// Inverse executes an operation in reverse
func (op *Operation) Inverse(input interface{}) (interface{}, error) {
	input2, err := op.inversePrepare(input.(*CoordXY))
	if err != nil {
		return nil, err
	}

	f := op.Info.inverse
	output, err := f(op, input2)

	output2, err := op.inverseFinalize(output.(*CoordLP))
	if err != nil {
		return nil, err
	}

	return output2, nil
}

func (op *Operation) String() string {
	b, err := json.MarshalIndent(op, "", " ")
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (op *Operation) initialize() error {

	projName, _ := op.ProjString.Args.GetAsString("proj")
	opInfo, ok := OperationInfoTable[projName]
	if !ok {
		return merror.New(merror.BadProjStringError)
	}

	op.Info = opInfo

	err := op.processDatum()
	if err != nil {
		return err
	}

	err = op.processEllipsoid()
	if err != nil {
		return err
	}

	/* Now that we have ellipse information check for WGS84 datum */
	if op.DatumType == DatumType3Param &&
		op.DatumParams[0] == 0.0 &&
		op.DatumParams[1] == 0.0 &&
		op.DatumParams[2] == 0.0 &&
		op.Ellipsoid.A == 6378137.0 &&
		math.Abs(op.Ellipsoid.Es-0.006694379990) < 0.000000000050 {
		/*WGS84/GRS80*/
		op.DatumType = DatumTypeWGS84
	}

	err = op.processMisc()
	if err != nil {
		return err
	}

	// do setup work specific to this operation
	// (locate_constructor)
	err = op.Info.setup(op)
	if err != nil {
		return err
	}

	return nil
}

func (op *Operation) processDatum() error {

	op.DatumType = DatumTypeUnknown

	datumName, ok := op.ProjString.Args.GetAsString("datum")
	if ok {

		datum, ok := DatumTable[datumName]
		if !ok {
			return merror.New(merror.NoSuchDatum)
		}

		// add the ellipse to the end of the list
		// TODO: move this into the ProjString processor?

		if datum.EllipseID != "" {
			op.ProjString.Args.Add(support.Pair{Key: "ellps", Value: datum.EllipseID})
		}
		if datum.DefinitionString != "" {
			op.ProjString.Args.AddList(datum.Definition)
		}
	}

	if op.ProjString.Args.ContainsKey("nadgrids") {
		op.DatumType = DatumTypeGridShift

	} else if op.ProjString.Args.ContainsKey("catalog") {
		op.DatumType = DatumTypeGridShift
		catalogName, ok := op.ProjString.Args.GetAsString("catalog")
		if !ok {
			return merror.New(merror.BadProjStringError)
		}
		op.CatalogName = catalogName
		datumDate, ok := op.ProjString.Args.GetAsString("sdate")
		if datumDate != "" {
			op.DatumDate = support.ParseDate(datumDate)
		}

	} else if op.ProjString.Args.ContainsKey("towgs84") {

		values, ok := op.ProjString.Args.GetAsFloats("towgs84")
		if !ok {
			return merror.New(merror.BadProjStringError)
		}

		if len(values) == 3 {
			op.DatumType = DatumType3Param

			op.DatumParams[0] = values[0]
			op.DatumParams[1] = values[1]
			op.DatumParams[2] = values[2]

		} else if len(values) == 7 {
			op.DatumType = DatumType7Param

			op.DatumParams[0] = values[0]
			op.DatumParams[1] = values[1]
			op.DatumParams[2] = values[2]
			op.DatumParams[3] = values[3]
			op.DatumParams[4] = values[4]
			op.DatumParams[5] = values[5]
			op.DatumParams[6] = values[6]

			// transform from arc seconds to radians
			op.DatumParams[3] = support.ConvertArcsecondsToRadians(op.DatumParams[3])
			op.DatumParams[4] = support.ConvertArcsecondsToRadians(op.DatumParams[4])
			op.DatumParams[5] = support.ConvertArcsecondsToRadians(op.DatumParams[5])

			// transform from parts per million to scaling factor
			op.DatumParams[6] = (op.DatumParams[6] / 1000000.0) + 1

			/* Note that pj_init() will later switch datum_type to
			   PJD_WGS84 if shifts are all zero, and ellipsoid is WGS84 or GRS80 */
		} else {
			return merror.New(merror.BadProjStringError)
		}
	}

	return nil
}

func (op *Operation) processEllipsoid() error {

	ellipsoid, err := NewEllipsoid(op)
	if err != nil {

		/* Didn't get an ellps, but doesn't need one: Get a free WGS84 */
		if op.NeedEllps {
			return merror.New(merror.BadProjStringError)
		}

		ellipsoid = &Ellipsoid{}
		ellipsoid.F = 1.0 / 298.257223563
		ellipsoid.AOrig = 6378137.0
		ellipsoid.A = 6378137.0
		ellipsoid.EsOrig = ellipsoid.F * (2 - ellipsoid.F)
		ellipsoid.Es = ellipsoid.F * (2 - ellipsoid.F)
	}

	ellipsoid.AOrig = ellipsoid.A
	ellipsoid.EsOrig = ellipsoid.Es

	err = ellipsoid.doCalcParams(ellipsoid.A, ellipsoid.Es)
	if err != nil {
		return err
	}

	op.Ellipsoid = ellipsoid

	return nil
}

func (op *Operation) readUnits(vertical bool) (float64, float64, error) {

	units := "units"
	toMeter := "toMeter"

	var to, from float64

	if vertical {
		units = "v" + units
		toMeter = "v" + toMeter
	}

	name, ok := op.ProjString.Args.GetAsString(units)
	var s string
	if ok {
		u, ok := UnitInfoTable[name]
		if !ok {
			return 0.0, 0.0, merror.New(merror.ErrUnknownUnit)
		}
		s = u.ToMetersS
	}

	if op.ProjString.Args.ContainsKey(toMeter) {
		s, _ = op.ProjString.Args.GetAsString(toMeter)
	}

	if s != "" {
		var factor float64
		var ratio = false

		/* ratio number? */
		if len(s) > 1 && s[0:1] == "1" && s[1:2] == "/" {
			ratio = true
			s = s[2:]
		}

		factor, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0.0, 0.0, merror.New(merror.BadProjStringError)
		}
		if (factor <= 0.0) || (1.0/factor == 0.0) {
			return 0.0, 0.0, merror.New(merror.ErrUnitFactorLessThanZero)
		}

		if ratio {
			to = 1.0 / factor
		} else {
			to = factor
		}

		from = 1.0 / op.FromMeter
	} else {
		to = 1.0
		from = 1.0
	}

	return to, from, nil
}

func (op *Operation) processMisc() error {

	/* Set PIN->geoc coordinate system */
	op.Geoc = (op.Ellipsoid.Es != 0.0 && op.ProjString.Args.ContainsKey("geoc"))

	/* Over-ranging flag */
	op.Over = op.ProjString.Args.ContainsKey("over")

	/* Vertical datum geoid grids */
	op.HasGeoidVgrids = op.ProjString.Args.ContainsKey("geoidgrids")

	/* Longitude center for wrapping */
	op.IsLongWrapSet = op.ProjString.Args.ContainsKey("lon_wrap")
	if op.IsLongWrapSet {
		op.LongWrapCenter, _ = op.ProjString.Args.GetAsFloat("lon_wrap")
		/* Don't accept excessive values otherwise we might perform badly */
		/* when correcting longitudes around it */
		/* The test is written this way to error on long_wrap_center "=" NaN */
		if !(math.Abs(op.LongWrapCenter) < 10.0*support.TwoPi) {
			return merror.New(merror.ErrLatOrLonExceededLimit)
		}
	}

	/* Axis orientation */
	if op.ProjString.Args.ContainsKey("axis") {
		axisLegal := "ewnsud"
		axisArg, _ := op.ProjString.Args.GetAsString("axis")
		if len(axisArg) != 3 {
			return merror.New(merror.ErrAxis)
		}

		if !strings.ContainsAny(axisArg[0:1], axisLegal) ||
			!strings.ContainsAny(axisArg[1:2], axisLegal) ||
			!strings.ContainsAny(axisArg[2:3], axisLegal) {
			return merror.New(merror.ErrAxis)
		}

		/* TODO: it would be nice to validate we don't have on axis repeated */
		op.Axis = axisArg
	}

	/* Central meridian */
	f, ok := op.ProjString.Args.GetAsFloat("lon_0")
	if ok {
		op.Lam0 = f
	}

	/* Central latitude */
	f, ok = op.ProjString.Args.GetAsFloat("lat_0")
	if ok {
		op.Phi0 = f
	}

	/* False easting and northing */
	f, ok = op.ProjString.Args.GetAsFloat("x_0")
	if ok {
		op.X0 = f
	}
	f, ok = op.ProjString.Args.GetAsFloat("y_0")
	if ok {
		op.Y0 = f
	}
	f, ok = op.ProjString.Args.GetAsFloat("z_0")
	if ok {
		op.Z0 = f
	}
	f, ok = op.ProjString.Args.GetAsFloat("t_0")
	if ok {
		op.T0 = f
	}

	/* General scaling factor */
	if op.ProjString.Args.ContainsKey("k_0") {
		op.K0, _ = op.ProjString.Args.GetAsFloat("k_0")
	} else if op.ProjString.Args.ContainsKey("k") {
		op.K0, _ = op.ProjString.Args.GetAsFloat("k")
	} else {
		op.K0 = 1.0
	}
	if op.K0 <= 0.0 {
		return merror.New(merror.ErrKLessThanZero)
	}

	/* Set units */
	to, from, err := op.readUnits(false)
	if err != nil {
		return err
	}
	op.ToMeter = to
	op.FromMeter = from

	to, from, err = op.readUnits(true)
	if err != nil {
		return err
	}
	op.VToMeter = to
	op.VFromMeter = from

	/* Prime meridian */
	name, ok := op.ProjString.Args.GetAsString("pm")
	if ok {
		var value string
		pm, ok := PrimeMeridianTable[name]
		if ok {
			value = pm.Definition
		} else {
			value = name
		}
		f, err = support.DMSToR(value)
		if err != nil {
			return err
		}
		op.FromGreenwich = f
	} else {
		op.FromGreenwich = 0.0
	}

	// TODO: geod_init(PIN->geod, PIN->a,  (1 - sqrt (1 - PIN->es)));

	return nil
}

func (op *Operation) forwardPrepare(lp *CoordLP) (*CoordLP, error) {

	if math.MaxFloat64 == lp.Lam {
		return nil, merror.New(merror.ErrCoordinateError)
	}

	// /* The helmert datum shift will choke unless it gets a sensible 4D coordinate */
	// if (HUGE_VAL==coo.v[2] && P->helmert) coo.v[2] = 0.0;
	// if (HUGE_VAL==coo.v[3] && P->helmert) coo.v[3] = 0.0;

	/* Check validity of angular input coordinates */
	if op.Left == IOUnitsAngular {

		/* check for latitude or longitude over-range */
		var t float64
		if lp.Phi < 0 {
			t = -lp.Phi - support.PiOverTwo
		} else {
			t = lp.Phi - support.PiOverTwo
		}
		if t > epsLat || lp.Lam > 10 || lp.Lam < -10 {
			return nil, merror.New(merror.ErrLatOrLonExceededLimit)
		}

		/* Clamp latitude to -90..90 degree range */
		if lp.Phi > support.PiOverTwo {
			lp.Phi = support.PiOverTwo
		}
		if lp.Phi < -support.PiOverTwo {
			lp.Phi = -support.PiOverTwo
		}

		/* If input latitude is geocentrical, convert to geographical */
		if op.Geoc {
			lp = GeocentricLatitude(op, DirectionInverse, lp)
		}

		/* Ensure longitude is in the -pi:pi range */
		if !op.Over {
			lp.Lam = support.Adjlon(lp.Lam)
		}

		//if (P->hgridshift) {
		//	coo = proj_trans (P->hgridshift, PJ_INV, coo);
		//}        else if (P->helmert || (P->cart_wgs84 != 0 && P->cart != 0)) {
		//    coo = proj_trans (P->cart_wgs84, PJ_FWD, coo); /* Go cartesian in WGS84 frame */
		//    if( P->helmert ) {
		//        coo = proj_trans (P->helmert,    PJ_INV, coo); /* Step into local frame */
		//	}
		//	coo = proj_trans (P->cart,       PJ_INV, coo); /* Go back to angular using local ellps */
		//}
		if lp.Lam == math.MaxFloat64 {
			return lp, nil
		}
		//if (P->vgridshift){
		//    coo = proj_trans (P->vgridshift, PJ_FWD, coo); /* Go orthometric from geometric */
		//}

		/* Distance from central meridian, taking system zero meridian into account */
		lp.Lam = (lp.Lam - op.FromGreenwich) - op.Lam0

		/* Ensure longitude is in the -pi:pi range */
		if !op.Over {
			lp.Lam = support.Adjlon(lp.Lam)
		}

		return lp, nil
	}

	/* We do not support gridshifts on cartesian input */
	//if (op.Left==IOUnitsCartesian && P->helmert) {
	//		return proj_trans (P->helmert, PJ_INV, coo);
	//}
	return lp, nil
}

func (op *Operation) forwardFinalize(coo *CoordXY) (*CoordXY, error) {

	switch op.Right {

	/* Handle false eastings/northings and non-metric linear units */
	case IOUnitsCartesian:
		panic(99)

		//if op.IsGeocentric {
		//    coo = proj_trans (P->cart, PJ_FWD, coo);
		//}

		//coo.X = op.FromMeter * (coo.X + P.X0);
		//coo.Y = op.FromMeter * (coo.Y + P.Y0);
		//coo.X = op.FromMeter * (coo.Z + P.Z0);

	/* Classic proj.4 functions return plane coordinates in units of the semimajor axis */
	case IOUnitsClassic:
		coo.X *= op.Ellipsoid.A
		coo.Y *= op.Ellipsoid.A
		fallthrough

	/* Falls through */ /* (<-- GCC warning silencer) */
	/* to continue processing in common with PJ_IO_UNITS_PROJECTED */
	case IOUnitsProjected:
		coo.X = op.FromMeter * (coo.X + op.X0)
		coo.Y = op.FromMeter * (coo.Y + op.Y0)
		///////////////////coo.Z = op.VFromMeter * (coo.Z + op.Z0)

	case IOUnitsWhatever:

	case IOUnitsAngular:
		panic(99)
		//if op.Left == IOUnitsAngular {
		//	break;
		//}

		/* adjust longitude to central meridian */
		//if !op.Over {
		//    coo.lpz.lam = adjlon(coo.lpz.lam);
		//}

		//if (P->vgridshift)
		//    coo = proj_trans (P->vgridshift, PJ_FWD, coo); /* Go orthometric from geometric */
		//if (coo.lp.lam==HUGE_VAL) {
		//	return coo;
		//}

		//if (P->hgridshift)
		//    coo = proj_trans (P->hgridshift, PJ_INV, coo);
		//else if (P->helmert || (P->cart_wgs84 != 0 && P->cart != 0)) {
		//    coo = proj_trans (P->cart_wgs84, PJ_FWD, coo); /* Go cartesian in WGS84 frame */
		//    if( P->helmert )
		//        coo = proj_trans (P->helmert,    PJ_INV, coo); /* Step into local frame */
		//    coo = proj_trans (P->cart,       PJ_INV, coo); /* Go back to angular using local ellps */
		//}
		//if (coo.lp.lam==HUGE_VAL) {
		//	return coo;
		//}

		/* If input latitude was geocentrical, convert back to geocentrical */
		//if op.Geoc {
		//	coo = GeocentricLatitude(op, DirectionForward, coo)
		//}

		/* Distance from central meridian, taking system zero meridian into account */
		//coo.lp.lam = coo.lp.lam + P->from_greenwich + P->lam0;

		/* adjust longitude to central meridian */
		//if (0==P->over) {
		//	coo.lpz.lam = adjlon(coo.lpz.lam);
		//}
	}

	//if (P->axisswap) {
	//    coo = proj_trans (P->axisswap, PJ_FWD, coo);
	//}

	return coo, nil
}

func (op *Operation) inversePrepare(coo *CoordXY) (*CoordXY, error) {
	if coo.X == math.MaxFloat64 {
		return nil, merror.New(merror.ErrInvalidXOrY)
	}

	///* The helmert datum shift will choke unless it gets a sensible 4D coordinate */
	//if (HUGE_VAL==coo.v[2] && P->helmert) coo.v[2] = 0.0;
	//if (HUGE_VAL==coo.v[3] && P->helmert) coo.v[3] = 0.0;

	//if (P->axisswap)
	//    coo = proj_trans (P->axisswap, PJ_INV, coo);

	/* Check validity of angular input coordinates */
	//if (INPUT_UNITS==PJ_IO_UNITS_ANGULAR) {
	//    double t;
	//
	//    /* check for latitude or longitude over-range */
	//    t = (coo.lp.phi < 0  ?  -coo.lp.phi  :  coo.lp.phi) - M_HALFPI;
	//    if (t > PJ_EPS_LAT  ||  coo.lp.lam > 10  ||  coo.lp.lam < -10) {
	//        proj_errno_set (P, PJD_ERR_LAT_OR_LON_EXCEED_LIMIT);
	//        return proj_coord_error ();
	//    }
	//
	//    /* Clamp latitude to -90..90 degree range */
	//    if (coo.lp.phi > M_HALFPI)
	//        coo.lp.phi = M_HALFPI;
	//    if (coo.lp.phi < -M_HALFPI)
	//        coo.lp.phi = -M_HALFPI;
	//
	//    /* If input latitude is geocentrical, convert to geographical */
	//    if (P->geoc)
	//        coo = proj_geocentric_latitude (P, PJ_INV, coo);
	//
	// /* Distance from central meridian, taking system zero meridian into account */
	//    coo.lp.lam = (coo.lp.lam + P->from_greenwich) - P->lam0;
	//
	//    /* Ensure longitude is in the -pi:pi range */
	//    if (0==P->over)
	//        coo.lp.lam = adjlon(coo.lp.lam);
	//
	//    if (P->hgridshift)
	//        coo = proj_trans (P->hgridshift, PJ_FWD, coo);
	//    else if (P->helmert || (P->cart_wgs84 != 0 && P->cart != 0)) {
	//        coo = proj_trans (P->cart,       PJ_FWD, coo); /* Go cartesian in local frame */
	//        if( P->helmert )
	//            coo = proj_trans (P->helmert,    PJ_FWD, coo); /* Step into WGS84 */
	//        coo = proj_trans (P->cart_wgs84, PJ_INV, coo); /* Go back to angular using WGS84 ellps */
	//    }
	//    if (coo.lp.lam==HUGE_VAL)
	//        return coo;
	//    if (P->vgridshift)
	//        coo = proj_trans (P->vgridshift, PJ_INV, coo); /* Go geometric from orthometric */
	//    return coo;
	//}

	/* Handle remaining possible input types */
	switch op.Left {

	case IOUnitsWhatever:
		return coo, nil

		/* de-scale and de-offset */
	case IOUnitsCartesian:
		coo.X = op.ToMeter*coo.X - op.X0
		coo.Y = op.ToMeter*coo.Y - op.Y0
		/////////////coo.Z = op.ToMeter*coo.Z - op.Z0

		//if (P->is_geocent)
		//    coo = proj_trans (P->cart, PJ_INV, coo);

		return coo, nil

	case IOUnitsProjected, IOUnitsClassic:
		coo.X = op.ToMeter*coo.X - op.X0
		coo.Y = op.ToMeter*coo.Y - op.Y0
		///////////coo.Z = op.VToMeter*coo.Z - op.Z0
		if op.Left == IOUnitsProjected {
			return coo, nil
		}

		/* Classic proj.4 functions expect plane coordinates in units of the semimajor axis  */
		/* Multiplying by ra, rather than dividing by a because the CalCOFI projection       */
		/* stomps on a and hence (apparently) depends on this to roundtrip correctly         */
		/* (CalCOFI avoids further scaling by stomping - but a better solution is possible)  */
		coo.X *= op.Ellipsoid.Ra
		coo.Y *= op.Ellipsoid.Ra
		return coo, nil
	}

	/* Should not happen, so we could return pj_coord_err here */
	return coo, nil
}

func (op *Operation) inverseFinalize(coo *CoordLP) (*CoordLP, error) {
	//if (coo.xyz.x == HUGE_VAL) {
	//    proj_errno_set (P, PJD_ERR_INVALID_X_OR_Y);
	//    return proj_coord_error ();
	//}

	if op.Right == IOUnitsAngular {

		if op.Left != IOUnitsAngular {
			/* Distance from central meridian, taking system zero meridian into account */
			coo.Lam = coo.Lam + op.FromGreenwich + op.Lam0

			/* adjust longitude to central meridian */
			if !op.Over {
				coo.Lam = support.Adjlon(coo.Lam)
			}

			//if (P->vgridshift)
			//    coo = proj_trans (P->vgridshift, PJ_INV, coo); /* Go geometric from orthometric */
			//if (coo.lp.lam==HUGE_VAL)
			//    return coo;
			//if (P->hgridshift)
			//    coo = proj_trans (P->hgridshift, PJ_FWD, coo);
			//else if (P->helmert || (P->cart_wgs84 != 0 && P->cart != 0)) {
			//    coo = proj_trans (P->cart,       PJ_FWD, coo); /* Go cartesian in local frame */
			//    if( P->helmert )
			//        coo = proj_trans (P->helmert,    PJ_FWD, coo); /* Step into WGS84 */
			//    coo = proj_trans (P->cart_wgs84, PJ_INV, coo); /* Go back to angular using WGS84 ellps */
			//}
			if coo.Lam == math.MaxFloat64 {
				return coo, nil
			}
		}

		/* If input latitude was geocentrical, convert back to geocentrical */
		if op.Geoc {
			coo = GeocentricLatitude(op, DirectionForward, coo)
		}
	}

	return coo, nil
}

// GeocentricLatitude converts geographical latitude to geocentric
// or the other way round if direction = PJ_INV
func GeocentricLatitude(op *Operation, direction DirectionType, lp *CoordLP) *CoordLP {
	/**************************************************************************************

		The conversion involves a call to the tangent function, which goes through the
		roof at the poles, so very close (the last centimeter) to the poles no
		conversion takes place and the input latitude is copied directly to the output.

		Fortunately, the geocentric latitude converges to the geographical at the
		poles, so the difference is negligible.

		For the spherical case, the geographical latitude equals the geocentric, and
		consequently, the input is copied directly to the output.
	**************************************************************************************/
	const limit = support.PiOverTwo - 1e-9

	res := lp

	if (lp.Phi > limit) || (lp.Phi < -limit) || (op.Ellipsoid.Es == 0) {
		return res
	}
	if direction == DirectionForward {
		res.Phi = math.Atan(op.Ellipsoid.OneEs * math.Tan(lp.Phi))
	} else {
		res.Phi = math.Atan(op.Ellipsoid.ROneEs * math.Tan(lp.Phi))
	}

	return res
}
