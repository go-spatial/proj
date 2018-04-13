package operations

import (
	"math"

	"github.com/go-spatial/proj4go/core"
)

func init() {
	core.RegisterOperation("aea",
		"Albers Equal Area",
		"\n\tConic Sph&Ell\n\tlat_1= lat_2=",
		core.CoordTypeLP, core.CoordTypeXY,
		NewAea,
	)
	core.RegisterOperation("leac",
		"Lambert Equal Area Conic",
		"\n\tConic, Sph&Ell\n\tlat_1= south",
		core.CoordTypeLP, core.CoordTypeXY,
		NewLeac)
}

// NewAea is
func NewAea(system *core.System) (core.IOperation, error) { return nil, nil }

// NewLeac is too
func NewLeac(system *core.System) (core.IOperation, error) { return nil, nil }

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

/*
 struct pj_opaque {
	 double  ec;
	 double  n;
	 double  c;
	 double  dd;
	 double  n2;
	 double  rho0;
	 double  rho;
	 double  phi1;
	 double  phi2;
	 double  *en;
	 int     ellips;
 };


 static XY e_forward (LP lp, PJ *P) {
	 XY xy = {0.0,0.0};
	 struct pj_opaque *Q = P->opaque;
	 Q->rho = Q->c - (Q->ellips ? Q->n * pj_qsfn(sin(lp.phi), P->e, P->one_es) : Q->n2 * sin(lp.phi));;
	 if (Q->rho < 0.) {
		 proj_errno_set(P, PJD_ERR_TOLERANCE_CONDITION);
		 return xy;
	 }
	 Q->rho = Q->dd * sqrt(Q->rho);
	 xy.x = Q->rho * sin( lp.lam *= Q->n );
	 xy.y = Q->rho0 - Q->rho * cos(lp.lam);
	 return xy;
 }


 static LP e_inverse (XY xy, PJ *P) {
	 LP lp = {0.0,0.0};
	 struct pj_opaque *Q = P->opaque;
	 if( (Q->rho = hypot(xy.x, xy.y = Q->rho0 - xy.y)) != 0.0 ) {
		 if (Q->n < 0.) {
			 Q->rho = -Q->rho;
			 xy.x = -xy.x;
			 xy.y = -xy.y;
		 }
		 lp.phi =  Q->rho / Q->dd;
		 if (Q->ellips) {
			 lp.phi = (Q->c - lp.phi * lp.phi) / Q->n;
			 if (fabs(Q->ec - fabs(lp.phi)) > TOL7) {
				 if ((lp.phi = phi1_(lp.phi, P->e, P->one_es)) == HUGE_VAL) {
					 proj_errno_set(P, PJD_ERR_TOLERANCE_CONDITION);
					 return lp;
				 }
			 } else
				 lp.phi = lp.phi < 0. ? -M_HALFPI : M_HALFPI;
		 } else if (fabs(lp.phi = (Q->c - lp.phi * lp.phi) / Q->n2) <= 1.)
			 lp.phi = asin(lp.phi);
		 else
			 lp.phi = lp.phi < 0. ? -M_HALFPI : M_HALFPI;
		 lp.lam = atan2(xy.x, xy.y) / Q->n;
	 } else {
		 lp.lam = 0.;
		 lp.phi = Q->n > 0. ? M_HALFPI : - M_HALFPI;
	 }
	 return lp;
 }



 static PJ *setup(PJ *P) {
	 double cosphi, sinphi;
	 int secant;
	 struct pj_opaque *Q = P->opaque;

	 P->inv = e_inverse;
	 P->fwd = e_forward;

	 if (fabs(Q->phi1 + Q->phi2) < EPS10)
		 return destructor(P, PJD_ERR_CONIC_LAT_EQUAL);
	 Q->n = sinphi = sin(Q->phi1);
	 cosphi = cos(Q->phi1);
	 secant = fabs(Q->phi1 - Q->phi2) >= EPS10;
	 if( (Q->ellips = (P->es > 0.))) {
		 double ml1, m1;

		 if (!(Q->en = pj_enfn(P->es)))
			 return destructor(P, 0);
		 m1 = pj_msfn(sinphi, cosphi, P->es);
		 ml1 = pj_qsfn(sinphi, P->e, P->one_es);
		 if (secant) { // secant cone
			 double ml2, m2;

			 sinphi = sin(Q->phi2);
			 cosphi = cos(Q->phi2);
			 m2 = pj_msfn(sinphi, cosphi, P->es);
			 ml2 = pj_qsfn(sinphi, P->e, P->one_es);
			 if (ml2 == ml1)
				 return destructor(P, 0);

			 Q->n = (m1 * m1 - m2 * m2) / (ml2 - ml1);
		 }
		 Q->ec = 1. - .5 * P->one_es * log((1. - P->e) /
			 (1. + P->e)) / P->e;
		 Q->c = m1 * m1 + Q->n * ml1;
		 Q->dd = 1. / Q->n;
		 Q->rho0 = Q->dd * sqrt(Q->c - Q->n * pj_qsfn(sin(P->phi0),
			 P->e, P->one_es));
	 } else {
		 if (secant) Q->n = .5 * (Q->n + sin(Q->phi2));
		 Q->n2 = Q->n + Q->n;
		 Q->c = cosphi * cosphi + Q->n2 * sinphi;
		 Q->dd = 1. / Q->n;
		 Q->rho0 = Q->dd * sqrt(Q->c - Q->n2 * sin(P->phi0));
	 }

	 return P;
 }


 PJ *PROJECTION(aea) {
	 struct pj_opaque *Q = pj_calloc (1, sizeof (struct pj_opaque));
	 if (0==Q)
		 return pj_default_destructor (P, ENOMEM);
	 P->opaque = Q;
	 P->destructor = destructor;

	 Q->phi1 = pj_param(P->ctx, P->params, "rlat_1").f;
	 Q->phi2 = pj_param(P->ctx, P->params, "rlat_2").f;
	 return setup(P);
 }


 PJ *PROJECTION(leac) {
	 struct pj_opaque *Q = pj_calloc (1, sizeof (struct pj_opaque));
	 if (0==Q)
		 return pj_default_destructor (P, ENOMEM);
	 P->opaque = Q;
	 P->destructor = destructor;

	 Q->phi2 = pj_param(P->ctx, P->params, "rlat_1").f;
	 Q->phi1 = pj_param(P->ctx, P->params, "bsouth").i ? - M_HALFPI: M_HALFPI;
	 return setup(P);
 }

*/
