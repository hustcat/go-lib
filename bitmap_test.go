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
