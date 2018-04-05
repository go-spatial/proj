package support_test

import (
	"math"
	"testing"

	"github.com/go-spatial/proj4go/support"
	"github.com/stretchr/testify/assert"
)

func TestDMSToDD(t *testing.T) {
	assert := assert.New(t)

	convert := func(d float64, m float64, s float64) float64 {
		dd := d + m/60.0 + s/3600.0
		return dd
	}

	assert.InDelta(15, convert(15, 0, 0), 1e-8)
	assert.InDelta(137, convert(135, 60, 3600), 1e-8)

	type Data struct {
		input    string
		expected float64
	}

	FAIL := math.MaxFloat32
	ok := convert(12, 34, 56.78)

	data := []Data{
		{`0.0d0m0s`, FAIL},
		{`0d0.0m0s`, FAIL},
		{`.0d0m0s`, FAIL},
		{`0d.0m0s`, FAIL},
		{`0d0m.0s`, FAIL},
		{`0d0m`, FAIL},
		{`0d`, FAIL},
		{`0d0m0sx`, FAIL},

		{`12d34m56.78s`, ok},
		{`+12d34m56.78s`, ok},
		{`-12d34m56.78s`, -ok},

		{`12d34m56.78se`, ok},
		{`12d34m56.78sE`, ok},
		{`12d34m56.78sw`, -ok},
		{`12d34m56.78sW`, -ok},
		{`12d34m56.78sn`, ok},
		{`12d34m56.78sN`, ok},
		{`12d34m56.78ss`, -ok},
		{`12d34m56.78sS`, -ok},

		{`12D34M56.78S`, ok},
		{`12°34'56.78"`, ok},

		{` 12 ° 34 ' 56.78 " e`, ok},

		{`15d0m0s`, convert(15, 0, 0)},
		{`135d0m0s`, convert(135, 0, 0)},
		{`134d60m0s`, convert(134, 60, 0)},
		{`134d0m3600s`, convert(134, 0, 3600)},
		{`134d1m3.3s`, convert(134, 1, 3.3)},
		{`-134d1'3.3s E`, -convert(134, 1, 3.3)},
		{`-134d1'3.3s W`, convert(134, 1, 3.3)},
	}

	for _, d := range data {
		actual, err := support.DMSToDD(d.input)
		if d.expected == FAIL {
			assert.Error(err)
		} else {
			assert.NoError(err)
			assert.InDelta(d.expected, actual, 1e-6, "Test %s", d.input)
		}
	}
}

func TestDDToR(t *testing.T) {
	assert := assert.New(t)

	r := support.DDToR(15.0)
	assert.InDelta(math.Pi/12.0, r, 1e-6)

	r = support.DDToR(135.0)
	assert.InDelta(math.Pi*0.75, r, 1e-6)
}

func TestDMSToR(t *testing.T) {
	assert := assert.New(t)

	r, err := support.DMSToR("15d0m0s")
	assert.NoError(err)
	assert.InDelta(math.Pi/12.0, r, 1e-6)

	r, err = support.DMSToR("135d0m0s")
	assert.NoError(err)
	assert.InDelta(math.Pi*0.75, r, 1e-6)
}

func TestArcseconds(t *testing.T) {
	assert := assert.New(t)

	r := support.ConvertArcsecondsToRadians(15.0 * 3600.0)
	assert.InDelta(math.Pi/12.0, r, 1e-6)
}
