package core

import (
	"encoding/json"

	"github.com/go-spatial/proj4go/merror"
	"github.com/go-spatial/proj4go/support"
)

// ProjString represents a "projection string", such as "+proj=utm +zone=11 +datum=WGS84"
// TODO: we don't support the "pipeline" or "step" keywords
type ProjString struct {
	Source string
	Args   *support.PairList
}

// NewProjString returns a new ProjString object representing the given string
func NewProjString(source string) (*ProjString, error) {
	ps := &ProjString{
		Source: source,
	}

	pairs, err := support.NewPairListFromString(source)
	if err != nil {
		return nil, err
	}

	ps.Args = pairs

	err = ps.validate()
	if err != nil {
		return nil, err
	}

	return ps, nil
}

func (ps *ProjString) String() string {
	b, err := json.MarshalIndent(ps, "", "    ")
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (ps *ProjString) validate() error {

	// TODO: we don't support +init or +pipeline
	if ps.Args.CountKey("init") > 0 {
		return merror.New(merror.BadProjStringError)
	}
	if ps.Args.CountKey("pipeline") > 0 {
		return merror.New(merror.BadProjStringError)
	}

	if ps.Args.CountKey("proj") != 1 {
		return merror.New(merror.BadProjStringError)
	}
	projName, ok := ps.Args.GetAsString("proj")
	if !ok || projName == "" {
		return merror.New(merror.ProjValueMissing)
	}

	return nil
}
