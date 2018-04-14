package operations

import (
	"math"

	"github.com/go-spatial/proj4go/core"
	"github.com/go-spatial/proj4go/merror"
	"github.com/go-spatial/proj4go/support"
)

func init() {
	core.RegisterOperation("merc",
		"Universal Transverse Mercator (UTM)",
		"\n\tCyl, Sph&Ell\n\tlat_ts=",
		core.OperationTypeConversion, core.CoordTypeLP, core.CoordTypeXY,
		NewMerc,
	)
}

const xeps10 = 1.e-10

// Merc implements core.IOperation and core.ConvertLPToXY
type Merc struct {
	core.OperationCommon
	isSphere bool
}

// NewMerc returns a new Merc
func NewMerc(system *core.System, desc *core.OperationDescription) (core.IOperation, error) {
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

	lp, err := merc.ForwardPrepare(lp)
	if err != nil {
		return nil, err
	}

	var xy *core.CoordXY
	if merc.isSphere {
		xy, err = merc.mercSphericalForward(lp)
	} else {
		xy, err = merc.mercEllipsoidalForward(lp)
	}
	if err != nil {
		return nil, err
	}

	xy, err = merc.ForwardFinalize(xy)
	if err != nil {
		return nil, err
	}

	return xy, nil
}

// Inverse goes backwards
func (merc *Merc) Inverse(xy *core.CoordXY) (*core.CoordLP, error) {

	xy, err := merc.InversePrepare(xy)
	if err != nil {
		return nil, err
	}

	var lp *core.CoordLP
	if merc.isSphere {
		lp, err = merc.mercSphericalInverse(xy)
	} else {
		lp, err = merc.mercEllipsoidalInverse(xy)
	}
	if err != nil {
		return nil, err
	}

	lp, err = merc.InverseFinalize(lp)
	if err != nil {
		return nil, err
	}

	return lp, nil
}

//---------------------------------------------------------------------

func (merc *Merc) mercEllipsoidalForward(lp *core.CoordLP) (*core.CoordXY, error) { /* Ellipsoidal, forward */
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

func (merc *Merc) mercSphericalForward(lp *core.CoordLP) (*core.CoordXY, error) { /* Spheroidal, forward */
	xy := &core.CoordXY{X: 0.0, Y: 0.0}

	P := merc.System

	if math.Abs(math.Abs(lp.Phi)-support.PiOverTwo) <= xeps10 {
		return xy, merror.New(merror.ErrToleranceCondition)
	}
	xy.X = P.K0 * lp.Lam
	xy.Y = P.K0 * math.Log(math.Tan(support.PiOverFour+.5*lp.Phi))
	return xy, nil
}

func (merc *Merc) mercEllipsoidalInverse(xy *core.CoordXY) (*core.CoordLP, error) { /* Ellipsoidal, inverse */
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

func (merc *Merc) mercSphericalInverse(xy *core.CoordXY) (*core.CoordLP, error) { /* Spheroidal, inverse */
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
