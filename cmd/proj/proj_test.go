package main_test

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"

	main "github.com/go-spatial/proj/cmd/proj"
	"github.com/stretchr/testify/assert"
)

func TestCmd(t *testing.T) {
	assert := assert.New(t)

	type testcase struct {
		args   string
		input  []float64
		output []float64
	}

	testcases := []testcase{
		{
			"proj -epsg 9999",
			[]float64{0.0, 0.0},
			nil,
		}, {
			"proj -epsg 9999 proj=merc",
			[]float64{0.0, 0.0},
			nil,
		}, {
			"proj-epsg 9999 -inverse",
			[]float64{0.0, 0.0},
			nil,
		}, {
			"proj -epsg 3395",
			[]float64{-77.625583, 38.833846},
			[]float64{-8641240.37, 4671101.60},
		}, {
			"proj +proj=utm +zone=32 +ellps=GRS80",
			[]float64{12.0, 55.0},
			[]float64{691875.63, 6098907.83},
		},
	}

	for _, tc := range testcases {

		s := fmt.Sprintf("%f %f", tc.input[0], tc.input[1])
		inBuf := bytes.NewBuffer([]byte(s))
		outBuf := &bytes.Buffer{}

		err := main.Main(inBuf, outBuf, strings.Fields(tc.args))

		if tc.output == nil {
			assert.Error(err, tc.args)
		} else {
			assert.NoError(err, tc.args)

			tokens := strings.Fields(string(outBuf.Bytes()))
			assert.Len(tokens, 2, tc.args)
			actual0, err := strconv.ParseFloat(tokens[0], 64)
			assert.NoError(err, tc.args)
			assert.InDelta(tc.output[0], actual0, 1.0e-2, tc.args)
			actual1, err := strconv.ParseFloat(tokens[1], 64)
			assert.NoError(err, tc.args)
			assert.InDelta(tc.output[1], actual1, 1.0e-2, tc.args)
		}
	}
}
