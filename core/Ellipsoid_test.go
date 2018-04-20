package core_test

import (
	"fmt"
	"testing"

	"github.com/go-spatial/proj/core"
	"github.com/go-spatial/proj/support"

	"github.com/stretchr/testify/assert"
)

func TestEllipsoid(t *testing.T) {
	assert := assert.New(t)

	ps, err := support.NewProjString("+proj=utm +zone=32 +ellps=GRS80")
	assert.NoError(err)
	sys, _, err := core.NewSystem(ps)
	assert.NoError(err)

	e := sys.Ellipsoid

	assert.Equal("GRS80", e.ID)

	s := fmt.Sprintf("%s", e)
	assert.True(len(s) > 1)
}
