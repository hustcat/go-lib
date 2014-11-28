// Copyright 2014 bitmap authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Author：Ye Yin<hustcat@gmail.com>

// NUMA CPU bitmap, used to affinity progress to some CPU
package bitmap

import "fmt"

const defaultSize = 24
const maxSize = 1024
const defaultNodeNum = 2

type NumaBitmap struct {
	bits []byte

	size uint

	userSize uint
	//NUMA node num
	nodeNum int
}

func NewNumaBitmap() *NumaBitmap {
	return NewNumaBitmapSize(defaultSize, defaultNodeNum)
}

func NewNumaBitmapSize(size uint, nodeNum int) *NumaBitmap {
	if size == 0 || size > maxSize {
		size = defaultSize
	}
	userSize := size
	size = align(size, 8)
	return &NumaBitmap{bits: make([]byte, size>>3), size: size, userSize: userSize, nodeNum: nodeNum}
}

func (b *NumaBitmap) SetBit(offset uint, value uint) bool {
	index, pos := offset/8, offset%8

	if b.size < offset {
		return false
	}

	if value == 0 {
		b.bits[index] &^= 0x01 << pos
	} else {
		b.bits[index] |= 0x01 << pos
	}

	return true
}

func (b *NumaBitmap) GetBit(offset uint) byte {
	index, pos := offset/8, offset%8

	if b.size < offset {
		return 0
	}

	return (b.bits[index] >> pos) & 0x01
}

// Get the offset of bits equal 1 all
func (b *NumaBitmap) Get1BitOffs() []uint {
	var (
		offset uint
		offs   []uint
	)

	maxNo := b.userSize

	offset = 0
	for index, line := range b.bits {
		for pos := 0; pos < 8; pos++ {

			offset = uint(index*8 + pos)
			if offset >= maxNo {
				goto OUT
			}

			if (line>>uint(pos))&0x01 != 0 {
				offs = append(offs, offset)
			}
		}
	}
OUT:
	return offs
}

// Get the offsets of bits equal 1 per Node
func (b *NumaBitmap) Get1BitOffsNuma(nodeNum uint) [][]uint {
	var (
		tmp     uint
		offset  uint
		curNode uint
	)

	maxNo := b.userSize
	//only surpport hyperthread CPU
	step := maxNo / (nodeNum * 2)

	//cpu cores, don't include hyperthread
	cpu := maxNo / 2

	curNode = 0

	offs := make([][]uint, nodeNum)

	offset = 0
	for index, line := range b.bits {
		for pos := 0; pos < 8; pos++ {
			offset = uint(index*8 + pos)
			if offset >= maxNo {
				goto OUT
			}

			//exlude hyperthread
			if offset > cpu {
				tmp = offset - cpu
			}

			curNode = tmp / step
			if (line>>uint(pos))&0x01 != 0 {
				offs[curNode] = append(offs[curNode], offset)
			}
		}
	}
OUT:
	return offs
}

/*
//Set num bit to 1 on Node(node), return bits index
func (b *NumaBitmap) Set1NumaBitNum(node uint, num uint) []uint {
	var (
		tmp     uint
		offs    []uint
		offset  uint
		curNode uint
	)

	maxNo := b.userSize
	//only surpport hyperthread CPU
	step := maxNo / (node * 2)

	//cpu cores, don't include hyperthread
	cpu := maxNo / 2

	curNode = 0

	offset = 0
	for index, line := range b.bits {
		for pos := 0; pos < 8; pos++ {
			offset = uint(index*8 + pos)
			if offset >= maxNo {
				break
			}

			//exlude hyperthread
			if offset > cpu {
				tmp = offset - cpu
			}

			curNode = tmp / step

			if curNode == node && ((line>>uint(pos))&0x01) == 0 {
				b.SetBit(offset, 1)
				offs = append(offs, offset)
			}
		}
	}
	return offs
}*/

func (b *NumaBitmap) String() string {

	offs := make([]uint, 0, b.userSize)

	var offset uint
	for offset = 0; offset < b.userSize; offset++ {
		if b.GetBit(offset) == 1 {
			offs = append(offs, offset)
		}
	}

	return fmt.Sprintf("%v", offs)
}

func align(n, align uint) uint {
	return (n + align - 1) & (^(align - 1))
}
