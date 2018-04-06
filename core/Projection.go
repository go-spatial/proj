package core

import (
	"encoding/json"
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

// Projection contains all the info needed to describe a projection.
type Projection struct {

	//
	// COORDINATE HANDLING
	//
	over         int  /* Over-range flag */
	geoc         int  /* Geocentric latitude flag */
	IsLatLong    bool /* proj=latlong ... not really a projection at all */
	IsGeocentric bool /* proj=geocent ... not really a projection at all */
	NeedEllps    bool /* 0 for operations that are purely cartesian */

	left  IOUnitsType /* Flags for input/output coordinate types */
	right IOUnitsType

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

// NewProjection returns a new Projection object
func NewProjection() (*Projection, error) {
	p := &Projection{}

	return p, nil
}

// Forward executes an operation
func (*Projection) Forward(interface{}) (interface{}, error) {
	return nil, nil
}

// Inverse executes an operation in reverse
func (*Projection) Inverse(interface{}) (interface{}, error) {
	return nil, nil
}

func (p *Projection) String() string {
	b, err := json.MarshalIndent(p, "", " ")
	if err != nil {
		panic(err)
	}

	return string(b)
}
