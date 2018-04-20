package core_test

import (
	"testing"

	"github.com/go-spatial/proj/support"

	"github.com/go-spatial/proj/core"
	"github.com/stretchr/testify/assert"

	// need to pull in the operations table entries
	_ "github.com/go-spatial/proj/operations"
)

func TestUtm(t *testing.T) {
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
	assert.InDelta(12.0, l, 1e-8)
	assert.InDelta(55.0, p, 1e-8)
}

func TestEtMerc(t *testing.T) {
	assert := assert.New(t)

	ps, err := support.NewProjString("+proj=etmerc +ellps=GRS80 +lat_1=0.5 +lat_2=2 +n=0.5 +zone=30")
	assert.NoError(err)

	sys, opx, err := core.NewSystem(ps)
	assert.NoError(err)
	assert.NotNil(sys)
	assert.NotNil(opx)
	assert.EqualValues(sys, opx.GetSystem())

	assert.True(opx.GetDescription().IsConvertLPToXY())
	op := opx.(core.IConvertLPToXY)

	input := &core.CoordLP{Lam: support.DDToR(2.0), Phi: support.DDToR(1.0)}
	output, err := op.Forward(input)
	assert.NoError(err)

	x, y := output.X, output.Y
	assert.InDelta(222650.796797586, x, 1e-8)
	assert.InDelta(110642.229411933, y, 1e-8)

	input2 := output
	output2, err := op.Inverse(input2)
	assert.NoError(err)

	l, p := output2.Lam, output2.Phi
	l = support.RToDD(l)
	p = support.RToDD(p)
	assert.InDelta(2.0, l, 1e-8)
	assert.InDelta(1.0, p, 1e-8)
}

func Test3395(t *testing.T) {

	assert := assert.New(t)

	ps, err := support.NewProjString("+proj=merc +lon_0=0 +k=1 +x_0=0 +y_0=0 +ellps=WGS84 +datum=WGS84 +units=m +no_defs")
	assert.NoError(err)

	t.Skip()

	sys, opx, err := core.NewSystem(ps)
	assert.NoError(err)
	assert.NotNil(sys)
	assert.NotNil(opx)
	assert.EqualValues(sys, opx.GetSystem())
}
