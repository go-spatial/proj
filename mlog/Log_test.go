package mlog_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
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

	ex1 := "[DEBUG] Log_test.go:65: debug 1"
	ex2 := "[LOG] Log_test.go:66: print 2"
	ex3 := "[ERROR] E"
	expected := ex1 + "\n" + ex2 + "\n" + ex3 + "\n"
	assert.Equal(expected, string(buf))
}
