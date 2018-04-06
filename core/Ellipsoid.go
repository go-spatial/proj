package core

import (
	"math"

	"github.com/go-spatial/proj4go/merror"
	"github.com/go-spatial/proj4go/support"
	"github.com/go-spatial/proj4go/tables"
)

// Ellipsoid represents an ellipsoid
type Ellipsoid struct {
	ID    string
	Major string
	Ell   string
	Name  string

	// ELLIPSOID

	DefSize           string
	DefShape          string
	DefSpherification string
	DefEllps          string

	/* The linear parameters */

	a  float64 /* semimajor axis (radius if eccentricity==0) */
	b  float64 /* semiminor axis */
	ra float64 /* 1/a */
	rb float64 /* 1/b */

	/* The eccentricities */

	alpha  float64 /* angular eccentricity */
	e      float64 /* first  eccentricity */
	es     float64 /* first  eccentricity squared */
	e2     float64 /* second eccentricity */
	e2s    float64 /* second eccentricity squared */
	e3     float64 /* third  eccentricity */
	e3s    float64 /* third  eccentricity squared */
	oneEs  float64 /* 1 - e^2 */
	rOneEs float64 /* 1/one_es */

	/* The flattenings */
	f   float64 /* first  flattening */
	f2  float64 /* second flattening */
	n   float64 /* third  flattening */
	rf  float64 /* 1/f  */
	rf2 float64 /* 1/f2 */
	rn  float64 /* 1/n  */
}

// EllipsoidTable is the global list of all the known datums
var EllipsoidTable map[string]*Ellipsoid

func init() {

	EllipsoidTable = map[string]*Ellipsoid{}

	for _, raw := range tables.RawEllipsoids {
		d := &Ellipsoid{
			ID:    raw.ID,
			Major: raw.Major,
			Ell:   raw.Ell,
			Name:  raw.Name,
		}

		EllipsoidTable[d.ID] = d
	}
}

func (ps *ProjString) processEllipsoid(P *Ellipsoid) error {

	var err error

	/* Specifying R overrules everything */
	if ps.Args.ContainsKey("R") {
		err = ellpsSize(ps, P)
		if err != nil {
			return err
		}
		err = pjCalcEllipsoidParams(P, P.a, 0)
		if err != nil {
			return err
		}
		return nil
	}

	/* If an ellps argument is specified, start by using that */
	err = ellpsEllps(ps, P)
	if err != nil {
		return err
	}

	/* We may overwrite the size */
	err = ellpsSize(ps, P)
	if err != nil {
		return err
	}

	/* We may also overwrite the shape */
	err = ellpsShape(ps, P)
	if err != nil {
		return err
	}

	/* When we're done with it, we compute all related ellipsoid parameters */
	err = pjCalcEllipsoidParams(P, P.a, P.es)
	if err != nil {
		return nil
	}

	/* And finally, we may turn it into a sphere */
	err = ellpsSpherification(ps, P)
	if err != nil {
		return err
	}

	//proj_log_debug (P, "pj_ellipsoid - final: a=%.3f f=1/%7.3f, errno=%d",
	//                  P->a,  P->f!=0? 1/P->f: 0,  proj_errno (P));
	//proj_log_debug (P, "pj_ellipsoid - final: %s %s %s %s",
	//                  P->def_size?           P->def_size: empty,
	//                P->def_shape?          P->def_shape: empty,
	//              P->def_spherification? P->def_spherification: empty,
	//            P->def_ellps?          P->def_ellps: empty            );

	/* success */
	return nil
}

