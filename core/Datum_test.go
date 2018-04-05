package core_test

import (
	"testing"

	"github.com/go-spatial/proj4go/core"
	"github.com/stretchr/testify/assert"
)

func TestDatumInit(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(core.DatumTable)
	assert.True(len(core.DatumTable) > 5)

	assert.Equal("NAD83", core.DatumTable["GRS80"].ID)
	assert.Equal("GRS80", core.DatumTable["GRS80"].EllipseID)
	assert.Equal("towgs84=0,0,0", core.DatumTable["GRS80"].DefinitionString)
	assert.Equal("North_American_Datum_1983", core.DatumTable["GRS80"].Comments)
}
