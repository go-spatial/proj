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
)

func TestOperationDescription(t *testing.T) {

	opDesc := core.OperationDescriptionTable["utm"]
	if opDesc == nil {
		t.Errorf("operaton description table for utm is nil")
	}
}