func pjCalcEllipsoidParams(P *Ellipsoid, a float64, es float64) error {

	P.a = a
	P.es = es

	/* Compute some ancillary ellipsoidal parameters */
	if P.e == 0 {
		P.e = math.Sqrt(P.es) /* eccentricity */
	}
	P.alpha = math.Asin(P.e) /* angular eccentricity */

	/* second eccentricity */
	P.e2 = math.Tan(P.alpha)
	P.e2s = P.e2 * P.e2

	/* third eccentricity */
	if 0 != P.alpha {
		P.e3 = math.Sin(P.alpha) / math.Sqrt(2-math.Sin(P.alpha)*math.Sin(P.alpha))
	} else {
		P.e3 = 0
	}

	P.e3s = P.e3 * P.e3

	/* flattening */
	if 0 == P.f {
		P.f = 1 - math.Cos(P.alpha) /* = 1 - sqrt (1 - PIN->es); */
	}
	P.rf = math.MaxFloat64
	if P.f != 0.0 {
		P.rf = 1.0 / P.f
	}

	/* second flattening */
	P.f2 = 0
	if math.Cos(P.alpha) != 0 {
		P.f2 = 1/math.Cos(P.alpha) - 1
	}
	P.rf2 = math.MaxFloat64
	if P.f2 != 0.0 {
		P.rf2 = 1 / P.f2
	}

	/* third flattening */
	P.n = math.Pow(math.Tan(P.alpha/2), 2)
	P.rn = math.MaxFloat64
	if P.n != 0.0 {
		P.rn = 1 / P.n
	}

	/* ...and a few more */
	if 0 == P.b {
		P.b = (1 - P.f) * P.a
	}
	P.rb = 1. / P.b
	P.ra = 1. / P.a

	P.oneEs = 1. - P.es
	if P.oneEs == 0. {
		return merror.New(merror.ErrEccentricityIsOne)
	}

	P.rOneEs = 1. / P.oneEs

	return nil
}

func ellpsEllps(ps *ProjString, P *Ellipsoid) error {

	/* Sail home if ellps=xxx is not specified */
	name, ok := ps.Args.GetAsString("ellps")
	if !ok {
		return nil
	}

	/* Then look up the right size and shape parameters from the builtin list */
	if name == "" {
		return merror.New(merror.ErrInvalidArg)
	}

	ellps, ok := EllipsoidTable[name]
	if !ok {
		return merror.New(merror.ErrUnknownEllpParam)
	}

	ellpsSize(ps, ellps)
	ellpsShape(ps, ellps)

	return nil
}

func ellpsSize(ps *ProjString, P *Ellipsoid) error {

	aWasSet := false

	/* A size parameter *must* be given, but may have been given as ellps prior */
	if P.a != 0.0 {
		aWasSet = true
	}

	/* Check which size key is specified */
	key := "R"
	value, ok := ps.Args.GetAsFloat("R")
	if !ok {
		key = "a"
		value, ok = ps.Args.GetAsFloat("a")
	}
	if !ok {
		if aWasSet {
			return nil
		}
		return merror.New(merror.ErrMajorAxisNotGiven)
	}

	P.DefSize = key
	P.a = value
	if P.a <= 0.0 {
		return merror.New(merror.ErrMajorAxisNotGiven)
	}
	if P.a == math.MaxFloat64 {
		return merror.New(merror.ErrMajorAxisNotGiven)
	}

	if key == "R" {
		P.es = 0
		P.f = 0
		P.e = 0
		P.rf = 0
		P.b = P.a
	}

	return nil
}

func ellpsShape(ps *ProjString, P *Ellipsoid) error {

	keys := []string{"rf", "f", "es", "e", "b"}

	/* Check which shape key is specified */
	var key string
	found := false
	var foundValue float64
	for _, key = range keys {
		value, ok := ps.Args.GetAsFloat(key)
		if ok {
			found = true
			foundValue = value
			break
		}
	}

	/* Not giving a shape parameter means selecting a sphere, unless shape */
	/* has been selected previously via ellps=xxx */
	if !found && P.es != 0 {
		return nil
	}
	if !found && P.es == 0 {
		P.es = 0
		P.f = 0
		P.b = P.a
		return nil
	}

	P.es = 0
	P.f = 0
	P.b = 0
	P.e = 0
	P.rf = 0

	switch key {

	/* reverse flattening, rf */
	case "rf":
		P.rf = foundValue
		if P.rf == math.MaxFloat64 {
			return merror.New(merror.ErrInvalidArg)
		}
		if P.rf == 0 {
			return merror.New(merror.ErrRevFlatteningIsZero)
		}
		P.f = 1 / P.rf
		P.es = 2*P.f - P.f*P.f

	/* flattening, f */
	case "f":
		P.f = foundValue
		if P.f == math.MaxFloat64 {
			return merror.New(merror.ErrInvalidArg)
		}
		if P.f == 0 {
			return merror.New(merror.ErrInvalidArg)
		}
		P.rf = 1 / P.f
		P.es = 2*P.f - P.f*P.f

	/* eccentricity squared, es */
	case "es":
		P.es = foundValue
		if P.es == math.MaxFloat64 {
			return merror.New(merror.ErrInvalidArg)
		}
		if P.es == 1 {
			return merror.New(merror.ErrEccentricityIsOne)
		}

	/* eccentricity, e */
	case "e":
		P.e = foundValue
		if P.e == math.MaxFloat64 {
			return merror.New(merror.ErrInvalidArg)
		}
		if P.e == 0 {
			return merror.New(merror.ErrInvalidArg)
		}
		if P.e == 1 {
			return merror.New(merror.ErrEccentricityIsOne)
		}
		P.es = P.e * P.e

	/* semiminor axis, b */
	case "b":
		P.b = foundValue
		if P.b == math.MaxFloat64 {
			return merror.New(merror.ErrInvalidArg)
		}
		if P.b == 0 {
			return merror.New(merror.ErrEccentricityIsOne)
		}
		if P.b == P.a {
			break
		}
		P.f = (P.a - P.b) / P.a
		P.es = 2*P.f - P.f*P.f

	default:
		return merror.New(merror.ErrInvalidArg)

	}

	if P.es < 0 {
		return merror.New(merror.ErrEsLessThanZero)
	}

	return nil
}

