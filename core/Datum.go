package core

import (
	"github.com/go-spatial/proj4go/support"
	"github.com/go-spatial/proj4go/tables"
)

// Datum contains information about one datum ellipse
type Datum struct {
	ID               string
	DefinitionString string
	EllipseID        string
	Comments         string
	Definition       *support.PairList
}

// DatumTable is the global list of all the known datums
var DatumTable map[string]*Datum

func init() {

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
}
