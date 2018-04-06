package core

// ForwardTransform is the function type of the forward operations
type ForwardTransform func(*Operation, interface{}) (interface{}, error)

// InverseTransform is the function type of the inverse operations
type InverseTransform func(*Operation, interface{}) (interface{}, error)

// SetupTransform is the function type of the setup/init fuction for this operation type
type SetupTransform func(*Operation) error

// OperationInfo stores the information about a particular kind of
// operation. It is populated from the "projections" package.
type OperationInfo struct {
	ID           string
	Description  string
	Description2 string
	Forward      ForwardTransform
	Inverse      InverseTransform
	Setup        SetupTransform
}

// RegisterOperation adds an OperationInfo entry to the OperationInfoTable
func RegisterOperation(
	id string,
	description string,
	description2 string,
	input CoordType,
	output CoordType,
	forward ForwardTransform,
	inverse InverseTransform,
	setup SetupTransform,
) {
	pi := &OperationInfo{
		ID:           id,
		Description:  description,
		Description2: description2,
		Forward:      forward,
		Inverse:      inverse,
		Setup:        setup,
	}

	_, ok := OperationInfoTable[id]
	if ok {
		panic(99)
	}
	OperationInfoTable[id] = pi
}
