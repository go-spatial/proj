package merror_test

import (
	"testing"

	"github.com/go-spatial/proj/merror"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	assert := assert.New(t)

	err1 := merror.New("errtest-%d", 1)
	assert.Error(err1)
	exp1 := "errtest-1 (from merror_test.TestError at merror_test.go:13)"
	assert.Equal(exp1, err1.Error())

	err2 := merror.Wrap(err1, "errtest-%d", 2)
	assert.Error(err2)
	exp2 := "errtest-2 (from merror_test.TestError at merror_test.go:18)"
	exp2 += " // Inner: " + exp1
	assert.Equal(exp2, err2.Error())

	err3 := merror.Wrap(err2)
	assert.Error(err3)
	exp3 := "wrapped error (from merror_test.TestError at merror_test.go:24)"
	exp3 += " // Inner: " + exp2
	assert.Equal(exp3, err3.Error())
}
