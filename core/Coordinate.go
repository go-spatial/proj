package core

// CoordType is the enum for the differetn intepretations of a Coordinate object
type CoordType int

// The coordinate type
const (
	CoordTypeXYZT = iota
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

// CoordXYZT is X,Y,Z,T
type CoordXYZT struct{ x, y, z, t float64 }

// CoordUVWT is U,V,W,T
type CoordUVWT struct{ u, v, w, t float64 }

// CoordLPZT is Lam,Phi,Z,T
type CoordLPZT struct{ lam, phi, z, t float64 }

// CoordOPK is Omega, Phi, Kappa (rotations)
type CoordOPK struct{ O, P, K float64 }

// CoordENU is East, North, Up
type CoordENU struct{ E, N, U float64 }

// CoordGEOD is geodesic length, fwd azi, rev azi
type CoordGEOD struct{ s, a1, a2 float64 }

// CoordUV is U,V
type CoordUV struct{ u, v float64 }

// CoordXY is X,Y
type CoordXY struct{ x, y float64 }

// CoordLP is Lam,Phi
type CoordLP struct{ lam, phi float64 }

// CoordXYZ is X,Y,Z
type CoordXYZ struct{ x, y, z float64 }

// CoordUVW is U,V,W
type CoordUVW struct{ u, v, w float64 }

// CoordLPZ is Lam, Phi, Z
type CoordLPZ struct{ lam, phi, z float64 }
