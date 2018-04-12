package core

// ForwardTransform is the function type of the forward operations
type ForwardTransform func(*System, interface{}) (interface{}, error)

// InverseTransform is the function type of the inverse operations
type InverseTransform func(*System, interface{}) (interface{}, error)

// SetupTransform is the function type of the setup/init fuction for this operation type
type SetupTransform func(*System) error

// Operation stores the information about a particular kind of
// operation. It is populated from the "projections" package.
type Operation struct {
	ID           string
	Description  string
	Description2 string
	forward      ForwardTransform
	inverse      InverseTransform
	setup        SetupTransform
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
	pi := &Operation{
		ID:           id,
		Description:  description,
		Description2: description2,
		forward:      forward,
		inverse:      inverse,
		setup:        setup,
	}

	_, ok := OperationTable[id]
	if ok {
		panic(99)
	}
	OperationTable[id] = pi
}
