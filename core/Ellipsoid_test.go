package core_test

import (
	"testing"

	"github.com/go-spatial/proj4go/core"
	"github.com/stretchr/testify/assert"
)

func TestEllipsoidInit(t *testing.T) {
	assert := assert.New(t)

	assert.True(len(core.EllipsoidTable) > 10)

	assert.Equal("danish", core.EllipsoidTable["danish"].ID)
	assert.Equal("a=6377019.2563", core.EllipsoidTable["danish"].Major)
	assert.Equal("rf=300.0", core.EllipsoidTable["danish"].Ell)
	assert.Equal("Andrae 1876 (Denmark, Iceland)", core.EllipsoidTable["danish"].Name)
}
