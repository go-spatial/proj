package api

import (
	"fmt"

	"github.com/go-spatial/proj4go/core"
	_ "github.com/go-spatial/proj4go/operations" // need to pull in the operations table
	"github.com/go-spatial/proj4go/support"
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

// NewProjector creates a Projector object for the given soure and destination systems.
func NewProjector(sourceEPSG, destEPSG int) (*Projector, error) {
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
	proj, err := NewProjector(source, dest)
	if err != nil {
		return nil, nil
	}
	return proj.Project(input)
}