func ellpsSpherification(ps *ProjString, P *Ellipsoid) error {

	/* series coefficients for calculating ellipsoid-equivalent spheres */
	const SIXTH = 1 / 6.
	const RA4 = 17 / 360.
	const RA6 = 67 / 3024.
	const RV4 = 5 / 72.
	const RV6 = 55 / 1296.

	keys := []string{"R_A", "R_V", "R_a", "R_g", "R_h", "R_lat_a", "R_lat_g"}

	var key string
	found := false
	for _, key = range keys {
		if ps.Args.ContainsKey(key) {
			found = true
			break
		}
	}

	/* No spherification specified? Then we're done */
	if !found {
		return nil
	}

	switch key {

	/* R_A - a sphere with same area as ellipsoid */
	case "R_A":
		P.a *= 1. - P.es*(SIXTH+P.es*(RA4+P.es*RA6))

	/* R_V - a sphere with same volume as ellipsoid */
	case "R_V":
		P.a *= 1. - P.es*(SIXTH+P.es*(RV4+P.es*RV6))

	/* R_a - a sphere with R = the arithmetic mean of the ellipsoid */
	case "R_a":
		P.a = (P.a + P.b) / 2

	/* R_g - a sphere with R = the geometric mean of the ellipsoid */
	case "R_g":
		P.a = math.Sqrt(P.a * P.b)

	/* R_h - a sphere with R = the harmonic mean of the ellipsoid */
	case "R_h":
		if P.a+P.b == 0 {
			return merror.New(merror.ErrToleranceCondition)
		}
		P.a = (2 * P.a * P.b) / (P.a + P.b)

		/* R_lat_a - a sphere with R = the arithmetic mean of the ellipsoid at given latitude */
		/* R_lat_g - a sphere with R = the geometric  mean of the ellipsoid at given latitude */
	case "R_lat_a", "R_lat_g":
		v, ok := ps.Args.GetAsString(key)
		if !ok {
			return merror.New(merror.ErrInvalidArg)
		}
		t, err := support.DMSToR(v)
		if err != nil {
			return err
		}
		if math.Abs(t) > support.PiOverTwo {
			return merror.New(merror.ErrRefRadLargerThan90)
		}
		t = math.Sin(t)
		t = 1 - P.es*t*t
		if key == "R_lat_a" { /* arithmetic */
			P.a *= (1. - P.es + t) / (2 * t * math.Sqrt(t))
		} else { /* geometric */
			P.a *= math.Sqrt(1-P.es) / t
		}

	default:
		return merror.New(merror.ErrInvalidArg)

	}

	/* Clean up the ellipsoidal parameters to reflect the sphere */
	P.es = 0
	P.e = 0
	P.f = 0
	P.rf = math.MaxFloat64
	P.b = P.a
	pjCalcEllipsoidParams(P, P.a, 0)

	return nil
}
