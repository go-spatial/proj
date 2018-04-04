package support

import (
	"log"
	"strconv"
	"strings"
)

// DMSToR converts a DMS string to radians
//
// We only support this format: n.nD n.n' n.n"
//
// TODO: the original dmstor() may support more, but we don't have any
// testcases
func DMSToR(input string) (float64, error) {

	p := strings.Replace(input, " ", "", -1)

	log.Printf("p=%s", p)

	dIndex := strings.IndexAny(p, "Dd")
	if dIndex == -1 {
		return 0, ErrInvalidDMS
	}
	dString := p[0:dIndex]
	p = p[dIndex+1:]
	log.Printf("dString=%s", dString)
	log.Printf("p=%s", p)

	mIndex := strings.Index(p, `'`)
	if mIndex == -1 {
		return 0, ErrInvalidDMS
	}
	mString := p[0:mIndex]
	p = p[mIndex+1:]
	log.Printf("mString=%s", mString)
	log.Printf("p=%s", p)

	sIndex := strings.Index(p, `"`)
	if sIndex == -1 {
		return 0, ErrInvalidDMS
	}
	sString := p[0:sIndex]
	p = p[sIndex+1:]
	log.Printf("sString=%s", sString)
	log.Printf("p=%s", p)

	sign := 1.0
	if p == "S" || p == "s" || p == "W" || p == "w" {
		sign = -sign
	}

	dValue, err := strconv.ParseFloat(dString, 64)
	if err != nil {
		return 0.0, err
	}
	mValue, err := strconv.ParseFloat(mString, 64)
	if err != nil {
		return 0.0, err
	}
	sValue, err := strconv.ParseFloat(sString, 64)
	if err != nil {
		return 0.0, err
	}

	degToRad := 0.017453292519943296
	r := (sign*dValue + mValue/60.0 + sValue/3600.0) * degToRad
	return r, nil
}
