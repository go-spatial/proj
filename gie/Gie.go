// Copyright (C) 2018, Michael P. Gerlek (Flaxen Consulting)
//
// Portions of this code were derived from the PROJ.4 software
// In keeping with the terms of the PROJ.4 project, this software
// is provided under the MIT-style license in `LICENSE.md` and may
// additionally be subject to the copyrights of the PROJ.4 authors.

package gie

import (
	"fmt"
	"io/ioutil"
	"strings"

	// need to pull in the operations table entries
	_ "github.com/go-spatial/proj/operations"
)

// These are the projections we know about. If the projection string has a
// "proj=" key whose valued is not in this list, the Gie will not try to
// execute the Command.
var supportedProjections = []string{
	"etmerc", "utm",
	"aea", "leac",
	"merc",
	"airy",
	"august",
}

// If the proj string has one of these keys, we won't execute the Command.
var unsupportedKeys = []string{
	"axis",
	"geoidgrids",
	"to_meter",
}

// If the Command is from this file and line, we won't execute the
// Command -- this acts as a way to shut off tests we don't like.
var skippedTests = []string{
	"ellipsoid.gie:64",
}

// Gie is the top-level object for the Gie test runner
//
// Gie manages reading and parsing the .gie files and then
// executing the commands and their testcases.
type Gie struct {
	dir      string
	files    []string
	Commands []*Command
}

// NewGie returns a new Gie object
func NewGie(dir string) (*Gie, error) {

	g := &Gie{
		dir:      dir,
		files:    []string{},
		Commands: []*Command{},
	}

	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, info := range infos {
		file := info.Name()
		if strings.HasSuffix(file, ".gie") {
			g.files = append(g.files, dir+"/"+file)
		}
	}
	return g, nil
}

// Parse reads the .gie files and creates the commands
func (g *Gie) Parse() error {
	for _, file := range g.files {
		p, err := NewParser(file)
		if err != nil {
			return err
		}
		g.Commands = append(g.Commands, p.Commands...)
	}

	return nil
}

// IsSupported returns true iff we support the projection and they aren't asking
// for proj strings things we don't do yet
func (g *Gie) IsSupported(cmd *Command) bool {

	if !g.isSupportedProjection(cmd) {
		return false
	}

	if g.hasUnsupportedKey(cmd) {
		return false
	}

	if g.isSkippedTest(cmd) {
		return false
	}

	return true
}

func (g *Gie) isSkippedTest(cmd *Command) bool {

	s := fmt.Sprintf("%s:%d", cmd.File, cmd.Line)
	for _, skippy := range skippedTests {
		if strings.HasSuffix(s, skippy) {
			return true
		}
	}
	return false
}

func (g *Gie) hasUnsupportedKey(cmd *Command) bool {

	for _, badkey := range unsupportedKeys {
		if strings.Contains(cmd.ProjString, badkey) {
			return true
		}
	}
	return false
}

func (g *Gie) isSupportedProjection(cmd *Command) bool {

	proj := cmd.ProjectionName()

	for _, sp := range supportedProjections {
		if proj == sp {
			return true
		}
	}

	return false
}
