package merror

// All the errors
var (
	BadProjStringError        = "bad proj string"
	ProjValueMissing          = "proj value missing in string"
	NoSuchDatum               = "no such datum"
	NotYetSupported           = "not yet supported" // TODO
	ErrMajorAxisNotGiven      = "ERR_MAJOR_AXIS_NOT_GIVEN"
	ErrEccentricityIsOne      = "ERR_ECCENTRICITY_IS_ONE"
	ErrInvalidArg             = "ERR_INVALID_ARG"
	ErrUnknownEllpParam       = "ERR_UNKNOWN_ELLP_PARAM"
	ErrRevFlatteningIsZero    = "ERR_REV_FLATTENING_IS_ZERO"
	ErrEsLessThanZero         = "ERR_ES_LESS_THAN_ZERO"
	ErrToleranceCondition     = "ERR_TOLERANCE_CONDITION"
	ErrRefRadLargerThan90     = "ERR_REF_RAD_LARGER_THAN_90"
	ErrInvalidDMS             = "ErrInvalidDMS"
	ErrEllipsoidUseRequired   = "ERR_ELLIPSOID_USE_REQUIRED"
	ErrInvalidUTMZone         = "ERR_INVALID_UTM_ZONE"
	ErrLatOrLonExceededLimit  = "ERR_LAT_OR_LON_EXCEED_LIMIT"
	ErrUnknownUnit            = "ERR_UNKNOWN_UNIT_ID"
	ErrUnitFactorLessThanZero = "ERR_UNIT_FACTOR_LESS_THAN_0"
	ErrAxis                   = "ERR_AXIS"
	ErrKLessThanZero          = "ERR_K_LESS_THAN_ZERO"
	ErrCoordinateError        = "ErrCoordinateError"
	ErrInvalidXOrY            = "ERR_INVALID_X_OR_Y"
)
