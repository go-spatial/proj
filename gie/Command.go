package gie

import (
	"math"
	"strconv"
	"strings"

	"github.com/go-spatial/proj4go/core"
	"github.com/go-spatial/proj4go/mlog"
	"github.com/go-spatial/proj4go/support"
)

type coord struct {
	a, b, c, d float64 // lam,phi or x,y
}

type testcase struct {
	inv    bool
	accept coord
	expect coord
}

// Command holds a set of tests as we build them up
type Command struct {
	proj            string
	delta           float64
	testcases       []testcase
	invFlag         bool
	completeFailure bool
}

// NewCommand returns a command
func NewCommand(ps string) *Command {
	c := &Command{
		proj:      ps,
		testcases: []testcase{},
	}
	//mlog.Printf("OPERATION: %s", ps)
	return c
}

func (c *Command) setDirection(s1 string) {
	if s1 == "inverse" {
		c.invFlag = true
	} else if s1 == "forward" {
		c.invFlag = false
	} else {
		panic(s1)
	}
}

func (c *Command) setAccept(s1, s2, s3, s4 string) {
	s1 = strings.Replace(s1, "_", "", -1)
	s2 = strings.Replace(s2, "_", "", -1)
	s3 = strings.Replace(s3, "_", "", -1)
	s4 = strings.Replace(s4, "_", "", -1)
	v1, err := strconv.ParseFloat(s1, 64)
	if err != nil {
		panic(err)
	}
	v2, err := strconv.ParseFloat(s2, 64)
	if err != nil {
		panic(err)
	}
	v3, err := strconv.ParseFloat(s3, 64)
	if err != nil {
		panic(err)
	}
	v4, err := strconv.ParseFloat(s4, 64)
	if err != nil {
		panic(err)
	}

	tc := testcase{
		accept: coord{v1, v2, v3, v4},
		inv:    c.invFlag,
	}

	c.testcases = append(c.testcases, tc)
}

func (c *Command) setExpectFailure() {
	n := len(c.testcases)
	if n == 0 {
		c.completeFailure = true
	} else {
		tc := &c.testcases[n-1]
		tc.expect = coord{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64, math.MaxFloat64}
	}
}

func (c *Command) setExpect(s1, s2, s3, s4 string) {
	s1 = strings.Replace(s1, "_", "", -1)
	s2 = strings.Replace(s2, "_", "", -1)
	s3 = strings.Replace(s3, "_", "", -1)
	s4 = strings.Replace(s4, "_", "", -1)
	v1, err := strconv.ParseFloat(s1, 64)
	if err != nil {
		panic(err)
	}
	v2, err := strconv.ParseFloat(s2, 64)
	if err != nil {
		panic(err)
	}
	v3, err := strconv.ParseFloat(s3, 64)
	if err != nil {
		panic(err)
	}
	v4, err := strconv.ParseFloat(s4, 64)
	if err != nil {
		panic(err)
	}

	tc := &c.testcases[len(c.testcases)-1]
	tc.expect = coord{v1, v2, v3, v4}
}

func (c *Command) setTolerance(s1, s2 string) {
	//mlog.Printf("TOLERANCE: %s %s", s1, s2)
	v, err := strconv.ParseFloat(s1, 64)
	if err != nil {
		panic(err)
	}

	c.delta = v

	switch s2 {
	case "*":
		c.delta /= 1.0
	case "cm":
		c.delta /= 100.0
	case "nm":
		c.delta /= 1.0e9
	case "um":
		c.delta /= 1.0e6
	case "mm":
		c.delta /= 1000.0
	case "m":
		c.delta /= 1.0
	default:
		panic(s2)
	}
}

func (c *Command) executeAll() error {

	ps, err := support.NewProjString(c.proj)
	if err != nil {
		return err
	}

	_, opx, err := core.NewSystem(ps)
	if err != nil {
		return err
	}

	op := opx.(core.IConvertLPToXY)

	for _, tc := range c.testcases {
		if !tc.inv {
			input := &core.CoordLP{Lam: support.DDToR(tc.accept.a), Phi: support.DDToR(tc.accept.b)}
			output, err := op.Forward(input)
			if err != nil {
				return err
			}

			x, y := output.X, output.Y
			check(tc.expect.a, x, c.delta)
			check(tc.expect.b, y, c.delta)
		}
	}

	return nil
}

func check(expect, actual, delta float64) {
	diff := math.Abs(expect - actual)
	if diff > delta {
		mlog.Printf("FAIL")
	}
}
