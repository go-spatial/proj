// Copyright (C) 2018, Michael P. Gerlek (Flaxen Consulting)
//
// Portions of this code were derived from the PROJ.4 software
// In keeping with the terms of the PROJ.4 project, this software
// is provided under the MIT-style license in `LICENSE.md` and may
// additionally be subject to the copyrights of the PROJ.4 authors.

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-spatial/proj"

	"github.com/go-spatial/proj/core"
	"github.com/go-spatial/proj/merror"
	"github.com/go-spatial/proj/mlog"
	"github.com/go-spatial/proj/support"
)

func main() {
	err := Main(os.Stdin, os.Stdout, os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
}

// Main is just a callable version of main(), for testing purposes
func Main(inS io.Reader, outS io.Writer, args []string) error {

	merror.ShowSource = false
	mlog.DebugEnabled = false
	mlog.InfoEnabled = false
	mlog.ErrorEnabled = false

	cli := flag.NewFlagSet(args[0], flag.ContinueOnError)
	cli.SetOutput(outS)

	verbose := cli.Bool("verbose", false, "enable logging")
	inverse := cli.Bool("inverse", false, "run the inverse transform")
	epsgDest := cli.Int("epsg", 0, "perform conversion from 4326 to given destination system")

	err := cli.Parse(args[1:])
	if err != nil {
		return err
	}
	projString := strings.Join(cli.Args(), " ")

	if *verbose {
		mlog.Printf("verbose: %t", *verbose)
		mlog.Printf("inverse: %t", *inverse)
		if *epsgDest == 0 {
			mlog.Printf("epsg: (not specified)")
		} else {
			mlog.Printf("epsg: %d", epsgDest)
		}
		if projString == "" {
			mlog.Printf("proj: (not specified)")
		} else {
			mlog.Printf("proj: %s", projString)
		}

		merror.ShowSource = true
		mlog.DebugEnabled = true
		mlog.InfoEnabled = true
		mlog.ErrorEnabled = true
	}

	if *epsgDest != 0 {
		if *inverse {
			return fmt.Errorf("-inverse not allowed with -epsg")
		}
		if projString != "" {
			return fmt.Errorf("projection string not allowed with -epsg")
		}
		input := make([]float64, 2)

		f := func(a, b float64) (float64, float64, error) {
			input[0] = a
			input[1] = b
			output, err := proj.Convert(proj.EPSGCode(*epsgDest), input)
			if err != nil {
				return 0.0, 0.0, err
			}
			return output[0], output[1], nil
		}

		return repl(inS, outS, f)
	}

	ps, err := support.NewProjString(projString)
	if err != nil {
		return err
	}

	_, opx, err := core.NewSystem(ps)
	if err != nil {
		return err
	}

	op := opx.(core.IConvertLPToXY)

	if !*inverse {

		f := func(a, b float64) (float64, float64, error) {
			input := &core.CoordLP{Lam: support.DDToR(a), Phi: support.DDToR(b)}
			output, err := op.Forward(input)
			if err != nil {
				return 0.0, 0.0, err
			}
			return output.X, output.Y, nil
		}
		return repl(inS, outS, f)
	}

	f := func(a, b float64) (float64, float64, error) {
		input := &core.CoordXY{X: a, Y: b}
		output, err := op.Inverse(input)
		if err != nil {
			return 0.0, 0.0, err
		}
		return support.RToDD(output.Lam), support.RToDD(output.Phi), nil
	}
	return repl(inS, outS, f)
}

type converter func(a, b float64) (float64, float64, error)

func repl(inS io.Reader, outS io.Writer, f converter) error {

	var a, b float64

	for {
		n, err := fmt.Fscanf(inS, "%f %f\n", &a, &b)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if n != 2 {
			return fmt.Errorf("error reading input")
		}

		c, d, err := f(a, b)
		if err != nil {
			return err
		}

		fmt.Fprintf(outS, "%f %f\n", c, d)
	}
}
