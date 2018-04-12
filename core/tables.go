package core

import (
	"github.com/go-spatial/proj4go/support"
	"github.com/go-spatial/proj4go/tables"
)

// DatumTable is the global list of all the known datums
var DatumTable map[string]*Datum

// EllipsoidTable is the global list of all the known datums
var EllipsoidTable map[string]*tables.RawEllipsoid

// PrimeMeridianTable is the global list of all the known datums
var PrimeMeridianTable map[string]*PrimeMeridian

// OperationTable is the global list of all the known operations
var OperationTable = map[string]*Operation{}

// UnitInfoTable is the global list of all the known units
var UnitInfoTable = map[string]*tables.RawUnit{}

func init() {

	//
	// Datum
	//
	DatumTable = map[string]*Datum{}

	for _, raw := range tables.RawDatums {
		pl, err := support.NewProjString(raw.DefinitionString)
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
	EllipsoidTable = map[string]*tables.RawEllipsoid{}

	for _, raw := range tables.RawEllipsoids {

		EllipsoidTable[raw.ID] = raw
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

	//
	// Units
	//
	UnitInfoTable = map[string]*tables.RawUnit{}

	for _, raw := range tables.RawUnits {
		UnitInfoTable[raw.ID] = raw
	}
}
