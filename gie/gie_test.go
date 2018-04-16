package gie_test

import (
	"testing"

	"github.com/go-spatial/proj4go/gie"
	"github.com/stretchr/testify/assert"
)

func TestGie(t *testing.T) {
	assert := assert.New(t)

	files := []string{
		"data/4D-API_cs2cs-style.gie",
		"data/DHDN_ETRS89.gie",
		"data/GDA.gie",
		"data/axisswap.gie",
		"data/builtins.gie",
		"data/deformation.gie",
		"data/ellipsoid.gie",
		"data/more_builtins.gie",
		"data/unitconvert.gie",
	}
	for _, f := range files {
		p, err := gie.NewParser(f)
		assert.NoError(err)
		assert.NotNil(p)
	}
}
