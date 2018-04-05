package support

import (
	"regexp"
	"strconv"

	"github.com/go-spatial/proj4go/merror"

	"github.com/go-spatial/proj4go/mlog"
)

// DMSToDD converts a degrees-minutes-seconds string to decimal-degrees
//
// Using an 8-part regexp, we support this format:
//    [+-] nnn [°Dd] nnn ['Mm] nnn.nnn ["Ss] [NnEeWwSs]
//
// TODO: the original dmstor() may support more, but the parsing code
// is messy and we don't have any testcases at this time.
func DMSToDD(input string) (float64, error) {

	mlog.Debugf("%s", input)

	deg := `\s*(-|\+)?\s*(\d+)\s*([°Dd])` // t1, t2, t3
	min := `\s*(\d+)\s*(['Mm])`           // t4, t5
	sec := `\s*(\d+\.?\d*)\s*(["Ss])`     // t6, t7
	news := `\s*([NnEeWwSs]?)`            // t8
	expr := "^" + deg + min + sec + news + "$"
	r := regexp.MustCompile(expr)

	tokens := r.FindStringSubmatch(input)
	if tokens == nil {
		return 0.0, merror.New(ErrInvalidArg)
	}

	sign := tokens[1]
	d := tokens[2]
	m := tokens[4]
	s := tokens[6]
	dir := tokens[8]

	df, err := strconv.ParseFloat(d, 64)
	if err != nil {
		return 0.0, err
	}
	mf, err := strconv.ParseFloat(m, 64)
	if err != nil {
		return 0.0, err
	}
	sf, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0, err
	}

	dd := df + mf/60.0 + sf/3600.0
	if sign == "-" {
		dd = -dd
	}

	if dir == "S" || dir == "s" || dir == "W" || dir == "w" {
		dd = -dd
	}

	return dd, nil
}

// DMSToR converts a DMS string to radians
func DMSToR(input string) (float64, error) {

	dd, err := DMSToDD(input)
	if err != nil {
		return 0.0, err
	}

	r := DDToR(dd)

	return r, nil
}

// DDToR converts decimal degrees to radians
func DDToR(deg float64) float64 {
	const degToRad = 0.017453292519943296
	r := deg * degToRad
	return r
}