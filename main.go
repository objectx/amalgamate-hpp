/*
 * Copyright (c) 2017. Masashi Fujita
 */
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

var (
	beVerbose = false
	progPath  string
	progName  string
	logger    *zap.SugaredLogger
)

func init() {
	var err error
	progPath, err = os.Executable()
	if err != nil {
		progPath = "amalgamate-hpp"
	}
	progName = filepath.Base(progPath)
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger = l.Sugar()
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <root header>\n", progName)
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.BoolVar(&beVerbose, "v", false, "Be verbose.")
	outputPath := flag.String("o", "-", "Output to this.")
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
	}
	err := AmalgamateToFile(*outputPath, flag.Arg(0))
	if err != nil {
		logger.Errorf("amalgamation failed (%v)", err)
		os.Exit(1)
	}
	os.Exit(0)
}

// Amalgamete performs *.hpp amalgamation starting from inputPath.
func AmalgamateToFile(outputPath string, inputPath string) error {
	if len(outputPath) == 0 || outputPath == "-" {
		return Amalgamate(os.Stdout, inputPath)
	}
	outDir := filepath.Dir(outputPath)
	tmpOut, err := ioutil.TempFile(outDir, "tmp-")
	if err != nil {
		return err
	}
	defer (func() {
		tmpOut.Close()
		_ = os.Remove(tmpOut.Name())
	})()
	err = Amalgamate(tmpOut, inputPath)
	if err != nil {
		return err
	}
	tmpOut.Close()
	err = os.Rename(tmpOut.Name(), outputPath)
	if err != nil {
		return err
	}
	return nil
}

func Amalgamate(output io.Writer, inputPath string) error {
	amalgamizer, err := NewAmalgamizer(output)
	err = amalgamizer.Apply(inputPath)
	if err != nil {
		return err
	}
	return nil
}
