package core

import (
	"encoding/json"

	"github.com/go-spatial/proj4go/merror"
	"github.com/go-spatial/proj4go/support"
)

// DatumType is the enum for the types of datums we support
type DatumType int

// All the DatumType constants
const (
	DatumTypeUnknown  DatumType = 0
	DatumType3Param             = 1
	DatumType7Param             = 2
	DatumTypeGridShif           = 3
	DatumTypeWGS84              = 4 /* WGS84 (or anything considered equivalent) */
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
	over         int  /* Over-range flag */
	geoc         int  /* Geocentric latitude flag */
	IsLatLong    bool /* proj=latlong ... not really a projection at all */
	IsGeocentric bool /* proj=geocent ... not really a projection at all */
	NeedEllps    bool /* 0 for operations that are purely cartesian */

	Left  IOUnitsType /* Flags for input/output coordinate types */
	Right IOUnitsType

	//
	// ELLIPSOID
	//
	E *Ellipsoid

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

	//int     has_geoid_vgrids;
	//struct _pj_gi **vgridlist_geoid;
	//int     vgridlist_geoid_count;

	//double  from_greenwich;            /* prime meridian offset (in radians) */
	//LongWrapCenter float64         /* 0.0 for -180 to 180, actually in radians*/
	//IsLongWrapSet  bool
	//Axis           string  /* Axis order, pj_transform/pj_adjust_axis */

	/* New Datum Shift Grid Catalogs */
	//char   *catalog_name;
	//struct _PJ_GridCatalog *catalog;

	//double  datum_date;                 /* TODO: Description needed */

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
	}

	return op, nil
}

// Forward executes an operation
func (op *Operation) Forward(input interface{}) (interface{}, error) {
	f := op.Info.Forward
	output, err := f(op, input)
	return output, err
}

// Inverse executes an operation in reverse
func (op *Operation) Inverse(input interface{}) (interface{}, error) {
	f := op.Info.Inverse
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

	err := op.Info.Setup(op)
	if err != nil {
		return err
	}

	err = op.processDatum()
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

		op.ProjString.Args.Add(support.Pair{Key: "ellps", Value: datum.EllipseID})
		op.ProjString.Args.AddList(datum.Definition)
	}

	_, ok = op.ProjString.Args.GetAsString("nadgrids")
	if ok {
		return merror.New(merror.NotYetSupported)
	}

	_, ok = op.ProjString.Args.GetAsString("catalog")
	if ok {
		return merror.New(merror.NotYetSupported)
	}

	values, ok := op.ProjString.Args.GetAsFloats("towgs84")
	if ok {
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

		} else {
			return merror.New(merror.BadProjStringError)
		}

		/* Note that pj_init() will later switch datum_type to
		   PJD_WGS84 if shifts are all zero, and ellipsoid is WGS84 or GRS80 */
	}

	return nil
}
