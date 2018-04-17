package support_test

import (
	"testing"

	"github.com/go-spatial/proj4go/support"
	"github.com/stretchr/testify/assert"
)

func TestUnitsTable(t *testing.T) {
	assert := assert.New(t)

	assert.True(len(support.UnitsTable) > 15)
	for key, value := range support.UnitsTable {
		assert.Equal(key, value.ID)
	}
}
