package merror

// All the errors
var (
	UnknownProjection               = "unknown projection: %s"
	UnknownEllipseParameter         = "unknown ellipse parameter: %s"
	UnsupportedProjectionString     = "unsupported projection string: %s"
	InvalidProjectionSyntax         = "invalid projection syntax: %s"
	ProjectionStringRequiresEllipse = "projection string requires ellipse"
	MajorAxisNotGiven               = "major axis not given"
	ReverseFlatteningIsZero         = "reverse flattening (rf) is zero"
	EccentricityIsOne               = "eccentricity is one"
	ToleranceCondition              = "tolerance condition error"

	ProjValueMissing          = "proj value missing in string"
	NoSuchDatum               = "no such datum"
	NotYetSupported           = "not yet supported" // TODO
	ErrInvalidArg             = "ERR_INVALID_ARG"
	ErrEsLessThanZero         = "ERR_ES_LESS_THAN_ZERO"
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
	ErrConicLatEqual          = "ErrConicLatEqual"
	ErrAeaSetupFailed         = "ErrAeaSetupFailed"
	ErrInvMlfn                = "ErrInvMlfn"
	ErrAeaProjString          = "ErrAeaProjString"
	ErrLatTSLargerThan90      = "ErrLatTSLargerThan90"
	ErrPhi2                   = "ErrPhi2"
)
