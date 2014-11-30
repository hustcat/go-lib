// Copyright 2014 bitmap authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Author：Ye Yin<hustcat@gmail.com>

package bitmap

import (
	"fmt"
	"testing"
)

func TestBitmapSet(t *testing.T) {
	bitmap := NewNumaBitmap()
	bitmap.SetBit(1, 1)
	bitmap.SetBit(10, 1)
	actual := bitmap.String()
	expected := "[1 10]"
	if expected != actual {
		t.Errorf("expected:%s, actual:%s", expected, actual)
	}
}

func TestGet1BitOffs(t *testing.T) {
	bitmap := NewNumaBitmap()
	bitmap.SetBit(1, 1)
	bitmap.SetBit(10, 1)
	bitmap.SetBit(15, 1)
	bitmap.SetBit(10, 0)
	bitmap.SetBit(10, 1)
	actual := bitmap.Get1BitOffs()
	expected := []uint{1, 10, 15}
	a := fmt.Sprintf("%v", expected)
	b := fmt.Sprintf("%v", actual)
	if a != b {
		t.Errorf("expected:%v, actual:%v", a, b)
	}
}

func TestGet1BitOffsNuma(t *testing.T) {
	bitmap := NewNumaBitmap()

	//node 0
	bitmap.SetBit(0, 1)
	bitmap.SetBit(5, 1)

	//node 1
	bitmap.SetBit(6, 1)
	bitmap.SetBit(11, 1)

	//node 0
	bitmap.SetBit(12, 1)
	bitmap.SetBit(17, 1)

	//node 1
	bitmap.SetBit(18, 1)
	bitmap.SetBit(23, 1)

	actual := bitmap.Get1BitOffsNuma(2)
	expected := [][]uint{
		[]uint{0, 5, 12, 17},
		[]uint{6, 11, 18, 23},
	}
	a := fmt.Sprintf("%v", expected)
	b := fmt.Sprintf("%v", actual)
	if a != b {
		t.Errorf("expected:%v, actual:%v", a, b)
	}
}
