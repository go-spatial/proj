package support_test

import (
	"testing"

	"github.com/go-spatial/proj/support"
	"github.com/stretchr/testify/assert"
)

func TestEllipsoidsTable(t *testing.T) {
	assert := assert.New(t)

	assert.True(len(support.EllipsoidsTable) > 40)
	for key, value := range support.EllipsoidsTable {
		assert.Equal(key, value.ID)
	}

	assert.Equal("danish", support.EllipsoidsTable["danish"].ID)
	assert.Equal("a=6377019.2563", support.EllipsoidsTable["danish"].Major)
	assert.Equal("rf=300.0", support.EllipsoidsTable["danish"].Ell)
	assert.Equal("Andrae 1876 (Denmark, Iceland)", support.EllipsoidsTable["danish"].Name)
}
