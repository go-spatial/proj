package core

import (
	"github.com/go-spatial/proj4go/merror"
	"github.com/go-spatial/proj4go/support"
)

// ProjString represents a projection string, such as "+proj=utm +zone=11 +datum=WGS84"
// TODO: we don't support the "pipeline" or "step" keywords
type ProjString struct {
	Source string
	Args   *support.PairList
}

// NewProjString returns a new ProjString object representing the given string
func NewProjString(source string) (*ProjString, error) {
	ps := &ProjString{
		Source: source,
	}

	pairs, err := support.NewPairListFromString(source)
	if err != nil {
		return nil, err
	}

	ps.Args = pairs

	err = ps.processInit()
	if err != nil {
		return nil, err
	}

	proj, ok := ps.Args.GetAsString("proj")
	if !ok || proj == "" {
		return nil, merror.New(merror.ProjValueMissing)
	}

	_, err = ps.processDatum()
	if err != nil {
		return nil, err
	}

	return ps, nil
}

func (ps *ProjString) processInit() error {

	numInit := ps.Args.CountKey("init")
	if numInit > 1 {
		return merror.New(merror.BadProjStringError)
	}

	// TODO: support "init" expansion
	if numInit != 0 {
		return merror.New(merror.NotYetSupported)
	}
	return nil
}

func (ps *ProjString) processDatum() (*Projection, error) {

	proj, err := NewProjection()
	if err != nil {
		return nil, err
	}

	proj.DatumType = DatumTypeUnknown

	datumName, ok := ps.Args.GetAsString("datum")
	if ok {

		datum, ok := DatumTable[datumName]
		if !ok {
			return nil, merror.New(merror.NoSuchDatum)
		}

		// add the ellipse to the end of the list

		ps.Args.Add(support.Pair{Key: "ellps", Value: datum.EllipseID})
		ps.Args.AddList(datum.Definition)
	}

	_, ok = ps.Args.GetAsString("nadgrids")
	if ok {
		return nil, merror.New(merror.NotYetSupported)
	}

	_, ok = ps.Args.GetAsString("catalog")
	if ok {
		return nil, merror.New(merror.NotYetSupported)
	}

	values, ok := ps.Args.GetAsFloats("towgs84")
	if ok {
		if len(values) == 3 {
			proj.DatumType = DatumType3Param

			proj.DatumParams[0] = values[0]
			proj.DatumParams[1] = values[1]
			proj.DatumParams[2] = values[2]

		} else if len(values) == 7 {
			proj.DatumType = DatumType7Param

			proj.DatumParams[0] = values[0]
			proj.DatumParams[1] = values[1]
			proj.DatumParams[2] = values[2]
			proj.DatumParams[3] = values[3]
			proj.DatumParams[4] = values[4]
			proj.DatumParams[5] = values[5]
			proj.DatumParams[6] = values[6]

			// transform from arc seconds to radians
			proj.DatumParams[3] = support.ConvertArcsecondsToRadians(proj.DatumParams[3])
			proj.DatumParams[4] = support.ConvertArcsecondsToRadians(proj.DatumParams[4])
			proj.DatumParams[5] = support.ConvertArcsecondsToRadians(proj.DatumParams[5])

			// transform from parts per million to scaling factor
			proj.DatumParams[6] = (proj.DatumParams[6] / 1000000.0) + 1

		} else {
			return nil, merror.New(merror.BadProjStringError)
		}

		/* Note that pj_init() will later switch datum_type to
		   PJD_WGS84 if shifts are all zero, and ellipsoid is WGS84 or GRS80 */
	}

	return proj, nil
}
