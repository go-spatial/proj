package operations

import (
	"math"

	"github.com/go-spatial/proj4go/core"
	"github.com/go-spatial/proj4go/merror"
	"github.com/go-spatial/proj4go/support"
)

func init() {
	core.RegisterOperation("aea",
		"Albers Equal Area",
		"\n\tConic Sph&Ell\n\tlat_1= lat_2=",
		core.OperationTypeConversion, core.CoordTypeLP, core.CoordTypeXY,
		NewAea,
	)
	core.RegisterOperation("leac",
		"Lambert Equal Area Conic",
		"\n\tConic, Sph&Ell\n\tlat_1= south",
		core.OperationTypeConversion, core.CoordTypeLP, core.CoordTypeXY,
		NewLeac)
}

// Aea implements core.IOperation and core.ConvertLPToXY
type Aea struct {
	core.OperationCommon
	isLambert bool

	// the "opaque" parts

	ec     float64
	n      float64
	c      float64
	dd     float64
	n2     float64
	rho0   float64
	rho    float64
	phi1   float64
	phi2   float64
	en     []float64
	ellips bool
}

// NewAea is
func NewAea(system *core.System, desc *core.OperationDescription) (core.IOperation, error) {
	xxx := &Aea{
		isLambert: false,
	}
	xxx.System = system

	err := xxx.aeaSetup(system)
	if err != nil {
		return nil, err
	}
	return xxx, nil
}

// NewLeac is too
func NewLeac(system *core.System, desc *core.OperationDescription) (core.IOperation, error) {
	xxx := &Aea{
		isLambert: true,
	}
	xxx.System = system

	err := xxx.leacSetup(system)
	if err != nil {
		return nil, err
	}
	return xxx, nil
}

//---------------------------------------------------------------------

// Forward goes forewards
func (aea *Aea) Forward(lp *core.CoordLP) (*core.CoordXY, error) {

	lp, err := aea.ForwardPrepare(lp)
	if err != nil {
		return nil, err
	}

	xy, err := aea.aeaForward(lp)
	if err != nil {
		return nil, err
	}

	xy, err = aea.ForwardFinalize(xy)
	if err != nil {
		return nil, err
	}

	return xy, nil
}

// Inverse goes backwards
func (aea *Aea) Inverse(xy *core.CoordXY) (*core.CoordLP, error) {

	xy, err := aea.InversePrepare(xy)
	if err != nil {
		return nil, err
	}

	lp, err := aea.aeaInverse(xy)
	if err != nil {
		return nil, err
	}

	lp, err = aea.InverseFinalize(lp)
	if err != nil {
		return nil, err
	}

	return lp, nil
}

//---------------------------------------------------------------------

const eps10 = 1.e-10
const tol7 = 1.e-7

/* determine latitude angle phi-1 */
const nIter = 15
const epsilon = 1.0e-7
const tol = 1.0e-10

func phi1(qs, Te, tOneEs float64) float64 {
	var i int
	var Phi, sinpi, cospi, con, com, dphi float64

	Phi = math.Asin(.5 * qs)
	if Te < epsilon {
		return (Phi)
	}
	i = nIter
	for {
		sinpi = math.Sin(Phi)
		cospi = math.Cos(Phi)
		con = Te * sinpi
		com = 1. - con*con
		dphi = .5 * com * com / cospi * (qs/tOneEs -
			sinpi/com + .5/Te*math.Log((1.-con)/
			(1.+con)))
		Phi += dphi
		i--
		if !(math.Abs(dphi) > tol && i != 0) {
			break
		}
	}
	if i != 0 {
		return Phi
	}
	return math.MaxFloat64
}

func (aea *Aea) localSetup(sys *core.System) error {
	var cosphi, sinphi float64
	var secant bool

	Q := aea
	P := aea.System
	PE := P.Ellipsoid

	if math.Abs(Q.phi1+Q.phi2) < eps10 {
		return merror.New(merror.ErrConicLatEqual)
	}
	sinphi = math.Sin(Q.phi1)
	Q.n = sinphi
	cosphi = math.Cos(Q.phi1)
	secant = math.Abs(Q.phi1-Q.phi2) >= eps10
	Q.ellips = (P.Ellipsoid.Es > 0.0)
	if Q.ellips {
		var ml1, m1 float64

		Q.en = support.Enfn(PE.Es)
		m1 = support.Msfn(sinphi, cosphi, PE.Es)
		ml1 = support.Qsfn(sinphi, PE.E, PE.OneEs)
		if secant { // secant cone
			var ml2, m2 float64

			sinphi = math.Sin(Q.phi2)
			cosphi = math.Cos(Q.phi2)
			m2 = support.Msfn(sinphi, cosphi, PE.Es)
			ml2 = support.Qsfn(sinphi, PE.E, PE.OneEs)
			if ml2 == ml1 {
				return merror.New(merror.ErrAeaSetupFailed)
			}

			Q.n = (m1*m1 - m2*m2) / (ml2 - ml1)
		}
		Q.ec = 1. - .5*PE.OneEs*math.Log((1.-PE.E)/
			(1.+PE.E))/PE.E
		Q.c = m1*m1 + Q.n*ml1
		Q.dd = 1. / Q.n
		Q.rho0 = Q.dd * math.Sqrt(Q.c-Q.n*support.Qsfn(math.Sin(P.Phi0),
			PE.E, PE.OneEs))
	} else {
		if secant {
			Q.n = .5 * (Q.n + math.Sin(Q.phi2))
		}
		Q.n2 = Q.n + Q.n
		Q.c = cosphi*cosphi + Q.n2*sinphi
		Q.dd = 1. / Q.n
		Q.rho0 = Q.dd * math.Sqrt(Q.c-Q.n2*math.Sin(P.Phi0))
	}

	return nil
}

