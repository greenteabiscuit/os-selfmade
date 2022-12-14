package main

import (
	"unsafe"
)

const (
	fbPhysAddr  uintptr = 0xa0000
	COL8_000000 uint16  = 0
	COL8_FF0000 uint16  = 1
	COL8_00FF00 uint16  = 2
	COL8_FFFF00 uint16  = 3
	COL8_0000FF uint16  = 4
	COL8_FF00FF uint16  = 5
	COL8_00FFFF uint16  = 6
	COL8_FFFFFF uint16  = 7
	COL8_C6C6C6 uint16  = 8
	COL8_840000 uint16  = 9
	COL8_008400 uint16  = 10
	COL8_848400 uint16  = 11
	COL8_000084 uint16  = 12
	COL8_840084 uint16  = 13
	COL8_008484 uint16  = 14
	COL8_848484 uint16  = 15
)

func main() {
	delay(1000)
	boxFill8(320, 20, 20, 120, 120, COL8_FF0000)
	boxFill8(320, 70, 50, 170, 150, COL8_00FF00)
	boxFill8(320, 120, 80, 220, 180, COL8_0000FF)

	delay(10000)
}

func boxFill8(xsize, x0, y0, x1, y1 int, color uint16) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			*(*uint16)(unsafe.Pointer(fbPhysAddr + uintptr(y*xsize) + uintptr(x))) = color
		}
	}
}

// delay implements a simple loop-based delay. The outer loop value is selected
// so that a reasonable delay is generated when running on virtualbox.
func delay(v int) {
	for i := 0; i < 684000; i++ {
		for j := 0; j < v; j++ {
		}
	}
}
