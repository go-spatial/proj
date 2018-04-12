package core

// ForwardTransform is the function type of the forward operations
type ForwardTransform func(*System, interface{}) (interface{}, error)

// InverseTransform is the function type of the inverse operations
type InverseTransform func(*System, interface{}) (interface{}, error)

// SetupTransform is the function type of the setup/init fuction for this operation type
type SetupTransform func(*System) error

// IOperation is for all operations
type IOperation interface {
	GetSystem() *System
	GetSignature() (CoordType, CoordType)
	ForwardAny(*System, *CoordAny) (*CoordAny, error)
	InverseAny(*System, *CoordAny) (*CoordAny, error)
	Setup(*System) error
}

// OperationCommon provides some common fields an implementation of IOperation will need
type OperationCommon struct {
	System     *System
	InputType  CoordType
	OutputType CoordType
}

// GetSystem returns the System
func (op OperationCommon) GetSystem() *System { return op.System }

// GetSignature returns the input type and output type
func (op OperationCommon) GetSignature() (CoordType, CoordType) { return op.InputType, op.OutputType }

// IConvertLPToXY is just for 2D LP->XY conversions
type IConvertLPToXY interface {
	Forward(*System, *CoordLP) (*CoordXY, error)
	Inverse(*System, *CoordXY) (*CoordLP, error)
}
