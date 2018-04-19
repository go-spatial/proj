package proj_test

import (
	"fmt"
	"testing"

	"github.com/go-spatial/proj"
	"github.com/stretchr/testify/assert"
)

var inputA = []float64{
	-0.127758, 51.507351, // London
	2.352222, 48.856614, // Paris
	12.496366, 41.902783, // Rome
}
var inputB = []float64{
	-77.625583, 38.833846, // mpg
}

type testcase struct {
	dest      proj.EPSGCode
	expectedA []float64
	expectedB []float64
}

var testcases = []testcase{
	{
		dest: proj.EPSG3395,
		expectedA: []float64{
			-14221.96, 6678068.96,
			261848.16, 6218371.80,
			1391089.10, 5117883.04,
		},
		expectedB: []float64{
			-8641240.37, 4671101.60,
		},
	},
	{
		dest: proj.EPSG3857,
		expectedA: []float64{
			-14221.96, 6711533.71,
			261848.16, 6250566.72,
			1391089.10, 5146427.91,
		},
		expectedB: []float64{
			-8641240.37, 4697899.31,
		},
	},
}

func TestConvert(t *testing.T) {
	assert := assert.New(t)

	for _, tc := range testcases {

		outputA, err := proj.Convert(tc.dest, inputA)
		assert.NoError(err)

		outputB, err := proj.Convert(tc.dest, inputB)
		assert.NoError(err)

		const tol = 1.0e-2

		for i := range tc.expectedA {
			tag := fmt.Sprintf("epsg:%d, input=A.%d", int(tc.dest), i)
			assert.InDelta(tc.expectedA[i], outputA[i], tol, tag)
			assert.InDelta(tc.expectedA[i], outputA[i], tol, tag)
		}
		for i := range tc.expectedB {
			tag := fmt.Sprintf("epsg:%d, input=B.%d", int(tc.dest), i)
			assert.InDelta(tc.expectedB[i], outputB[i], tol, tag)
			assert.InDelta(tc.expectedB[i], outputB[i], tol, tag)
		}
	}
}
