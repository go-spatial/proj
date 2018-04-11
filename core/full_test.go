package core_test

import (
	"testing"

	"github.com/go-spatial/proj4go/mlog"
	"github.com/go-spatial/proj4go/support"

	"github.com/go-spatial/proj4go/core"
	_ "github.com/go-spatial/proj4go/operations"
	"github.com/stretchr/testify/assert"
)

func TestFull(t *testing.T) {
	assert := assert.New(t)

	ps, err := core.NewProjString("+proj=utm +zone=32 +ellps=GRS80")
	assert.NoError(err)

	op, err := core.NewOperation(ps)
	assert.NoError(err)

	//mlog.Printv(op)

	// TODO: convert to radians, the internal format
	// 55d N, 12d E (lon lat) (lam phi)
	input := &core.CoordLP{Lam: support.DDToR(12.0), Phi: support.DDToR(55.0)}
	mlog.Printv(input)
	output, err := op.Forward(input)
	assert.NoError(err)

	x, y := output.(*core.CoordXY).X, output.(*core.CoordXY).Y
	assert.InDelta(691875.63, x, 1e-2)
	assert.InDelta(6098907.83, y, 1e-2)

	output, err = op.Inverse(output)
	assert.NoError(err)

	l, p := output.(*core.CoordLP).Lam, output.(*core.CoordLP).Phi
	assert.InDelta(12.0, l, 1e-6)
	assert.InDelta(55.0, p, 1e-6)
}
