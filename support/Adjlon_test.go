package support_test

import (
	"math"
	"testing"

	"github.com/go-spatial/proj/support"
	"github.com/stretchr/testify/assert"
)

func TestAdjlon(t *testing.T) {
	assert := assert.New(t)

	assert.InDelta(math.Pi*0.5, support.Adjlon(math.Pi*2.5), 1.0e-8)
}
