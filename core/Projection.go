package core

// DatumType is the enum for the types of datums we support
type DatumType string

// All the DatumType constants
const (
	DatumTypeUnknown DatumType = "UNKNOWN"
	DatumType3Param            = "3PARAM"
	DatumType7Param            = "7PARAM"
)

// Projection contains all the info needed to describe a projection.
type Projection struct {
	IsLatLong      bool
	IsGeocentric   bool
	IsLongWrapSet  bool
	LongWrapCenter float64
	Axis           string
	DatumType      DatumType
	DatumParams    [7]float64

	// gridlist, gridlist_count
	// vgridlist_geoid, vgtridlist_geoid_count

	DefSize           string
	DefShape          string
	DefSpherification string
	DefEllps          string

	// ELLIPSOID

	/* The linear parameters */

	a  float64 /* semimajor axis (radius if eccentricity==0) */
	b  float64 /* semiminor axis */
	ra float64 /* 1/a */
	rb float64 /* 1/b */

	/* The eccentricities */

	alpha  float64 /* angular eccentricity */
	e      float64 /* first  eccentricity */
	es     float64 /* first  eccentricity squared */
	e2     float64 /* second eccentricity */
	e2s    float64 /* second eccentricity squared */
	e3     float64 /* third  eccentricity */
	e3s    float64 /* third  eccentricity squared */
	oneEs  float64 /* 1 - e^2 */
	rOneEs float64 /* 1/one_es */

	/* The flattenings */
	f   float64 /* first  flattening */
	f2  float64 /* second flattening */
	n   float64 /* third  flattening */
	rf  float64 /* 1/f  */
	rf2 float64 /* 1/f2 */
	rn  float64 /* 1/n  */

	/* This one's for GRS80 */
	J float64 /* "Dynamic form factor" */

	esOrig, aOrig float64 /* es and a before any +proj related adjustment */

}

// NewProjection returns a new Projection object
func NewProjection() (*Projection, error) {
	p := &Projection{}

	return p, nil
}
