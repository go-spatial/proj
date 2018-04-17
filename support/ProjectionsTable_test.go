package support_test

import (
	"testing"

	"github.com/go-spatial/proj4go/support"
	"github.com/stretchr/testify/assert"
)

func TestProjectionsTable(t *testing.T) {
	assert := assert.New(t)

	assert.True(len(support.ProjectionsTable) > 150)
	for key, value := range support.ProjectionsTable {
		assert.Equal(key, value.ID)
	}
}
