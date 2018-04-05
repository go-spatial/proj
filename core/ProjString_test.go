package core_test

import (
	"testing"

	"github.com/go-spatial/proj4go/core"
	"github.com/stretchr/testify/assert"
)

func TestProjString(t *testing.T) {
	assert := assert.New(t)

	ps, err := core.NewProjString("")
	assert.Error(err)

	ps, err = core.NewProjString("  +proj=P99 +k1=a   +k2=b    \t  ")
	assert.NoError(err)
	assert.Equal(3, ps.Args.Len())

	// only 1 "init" allowed
	{
		ps, err = core.NewProjString("init=foo proj=foo init=foo")
		assert.Error(err)
	}

	// "proj" may not be empty
	{
		ps, err = core.NewProjString("proj")
		assert.Error(err)
	}
}
