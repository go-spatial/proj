package support

import "math"

// Msfn is to "determine constant small m"
func Msfn(sinphi, cosphi, es float64) float64 {
	return (cosphi / math.Sqrt(1.-es*sinphi*sinphi))
}
