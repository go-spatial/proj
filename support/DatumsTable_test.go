package support_test

import (
	"testing"

	"github.com/go-spatial/proj4go/support"
	"github.com/stretchr/testify/assert"
)

func TestDatumsTable(t *testing.T) {
	assert := assert.New(t)

	assert.True(len(support.DatumsTable) > 8)

	for key, value := range support.DatumsTable {
		assert.True(key == value.ID || key == value.EllipseID)
	}

	assert.Equal("NAD83", support.DatumsTable["GRS80"].ID)
	assert.Equal("GRS80", support.DatumsTable["GRS80"].EllipseID)
	assert.Equal("towgs84=0,0,0", support.DatumsTable["GRS80"].DefinitionString)
	assert.Equal("North_American_Datum_1983", support.DatumsTable["GRS80"].Comments)
}
