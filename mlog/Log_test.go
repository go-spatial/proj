// Copyright (C) 2018, Michael P. Gerlek (Flaxen Consulting)
//
// Portions of this code were derived from the PROJ.4 software
// In keeping with the terms of the PROJ.4 project, this software
// is provided under the MIT-style license in `LICENSE.md` and may
// additionally be subject to the copyrights of the PROJ.4 authors.

package mlog_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"
	"testing"

	"github.com/go-spatial/proj/mlog"
	"github.com/stretchr/testify/assert"
)

func redirectStderr(f *os.File) int {
	savedFd, err := syscall.Dup(int(os.Stderr.Fd()))
	if err != nil {
		panic(err)
	}
	err = syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		panic(err)
	}
	return savedFd
}

func unredirectStderr(savedFd int) {
	err := syscall.Dup2(savedFd, int(os.Stderr.Fd()))
	if err != nil {
		log.Fatalf("Failed to redirect stderr to file: %v", err)
	}
	syscall.Close(savedFd)
}

func TestLogger(t *testing.T) {
	assert := assert.New(t)

	tmpfile, err := ioutil.TempFile("", "mlog-test")
	if err != nil {
		panic(err)
	}
	defer func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}()

	savedFd := redirectStderr(tmpfile)

	// save log state
	oldDebug := mlog.DebugEnabled
	oldInfo := mlog.InfoEnabled
	oldError := mlog.ErrorEnabled

	func() {
		// the defer is inside a lambda so that it always gets
		// called immediately after the log stmts run, even if
		// they crash -- otherwise, we'd have lost our stderr!
		defer unredirectStderr(savedFd)

		mlog.DebugEnabled = true
		mlog.InfoEnabled = true
		mlog.ErrorEnabled = true

		mlog.Debugf("debug %d", 1)
		mlog.Printf("print %s", "2")
		e := fmt.Errorf("E")
		mlog.Error(e)
		x := "yow"
		mlog.Printv(x)

		mlog.DebugEnabled = false
		mlog.InfoEnabled = false
		mlog.ErrorEnabled = false

		mlog.Debugf("nope")
		mlog.Printf("nope")
		mlog.Error(e)
	}()

	// restore log state
	mlog.DebugEnabled = oldDebug
	mlog.InfoEnabled = oldInfo
	mlog.ErrorEnabled = oldError

	tmpfile.Seek(0, io.SeekStart)
	buf := make([]byte, 1024)
	n, err := tmpfile.Read(buf)
	assert.NoError(err)
	buf = buf[0:n]

	ex := []string{
		"[DEBUG] Log_test.go:66: debug 1",
		"[LOG] Log_test.go:67: print 2",
		"[ERROR] E",
		"[LOG] Log_test.go:71: \"yow\"",
	}
	expected := strings.Join(ex, "\n") + "\n"
	assert.Equal(expected, string(buf))
}
