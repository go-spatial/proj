package core

// NewFuncType is the type of the function which creates an operation-specific object
type NewFuncType func(*System) (IOperation, error)

// OperationDescription stores the information about a particular kind of
// operation. It is populated from each op in the "operations" package
// into the global table.
type OperationDescription struct {
	ID           string
	Description  string
	Description2 string
	InputType    CoordType
	OutputType   CoordType
	NewFunc      NewFuncType
}

// RegisterOperation adds an OperationDescription entry to the OperationDescriptionTable
func RegisterOperation(
	id string,
	description string,
	description2 string,
	input CoordType,
	output CoordType,
	newFunc NewFuncType,
) {
	pi := &OperationDescription{
		ID:           id,
		Description:  description,
		Description2: description2,
		InputType:    input,
		OutputType:   output,
		NewFunc:      newFunc,
	}

	_, ok := OperationDescriptionTable[id]
	if ok {
		panic(99)
	}
	OperationDescriptionTable[id] = pi
}
