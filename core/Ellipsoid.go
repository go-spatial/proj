package core

import (
	"encoding/json"
	"math"

	"github.com/go-spatial/proj4go/merror"
	"github.com/go-spatial/proj4go/support"
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

	A  float64 /* semimajor axis (radius if eccentricity==0) */
	B  float64 /* semiminor axis */
	Ra float64 /* 1/a */
	Rb float64 /* 1/b */

	/* The eccentricities */

	Alpha  float64 /* angular eccentricity */
	E      float64 /* first  eccentricity */
	Es     float64 /* first  eccentricity squared */
	E2     float64 /* second eccentricity */
	E2s    float64 /* second eccentricity squared */
	E3     float64 /* third  eccentricity */
	E3s    float64 /* third  eccentricity squared */
	OneEs  float64 /* 1 - e^2 */
	ROneEs float64 /* 1/one_es */

	/* The flattenings */
	F   float64 /* first  flattening */
	F2  float64 /* second flattening */
	N   float64 /* third  flattening */
	Rf  float64 /* 1/f  */
	Rf2 float64 /* 1/f2 */
	Rn  float64 /* 1/n  */

	/* This one's for GRS80 */
	J float64 /* "Dynamic form factor" */

	EsOrig, AOrig float64 /* es and a before any +proj related adjustment */
}

// NewEllipsoid creates an Ellipsoid and initializes it from the proj string
func NewEllipsoid(op *System) (*Ellipsoid, error) {
	ellipsoid := &Ellipsoid{}

	err := ellipsoid.initialize(op)
	if err != nil {
		return nil, err
	}

	return ellipsoid, nil
}

