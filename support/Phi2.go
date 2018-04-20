package support

import (
	"math"

	"github.com/go-spatial/proj/merror"
)

const tol = 1.0e-10
const nIter = 15

// Phi2 is to "determine latitude angle phi-2"
func Phi2(ts, e float64) (float64, error) {
	var eccnth, Phi, con float64
	var i int

	eccnth = .5 * e
	Phi = PiOverTwo - 2.*math.Atan(ts)
	i = nIter
	for {

		con = e * math.Sin(Phi)
		dphi := PiOverTwo - 2.*math.Atan(ts*math.Pow((1.-con)/(1.+con), eccnth)) - Phi
		Phi += dphi
		i--
		if math.Abs(dphi) > tol && i != 0 {
			continue
		}
		break
	}
	if i <= 0 {
		return 0.0, merror.New(merror.Phi2)
	}
	return Phi, nil
}
