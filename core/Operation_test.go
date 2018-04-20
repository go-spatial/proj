package core_test

import (
	"fmt"
	"testing"

	"github.com/go-spatial/proj/core"
	"github.com/go-spatial/proj/support"

	"github.com/stretchr/testify/assert"
)

func TestOperation(t *testing.T) {
	assert := assert.New(t)

	ps, err := support.NewProjString("+proj=utm +zone=32 +ellps=GRS80")
	assert.NoError(err)
	_, opx, err := core.NewSystem(ps)
	assert.NoError(err)

	s := fmt.Sprintf("%s", opx)
	assert.True(len(s) > 1)

	id := opx.GetDescription().ID
	assert.Equal("utm", id)
}
