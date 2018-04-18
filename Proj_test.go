package proj_test

import (
	"testing"

	"github.com/go-spatial/proj"
	"github.com/stretchr/testify/assert"
)

func Test3395(t *testing.T) {
	assert := assert.New(t)

	var err error

	source := 4326
	dest := 3395
	inputA := []float64{2, 1, 10, 20, 30, 40}
	inputB := []float64{2, 1}

	var output1A, output1B, output2A, output2B []float64

	// *** first way
	{
		output1A, err = proj.Project(source, dest, inputA)
		assert.NoError(err)

		output1B, err = proj.Project(source, dest, inputB)
		assert.NoError(err)
	}

	// *** second way, which doesn't need to build the coordinate system object a second time
	{
		proj, err := proj.New(source, dest)
		assert.NoError(err)

		output2A, err = proj.Project(inputA)
		assert.NoError(err)

		output2B, err = proj.Project(inputB)
		assert.NoError(err)
	}

	expectedA := []float64{
		222638.98, 110579.97,
		1113194.91, 2258423.65,
		3339584.72, 4838471.40,
	}
	expectedB := []float64{
		222638.98, 110579.97,
	}

	for i := range expectedA {
		assert.InDelta(expectedA[i], output1A[i], 1e-2)
		assert.InDelta(expectedA[i], output2A[i], 1e-2)
	}
	for i := range expectedB {
		assert.InDelta(expectedB[i], output1B[i], 1e-2)
		assert.InDelta(expectedB[i], output2B[i], 1e-2)
	}
}

func Test3857(t *testing.T) {
	// TODO
}
