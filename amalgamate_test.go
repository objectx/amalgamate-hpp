package main

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAmalgamizer_Apply(t *testing.T) {
	Convey("Test Apply", t, func() {
		var outbuf bytes.Buffer
		a, err := NewAmalgamizer(&outbuf)
		So(err, ShouldBeNil)
		Convey("GIVEN: A Amalgamizer", func() {
			err = a.Apply("testdata/test.hpp")
			So(err, ShouldBeNil)
			Convey("WHEN: Reading a test data", func() {
				Convey("THEN: Should match", func() {
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
					So(outbuf.String(), ShouldEqual, expected)
				})
			})
		})
	})
}
