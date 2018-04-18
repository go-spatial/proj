package support

import (
	"math"
)

// Tsfn is to "determine small t"
func Tsfn(phi, sinphi, e float64) float64 {
	sinphi *= e

	/* avoid zero division, fail gracefully */
	denominator := 1.0 + sinphi
	if denominator == 0.0 {
		return math.MaxFloat64
	}

	return (math.Tan(.5*(PiOverTwo-phi)) /
		math.Pow((1.-sinphi)/(denominator), .5*e))
}
