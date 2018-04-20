// Copyright (C) 2018, Michael P. Gerlek (Flaxen Consulting)
//
// Portions of this code were derived from the PROJ.4 software
// In keeping with the terms of the PROJ.4 project, this software
// is provided under the MIT-style license in `LICENSE.md` and may
// additionally be subject to the copyrights of the PROJ.4 authors.

package gie

import (
	"bufio"
	"os"
	"strings"
)

// Parser reads the .gie files and returns the list of commands
// it constructed
type Parser struct {
	lines    []string
	Commands []*Command
	fname    string
	lineNum  int
}

// NewParser creates a Parser object and runs the parser
//
// The parser runs one line at time. It's not quite clear
// what the official .gie format is supposed to be, so we
// use a very dumb but effective approach.
func NewParser(fname string) (*Parser, error) {

	lines, err := readLines(fname)
	if err != nil {
		return nil, err
	}

	p := &Parser{
		lines:    lines,
		Commands: []*Command{},
		fname:    fname,
		lineNum:  0,
	}

	for len(p.lines) > 0 {
		//mlog.Printf("%s: %s", fname, p.lines[0])

		if p.removeBlank() {
			continue
		}
		if p.removeGie() {
			continue
		}
		if p.removeTitleBlock() {
			continue
		}
		if p.removeSingleDashed() {
			continue
		}
		if p.removeComment() {
			continue
		}
		if p.doCommand() {
			continue
		}

		// anything left here has to be part of a comment section,
		// so we throw it away
		p.pop()
	}

	return p, nil
}

func (p *Parser) doCommand() bool {
	s := p.lines[0]
	idx := strings.Index(s, "#")
	if idx >= 0 {
		s = s[:idx]
	}

	tokens := strings.Fields(s)

	switch tokens[0] {
	case "operation":
		return p.doCommandOperation(tokens[1:])
	case "accept":
		return p.doCommandAccept(tokens[1:])
	case "expect":
		return p.doCommandExpect(tokens[1:])
	case "roundtrip":
		return p.doCommandRoundtrip(tokens[1:])
	case "banner":
		return p.doCommandBanner(tokens[1:])
	case "verbose":
		return p.doCommandVerbose(tokens[1:])
	case "direction":
		return p.doCommandDirection(tokens[1:])
	case "tolerance":
		return p.doCommandTolerance(tokens[1:])
	case "ignore":
		return p.doCommandIgnore(tokens[1:])
	case "builtins":
		return p.doCommandBuiltins(tokens[1:])
	case "echo":
		return p.doCommandEcho(tokens[1:])
	case "skip":
		return p.doCommandSkip(tokens[1:])
	}

	return false
}

func (p *Parser) doCommandOperation(tokens []string) bool {
	ss := strings.Join(tokens, " ")
	p.pop()
	// naive assumption: if the next line starts with whitespace, it's a
	// continuation of the command
	for len(p.lines) > 0 && (p.lines[0] == "" || p.lines[0][0:1] == " " || p.lines[0][0:1] == "\t") {
		ss += p.lines[0]
		p.pop()
	}
	cmd := NewCommand(p.fname, p.lineNum, ss)
	p.Commands = append(p.Commands, cmd)
	return true
}

func (p *Parser) doCommandTolerance(tokens []string) bool {
	cmd := p.Commands[len(p.Commands)-1]
	numTokens := len(tokens)

	if numTokens == 2 {
		cmd.setTolerance(tokens[0], tokens[1])
		p.pop()
		return true
	} else if numTokens == 1 {
		cmd.setTolerance(tokens[0], "*")
		p.pop()
		return true
	}
	panic(55)
}

func (p *Parser) doCommandAccept(tokens []string) bool {
	cmd := p.Commands[len(p.Commands)-1]
	numTokens := len(tokens)

	switch {
	case numTokens == 2:
		cmd.setAccept(tokens[0], tokens[1], "0.0", "0.0")
		p.pop()
		return true
	case numTokens == 3:
		cmd.setAccept(tokens[0], tokens[1], tokens[2], "0.0")
		p.pop()
		return true
	case numTokens == 4:
		cmd.setAccept(tokens[0], tokens[1], tokens[2], tokens[3])
		p.pop()
		return true
	}
	panic(44)
}

func (p *Parser) doCommandExpect(tokens []string) bool {
	cmd := p.Commands[len(p.Commands)-1]
	numTokens := len(tokens)

	switch {
	case tokens[0] == "failure":
		cmd.setExpectFailure()
		p.pop()
		return true
	case numTokens == 2:
		cmd.setExpect(tokens[0], tokens[1], "0.0", "0.0")
		p.pop()
		return true
	case numTokens == 3:
		cmd.setExpect(tokens[0], tokens[1], tokens[2], "0.0")
		p.pop()
		return true
	case numTokens == 4:
		cmd.setExpect(tokens[0], tokens[1], tokens[2], tokens[3])
		p.pop()
		return true
	}
	panic(33)
}

func (p *Parser) doCommandDirection(tokens []string) bool {
	cmd := p.Commands[len(p.Commands)-1]
	numTokens := len(tokens)

	if numTokens != 1 {
		panic(22)
	}
	cmd.setDirection(tokens[0])
	p.pop()
	return true
}

func (p *Parser) doCommandRoundtrip(tokens []string) bool {
	cmd := p.Commands[len(p.Commands)-1]
	numTokens := len(tokens)

	if numTokens == 1 {
		cmd.setRoundtrip(tokens[0], "1.0", "*")
		p.pop()
		return true
	}
	if numTokens == 3 {
		cmd.setRoundtrip(tokens[0], tokens[1], tokens[2])
		p.pop()
		return true
	}
	panic(77)
}

func (p *Parser) doCommandBanner(tokens []string) bool {
	panic("banner")
}

func (p *Parser) doCommandVerbose(tokens []string) bool {
	panic("verbose")
}

func (p *Parser) doCommandIgnore(tokens []string) bool {
	p.pop()
	return true
}

func (p *Parser) doCommandBuiltins(tokens []string) bool {
	p.pop()
	return true
}

func (p *Parser) doCommandEcho(tokens []string) bool {
	panic("echo")
}

func (p *Parser) doCommandSkip(tokens []string) bool {
	panic("skip")
}

func (p *Parser) removeComment() bool {
	s := strings.TrimSpace(p.lines[0])
	if s[0:1] == "#" {
		p.pop()
		return true
	}
	return false
}

func (p *Parser) removeSingleDashed() bool {
	if p.isDelimiter("-") {
		p.pop()
		return true
	}
	return false
}

func (p *Parser) removeGie() bool {
	s := strings.TrimSpace(p.lines[0])
	if s == "<gie>" || s == "</gie>" {
		p.pop()
		return true
	}
	return false
}

func (p *Parser) removeTitleBlock() bool {
	if !p.isDelimiter("=") {
		return false
	}

	p.pop()
	for !p.isDelimiter("=") {
		p.pop()
	}
	p.pop()
	return true
}

func (p *Parser) removeBlank() bool {
	s := strings.TrimSpace(p.lines[0])
	if len(s) == 0 {
		p.pop()
		return true
	}
	return false
}

func (p *Parser) pop() {
	p.lines = p.lines[1:]
	p.lineNum++
}

func (p *Parser) isDelimiter(c string) bool {
	s := strings.TrimSpace(p.lines[0])
	s = strings.Replace(s, c, "", -1)
	return len(s) == 0
}

func readLines(fname string) ([]string, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
