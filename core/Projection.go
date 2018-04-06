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

	E *Ellipsoid

	/* This one's for GRS80 */
	J float64 /* "Dynamic form factor" */

	esOrig, aOrig float64 /* es and a before any +proj related adjustment */

}

// NewProjection returns a new Projection object
func NewProjection() (*Projection, error) {
	p := &Projection{}

	return p, nil
}
