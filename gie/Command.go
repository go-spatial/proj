package gie

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/go-spatial/proj/core"
	"github.com/go-spatial/proj/mlog"
	"github.com/go-spatial/proj/support"
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
	ProjString      string
	tolerance       float64
	testcases       []testcase
	invFlag         bool
	completeFailure bool
	File            string
	Line            int
	roundtripCount  int
	roundtripDelta  float64
}

// NewCommand returns a command
func NewCommand(file string, line int, ps string) *Command {
	c := &Command{
		ProjString: ps,
		testcases:  []testcase{},
		File:       file,
		Line:       line,
		tolerance:  0.5 * unitsValue("mm"),
	}
	//mlog.Printf("OPERATION: %s", ps)
	return c
}

// ProjectionName returns the name of the projection used in this test
func (c *Command) ProjectionName() string {
	s := c.ProjString
	for {
		t := strings.Replace(s, "\t", " ", -1)
		t = strings.Replace(t, "  ", " ", -1)
		t = strings.Replace(t, " =", "=", -1)
		t = strings.Replace(t, "= ", "=", -1)
		if s == t {
			break
		}
		s = t
	}

	toks := strings.Fields(s)
	for _, tok := range toks {
		if tok[0:1] == "+" {
			tok = tok[1:]
		}
		if strings.HasPrefix(tok, "proj=") {
			return tok[5:]
		}
	}
	return "UNKNOWN"
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

func (c *Command) setRoundtrip(s1, s2, s3 string) {
	count, err := strconv.Atoi(s1)
	if err != nil {
		panic(err)
	}
	v, err := strconv.ParseFloat(s2, 64)
	if err != nil {
		panic(err)
	}
	delta := v / unitsValue(s3)

	c.roundtripCount = count
	c.roundtripDelta = delta
}

func (c *Command) setTolerance(s1, s2 string) {
	//mlog.Printf("TOLERANCE: %s %s", s1, s2)
	v, err := strconv.ParseFloat(s1, 64)
	if err != nil {
		panic(err)
	}

	c.tolerance = v / unitsValue(s2)
}

func unitsValue(s string) float64 {
	switch s {
	case "*":
		return 1.0
	case "cm":
		return 100.0
	case "nm":
		return 1.0e9
	case "um":
		return 1.0e6
	case "mm":
		return 1000.0
	case "m":
		return 1.0
	}
	panic(s)
}

// Execute runs the tests
func (c *Command) Execute() error {

	ps, err := support.NewProjString(c.ProjString)
	if err != nil {
		if c.completeFailure {
			return nil
		}
		return err
	}

	_, opx, err := core.NewSystem(ps)
	if err != nil {
		if c.completeFailure {
			return nil
		}
		return err
	}

	op := opx.(core.IConvertLPToXY)

	for _, tc := range c.testcases {

		if !tc.inv {

			if c.roundtripCount == 0 {
				_, _, err = c.executeForwardOnce(
					tc.accept.a, tc.accept.b,
					tc.expect.a, tc.expect.b,
					op, c.tolerance)
			} else {
				err = c.executeRoundtrip(
					tc.accept.a, tc.accept.b,
					tc.expect.a, tc.expect.b,
					op, c.roundtripDelta, c.roundtripCount)
			}
		} else {
			if c.roundtripCount == 0 {
				_, _, err = c.executeInverseOnce(
					tc.accept.a, tc.accept.b,
					tc.expect.a, tc.expect.b,
					op, c.roundtripDelta)
			} else {
				// roundtrips are always done from the Forward funcs
				panic(9)
			}
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Command) executeForwardOnce(
	in1, in2, out1, out2 float64,
	op core.IConvertLPToXY,
	tolerance float64) (float64, float64, error) {

	input := &core.CoordLP{Lam: support.DDToR(in1), Phi: support.DDToR(in2)}
	output, err := op.Forward(input)
	if err != nil {
		return 0, 0, err
	}

	x, y := output.X, output.Y
	ok1 := check(out1, x, c.tolerance)
	ok2 := check(out2, y, c.tolerance)
	if !ok1 || !ok2 {
		return 0, 0, fmt.Errorf("delta failed")
	}

	return x, y, nil
}

func (c *Command) executeInverseOnce(
	in1, in2, out1, out2 float64,
	op core.IConvertLPToXY,
	tolerance float64) (float64, float64, error) {

	input := &core.CoordXY{X: in1, Y: in2}
	output, err := op.Inverse(input)
	if err != nil {
		return 0, 0, err
	}

	lam, phi := support.RToDD(output.Lam), support.RToDD(output.Phi)
	ok1 := check(out1, lam, c.tolerance)
	ok2 := check(out2, phi, c.tolerance)
	if !ok1 || !ok2 {
		return 0, 0, fmt.Errorf("delta failed")
	}

	return lam, phi, nil
}

func (c *Command) executeRoundtrip(
	in1, in2, out1, out2 float64,
	op core.IConvertLPToXY,
	tolerance float64,
	count int) error {

	for i := 0; i < count; i++ {

		x, y, err := c.executeForwardOnce(in1, in2, out1, out2, op, tolerance)
		if err != nil {
			return err
		}
		lam, phi, err := c.executeInverseOnce(x, y, in1, in1, op, tolerance)
		if err != nil {
			return err
		}

		in1, in2 = lam, phi
	}

	return nil
}

func check(expect, actual, tolerance float64) bool {

	diff := math.Abs(expect - actual)

	if diff > tolerance {
		mlog.Printf("TEST FAILED")
		mlog.Printf("expected:  %f", expect)
		mlog.Printf("actual:    %f", actual)
		mlog.Printf("tolerance: %f", tolerance)
		mlog.Printf("diff:      %f", diff)
		return false
	}
	return true
}
