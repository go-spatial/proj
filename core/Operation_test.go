// Copyright (C) 2018, Michael P. Gerlek (Flaxen Consulting)
//
// Portions of this code were derived from the PROJ.4 software
// In keeping with the terms of the PROJ.4 project, this software
// is provided under the MIT-style license in `LICENSE.md` and may
// additionally be subject to the copyrights of the PROJ.4 authors.

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
