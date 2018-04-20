// Copyright (C) 2018, Michael P. Gerlek (Flaxen Consulting)
//
// Portions of this code were derived from the PROJ.4 software
// In keeping with the terms of the PROJ.4 project, this software
// is provided under the MIT-style license in `LICENSE.md` and may
// additionally be subject to the copyrights of the PROJ.4 authors.

package support_test

import (
	"math"
	"testing"

	"github.com/go-spatial/proj/support"
	"github.com/stretchr/testify/assert"
)

func TestAasincos(t *testing.T) {
	assert := assert.New(t)

	assert.InDelta(math.Asin(0.5), support.Aasin(0.5), 1.0e-8)
	assert.InDelta(math.Asin(0.5), support.Aasin(0.5), 1.0e-8)
	assert.InDelta(math.Sqrt(0.5), support.Asqrt(0.5), 1.0e-8)
	assert.InDelta(math.Atan2(0.5, 0.5), support.Aatan2(0.5, 0.5), 1.0e-8)
}
