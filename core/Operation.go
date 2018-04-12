package core

import (
	"math"

	"github.com/go-spatial/proj4go/merror"
	"github.com/go-spatial/proj4go/support"
)

// IOperation is for all operations
type IOperation interface {
	GetSystem() *System
	GetSignature() (CoordType, CoordType)
	ForwardAny(*CoordAny) (*CoordAny, error)
	InverseAny(*CoordAny) (*CoordAny, error)
	Setup() error
}

// OperationCommon provides some common fields an implementation of IOperation will need
type OperationCommon struct {
	System     *System
	InputType  CoordType
	OutputType CoordType
}

// GetSystem returns the System
func (opcommon OperationCommon) GetSystem() *System { return opcommon.System }

// GetSignature returns the input type and output type
func (opcommon OperationCommon) GetSignature() (CoordType, CoordType) { 
	return opcommon.InputType, opcommon.OutputType
 }

// IConvertLPToXY is just for 2D LP->XY conversions
type IConvertLPToXY interface {
	Forward(*CoordLP) (*CoordXY, error)
	Inverse(*CoordXY) (*CoordLP, error)
}

//---------------------------------------------------------------------

// ForwardPrepare is called just before calling Forward()
func (opcommon *OperationCommon) ForwardPrepare(lp *CoordLP) (*CoordLP, error) {

	op := opcommon.GetSystem()

	if math.MaxFloat64 == lp.Lam {
		return nil, merror.New(merror.ErrCoordinateError)
	}

	// /* The helmert datum shift will choke unless it gets a sensible 4D coordinate */
	// if (HUGE_VAL==coo.v[2] && P->helmert) coo.v[2] = 0.0;
	// if (HUGE_VAL==coo.v[3] && P->helmert) coo.v[3] = 0.0;

	/* Check validity of angular input coordinates */
	if op.Left == IOUnitsAngular {

		/* check for latitude or longitude over-range */
		var t float64
		if lp.Phi < 0 {
			t = -lp.Phi - support.PiOverTwo
		} else {
			t = lp.Phi - support.PiOverTwo
		}
		if t > epsLat || lp.Lam > 10 || lp.Lam < -10 {
			return nil, merror.New(merror.ErrLatOrLonExceededLimit)
		}

		/* Clamp latitude to -90..90 degree range */
		if lp.Phi > support.PiOverTwo {
			lp.Phi = support.PiOverTwo
		}
		if lp.Phi < -support.PiOverTwo {
			lp.Phi = -support.PiOverTwo
		}

		/* If input latitude is geocentrical, convert to geographical */
		if op.Geoc {
			lp = GeocentricLatitude(op, DirectionInverse, lp)
		}

		/* Ensure longitude is in the -pi:pi range */
		if !op.Over {
			lp.Lam = support.Adjlon(lp.Lam)
		}

		//if (P->hgridshift) {
		//	coo = proj_trans (P->hgridshift, PJ_INV, coo);
		//}        else if (P->helmert || (P->cart_wgs84 != 0 && P->cart != 0)) {
		//    coo = proj_trans (P->cart_wgs84, PJ_FWD, coo); /* Go cartesian in WGS84 frame */
		//    if( P->helmert ) {
		//        coo = proj_trans (P->helmert,    PJ_INV, coo); /* Step into local frame */
		//	}
		//	coo = proj_trans (P->cart,       PJ_INV, coo); /* Go back to angular using local ellps */
		//}
		if lp.Lam == math.MaxFloat64 {
			return lp, nil
		}
		//if (P->vgridshift){
		//    coo = proj_trans (P->vgridshift, PJ_FWD, coo); /* Go orthometric from geometric */
		//}

		/* Distance from central meridian, taking system zero meridian into account */
		lp.Lam = (lp.Lam - op.FromGreenwich) - op.Lam0

		/* Ensure longitude is in the -pi:pi range */
		if !op.Over {
			lp.Lam = support.Adjlon(lp.Lam)
		}

		return lp, nil
	}

	/* We do not support gridshifts on cartesian input */
	//if (op.Left==IOUnitsCartesian && P->helmert) {
	//		return proj_trans (P->helmert, PJ_INV, coo);
	//}
	return lp, nil
}

