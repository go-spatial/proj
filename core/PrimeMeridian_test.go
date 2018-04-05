package core_test

import (
	"testing"

	"github.com/go-spatial/proj4go/core"
	"github.com/stretchr/testify/assert"
)

func TestPrimeMeridianInit(t *testing.T) {
	assert := assert.New(t)

	assert.True(len(core.PrimeMeridianTable) > 10)

	assert.Equal("paris", core.PrimeMeridianTable["paris"].ID)
	assert.Equal("2d20'14.025\"E", core.PrimeMeridianTable["paris"].Definition)
}
