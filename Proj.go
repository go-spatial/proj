package proj

import (
	"fmt"

	"github.com/go-spatial/proj/core"
	"github.com/go-spatial/proj/support"

	// need to pull in the operations table entries
	_ "github.com/go-spatial/proj/operations"
)

// Projector is used to perform a conversion from a lon/lat (degree) system
// to a projected (meters) system
type Projector struct {
	pair       pair
	projString *support.ProjString
	system     *core.System
	operation  core.IOperation
	converter  core.IConvertLPToXY
}

// New creates a Projector object for the given soure and destination systems.
func New(sourceEPSG, destEPSG int) (*Projector, error) {
	pair := pair{sourceEPSG, destEPSG}
	str, ok := knownProjections[pair]
	if !ok {
		return nil, fmt.Errorf("epsg code pair is not a supported projection")
	}

	ps, err := support.NewProjString(str)
	if err != nil {
		return nil, err
	}

	sys, opx, err := core.NewSystem(ps)
	if err != nil {
		return nil, err
	}

	if !opx.GetDescription().IsConvertLPToXY() {
		return nil, fmt.Errorf("projection type is not supported")
	}

	proj := &Projector{
		pair:       pair,
		projString: ps,
		system:     sys,
		operation:  opx,
		converter:  opx.(core.IConvertLPToXY),
	}

	return proj, nil
}

// Project performs the projection on the given input points
func (proj *Projector) Project(input []float64) ([]float64, error) {
	if proj.converter == nil {
		return nil, fmt.Errorf("projector not initialized")
	}

	if len(input)%2 != 0 {
		return nil, fmt.Errorf("input array of lon/lat values must be an even number")
	}

	output := make([]float64, len(input))

	lp := core.CoordLP{}

	for i := 0; i < len(input); i += 2 {
		lp.Lam = support.DDToR(input[i])
		lp.Phi = support.DDToR(input[i+1])
		xy, err := proj.converter.Forward(&lp)
		if err != nil {
			return nil, err
		}
		output[i] = xy.X
		output[i+1] = xy.Y
	}

	return output, nil
}

// Project performs the above functions in one step
func Project(source int, dest int, input []float64) ([]float64, error) {
	proj, err := New(source, dest)
	if err != nil {
		return nil, nil
	}
	return proj.Project(input)
}

//---------------------------------------------------------------------

type pair struct {
	source, dest int
}

type pairTable map[pair]string

var knownProjections = pairTable{
	pair{4326, 3395}: "+proj=merc +lon_0=0 +k=1 +x_0=0 +y_0=0 +datum=WGS84", // TODO: support +units=m +no_defs
	pair{4326, 3857}: "+proj=merc +a=6378137 +b=6378137 +lat_ts=0.0 +lon_0=0.0 +x_0=0.0 +y_0=0 +k=1.0 +units=m +nadgrids=@null +wktext +no_defs",
	// TODO: add (3857,3395)?
}
