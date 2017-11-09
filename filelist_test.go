package main

import (
	"sort"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/convey"
	"github.com/leanovate/gopter/gen"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFileList_Register(t *testing.T) {
	Convey("Test Register", t, func() {
		Convey("GIVEN: A FileList", func() {
			var fl FileList
			Convey("WHEN: Add an item", func() {
				idx := fl.Register("a")
				Convey("THEN: Index should be 0", func() {
					So(idx, ShouldEqual, 0)
				})
				Convey("AND WHEN: Add another item", func() {
					idx2 := fl.Register("b")
					Convey("THEN: Index should be 1", func() {
						So(idx2, ShouldEqual, 1)
						Convey("AND THEN: Contents should match", func() {
							So(fl.Items(), ShouldResemble, []string{"a", "b"})
						})
					})
					Convey("AND WHEN: Add same item", func() {
						idx3 := fl.Register("a")
						Convey("THEN: Index should be 0", func() {
							So(idx3, ShouldEqual, 0)
						})
					})
				})
			})
		})
	})
}

func TestFileList_Register2(t *testing.T) {
	Convey("Test FileList properties", t, func() {
		prop := func(paths []string) bool {
			//t.Log("paths ", paths)
			var fl1 FileList
			for _, p := range paths {
				fl1.Register(p)
			}
			tmp := append([]string{}, fl1.Items()...)
			sort.Reverse(sort.StringSlice(paths))
			for _, p := range paths {
				fl1.Register(p)
			}
			if len(tmp) != len(fl1.Items()) {
				return false
			}
			for i, v := range fl1.Items() {
				if tmp[i] != v {
					return false
				}
			}
			return true
		}
		So(prop, convey.ShouldSucceedForAll, genItems())
	})
}

func TestFileList_FindIndex(t *testing.T) {
	prop := func(items []string) bool {
		var fl FileList
		// t.Log("items =", items)
		registerItems(&fl, items)

		for i, v := range items {
			idx := fl.FindIndex(v)
			if idx < 0 {
				return false
			}
			//t.Log("i =", i, " v = ", v, " idx =", idx)
			if i < idx {
				t.Log("i =", i, " v =", v, " idx =", idx)
				return false
			}
		}
		return true
	}
	Convey("Testing FindIndex", t, func() {
		So(prop, convey.ShouldSucceedForAll, genItems())
	})
}

func registerItems(fl *FileList, items []string) {
	for _, v := range items {
		fl.Register(v)
	}
}

func genItems() gopter.Gen {
	return gen.SliceOf(
		gen.Weighted([]gen.WeightedGen{
			{1, gen.AlphaString()},
			{3, gen.OneConstOf("YES", "NO", "TRUE", "FALSE")},
		}))
}
