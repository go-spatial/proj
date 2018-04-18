package support

import (
	"math"

	"github.com/go-spatial/proj4go/mlog"
)

/* arc sin, cosine, tan2 and sqrt that will NOT fail */

const oneTol = 1.00000000000001
const aTol = 1e-50

// Aasin is asin w/ error catching
func Aasin(v float64) float64 {

	av := math.Abs(v)

	if av >= 1. {
		if av > oneTol {
			// TODO: we are supposed to signal an error, but not actually fail
			mlog.Printf("error signal in aasin()")
		}
		if v < 0. {
			return -math.Pi / 2.0
		}
		return math.Pi / 2.0
	}
	return math.Asin(v)
}

// Aacos is acos w/ error catching
func Aacos(v float64) float64 {

	av := math.Abs(v)

	if av >= 1. {
		if av > oneTol {
			// TODO: we are supposed to signal an error, but not actually fail
			mlog.Printf("error signal in aacos()")
		}
		if v < 0. {
			return math.Pi
		}
		return 0.
	}
	return math.Acos(v)
}

// Asqrt is sqrt w/ error catching
func Asqrt(v float64) float64 {
	if v <= 0 {
		return 0.0
	}
	return math.Sqrt(v)
}

// Aatan2 is atan2 w/ error catching
func Aatan2(n, d float64) float64 {
	if math.Abs(n) < aTol && math.Abs(d) < aTol {
		return 0.0
	}
	return math.Atan2(n, d)
}
