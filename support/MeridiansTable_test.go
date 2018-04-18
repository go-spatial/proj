package support_test

import (
	"testing"

	"github.com/go-spatial/proj/support"
	"github.com/stretchr/testify/assert"
)

func TestMeridiansTable(t *testing.T) {
	assert := assert.New(t)

	assert.True(len(support.MeridiansTable) > 10)
	for key, value := range support.MeridiansTable {
		assert.Equal(key, value.ID)
	}

	assert.Equal("paris", support.MeridiansTable["paris"].ID)
	assert.Equal("2d20'14.025\"E", support.MeridiansTable["paris"].Definition)
}
