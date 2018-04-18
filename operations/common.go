package operations

type mode int

const (
	modeNPole mode = 0
	modeSPole      = 1
	modeEquit      = 2
	modeObliq      = 3
)

const tol7 = 1.e-7
const tol10 = 1.0e-10
const tol14 = 1.0e-14

const eps7 = 1.0e-7
const eps10 = 1.e-10
