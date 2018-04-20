// Copyright (C) 2018, Michael P. Gerlek (Flaxen Consulting)
//
// Portions of this code were derived from the PROJ.4 software
// In keeping with the terms of the PROJ.4 project, this software
// is provided under the MIT-style license in `LICENSE.md` and may
// additionally be subject to the copyrights of the PROJ.4 authors.

package support_test

import (
	"testing"

	"github.com/go-spatial/proj/support"
	"github.com/stretchr/testify/assert"
)

func TestUnitsTable(t *testing.T) {
	assert := assert.New(t)

	assert.True(len(support.UnitsTable) > 15)
	for key, value := range support.UnitsTable {
		assert.Equal(key, value.ID)
	}
}