// Forward goes frontwords
func (aea *Aea) aeaForward(lp *core.CoordLP) (*core.CoordXY, error) {
	xy := &core.CoordXY{X: 0.0, Y: 0.0}
	Q := aea
	PE := aea.System.Ellipsoid

	var t float64
	if Q.ellips {
		t = Q.n * support.Qsfn(math.Sin(lp.Phi), PE.E, PE.OneEs)
	} else {
		t = Q.n2 * math.Sin(lp.Phi)
	}
	Q.rho = Q.c - t
	if Q.rho < 0. {
		return xy, merror.New(merror.ErrToleranceCondition)
	}
	Q.rho = Q.dd * math.Sqrt(Q.rho)
	lp.Lam *= Q.n
	xy.X = Q.rho * math.Sin(lp.Lam)
	xy.Y = Q.rho0 - Q.rho*math.Cos(lp.Lam)
	return xy, nil
}

// Inverse goes backwards
func (aea *Aea) aeaInverse(xy *core.CoordXY) (*core.CoordLP, error) {

	lp := &core.CoordLP{Lam: 0.0, Phi: 0.0}
	Q := aea
	PE := aea.System.Ellipsoid

	xy.Y = Q.rho0 - xy.Y
	Q.rho = math.Hypot(xy.X, xy.Y)
	if Q.rho != 0.0 {
		if Q.n < 0. {
			Q.rho = -Q.rho
			xy.X = -xy.X
			xy.Y = -xy.Y
		}
		lp.Phi = Q.rho / Q.dd
		if Q.ellips {
			lp.Phi = (Q.c - lp.Phi*lp.Phi) / Q.n
			if math.Abs(Q.ec-math.Abs(lp.Phi)) > tol7 {
				lp.Phi = phi1(lp.Phi, PE.E, PE.OneEs)
				if lp.Phi == math.MaxFloat64 {
					return lp, merror.New(merror.ErrToleranceCondition)
				}
			} else {
				if lp.Phi < 0. {
					lp.Phi = -support.PiOverTwo
				} else {
					lp.Phi = support.PiOverTwo
				}
			}
		} else {
			lp.Phi = (Q.c - lp.Phi*lp.Phi) / Q.n2
			if math.Abs(lp.Phi) <= 1. {
				lp.Phi = math.Asin(lp.Phi)
			} else {
				if lp.Phi < 0. {
					lp.Phi = -support.PiOverTwo
				} else {
					lp.Phi = support.PiOverTwo
				}
			}
		}
		lp.Lam = math.Atan2(xy.X, xy.Y) / Q.n
	} else {
		lp.Lam = 0.
		if Q.n > 0. {
			lp.Phi = support.PiOverTwo
		} else {
			lp.Phi = -support.PiOverTwo
		}
	}
	return lp, nil
}

func (aea *Aea) aeaSetup(sys *core.System) error {

	lat1, ok := aea.System.ProjString.GetAsFloat("lat_1")
	if !ok {
		lat1 = 0.0
	}
	lat2, ok := aea.System.ProjString.GetAsFloat("lat_2")
	if !ok {
		lat2 = 0.0
	}

	aea.phi1 = support.DDToR(lat1)
	aea.phi2 = support.DDToR(lat2)

	return aea.localSetup(aea.System)
}

func (aea *Aea) leacSetup(sys *core.System) error {

	lat1, ok := aea.System.ProjString.GetAsFloat("lat_1")
	if !ok {
		lat1 = 0.0
	}

	south := -support.PiOverTwo
	_, ok = aea.System.ProjString.GetAsInt("south")
	if !ok {
		south = support.PiOverTwo
	}

	aea.phi2 = support.DDToR(lat1)
	aea.phi1 = south

	return aea.localSetup(aea.System)
}
