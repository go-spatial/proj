// Copyright (C) 2018, Michael P. Gerlek (Flaxen Consulting)
//
// Portions of this code were derived from the PROJ.4 software
// In keeping with the terms of the PROJ.4 project, this software
// is provided under the MIT-style license in `LICENSE.md` and may
// additionally be subject to the copyrights of the PROJ.4 authors.

package gie_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/go-spatial/proj/gie"
	"github.com/stretchr/testify/assert"
)

// must be run from the "./proj/gie" directory, so it can access the "gie_data" directory
func TestGie(t *testing.T) {
	assert := assert.New(t)

	g, err := gie.NewGie("./gie_data")
	assert.NoError(err)

	err = g.Parse()
	assert.NoError(err)

	total := 0
	actual := 0
	passed := 0
	failed := 0

	for _, command := range g.Commands {
		total++
		tag := fmt.Sprintf("%s:%d", command.File, command.Line)

		if g.IsSupported(command) {
			actual++

			err = command.Execute()
			assert.NoError(err, tag)

			if err != nil {
				failed++
			} else {
				passed++
			}
		}
	}

	log.Printf("total:  %d", total)
	log.Printf("actual: %d", actual)
	log.Printf("passed: %d", passed)
	log.Printf("failed: %d", failed)
}
