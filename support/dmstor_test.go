package support_test

import (
	"math"
	"testing"

	"github.com/go-spatial/proj4go/support"
	"github.com/stretchr/testify/assert"
)

func TestDMSToR(t *testing.T) {
	assert := assert.New(t)

	convert := func(d float64, m float64, s float64, news string) float64 {
		deg := d + m/60.0 + s/3600.0
		rad := deg * (math.Pi / 180.0)
		if news == "W" || news == "w" || news == "S" || news == "s" {
			rad = -rad
		}
		return rad
	}

	assert.InDelta(math.Pi/12.0, convert(15, 0, 0, ""), 0.00001)
	assert.InDelta(math.Pi*0.75, convert(135, 0, 0, ""), 0.00001)

	type Data struct {
		input    string
		expected float64
	}
	data := []Data{
		{`0d0'0"`, convert(0, 0, 0, "")},
		{`15d0'0"`, convert(15, 0, 0, "")},
		{`135d0'0"`, convert(135, 0, 0, "")},
		{`134d60'0"`, convert(134, 60, 0, "")},
		{`134d0'3600"`, convert(134, 0, 3600, "")},
		{`134.1d1.2'3.3"`, convert(134.1, 1.2, 3.3, "")},
		{`-134.1d1.2'3.3"`, convert(-134.1, 1.2, 3.3, "")},
	}

	for _, d := range data {
		actual, err := support.DMSToR(d.input)
		assert.NoError(err)
		assert.InDelta(d.expected, actual, 1e-6)
		//log.Printf("OUT: %f", r)
	}
}
