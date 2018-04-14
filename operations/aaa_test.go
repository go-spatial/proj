package operations_test

import (
	"fmt"
	"testing"

	"github.com/go-spatial/proj4go/core"
	"github.com/go-spatial/proj4go/mlog"
	_ "github.com/go-spatial/proj4go/operations"
	"github.com/go-spatial/proj4go/support"
	"github.com/stretchr/testify/assert"
)

type data struct {
	proj     string
	lam, phi float64
	x, y     float64
}

var testdata = []data{
	{"+proj=utm +ellps=GRS80  +lat_1=0.5 +lat_2=2 +n=0.5 +zone=30", 2, 1, 1057002.405491298, 110955.141175949},
	{"+proj=etmerc +ellps=GRS80 +ellps=GRS80 +lat_1=0.5 +lat_2=2 +n=0.5 +zone=30", 2, 1, 222650.796797586, 110642.229411933},
	{"+proj=aea +ellps=GRS80 +lat_1=0 +lat_2=2", 2.0, 1.0, 222571.608757106, 110653.326743030},
	{"+proj=leac +ellps=GRS80 +lat_1=0 +lat_2=2", 2.0, 1.0, 220685.140542979, 112983.500889396},
	{"+proj=merc +ellps=GRS80 +lat_1=0.5 +lat_2=2", 2, 1, 222638.981586547, 110579.965218250},
}

func TestConvertLPToXY(t *testing.T) {
	assert := assert.New(t)

	for i, d := range testdata {

		mssg := fmt.Sprintf("%d: %s", i, d.proj)
		mlog.Printf(mssg)

		ps, err := support.NewProjString(d.proj)
		assert.NoError(err)

		sys, opx, err := core.NewSystem(ps)
		assert.NoError(err)
		assert.NotNil(sys)
		assert.NotNil(opx)
		assert.EqualValues(sys, opx.GetSystem())

		op := opx.(core.IConvertLPToXY)

		input := &core.CoordLP{Lam: support.DDToR(d.lam), Phi: support.DDToR(d.phi)}
		output, err := op.Forward(input)
		assert.NoError(err)

		x, y := output.X, output.Y
		assert.InDelta(d.x, x, 1e-2, mssg)
		assert.InDelta(d.y, y, 1e-2, mssg)

		input2 := output
		output2, err := op.Inverse(input2)
		assert.NoError(err)

		l, p := output2.Lam, output2.Phi
		l = support.RToDD(l)
		p = support.RToDD(p)
		//assert.InDelta(d.lam, l, 1e-8, mssg)
		//assert.InDelta(d.phi, p, 1e-8, mssg)
	}
}
