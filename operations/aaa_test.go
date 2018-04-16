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

type coord struct {
	a, b float64 // lam,phi or x,y
}

type testcase struct {
	accept []float64
	expect []float64
}

type data struct {
	proj  string
	delta float64
	fwd   []testcase
	inv   []testcase
}

var testdata = []data{
	{
		proj:  "+proj=aea   +ellps=GRS80  +lat_1=0 +lat_2=2",
		delta: 0.1 * 0.001,
		fwd: []testcase{
			{
				accept: coord{2, 1},
				expect: coord{1057002.405491298, 110955.141175949},
			}, {
				accept: {2, 1},
				expect: {222571.608757106, 110653.326743030},
			}, {
				accept: {2, -1},
				expect: {222706.306508391, -110484.267144400},
			}, {
				accept: {-2, 1},
				expect: {-222571.608757106, 110653.326743030},
			}, {
				accept: {-2, -1},
				expect: {-222706.306508391, -110484.267144400},
			},
		},
		//			direction inverse
		//			accept  200 100
		//			expect  0.001796631 0.000904369
		//			accept  200 -100
		//			expect  0.001796630 -0.000904370
		//			accept  -200 100
		//			expect  -0.001796631 0.000904369
		//			accept  -200 -100
		//			expect  -0.001796630 -0.000904370
		//				},
		inv: []testcase{},
	},
	//{"+proj=etmerc +ellps=GRS80 +ellps=GRS80 +lat_1=0.5 +lat_2=2 +n=0.5 +zone=30", 2, 1, 222650.796797586, 110642.229411933},
	//{"+proj=aea +ellps=GRS80 +lat_1=0 +lat_2=2", 2.0, 1.0, 222571.608757106, 110653.326743030},
	//{"+proj=leac +ellps=GRS80 +lat_1=0 +lat_2=2", 2.0, 1.0, 220685.140542979, 112983.500889396},
	//{"+proj=merc +ellps=GRS80 +lat_1=0.5 +lat_2=2", 2, 1, 222638.981586547, 110579.965218250},
}

func TestConvert(t *testing.T) {
	assert := assert.New(t)

	for i, td := range testdata {

		mssg := fmt.Sprintf("%d: %s ", i, td.proj)
		mlog.Printf(mssg)

		ps, err := support.NewProjString(td.proj)
		assert.NoError(err)

		sys, opx, err := core.NewSystem(ps)
		assert.NoError(err)
		assert.NotNil(sys)
		assert.NotNil(opx)
		assert.EqualValues(sys, opx.GetSystem())

		op := opx.(core.IConvertLPToXY)

		for _, tc := range td.fwd {
			input := &core.CoordLP{Lam: support.DDToR(tc.in1), Phi: support.DDToR(tc.in2)}
			output, err := op.Forward(input)
			assert.NoError(err)

			x, y := output.X, output.Y
			assert.InDelta(tc.out1, x, td.delta, mssg+"fwd")
			assert.InDelta(tc.out2, y, td.delta, mssg+"fwd")
		}

		for _, tc := range td.inv {
			input := &core.CoordXY{X: tc.in1, Y: tc.in2}
			output, err := op.Inverse(input)
			assert.NoError(err)

			l, p := output.Lam, output.Phi
			assert.InDelta(tc.out1, support.RToDD(l), td.delta, mssg+"inv")
			assert.InDelta(tc.out2, support.RToDD(p), td.delta, mssg+"inv")
		}
	}
}
