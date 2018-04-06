package core

// ForwardTransform is the function type of the forward operations
type ForwardTransform interface{}

// InverseTransform is the function type of the inverse operations
type InverseTransform interface{}

// SetupTransform is the function type of the setup/init fuction for this operation type
type SetupTransform func() interface{}

// ProjectionInfo stores the information about a particular kind of
// projection. It is populated from the "projections" package.
type ProjectionInfo struct {
	ID           string
	Description  string
	Description2 string
	Forward      ForwardTransform
	Inverse      InverseTransform
	Setup        SetupTransform
}

// RegisterProjection adds a ProjectionInfo entry to the ProjectionTable
func RegisterProjection(
	id string,
	description string,
	description2 string,
	input CoordType,
	output CoordType,
	forward ForwardTransform,
	inverse InverseTransform,
	setup SetupTransform,
) {
	pi := &ProjectionInfo{
		ID:           id,
		Description:  description,
		Description2: description2,
		Forward:      forward,
		Inverse:      inverse,
		Setup:        setup,
	}

	_, ok := ProjectionTable[id]
	if !ok {
		panic(99)
	}
	ProjectionTable[id] = pi
}
