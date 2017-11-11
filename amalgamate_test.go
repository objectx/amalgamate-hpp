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
