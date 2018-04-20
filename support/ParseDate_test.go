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

func TestParseDate(t *testing.T) {
	assert := assert.New(t)

	assert.InDelta(0.0, support.ParseDate("yow"), 1.0e-8)
	assert.InDelta(1999.5, support.ParseDate("1999.50"), 1.0e-8)
	assert.InDelta(2000.0, support.ParseDate("1999.99999999"), 1.0e-8)
	assert.InDelta(1999.0+(12.0*31.0-1.0)/(12.0*31.0), support.ParseDate("1999-12-31"), 1.0e-8)
	assert.InDelta(1999.0+(6.0*31.0)/(12.0*31.0), support.ParseDate("1999-07-01"), 1.0e-8)
}
