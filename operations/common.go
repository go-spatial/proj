package operations

type mode int

const (
	modeNPole mode = 0
	modeSPole      = 1
	modeEquit      = 2
	modeObliq      = 3
)

const tol14 = 1.0e-14
