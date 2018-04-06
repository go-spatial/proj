package core

import (
	"github.com/go-spatial/proj4go/support"
	"github.com/go-spatial/proj4go/tables"
)

// DatumTable is the global list of all the known datums
var DatumTable map[string]*Datum

// EllipsoidTable is the global list of all the known datums
var EllipsoidTable map[string]*Ellipsoid

// PrimeMeridianTable is the global list of all the known datums
var PrimeMeridianTable map[string]*PrimeMeridian

// ProjectionTable is the global list of all the known projections
var ProjectionTable map[string]*ProjectionInfo

func init() {

	//
	// Datum
	//
	DatumTable = map[string]*Datum{}

	for _, raw := range tables.RawDatums {
		pl, err := support.NewPairListFromString(raw.DefinitionString)
		if err != nil {
			panic(err)
		}

		d := &Datum{
			ID:               raw.ID,
			DefinitionString: raw.DefinitionString,
			EllipseID:        raw.EllipseID,
			Comments:         raw.Comments,
			Definition:       pl,
		}

		DatumTable[d.EllipseID] = d
	}

	//
	// Ellipsoid
	//
	EllipsoidTable = map[string]*Ellipsoid{}

	for _, raw := range tables.RawEllipsoids {
		d := &Ellipsoid{
			ID:    raw.ID,
			Major: raw.Major,
			Ell:   raw.Ell,
			Name:  raw.Name,
		}

		EllipsoidTable[d.ID] = d
	}

	//
	// PrimeMeridian
	//
	PrimeMeridianTable = map[string]*PrimeMeridian{}

	for _, raw := range tables.RawPrimeMeridians {
		pm := &PrimeMeridian{
			ID:         raw.ID,
			Definition: raw.Definition,
		}
		PrimeMeridianTable[pm.ID] = pm
	}
}
