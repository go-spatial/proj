package core

import (
	"math"

	"github.com/go-spatial/proj4go/merror"
	"github.com/go-spatial/proj4go/support"
)

func (ps *ProjString) processEllipsoid(P *Projection) error {

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

func ellpsSize(ps *ProjString, P *Projection) error {

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

func pjCalcEllipsoidParams(P *Projection, a float64, es float64) error {

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

func ellpsEllps(ps *ProjString, P *Projection) error {

	/* Sail home if ellps=xxx is not specified */
	name, ok := ps.Args.GetAsString("ellps")
	if !ok {
		return nil
	}

	/* Then look up the right size and shape parameters from the builtin list */
	if name == "" {
		return merror.New(merror.ErrInvalidArg)
	}

	ellps := pjFindEllps(name)
	if ellps == nil {
		return merror.New(merror.ErrUnknownEllpParam)
	}

	ellpsSize(ps, P)
	ellpsShape(ps, P)

	return nil
}

func ellpsShape(ps *ProjString, P *Projection) error {

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

func ellpsSpherification(ps *ProjString, P *Projection) error {

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

//-----------------------------------------------------------------

// Ellipsoid represents the parameters of an ellipsoid
type Ellipsoid struct {
	id    string
	major string
	ell   string
	name  string
}

// EllipsoidTable is an array of Ellipsoid objects
type EllipsoidTable []*Ellipsoid

// Ellipsoids is the global table of ellipsoids
var Ellipsoids = &EllipsoidTable{
	{"MERIT", "a=6378137.0", "rf=298.257", "MERIT 1983"},
	{"SGS85", "a=6378136.0", "rf=298.257", "Soviet Geodetic System 85"},
	{"GRS80", "a=6378137.0", "rf=298.257222101", "GRS 1980(IUGG, 1980)"},
	{"IAU76", "a=6378140.0", "rf=298.257", "IAU 1976"},
	{"airy", "a=6377563.396", "b=6356256.910", "Airy 1830"},
	{"APL4.9", "a=6378137.0.", "rf=298.25", "Appl. Physics. 1965"},
	{"NWL9D", "a=6378145.0.", "rf=298.25", "Naval Weapons Lab., 1965"},
	{"mod_airy", "a=6377340.189", "b=6356034.446", "Modified Airy"},
	{"andrae", "a=6377104.43", "rf=300.0", "Andrae 1876 (Den., Iclnd.)"},
	{"danish", "a=6377019.2563", "rf=300.0", "Andrae 1876 (Denmark, Iceland)"},
	{"aust_SA", "a=6378160.0", "rf=298.25", "Australian Natl & S. Amer. 1969"},
	{"GRS67", "a=6378160.0", "rf=298.2471674270", "GRS 67(IUGG 1967)"},
	{"GSK2011", "a=6378136.5", "rf=298.2564151", "GSK-2011"},
	{"bessel", "a=6377397.155", "rf=299.1528128", "Bessel 1841"},
	{"bess_nam", "a=6377483.865", "rf=299.1528128", "Bessel 1841 (Namibia)"},
	{"clrk66", "a=6378206.4", "b=6356583.8", "Clarke 1866"},
	{"clrk80", "a=6378249.145", "rf=293.4663", "Clarke 1880 mod."},
	{"clrk80ign", "a=6378249.2", "rf=293.4660212936269", "Clarke 1880 (IGN)."},
	{"CPM", "a=6375738.7", "rf=334.29", "Comm. des Poids et Mesures 1799"},
	{"delmbr", "a=6376428.", "rf=311.5", "Delambre 1810 (Belgium)"},
	{"engelis", "a=6378136.05", "rf=298.2566", "Engelis 1985"},
	{"evrst30", "a=6377276.345", "rf=300.8017", "Everest 1830"},
	{"evrst48", "a=6377304.063", "rf=300.8017", "Everest 1948"},
	{"evrst56", "a=6377301.243", "rf=300.8017", "Everest 1956"},
	{"evrst69", "a=6377295.664", "rf=300.8017", "Everest 1969"},
	{"evrstSS", "a=6377298.556", "rf=300.8017", "Everest (Sabah & Sarawak)"},
	{"fschr60", "a=6378166.", "rf=298.3", "Fischer (Mercury Datum) 1960"},
	{"fschr60m", "a=6378155.", "rf=298.3", "Modified Fischer 1960"},
	{"fschr68", "a=6378150.", "rf=298.3", "Fischer 1968"},
	{"helmert", "a=6378200.", "rf=298.3", "Helmert 1906"},
	{"hough", "a=6378270.0", "rf=297.", "Hough"},
	{"intl", "a=6378388.0", "rf=297.", "International 1909 (Hayford)"},
	{"krass", "a=6378245.0", "rf=298.3", "Krassovsky, 1942"},
	{"kaula", "a=6378163.", "rf=298.24", "Kaula 1961"},
	{"lerch", "a=6378139.", "rf=298.257", "Lerch 1979"},
	{"mprts", "a=6397300.", "rf=191.", "Maupertius 1738"},
	{"new_intl", "a=6378157.5", "b=6356772.2", "New International 1967"},
	{"plessis", "a=6376523.", "b=6355863.", "Plessis 1817 (France)"},
	{"PZ90", "a=6378136.0", "rf=298.25784", "PZ-90"},
	{"SEasia", "a=6378155.0", "b=6356773.3205", "Southeast Asia"},
	{"walbeck", "a=6376896.0", "b=6355834.8467", "Walbeck"},
	{"WGS60", "a=6378165.0", "rf=298.3", "WGS 60"},
	{"WGS66", "a=6378145.0", "rf=298.25", "WGS 66"},
	{"WGS72", "a=6378135.0", "rf=298.26", "WGS 72"},
	{"WGS84", "a=6378137.0", "rf=298.257223563", "WGS 84"},
	{"sphere", "a=6370997.0", "b=6370997.0", "Normal Sphere (r=6370997)"},
}

func pjFindEllps(name string) *Ellipsoid {
	for _, e := range *Ellipsoids {
		if name == e.id {
			return e
		}
	}
	return nil
}
