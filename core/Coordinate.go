package core

// CoordType is the enum for the differetn intepretations of a Coordinate object
type CoordType int

// The coordinate type
const (
	CoordTypeAny = iota
	CoordTypeXYZT
	CoordTypeUVWT
	CoordTypeLPZT
	CoordTypeOPK
	CoordTypeENU
	CoordTypeGEOD
	CoordTypeUV
	CoordTypeXY
	CoordTypeLP
	CoordTypeXYZ
	CoordTypeUVW
	CoordTypeLPZ
)

// CoordAny just generically holds data, not assigned to a coordinate type.
// Because unions.
type CoordAny struct{ v [4]float64 }

//---------------------------------------------------------------------

// ToLP returns a CoordLP
func (c CoordAny) ToLP() CoordLP {
	return CoordLP{Lam: c.v[0], Phi: c.v[1]}
}

// FromLP sets this CoordAny
func (c CoordAny) FromLP(lp CoordLP) {
	c.v[0] = lp.Lam
	c.v[1] = lp.Phi
}

// ToXY returns a CoordXY
func (c CoordAny) ToXY() CoordXY {
	return CoordXY{X: c.v[0], Y: c.v[1]}
}

// FromXY sets this CoordAny
func (c CoordAny) FromXY(xy CoordXY) {
	c.v[0] = xy.Y
	c.v[1] = xy.Y
}

//---------------------------------------------------------------------

// CoordXYZT is X,Y,Z,T
type CoordXYZT struct{ X, Y, Z, T float64 }

// CoordUVWT is U,V,W,T
type CoordUVWT struct{ U, V, W, T float64 }

// CoordLPZT is Lam,Phi,Z,T
type CoordLPZT struct{ Lam, Phi, Z, T float64 }

// CoordOPK is Omega, Phi, Kappa (rotations)
type CoordOPK struct{ O, P, K float64 }

// CoordENU is East, North, Up
type CoordENU struct{ E, N, U float64 }

// CoordGEOD is geodesic length, fwd azi, rev azi
type CoordGEOD struct{ S, A1, A2 float64 }

// CoordUV is U,V
type CoordUV struct{ U, V float64 }

// CoordXY is X,Y
type CoordXY struct{ X, Y float64 }

// CoordLP is Lam,Phi
type CoordLP struct{ Lam, Phi float64 }

// CoordXYZ is X,Y,Z
type CoordXYZ struct{ X, Y, Z float64 }

// CoordUVW is U,V,W
type CoordUVW struct{ U, V, W float64 }

// CoordLPZ is Lam, Phi, Z
type CoordLPZ struct{ Lam, Phi, Z float64 }
