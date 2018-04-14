package operations

import (
	"math"

	"github.com/go-spatial/proj4go/core"
	"github.com/go-spatial/proj4go/merror"
	"github.com/go-spatial/proj4go/support"
)

func init() {
	core.RegisterConvertLPToXY("merc",
		"Universal Transverse Mercator (UTM)",
		"\n\tCyl, Sph&Ell\n\tlat_ts=",
		NewMerc,
	)
}

const xeps10 = 1.e-10

// Merc implements core.IOperation and core.ConvertLPToXY
type Merc struct {
	core.Operation
	isSphere bool
}

// NewMerc returns a new Merc
func NewMerc(system *core.System, desc *core.OperationDescription) (core.IConvertLPToXY, error) {
	xxx := &Merc{
		isSphere: false,
	}
	xxx.System = system

	err := xxx.mercSetup(system)
	if err != nil {
		return nil, err
	}
	return xxx, nil
}

// Forward goes forewards
func (merc *Merc) Forward(lp *core.CoordLP) (*core.CoordXY, error) {

	if merc.isSphere {
		return merc.sphericalForward(lp)
	}
	return merc.ellipsoidalForward(lp)
}

// Inverse goes backwards
func (merc *Merc) Inverse(xy *core.CoordXY) (*core.CoordLP, error) {

	if merc.isSphere {
		return merc.sphericalInverse(xy)
	}
	return merc.ellipsoidalInverse(xy)
}

//---------------------------------------------------------------------

func (merc *Merc) ellipsoidalForward(lp *core.CoordLP) (*core.CoordXY, error) { /* Ellipsoidal, forward */
	xy := &core.CoordXY{X: 0.0, Y: 0.0}

	P := merc.System
	PE := merc.System.Ellipsoid

	if math.Abs(math.Abs(lp.Phi)-support.PiOverTwo) <= xeps10 {
		return xy, merror.New(merror.ErrToleranceCondition)
	}
	xy.X = P.K0 * lp.Lam
	xy.Y = -P.K0 * math.Log(support.Tsfn(lp.Phi, math.Sin(lp.Phi), PE.E))
	return xy, nil
}

func (merc *Merc) sphericalForward(lp *core.CoordLP) (*core.CoordXY, error) { /* Spheroidal, forward */
	xy := &core.CoordXY{X: 0.0, Y: 0.0}

	P := merc.System

	if math.Abs(math.Abs(lp.Phi)-support.PiOverTwo) <= xeps10 {
		return xy, merror.New(merror.ErrToleranceCondition)
	}
	xy.X = P.K0 * lp.Lam
	xy.Y = P.K0 * math.Log(math.Tan(support.PiOverFour+.5*lp.Phi))
	return xy, nil
}

func (merc *Merc) ellipsoidalInverse(xy *core.CoordXY) (*core.CoordLP, error) { /* Ellipsoidal, inverse */
	lp := &core.CoordLP{Lam: 0.0, Phi: 0.0}

	P := merc.System
	PE := merc.System.Ellipsoid
	var err error

	lp.Phi, err = support.Phi2(math.Exp(-xy.Y/P.K0), PE.E)
	if err != nil {
		return nil, err
	}
	if lp.Phi == math.MaxFloat64 {
		return lp, merror.New(merror.ErrToleranceCondition)
	}
	lp.Lam = xy.X / P.K0
	return lp, nil
}

func (merc *Merc) sphericalInverse(xy *core.CoordXY) (*core.CoordLP, error) { /* Spheroidal, inverse */
	lp := &core.CoordLP{Lam: 0.0, Phi: 0.0}

	P := merc.System

	lp.Phi = support.PiOverTwo - 2.*math.Atan(math.Exp(-xy.Y/P.K0))
	lp.Lam = xy.X / P.K0
	return lp, nil
}

func (merc *Merc) mercSetup(sys *core.System) error {
	var phits float64

	ps := merc.System.ProjString

	isPhits := ps.ContainsKey("lat_ts")
	if isPhits {
		phits, _ = ps.GetAsFloat("lat_ts")
		phits = support.DDToR(phits)
		phits = math.Abs(phits)
		if phits >= support.PiOverTwo {
			return merror.New(merror.ErrLatTSLargerThan90)
		}
	}

	P := merc.System
	PE := merc.System.Ellipsoid

	if PE.Es != 0.0 { /* ellipsoid */
		merc.isSphere = false
		if isPhits {
			P.K0 = support.Msfn(math.Sin(phits), math.Cos(phits), PE.Es)
		}
	} else { /* sphere */
		merc.isSphere = true
		if isPhits {
			P.K0 = math.Cos(phits)
		}
	}

	return nil
}