func (e *Ellipsoid) String() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (e *Ellipsoid) initialize(op *System) error {

	ps := op.ProjString

	/* Specifying R overrules everything */
	if ps.ContainsKey("R") {

		err := e.doSize(ps)
		if err != nil {
			return err
		}
		err = e.doCalcParams(e.A, 0)
		if err != nil {
			return err
		}
		return nil
	}

	/* If an ellps argument is specified, start by using that */
	err := e.doEllps(op.ProjString)
	if err != nil {
		return err
	}

	/* We may overwrite the size */
	err = e.doSize(op.ProjString)
	if err != nil {
		return err
	}

	/* We may also overwrite the shape */
	err = e.doShape(op.ProjString)
	if err != nil {
		return err
	}

	/* When we're done with it, we compute all related ellipsoid parameters */
	err = e.doCalcParams(e.A, e.Es)
	if err != nil {
		return nil
	}

	/* And finally, we may turn it into a sphere */
	err = e.doSpherification(op.ProjString)
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

func (e *Ellipsoid) doCalcParams(a float64, es float64) error {

	P := e

	P.A = a
	P.Es = es

	/* Compute some ancillary ellipsoidal parameters */
	if P.E == 0 {
		P.E = math.Sqrt(P.Es) /* eccentricity */
	}
	P.Alpha = math.Asin(P.E) /* angular eccentricity */

	/* second eccentricity */
	P.E2 = math.Tan(P.Alpha)
	P.E2s = P.E2 * P.E2

	/* third eccentricity */
	if 0 != P.Alpha {
		P.E3 = math.Sin(P.Alpha) / math.Sqrt(2-math.Sin(P.Alpha)*math.Sin(P.Alpha))
	} else {
		P.E3 = 0
	}

	P.E3s = P.E3 * P.E3

	/* flattening */
	if 0 == P.F {
		P.F = 1 - math.Cos(P.Alpha) /* = 1 - sqrt (1 - PIN->es); */
	}
	P.Rf = math.MaxFloat64
	if P.F != 0.0 {
		P.Rf = 1.0 / P.F
	}

	/* second flattening */
	P.F2 = 0
	if math.Cos(P.Alpha) != 0 {
		P.F2 = 1/math.Cos(P.Alpha) - 1
	}
	P.Rf2 = math.MaxFloat64
	if P.F2 != 0.0 {
		P.Rf2 = 1 / P.F2
	}

	/* third flattening */
	P.N = math.Pow(math.Tan(P.Alpha/2), 2)
	P.Rn = math.MaxFloat64
	if P.N != 0.0 {
		P.Rn = 1 / P.N
	}

	/* ...and a few more */
	if 0 == P.B {
		P.B = (1 - P.F) * P.A
	}
	P.Rb = 1. / P.B
	P.Ra = 1. / P.A

	P.OneEs = 1. - P.Es
	if P.OneEs == 0. {
		return merror.New(merror.ErrEccentricityIsOne)
	}

	P.ROneEs = 1. / P.OneEs

	return nil
}

func (e *Ellipsoid) doEllps(ps *support.ProjString) error {

	/* Sail home if ellps=xxx is not specified */
	name, ok := ps.GetAsString("ellps")
	if !ok {
		return nil
	}

	/* Then look up the right size and shape parameters from the builtin list */
	if name == "" {
		return merror.New(merror.ErrInvalidArg)
	}

	ellps, ok := EllipsoidTable[name]
	if !ok {
		return merror.New(merror.UnknownEllipseParameter, name)
	}

	e.ID = ellps.ID
	e.Major = ellps.Major
	e.Ell = ellps.Ell
	e.Name = ellps.Name

	pl, err := support.NewProjString(ellps.Ell + " " + ellps.Major)
	if err != nil {
		panic(err)
	}
	ps.AddList(pl)

	err = e.doSize(ps)
	if err != nil {
		return err
	}

	err = e.doShape(ps)
	if err != nil {
		return err
	}

	return nil
}

func (e *Ellipsoid) doSize(ps *support.ProjString) error {

	P := e

	aWasSet := false

	/* A size parameter *must* be given, but may have been given as ellps prior */
	if P.A != 0.0 {
		aWasSet = true
	}

	/* Check which size key is specified */
	key := "R"
	value, ok := ps.GetAsFloat("R")
	if !ok {
		key = "a"
		value, ok = ps.GetAsFloat("a")
	}
	if !ok {
		if aWasSet {
			return nil
		}
		return merror.New(merror.MajorAxisNotGiven)
	}

	P.DefSize = key
	P.A = value
	if P.A <= 0.0 {
		return merror.New(merror.MajorAxisNotGiven)
	}
	if P.A == math.MaxFloat64 {
		return merror.New(merror.MajorAxisNotGiven)
	}

	if key == "R" {
		P.Es = 0
		P.F = 0
		P.E = 0
		P.Rf = 0
		P.B = P.A
	}

	return nil
}

func (e *Ellipsoid) doShape(ps *support.ProjString) error {

	P := e

	keys := []string{"rf", "f", "es", "e", "b"}

	/* Check which shape key is specified */
	var key string
	found := false
	var foundValue float64
	for _, key = range keys {
		value, ok := ps.GetAsFloat(key)
		if ok {
			found = true
			foundValue = value
			break
		}
	}

	/* Not giving a shape parameter means selecting a sphere, unless shape */
	/* has been selected previously via ellps=xxx */
	if !found && P.Es != 0 {
		return nil
	}
	if !found && P.Es == 0 {
		P.Es = 0
		P.F = 0
		P.B = P.A
		return nil
	}

	P.Es = 0
	P.F = 0
	P.B = 0
	P.E = 0
	P.Rf = 0

	switch key {

	/* reverse flattening, rf */
	case "rf":
		P.Rf = foundValue
		if P.Rf == math.MaxFloat64 {
			return merror.New(merror.ErrInvalidArg)
		}
		if P.Rf == 0 {
			return merror.New(merror.ReverseFlatteningIsZero)
		}
		P.F = 1 / P.Rf
		P.Es = 2*P.F - P.F*P.F

	/* flattening, f */
	case "f":
		P.F = foundValue
		if P.F == math.MaxFloat64 {
			return merror.New(merror.ErrInvalidArg)
		}
		if P.F == 0 {
			return merror.New(merror.ErrInvalidArg)
		}
		P.Rf = 1 / P.F
		P.Es = 2*P.F - P.F*P.F

	/* eccentricity squared, es */
	case "es":
		P.Es = foundValue
		if P.Es == math.MaxFloat64 {
			return merror.New(merror.ErrInvalidArg)
		}
		if P.Es == 1 {
			return merror.New(merror.ErrEccentricityIsOne)
		}

	/* eccentricity, e */
	case "e":
		P.E = foundValue
		if P.E == math.MaxFloat64 {
			return merror.New(merror.ErrInvalidArg)
		}
		if P.E == 0 {
			return merror.New(merror.ErrInvalidArg)
		}
		if P.E == 1 {
			return merror.New(merror.ErrEccentricityIsOne)
		}
		P.Es = P.E * P.E

	/* semiminor axis, b */
	case "b":
		P.B = foundValue
		if P.B == math.MaxFloat64 {
			return merror.New(merror.ErrInvalidArg)
		}
		if P.B == 0 {
			return merror.New(merror.ErrEccentricityIsOne)
		}
		if P.B == P.A {
			break
		}
		P.F = (P.A - P.B) / P.A
		P.Es = 2*P.F - P.F*P.F

	default:
		return merror.New(merror.ErrInvalidArg)

	}

	if P.Es < 0 {
		return merror.New(merror.ErrEsLessThanZero)
	}

	return nil
}

func (e *Ellipsoid) doSpherification(ps *support.ProjString) error {

	P := e

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
		if ps.ContainsKey(key) {
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
		P.A *= 1. - P.Es*(SIXTH+P.Es*(RA4+P.Es*RA6))

	/* R_V - a sphere with same volume as ellipsoid */
	case "R_V":
		P.A *= 1. - P.Es*(SIXTH+P.Es*(RV4+P.Es*RV6))

	/* R_a - a sphere with R = the arithmetic mean of the ellipsoid */
	case "R_a":
		P.A = (P.A + P.B) / 2

	/* R_g - a sphere with R = the geometric mean of the ellipsoid */
	case "R_g":
		P.A = math.Sqrt(P.A * P.B)

	/* R_h - a sphere with R = the harmonic mean of the ellipsoid */
	case "R_h":
		if P.A+P.B == 0 {
			return merror.New(merror.ErrToleranceCondition)
		}
		P.A = (2 * P.A * P.B) / (P.A + P.B)

		/* R_lat_a - a sphere with R = the arithmetic mean of the ellipsoid at given latitude */
		/* R_lat_g - a sphere with R = the geometric  mean of the ellipsoid at given latitude */
	case "R_lat_a", "R_lat_g":
		v, ok := ps.GetAsString(key)
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
		t = 1 - P.Es*t*t
		if key == "R_lat_a" { /* arithmetic */
			P.A *= (1. - P.Es + t) / (2 * t * math.Sqrt(t))
		} else { /* geometric */
			P.A *= math.Sqrt(1-P.Es) / t
		}

	default:
		return merror.New(merror.ErrInvalidArg)

	}

	/* Clean up the ellipsoidal parameters to reflect the sphere */
	P.Es = 0
	P.E = 0
	P.F = 0
	P.Rf = math.MaxFloat64
	P.B = P.A

	err := e.doCalcParams(P.A, 0)
	if err != nil {
		return err
	}

	return nil
}
