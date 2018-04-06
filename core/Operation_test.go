package core_test

import (
	"testing"

	"github.com/go-spatial/proj4go/core"
	"github.com/stretchr/testify/assert"
)

func TestOperation(t *testing.T) {
	assert := assert.New(t)

	ps, err := core.NewProjString("+proj=utm +zone=32 +ellps=GRS80")
	assert.NoError(err)
	assert.NotNil(ps)

	op, err := core.NewOperation(ps)
	assert.NoError(err)
	assert.NotNil(op)
}
