// Copyright (c) 2017 Masashi Fujita

package main

import (
	"io"
	"path/filepath"
)

type Amalgamizer struct {
	output         io.Writer
	systemIncludes IncludeFileList
	localIncludes  IncludeFileList
	sourceRoot     string
}

func NewAmalgamizer(out io.Writer) (*Amalgamizer, error) {
	return &Amalgamizer{
		output: out,
	}, nil
}

func (a *Amalgamizer) Apply(inputPath string) error {
	a.sourceRoot = filepath.Dir(inputPath)
	return nil
}

type IncludeFileList struct {
	files []string
	dict  map[string]int
}

func (l *IncludeFileList) Items() []string {
	return l.files
}

func (l *IncludeFileList) FindIndex(p string) int {
	if idx, ok := l.dict[p]; ok {
		return idx
	}
	return -1
}

// Register records `filePath` as a include file
func (l *IncludeFileList) Register(filePath string) int {
	idx, ok := l.dict[filePath]
	if ok {
		return idx // Already appeard
	}
	idx = len(l.files)
	l.dict[filePath] = idx
	l.files = append(l.files, filePath)
	return idx
}
