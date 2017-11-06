/*
 * Copyright (c) 2017. Masashi Fujita
 */
package main

import (
	"flag"
	"fmt"
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
	if flag.NArg() < 2 {
		flag.Usage()
	}
	err := Amalgamate(*outputPath, flag.Arg(0))
	if err != nil {
		logger.Errorf("amalgamation failed (%v)", err)
		os.Exit(1)
	}
	os.Exit(0)
}

// Amalgamete performs *.hpp amalgamation starting from inputPath.
func Amalgamate(outputPath string, inputPath string) error {
	return nil
}
