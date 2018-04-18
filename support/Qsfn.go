package support

import (
	"math"
)

const epsilon = 1.0e-7

// Qsfn is ..?
func Qsfn(sinphi, e, oneEs float64) float64 {
	var con, div1, div2 float64

	if e >= epsilon {
		con = e * sinphi
		div1 = 1.0 - con*con
		div2 = 1.0 + con

		/* avoid zero division, fail gracefully */
		if div1 == 0.0 || div2 == 0.0 {
			return math.MaxFloat64
		}

		return (oneEs * (sinphi/div1 - (.5/e)*math.Log((1.-con)/div2)))
	}
	return (sinphi + sinphi)
}
