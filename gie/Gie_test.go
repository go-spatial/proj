package gie_test

import (
	"testing"

	"github.com/go-spatial/proj/gie"

	"github.com/stretchr/testify/assert"
)

func TestGie(t *testing.T) {
	assert := assert.New(t)

	g, err := gie.NewGie(".")
	assert.NoError(err)
	assert.NotNil(g)
}
