package core


// OperationType is the enum for the different kinds of conversions and transforms
type OperationType int

// The operation type
const (
	OperationTypeInvalid OperationType = iota
	OperationTypeConversion
	OperationTypeTransformation
)

// IOperation is for all the operation
type IOperation interface {
	GetSystem() *System
	GetDescription() *OperationDescription
}

// Operation is for all operations
type Operation struct {
	System      *System
	Description *OperationDescription
}

// GetSystem returns the system
func (op *Operation) GetSystem() *System {
	return op.System
}

// GetDescription returns the descr
func (op *Operation) GetDescription() *OperationDescription {
	return op.Description
}

// GetSignature returns the operation type, the input type, and the output type
func (op *Operation) GetSignature() (OperationType, CoordType, CoordType) {
	d := op.Description
	return d.OperationType, d.InputType, d.OutputType
}
