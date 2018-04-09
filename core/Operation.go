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
	f := op.Info.forward
	output, err := f(op, input)
	return output, err
}

// Inverse executes an operation in reverse
func (op *Operation) Inverse(input interface{}) (interface{}, error) {
	f := op.Info.inverse
	output, err := f(op, input)
	return output, err
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
