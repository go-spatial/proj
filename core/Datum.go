package core

import (
	"github.com/go-spatial/proj4go/support"
)

// Datum contains information about one datum ellipse
type Datum struct {
	ID               string
	DefinitionString string
	EllipseID        string
	Comments         string
	Definition       *support.PairList
}