// ForwardFinalize is called just after calling Forward()
func (opcommon *OperationCommon) ForwardFinalize(coo *CoordXY) (*CoordXY, error) {

	op := opcommon.GetSystem()

	switch op.Right {

	/* Handle false eastings/northings and non-metric linear units */
	case IOUnitsCartesian:
		panic(99)

		//if op.IsGeocentric {
		//    coo = proj_trans (P->cart, PJ_FWD, coo);
		//}

		//coo.X = op.FromMeter * (coo.X + P.X0);
		//coo.Y = op.FromMeter * (coo.Y + P.Y0);
		//coo.X = op.FromMeter * (coo.Z + P.Z0);

	/* Classic proj.4 functions return plane coordinates in units of the semimajor axis */
	case IOUnitsClassic:
		coo.X *= op.Ellipsoid.A
		coo.Y *= op.Ellipsoid.A
		fallthrough

	/* Falls through */ /* (<-- GCC warning silencer) */
	/* to continue processing in common with PJ_IO_UNITS_PROJECTED */
	case IOUnitsProjected:
		coo.X = op.FromMeter * (coo.X + op.X0)
		coo.Y = op.FromMeter * (coo.Y + op.Y0)
		///////////////////coo.Z = op.VFromMeter * (coo.Z + op.Z0)

	case IOUnitsWhatever:

	case IOUnitsAngular:
		panic(99)
		//if op.Left == IOUnitsAngular {
		//	break;
		//}

		/* adjust longitude to central meridian */
		//if !op.Over {
		//    coo.lpz.lam = adjlon(coo.lpz.lam);
		//}

		//if (P->vgridshift)
		//    coo = proj_trans (P->vgridshift, PJ_FWD, coo); /* Go orthometric from geometric */
		//if (coo.lp.lam==HUGE_VAL) {
		//	return coo;
		//}

		//if (P->hgridshift)
		//    coo = proj_trans (P->hgridshift, PJ_INV, coo);
		//else if (P->helmert || (P->cart_wgs84 != 0 && P->cart != 0)) {
		//    coo = proj_trans (P->cart_wgs84, PJ_FWD, coo); /* Go cartesian in WGS84 frame */
		//    if( P->helmert )
		//        coo = proj_trans (P->helmert,    PJ_INV, coo); /* Step into local frame */
		//    coo = proj_trans (P->cart,       PJ_INV, coo); /* Go back to angular using local ellps */
		//}
		//if (coo.lp.lam==HUGE_VAL) {
		//	return coo;
		//}

		/* If input latitude was geocentrical, convert back to geocentrical */
		//if op.Geoc {
		//	coo = GeocentricLatitude(op, DirectionForward, coo)
		//}

		/* Distance from central meridian, taking system zero meridian into account */
		//coo.lp.lam = coo.lp.lam + P->from_greenwich + P->lam0;

		/* adjust longitude to central meridian */
		//if (0==P->over) {
		//	coo.lpz.lam = adjlon(coo.lpz.lam);
		//}
	}

	//if (P->axisswap) {
	//    coo = proj_trans (P->axisswap, PJ_FWD, coo);
	//}

	return coo, nil
}

