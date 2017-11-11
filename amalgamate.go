// Copyright (c) 2017 Masashi Fujita

package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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
	rxPragmaOnce    = regexp.MustCompile(`^\s*#\s*pragma\s+once\b`)
)

// Amalgamizer amalgamates supplied *.hpp into single *.hpp file.
// Assumes target contains following structure:
// /* preamble */
// #ifndef GUARD
// #define GUARD 1
// /*
//  * Body
//  */
// #endif /* GUARD */
// /* postamble */
//
type Amalgamizer struct {
	output         io.Writer
	systemIncludes FileList
	localIncludes  FileList
	sourceRoot     string
	contexts       []*readContext
}

type readContext struct {
	input io.ReadCloser
	state int
}

func newReadContext(filePath string) (*readContext, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open \"%s\" for reading", filePath)
	}
	return &readContext{
		input: f,
		state: statePreamble,
	}, nil
}

func (ctx *readContext) Close() {
	ctx.input.Close()
}

func NewAmalgamizer(out io.Writer) (*Amalgamizer, error) {
	return &Amalgamizer{
		output: out,
	}, nil
}

func (a *Amalgamizer) Clear() {
	a.sourceRoot = ""
	a.systemIncludes.Clear()
	a.localIncludes.Clear()
	a.contexts = nil
}

func (a *Amalgamizer) Apply(inputPath string) error {
	a.Clear()
	a.sourceRoot = filepath.Dir(inputPath)
	r, err := a.applyInternal(inputPath)
	if err != nil {
		return err
	}
	_, err = io.Copy(a.output, &r.preamble)
	if err != nil {
		return err
	}
	fmt.Fprintln(a.output, "#pragma once")
	b := r.body.Bytes()
	guard := computeGuard(b, inputPath)
	fmt.Fprintf(a.output, "#ifndef %s\n", guard)
	fmt.Fprintf(a.output, "#define %s\t1\n", guard)
	for _, h := range a.systemIncludes.Items() {
		fmt.Fprintf(a.output, "#include <%s>\n", h)
	}
	_, err = a.output.Write(b)
	if err != nil {
		return err
	}
	fmt.Fprintf(a.output, "#endif\t/* %s */\n", guard)
	_, err = io.Copy(a.output, &r.postamble)
	if err != nil {
		return err
	}
	return nil
}

func (a *Amalgamizer) applyInternal(inputPath string) (*includeResult, error) {
	ctx, err := newReadContext(inputPath)
	if err != nil {
		return nil, err
	}
	defer ctx.Close()
	a.contexts = append(a.contexts, ctx)
	defer (func() {
		a.contexts = a.contexts[:len(a.contexts)-1]
	})()
	scanner := bufio.NewScanner(ctx.input)
	var guard string
	var result includeResult
	for scanner.Scan() {
		txt := scanner.Text()
		switch ctx.state {
		case statePreamble:
			if g := findGuardIf(txt); 0 < len(g) {
				ctx.state = stateGuardOpen
				guard = g
				break
			}
			if findPragmaOnce(txt) {
				break
			}
			fmt.Fprintln(&result.preamble, txt)
		case stateGuardOpen:
			if g := findGuardDefine(txt); 0 < len(g) {
				if g == guard {
					ctx.state = stateBody
					break
				}
			}
			return nil, errors.Errorf("missing #define %s just after the #if", guard)
		case stateBody:
			if g := findGuardEndif(txt); 0 < len(g) {
				if g == guard {
					ctx.state = statePostamble
					break
				}
			}
			if inc := findLocalInclude(txt); 0 < len(inc) {
				logger.Debugf("Local include: %s", inc)
				if a.localIncludes.FindIndex(inc) < 0 {
					// Newly found local include file.
					// Expand to here.
					a.localIncludes.Register(inc)
					r, err := a.applyInternal(filepath.Join(a.sourceRoot, inc))
					if err != nil {
						return nil, err
					}
					err = r.WriteTo(&result.body)
					if err != nil {
						return nil, err
					}
				}
				break
			}
			if inc := findSystemInclude(txt); 0 < len(inc) {
				logger.Debugf("System include: %s", inc)
				a.systemIncludes.Register(inc)
				break
			}
			fmt.Fprintln(&result.body, txt)
		case statePostamble:
			fmt.Fprintln(&result.postamble, txt)
		default:
			panic(fmt.Sprintf("unexpected state %v", ctx.state))
		}
	}
	return &result, nil
}

func findPragmaOnce(line string) bool {
	return rxPragmaOnce.FindString(line) != ""
}

func findGuardIf(line string) string {
	if m := rxGuardIf.FindStringSubmatch(line); m != nil {
		return m[1]
	}
	return ""
}

func findGuardDefine(line string) string {
	if m := rxGuardDefine.FindStringSubmatch(line); m != nil {
		return m[1]
	}
	return ""
}

func findGuardEndif(line string) string {
	if m := rxGuardEndif.FindStringSubmatch(line); m != nil {
		return m[1]
	}
	return ""
}

func findLocalInclude(line string) string {
	if m := rxLocalInclude.FindStringSubmatch(line); m != nil {
		return m[1]
	}
	return ""
}

func findSystemInclude(line string) string {
	if m := rxSystemInclude.FindStringSubmatch(line); m != nil {
		return m[1]
	}
	return ""
}

func computeGuard(b []byte, filePath string) string {
	mapper := func(r rune) rune {
		switch {
		case 'A' <= r && r <= 'Z', 'a' <= r && r <= 'z', '0' <= r && r <= '9':
			break
		default:
			r = '_'
		}
		return r
	}
	g := strings.Map(mapper, filepath.Base(filePath))
	return fmt.Sprintf("%s__%x", g, sha256.Sum256(b))
}

type includeResult struct {
	preamble  bytes.Buffer
	body      bytes.Buffer
	postamble bytes.Buffer
}

func (r *includeResult) WriteTo(w io.Writer) error {
	var err error
	_, err = io.Copy(w, &r.preamble)
	_, err = io.Copy(w, &r.body)
	_, err = io.Copy(w, &r.postamble)
	return err
}
