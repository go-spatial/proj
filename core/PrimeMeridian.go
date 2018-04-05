package core

import "github.com/go-spatial/proj4go/tables"

// PrimeMeridian contains information about a prime meridian
type PrimeMeridian struct {
	ID         string
	Definition string
}

// PrimeMeridianTable is the global list of all the known datums
var PrimeMeridianTable map[string]*PrimeMeridian

func init() {
	PrimeMeridianTable = map[string]*PrimeMeridian{}

	for _, raw := range tables.RawPrimeMeridians {
		pm := &PrimeMeridian{
			ID:         raw.ID,
			Definition: raw.Definition,
		}
		PrimeMeridianTable[pm.ID] = pm
	}
}