// InversePrepare is called just before calling Inverse()
func (opcommon *OperationCommon) InversePrepare(coo *CoordXY) (*CoordXY, error) {

	op := opcommon.GetSystem()

	if coo.X == math.MaxFloat64 {
		return nil, merror.New(merror.ErrInvalidXOrY)
	}

	///* The helmert datum shift will choke unless it gets a sensible 4D coordinate */
	//if (HUGE_VAL==coo.v[2] && P->helmert) coo.v[2] = 0.0;
	//if (HUGE_VAL==coo.v[3] && P->helmert) coo.v[3] = 0.0;

	//if (P->axisswap)
	//    coo = proj_trans (P->axisswap, PJ_INV, coo);

	/* Check validity of angular input coordinates */
	//if (INPUT_UNITS==PJ_IO_UNITS_ANGULAR) {
	//    double t;
	//
	//    /* check for latitude or longitude over-range */
	//    t = (coo.lp.phi < 0  ?  -coo.lp.phi  :  coo.lp.phi) - M_HALFPI;
	//    if (t > PJ_EPS_LAT  ||  coo.lp.lam > 10  ||  coo.lp.lam < -10) {
	//        proj_errno_set (P, PJD_ERR_LAT_OR_LON_EXCEED_LIMIT);
	//        return proj_coord_error ();
	//    }
	//
	//    /* Clamp latitude to -90..90 degree range */
	//    if (coo.lp.phi > M_HALFPI)
	//        coo.lp.phi = M_HALFPI;
	//    if (coo.lp.phi < -M_HALFPI)
	//        coo.lp.phi = -M_HALFPI;
	//
	//    /* If input latitude is geocentrical, convert to geographical */
	//    if (P->geoc)
	//        coo = proj_geocentric_latitude (P, PJ_INV, coo);
	//
	// /* Distance from central meridian, taking system zero meridian into account */
	//    coo.lp.lam = (coo.lp.lam + P->from_greenwich) - P->lam0;
	//
	//    /* Ensure longitude is in the -pi:pi range */
	//    if (0==P->over)
	//        coo.lp.lam = adjlon(coo.lp.lam);
	//
	//    if (P->hgridshift)
	//        coo = proj_trans (P->hgridshift, PJ_FWD, coo);
	//    else if (P->helmert || (P->cart_wgs84 != 0 && P->cart != 0)) {
	//        coo = proj_trans (P->cart,       PJ_FWD, coo); /* Go cartesian in local frame */
	//        if( P->helmert )
	//            coo = proj_trans (P->helmert,    PJ_FWD, coo); /* Step into WGS84 */
	//        coo = proj_trans (P->cart_wgs84, PJ_INV, coo); /* Go back to angular using WGS84 ellps */
	//    }
	//    if (coo.lp.lam==HUGE_VAL)
	//        return coo;
	//    if (P->vgridshift)
	//        coo = proj_trans (P->vgridshift, PJ_INV, coo); /* Go geometric from orthometric */
	//    return coo;
	//}

	/* Handle remaining possible input types */
	switch op.Right {

	case IOUnitsWhatever:
		return coo, nil

		/* de-scale and de-offset */
	case IOUnitsCartesian:
		coo.X = op.ToMeter*coo.X - op.X0
		coo.Y = op.ToMeter*coo.Y - op.Y0
		/////////////coo.Z = op.ToMeter*coo.Z - op.Z0

		//if (P->is_geocent)
		//    coo = proj_trans (P->cart, PJ_INV, coo);

		return coo, nil

	case IOUnitsProjected, IOUnitsClassic:

		coo.X = op.ToMeter*coo.X - op.X0
		coo.Y = op.ToMeter*coo.Y - op.Y0
		///////////coo.Z = op.VToMeter*coo.Z - op.Z0
		if op.Right == IOUnitsProjected {
			return coo, nil
		}

		/* Classic proj.4 functions expect plane coordinates in units of the semimajor axis  */
		/* Multiplying by ra, rather than dividing by a because the CalCOFI projection       */
		/* stomps on a and hence (apparently) depends on this to roundtrip correctly         */
		/* (CalCOFI avoids further scaling by stomping - but a better solution is possible)  */
		coo.X *= op.Ellipsoid.Ra
		coo.Y *= op.Ellipsoid.Ra
		return coo, nil
	}

	/* Should not happen, so we could return pj_coord_err here */
	return coo, nil
}

// InverseFinalize is called just after calling Inverse()
func (opcommon *OperationCommon) InverseFinalize(coo *CoordLP) (*CoordLP, error) {

	op := opcommon.GetSystem()

	//if (coo.xyz.x == HUGE_VAL) {
	//    proj_errno_set (P, PJD_ERR_INVALID_X_OR_Y);
	//    return proj_coord_error ();
	//}

	if op.Left == IOUnitsAngular {

		if op.Right != IOUnitsAngular {
			/* Distance from central meridian, taking system zero meridian into account */
			coo.Lam = coo.Lam + op.FromGreenwich + op.Lam0

			/* adjust longitude to central meridian */
			if !op.Over {
				coo.Lam = support.Adjlon(coo.Lam)
			}

			//if (P->vgridshift)
			//    coo = proj_trans (P->vgridshift, PJ_INV, coo); /* Go geometric from orthometric */
			//if (coo.lp.lam==HUGE_VAL)
			//    return coo;
			//if (P->hgridshift)
			//    coo = proj_trans (P->hgridshift, PJ_FWD, coo);
			//else if (P->helmert || (P->cart_wgs84 != 0 && P->cart != 0)) {
			//    coo = proj_trans (P->cart,       PJ_FWD, coo); /* Go cartesian in local frame */
			//    if( P->helmert )
			//        coo = proj_trans (P->helmert,    PJ_FWD, coo); /* Step into WGS84 */
			//    coo = proj_trans (P->cart_wgs84, PJ_INV, coo); /* Go back to angular using WGS84 ellps */
			//}
			if coo.Lam == math.MaxFloat64 {
				return coo, nil
			}
		}

		/* If input latitude was geocentrical, convert back to geocentrical */
		if op.Geoc {
			coo = GeocentricLatitude(op, DirectionForward, coo)
		}
	}

	return coo, nil
}
