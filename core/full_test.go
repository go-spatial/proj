package core_test

import (
	"testing"

	"github.com/go-spatial/proj4go/support"

	"github.com/go-spatial/proj4go/core"
	_ "github.com/go-spatial/proj4go/operations"
	"github.com/stretchr/testify/assert"
)

func TestFull(t *testing.T) {
	assert := assert.New(t)

	ps, err := support.NewProjString("+proj=utm +zone=32 +ellps=GRS80")
	assert.NoError(err)

	sys, opx, err := core.NewSystem(ps)
	assert.NoError(err)
	assert.NotNil(sys)
	assert.NotNil(opx)
	assert.EqualValues(sys, opx.GetSystem())

	op := opx.(core.IConvertLPToXY)

	// 55d N, 12d E (lon lat) (lam phi)
	input := &core.CoordLP{Lam: support.DDToR(12.0), Phi: support.DDToR(55.0)}
	output, err := op.Forward(input)
	assert.NoError(err)

	x, y := output.X, output.Y
	assert.InDelta(691875.63, x, 1e-2)
	assert.InDelta(6098907.83, y, 1e-2)

	input2 := output
	output2, err := op.Inverse(input2)
	assert.NoError(err)

	l, p := output2.Lam, output2.Phi
	l = support.RToDD(l)
	p = support.RToDD(p)
	assert.InDelta(12.0, l, 1e-6)
	assert.InDelta(55.0, p, 1e-6)
}
