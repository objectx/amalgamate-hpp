// Copyright (c) 2017 Masashi Fujita

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
)

const (
	// State for reading *.hpp
	statePreamble = iota
	stateGuardOpen
	stateBody
	statePostamble
)

var (
	rxGuardIf       = regexp.MustCompile(`^\s*#\s*ifndef\s+(\S+)`)
	rxGuardDefine   = regexp.MustCompile(`^\s*#\s*define\s+(\S+)`)
	rxGuardEndif    = regexp.MustCompile(`^\s*#\s*endif\s+/\*\s+(\S+)\s+\*/`)
	rxLocalInclude  = regexp.MustCompile(`^\s*#\s*include\s+"([^"]+)"`)
	rxSystemInclude = regexp.MustCompile(`^\s*#\s*include\s+<([^>]+)>`)
)

// Amalgamizer amalgamates supplied *.hpp into single *.hpp file.
type Amalgamizer struct {
	output         io.Writer
	systemIncludes FileList
	localIncludes  FileList
	sourceRoot     string
	scanner        *bufio.Scanner
	state          int
}

func NewAmalgamizer(out io.Writer) (*Amalgamizer, error) {
	return &Amalgamizer{
		output: out,
		state:  statePreamble,
	}, nil
}

func (a *Amalgamizer) Clear() {
	a.sourceRoot = ""
	a.systemIncludes.Clear()
	a.localIncludes.Clear()
	a.scanner = nil
	a.state = statePreamble
}

func (a *Amalgamizer) Apply(inputPath string) error {
	a.Clear()
	a.sourceRoot = filepath.Dir(inputPath)
	f, err := os.Open(inputPath)
	if err != nil {
		return errors.Wrapf(err, "failed to open %s as input", inputPath)
	}
	defer f.Close()
	a.scanner = bufio.NewScanner(f)
	var guard string
	for a.scanner.Scan() {
		txt := a.scanner.Text()
		switch a.state {
		case statePreamble:
			m := rxGuardIf.FindStringSubmatch(txt)
			if m != nil {
				a.state = stateGuardOpen
				guard = m[1]
			} else {
				fmt.Fprintln(a.output, txt)
			}
		case stateGuardOpen:
			m := rxGuardDefine.FindStringSubmatch(txt)
			if m != nil {
				if m[1] == guard {
					a.state = stateBody
					break
				}
			}
			logger.Debugf("m = %v", m)
			return errors.Errorf("missing #define %s just after the #if", guard)
		case stateBody:
			m := rxGuardEndif.FindStringSubmatch(txt)
			if m != nil {
				if m[1] == guard {
					a.state = statePreamble
					break
				}
			}
			m = rxLocalInclude.FindStringSubmatch(txt)
			if m != nil {
				f := m[1]
				logger.Debugf("Local include: %s", f)
				if a.localIncludes.FindIndex(f) < 0 {
					// Newly found local include file.
					// Expand to here.
					a.localIncludes.Register(f)
				}
				break
			}
			m = rxSystemInclude.FindStringSubmatch(txt)
			if m != nil {
				f := m[1]
				logger.Debugf("System include: %s", f)
				a.systemIncludes.Register(f)
				break
			}
			fmt.Fprintln(a.output, txt)
		case statePostamble:
			fmt.Fprintln(a.output, txt)
		default:
			panic(fmt.Sprintf("unexpected state %v", a.state))
		}
	}
	return nil
}
