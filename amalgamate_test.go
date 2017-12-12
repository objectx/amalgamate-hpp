package main

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/ToQoz/gopwt"
	"github.com/ToQoz/gopwt/assert"
)

func TestAmalgamizer_Apply(t *testing.T) {
	var outbuf bytes.Buffer
	a, err := NewAmalgamizer(&outbuf)
	assert.OK(t, err == nil, "Should success")
	err = a.Apply("testdata/test.hpp")
	assert.OK(t, err == nil, "Should success")
	expected := `/*
 * Preamble test
 */
#pragma once
#ifndef test_hpp__c24d74e8b3ea1d81019cf4a6cd7c200bac159d06697a1a7de46b535042e78954
#define test_hpp__c24d74e8b3ea1d81019cf4a6cd7c200bac159d06697a1a7de46b535042e78954	1
#include <sys/types.h>

/*
 * Preamble child
 */

/*
 * Preamble child2
 */

/* Body of the child2 */


/* Body of the child */


#ifndef TEST
mokeke
#else
gugugu
#endif

#endif	/* test_hpp__c24d74e8b3ea1d81019cf4a6cd7c200bac159d06697a1a7de46b535042e78954 */
/*
 * Postamble test
 */
`
	assert.OK(t, outbuf.String() == expected)
}

func TestMain(m *testing.M) {
	flag.Parse()
	gopwt.Empower()
	os.Exit(m.Run())
}
