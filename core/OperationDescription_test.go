package core_test

import (
	"testing"

	"github.com/go-spatial/proj/core"

	"github.com/stretchr/testify/assert"
)

func TestOperationDescription(t *testing.T) {
	assert := assert.New(t)

	opDesc := core.OperationDescriptionTable["utm"]
	assert.NotNil(opDesc)
}
