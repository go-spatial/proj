package core_test

import (
	"testing"

	"github.com/go-spatial/proj4go/core"
	"github.com/stretchr/testify/assert"
)

func TestProjection(t *testing.T) {
	assert := assert.New(t)

	p, err := core.NewProjection()
	assert.NoError(err)
	assert.NotNil(p)
}
