package core

// CreatorFuncType is the type of the function which creates an operation-specific object
//
// This returns an IOperation interface, which can be downcasted to something that implements
// the operation's specific signature, such as an IConvertLPToXY.
type CreatorFuncType func(*System, *OperationDescription) (IOperation, error)

// OperationDescription stores the information about a particular kind of
// operation. It is populated from each op in the "operations" package
// into the global table.
type OperationDescription struct {
	ID            string
	Description   string
	Description2  string
	OperationType OperationType
	InputType     CoordType
	OutputType    CoordType
	creatorFunc   CreatorFuncType
}

// RegisterOperation adds an OperationDescription entry to the OperationDescriptionTable
func RegisterOperation(
	id string,
	description string,
	description2 string,
	operationType OperationType,
	inputType CoordType,
	outputType CoordType,
	creatorFunc CreatorFuncType,
) {
	pi := &OperationDescription{
		ID:            id,
		Description:   description,
		Description2:  description2,
		OperationType: operationType,
		InputType:     inputType,
		OutputType:    outputType,
		creatorFunc:   creatorFunc,
	}

	_, ok := OperationDescriptionTable[id]
	if ok {
		panic(99)
	}
	OperationDescriptionTable[id] = pi
}

// Create returns a new object of the specific operation type, e.g. an operations.EtMerc
func (desc *OperationDescription) Create(sys *System) (IOperation, error) {
	f := desc.creatorFunc

	specificOperation, err := f(sys, desc)
	if err != nil {
		return nil, err
	}

	return specificOperation, nil
}
