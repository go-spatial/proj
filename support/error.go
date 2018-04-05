package support

// All the errors
var (
	BadProjStringError     = "bad proj string"
	ProjValueMissing       = "proj value missing in string"
	NoSuchDatum            = "no such datum"
	NotYetSupported        = "not yet supported" // TODO
	ErrMajorAxisNotGiven   = "ERR_MAJOR_AXIS_NOT_GIVEN"
	ErrEccentricityIsOne   = "ERR_ECCENTRICITY_IS_ONE"
	ErrInvalidArg          = "ERR_INVALID_ARG"
	ErrUnknownEllpParam    = "ERR_UNKNOWN_ELLP_PARAM"
	ErrRevFlatteningIsZero = "ERR_REV_FLATTENING_IS_ZERO"
	ErrEsLessThanZero      = "ERR_ES_LESS_THAN_ZERO"
	ErrToleranceCondition  = "ERR_TOLERANCE_CONDITION"
	ErrRefRadLargerThan90  = "ERR_REF_RAD_LARGER_THAN_90"
	ErrInvalidDMS          = "ErrInvalidDMS"
)