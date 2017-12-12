package main

import (
	"reflect"
	"sort"
	"testing"

	"github.com/ToQoz/gopwt/assert"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestFileList_Register(t *testing.T) {
	var fl FileList
	idx := fl.Register("a")
	assert.OK(t, idx == 0)
	idx2 := fl.Register("b")
	assert.OK(t, idx2 == 1)
	assert.OK(t, reflect.DeepEqual(fl.Items(), []string{"a", "b"}))
	idx3 := fl.Register("a")
	assert.OK(t, idx3 == 0)
}

func TestFileList_Register2(t *testing.T) {
	properties := gopter.NewProperties(nil)
	{
		cond := func(paths []string) bool {
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
		properties.Property("FileList property", prop.ForAll(cond, genItems()))
	}
	properties.TestingRun(t)
}

func TestFileList_FindIndex(t *testing.T) {
	properties := gopter.NewProperties(nil)
	cond := func(items []string) bool {
		var fl FileList
		// t.Log("items =", items)
		registerItems(&fl, items)

		for i, v := range items {
			idx := fl.FindIndex(v)
			if idx < 0 {
				return false
			}
			// t.Log("i =", i, " v = ", v, " idx =", idx)
			if i < idx {
				t.Log("i =", i, " v =", v, " idx =", idx)
				return false
			}
		}
		return true
	}
	properties.Property("FileIndex property", prop.ForAll(cond, genItems()))
	properties.TestingRun(t)
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
