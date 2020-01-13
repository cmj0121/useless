/* Copyright (C) 2019-2020 cmj. All right reserved. */
package tree

import (
	"testing"
)

func TestSimpleType(t *testing.T) {
	cases := map[string]interface{}{
		"bool":       true,
		"int":        123,
		"uint":       uint(123),
		"int8":       int8(1),
		"uint8":      uint8(1),
		"int16":      int16(1),
		"uint16":     uint16(1),
		"int32":      int32(1),
		"uint32":     uint32(1),
		"int64":      int64(1),
		"uint64":     uint64(1),
		"uintptr":    uintptr(1),
		"float32":    float32(1.23),
		"float64":    float64(1.23),
		"complex64":  complex64(123),
		"complex128": complex128(123),
		"string":     "abc",
		"struct":     struct{}{},
	}

	for key, value := range cases {
		tree := New(value)

		if key != tree.String() {
			/* The format is not match */
			t.Errorf("Case %v not match - %s <> %s", value, tree.String(), key)
		}
	}
}

func TestSliceMapChan(t *testing.T) {
	x, y, z := 1, 2, 3
	var c1 chan int
	var c2 <-chan int
	var c3 chan<- int

	cases := map[string]interface{}{
		"slice int":  []int{1, 2, 3},
		"slice *int": []*int{&x, &y, &z},
		"map int:string": map[int]string{
			1: "x",
			2: "z",
		},
		"map string:*int": map[string]*int{
			"x": &x,
			"y": &y,
		},
		"chan int":   c1,
		"<-chan int": c2,
		"chan<- int": c3,
	}

	for key, value := range cases {
		if tree := New(value); tree.String() != key {
			t.Errorf("Case %v not match - %s <> %s", value, tree.String(), key)
		}
	}
}
