package core_test

import (
	"testing"

	"github.com/go-spatial/proj4go/core"
	"github.com/stretchr/testify/assert"
)

func TestFull(t *testing.T) {
	assert := assert.New(t)

	ps, err := core.NewProjString("+proj=utm +zone=32 +ellps=GRS80")
	assert.NoError(err)

	proj := ps.Projection

	// TODO: convert to radians, the internal format
	// 55d N, 12d E
	input := core.CoordLP{12.0, 55.0}

	output, err := proj.Forward(input)
	assert.NoError(err)

	x, y := output.ToDegrees()
	assert.InDelta(32.1, x, 1e-6)
	assert.InDelta(65.4, x, 1e-6)

	output, err = proj.Inverse(output)
	x, y = output.ToDegrees()
	assert.InDelta(12.0, x, 1e-6)
	assert.InDelta(55.0, x, 1e-6)
}
