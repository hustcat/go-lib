// Copyright 2014 bitmap authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Author：Ye Yin<hustcat@gmail.com>

// NUMA CPU bitmap, used to affinity progress to some CPU.
// CPU must be hypethreaded, and CPU number look like as follows:
// [node0, node1, ... , node0, node1, ...]
// For example:
// node0: [0,1,2,3,4,5,12,13,14,15,16,17]
// node1: [6,7,8,9,10,11,18,19,20,21,22,23]

package bitmap

import "fmt"

//24 cores
const defaultSize = 24
const maxSize = 1024

//two nodes
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

func (b *NumaBitmap) SetBit(offset uint, value uint) error {
	index, pos := offset/8, offset%8

	if b.userSize <= offset {
		return fmt.Errorf("offset: %d is out of range %d", offset, b.userSize)
	}

	if value == 0 {
		b.bits[index] &^= 0x01 << pos
	} else {
		b.bits[index] |= 0x01 << pos
	}

	return nil
}

func (b *NumaBitmap) GetBit(offset uint) (byte, error) {
	index, pos := offset/8, offset%8

	if b.userSize <= offset {
		return 0, fmt.Errorf("offset: %d is out of range %d", offset, b.userSize)
	}

	return (b.bits[index] >> pos) & 0x01, nil
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
				//err = fmt.Errorf("offset: %d is out of range %d", offset, maxNo)
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
func (b *NumaBitmap) Get1BitOffsNuma(nodeNum uint) ([][]uint, error) {
	var (
		tmp     uint
		offset  uint
		curNode uint
		err     error
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
				//err = fmt.Errorf("offset: %d is out of range %d", offset, maxNo)
				goto OUT
			}

			//exlude hyperthread
			if offset >= cpu {
				tmp = offset - cpu
			} else {
				tmp = offset
			}

			curNode = tmp / step
			if curNode >= nodeNum {
				err = fmt.Errorf("Node index out of range, curNode: %d, offset: %d, tmp: %d", curNode, offset, tmp)
				goto OUT
			}
			if (line>>uint(pos))&0x01 != 0 {
				offs[curNode] = append(offs[curNode], offset)
			}
		}
	}
OUT:
	return offs, err
}

// Get the offset of bits equal 0 all
func (b *NumaBitmap) Get0BitOffs() []uint {
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
				//err = fmt.Errorf("offset: %d is out of range %d", offset, maxNo)
				goto OUT
			}

			if (line>>uint(pos))&0x01 == 0 {
				offs = append(offs, offset)
			}
		}
	}
OUT:
	return offs
}

// Get the offsets of bits equal 0 per Node
func (b *NumaBitmap) Get0BitOffsNuma(nodeNum uint) ([][]uint, error) {
	var (
		tmp     uint
		offset  uint
		curNode uint
		err     error
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
				//err = fmt.Errorf("offset: %d is out of range %d", offset, maxNo)
				goto OUT
			}

			//exlude hyperthread
			if offset >= cpu {
				tmp = offset - cpu
			} else {
				tmp = offset
			}

			curNode = tmp / step
			if curNode >= nodeNum {
				err = fmt.Errorf("Node index out of range, curNode: %d, offset: %d, tmp: %d", curNode, offset, tmp)
				goto OUT
			}
			if (line>>uint(pos))&0x01 == 0 {
				offs[curNode] = append(offs[curNode], offset)
			}
		}
	}
OUT:
	return offs, err
}

func (b *NumaBitmap) String() string {

	offs := make([]uint, 0, b.userSize)

	var offset uint
	for offset = 0; offset < b.userSize; offset++ {
		if v, _ := b.GetBit(offset); v == 1 {
			offs = append(offs, offset)
		}
	}

	return fmt.Sprintf("%v", offs)
}

func align(n, align uint) uint {
	return (n + align - 1) & (^(align - 1))
}
