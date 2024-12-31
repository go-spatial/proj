// Copyright (C) 2018, Michael P. Gerlek (Flaxen Consulting)
//
// Portions of this code were derived from the PROJ.4 software
// In keeping with the terms of the PROJ.4 project, this software
// is provided under the MIT-style license in `LICENSE.md` and may
// additionally be subject to the copyrights of the PROJ.4 authors.

package mlog_test

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"testing"

	"github.com/go-spatial/proj/mlog"
)

func TestLogger(t *testing.T) {
	//assert := assert.New(t)

	tmpfile, err := os.CreateTemp("", "mlog-test")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
		return
	}
	defer tmpfile.Close()
	tempFilename := tmpfile.Name()
	t.Cleanup(func() { os.Remove(tempFilename) })

	log := mlog.NewLoggerSingleOutput(tmpfile)
	log.EnableDebug()
	log.EnableError()
	log.EnableInfo()

	// the following is put in an inlined lambda, so that
	// we have a place to put the defer: we need it to always get
	// called immediately after the log stmts run, even if
	// they crash -- otherwise, we'd have lost our stderr!

	log.Debugf("debug %d", 1)
	log.Printf("print %s", "2")
	e := fmt.Errorf("E")
	log.Error(e)
	x := "yow"
	log.Printv(x)

	log.DisableDebug()
	log.DisableError()
	log.DisableInfo()

	log.Debugf("nope")
	log.Printf("nope")
	log.Error(e)

	tmpfile.Seek(0, io.SeekStart)

	expectedLines := []string{
		`\[DEBUG\] Log_test.go:\d+: debug 1`,
		`\[LOG\] Log_test.go:\d+: print 2`,
		`\[ERROR\] E`,
		`\[LOG\] Log_test.go:\d+: "yow"`,
	}

	scanner := bufio.NewScanner(tmpfile)
	count := 0
	for scanner.Scan() {
		scanner.Text()
		if count > len(expectedLines) {
			t.Errorf("Found too many lines, expected %d, got %d", len(expectedLines), count)
			return
		}
		txt := scanner.Text()
		m, err := regexp.MatchString(expectedLines[count], txt)
		if err != nil {
			t.Errorf("error failed to match regexp[%s]: error: %v", expectedLines[count], err)
			return
		}
		if !m {
			t.Errorf("failed to match regexp[%s]: %s", expectedLines[count], txt)
			return
		}
		count++
	}
}
