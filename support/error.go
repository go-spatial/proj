package support

// Error is our rich error class
// implements "error"
type Error struct {
	Code string
}

func (e Error) Error() string {
	return e.Code
}

// All the errors
var (
	BadProjStringError     = Error{Code: "bad proj string"}
	ProjValueMissing       = Error{Code: "proj value missing in string"}
	NoSuchDatum            = Error{Code: "no such datum"}
	NotYetSupported        = Error{Code: "not yet supported"}
	ErrMajorAxisNotGiven   = Error{Code: "ERR_MAJOR_AXIS_NOT_GIVEN"}
	ErrEccentricityIsOne   = Error{Code: "ERR_ECCENTRICITY_IS_ONE"}
	ErrInvalidArg          = Error{Code: "ERR_INVALID_ARG"}
	ErrUnknownEllpParam    = Error{Code: "ERR_UNKNOWN_ELLP_PARAM"}
	ErrRevFlatteningIsZero = Error{Code: "ERR_REV_FLATTENING_IS_ZERO"}
	ErrEsLessThanZero      = Error{Code: "ERR_ES_LESS_THAN_ZERO"}
	ErrToleranceCondition  = Error{Code: "ERR_TOLERANCE_CONDITION"}
	ErrRefRadLargerThan90  = Error{Code: "ERR_REF_RAD_LARGER_THAN_90"}
	ErrInvalidDMS          = Error{Code: "ErrInvalidDMS"}
)
