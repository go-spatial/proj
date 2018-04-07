package core_test

import (
	"testing"

	"github.com/go-spatial/proj4go/core"
	"github.com/go-spatial/proj4go/mlog"
	"github.com/stretchr/testify/assert"
)

func TestFull(t *testing.T) {
	assert := assert.New(t)

	ps, err := core.NewProjString("+proj=utm +zone=32 +ellps=GRS80")
	assert.NoError(err)
	mlog.Printf("%s", ps)

	op, err := core.NewOperation(ps)
	assert.NoError(err)
	mlog.Printf("%s", op)
	mlog.Printf("%s", op.Ellipsoid)

	// TODO: convert to radians, the internal format
	// 55d N, 12d E (lon lat) (lam phi)
	input := core.CoordLP{Lam: 12.0, Phi: 55.0}
	t.Skip()
	output, err := op.Forward(input)
	assert.NoError(err)

	e, n := output.(core.CoordENU).E, output.(core.CoordENU).N
	assert.InDelta(691875.63, e, 1e-6)
	assert.InDelta(6098907.83, n, 1e-6)

	output, err = op.Inverse(output)
	l, p := output.(core.CoordLP).Lam, output.(core.CoordLP).Phi
	assert.InDelta(12.0, l, 1e-6)
	assert.InDelta(55.0, p, 1e-6)
}
