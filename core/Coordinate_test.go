package core_test

import (
	"testing"

	"github.com/go-spatial/proj/core"
	"github.com/stretchr/testify/assert"
)

func TestCoordinate(t *testing.T) {
	assert := assert.New(t)

	{
		any := &core.CoordAny{V: [4]float64{1.0, 2.0, 3.0, 4.0}}
		lp := any.ToLP()
		assert.Equal(1.0, lp.Lam)
		assert.Equal(2.0, lp.Phi)
		lp.Lam = 10.0
		lp.Phi = 20.0
		any.FromLP(lp)
		assert.Equal(10.0, any.V[0])
		assert.Equal(20.0, any.V[1])
		assert.Equal(3.0, any.V[2])
		assert.Equal(4.0, any.V[3])
	}
	{
		any := &core.CoordAny{V: [4]float64{1.0, 2.0, 3.0, 4.0}}
		xy := any.ToXY()
		assert.Equal(1.0, xy.X)
		assert.Equal(2.0, xy.Y)
		xy.X = 10.0
		xy.Y = 20.0
		any.FromXY(xy)
		assert.Equal(10.0, any.V[0])
		assert.Equal(20.0, any.V[1])
		assert.Equal(3.0, any.V[2])
		assert.Equal(4.0, any.V[3])
	}
}
