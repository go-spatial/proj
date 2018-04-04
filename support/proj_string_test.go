package support_test

import (
	"testing"

	"github.com/go-spatial/proj4go/support"
	"github.com/stretchr/testify/assert"
)

func TestProjString(t *testing.T) {
	assert := assert.New(t)

	ps, err := support.NewProjString("")
	assert.Error(err)

	ps, err = support.NewProjString("k1=a=b")
	assert.Error(err)

	kv0 := support.Pair{Key: "k0", Value: ""}
	kv1 := support.Pair{Key: "k1", Value: "a"}
	kv2 := support.Pair{Key: "k2", Value: "b"}
	//init1 := support.Pair{Key: "init", Value: "foo"}
	proj1 := support.Pair{Key: "proj", Value: "P99"}

	pl1 := support.NewPairList()
	pl1.Add(proj1)
	pl1.Add(kv0)

	pl2 := support.NewPairList()
	pl2.Add(kv1)
	pl2.Add(proj1)
	pl2.Add(kv2)

	pl3 := support.NewPairList()
	pl3.Add(kv1)
	pl3.Add(kv2)
	pl3.Add(proj1)

	pl4 := support.NewPairList()
	pl4.Add(proj1)
	pl4.Add(kv1)
	pl4.Add(kv2)

	ps, err = support.NewProjString("proj=P99 k0= ")
	assert.NoError(err)
	assert.Equal(pl1, ps.Args)

	ps, err = support.NewProjString("k1=a proj=P99 k2=b")
	assert.NoError(err)
	assert.Equal(pl2, ps.Args)

	ps, err = support.NewProjString("  k1=a   k2=b  proj=P99  \t  ")
	assert.NoError(err)
	assert.Equal(pl3, ps.Args)

	ps, err = support.NewProjString("  +proj=P99 +k1=a   +k2=b    \t  ")
	assert.NoError(err)
	assert.Equal(pl4, ps.Args)

	ps, err = support.NewProjString("init=foo") // TODO
	assert.Error(err)
	//assert.Equal([]support.KeyValue{init1}, ps.Args)

	ps, err = support.NewProjString("init=foo init=foo")
	assert.Error(err)
}
