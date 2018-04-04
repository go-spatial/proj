package support

// ConvertArcsecondsToRadians converts from arc secs to rads
func ConvertArcsecondsToRadians(s float64) float64 {
	// Pi/180/3600
	r := s * 4.84813681109535993589914102357e-6
	return r
}

// some more useful math constants and aliases */
const (
	mPi          = 3.14159265358979323846
	mPiOverTwo   = 1.57079632679489661923
	mPiOverFour  = 0.78539816339744830962
	mTwoOverPi   = 0.63661977236758134308
	mPiHalfPi    = 4.71238898038468985769 /* 1.5*pi */
	mTwoPI       = 6.28318530717958647693 /* 2*pi */
	mTwoPiHalfPi = 7.85398163397448309616 /* 2.5*pi */
)
