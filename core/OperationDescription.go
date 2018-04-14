package core

import (
	"github.com/go-spatial/proj4go/merror"
)

// ConvertLPToXYCreatorFuncType is the type of the function which creates an operation-specific object
//
// This kind of function, when executed, creates an operation-specific type
// which implements IConvertLPToXY.
type ConvertLPToXYCreatorFuncType func(*System, *OperationDescription) (IConvertLPToXY, error)

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
	creatorFunc   interface{} // for now, this will always be a ConvertLPToXYCreatorFuncType
}

// RegisterConvertLPToXY adds an OperationDescription entry to the OperationDescriptionTable
func RegisterConvertLPToXY(
	id string,
	description string,
	description2 string,
	creatorFunc ConvertLPToXYCreatorFuncType,
) {
	pi := &OperationDescription{
		ID:            id,
		Description:   description,
		Description2:  description2,
		OperationType: OperationTypeConversion,
		InputType:     CoordTypeLP,
		OutputType:    CoordTypeXY,
		creatorFunc:   creatorFunc,
	}

	_, ok := OperationDescriptionTable[id]
	if ok {
		panic(99)
	}
	OperationDescriptionTable[id] = pi
}

// CreateOperation returns a new object of the specific operation type, e.g. an operations.EtMerc
func (desc *OperationDescription) CreateOperation(sys *System) (IOperation, error) {

	if desc.IsConvertLPToXY() {
		return NewConvertLPToXY(sys, desc)
	}

	return nil, merror.New(merror.NotYetSupported)
}

// IsConvertLPToXY returns true iff the operation can be casted to an IConvertLPToXY
func (desc *OperationDescription) IsConvertLPToXY() bool {
	return desc.OperationType == OperationTypeConversion &&
		desc.InputType == CoordTypeLP &&
		desc.OutputType == CoordTypeXY
}
