// Copyright (C) 2018, Michael P. Gerlek (Flaxen Consulting)
//
// Portions of this code were derived from the PROJ.4 software
// In keeping with the terms of the PROJ.4 project, this software
// is provided under the MIT-style license in `LICENSE.md` and may
// additionally be subject to the copyrights of the PROJ.4 authors.

package core_test

import (
	"testing"

	"github.com/go-spatial/proj/core"
	"github.com/go-spatial/proj/support"
	"github.com/stretchr/testify/assert"
)

func TestSystem(t *testing.T) {
	assert := assert.New(t)

	ps, err := support.NewProjString("+proj=utm +zone=32 +ellps=GRS80")
	assert.NoError(err)
	assert.NotNil(ps)

	sys, op, err := core.NewSystem(ps)
	assert.NoError(err)
	assert.NotNil(sys)
	assert.NotNil(op)

	assert.NotEqual("", sys.String())
}

func TestProjStringValidation(t *testing.T) {
	assert := assert.New(t)

	_, err := support.NewProjString("")
	assert.NoError(err)

	ps, err := support.NewProjString("  +proj=P99 +k1=a   +k2=b    \t  ")
	assert.NoError(err)
	assert.Equal(3, ps.Len())

	err = core.ValidateProjStringContents(ps)
	assert.NoError(err)

	// only 1 "init" allowed
	{
		ps, err = support.NewProjString("init=foo proj=foo init=foo")
		assert.NoError(err)
		err = core.ValidateProjStringContents(ps)
		assert.Error(err)
	}

	// "proj" may not be empty
	{
		ps, err = support.NewProjString("proj")
		assert.NoError(err)
		err = core.ValidateProjStringContents(ps)
		assert.Error(err)
	}
}
