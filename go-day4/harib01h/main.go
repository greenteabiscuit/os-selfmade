package main

import (
	"unsafe"
)

const (
	fbPhysAddr  uintptr = 0xa0000
	COL8_000000 uint16  = 0 // black
	COL8_FF0000 uint16  = 1 // blue
	COL8_00FF00 uint16  = 2 // green
	COL8_FFFF00 uint16  = 3
	COL8_0000FF uint16  = 4
	COL8_FF00FF uint16  = 5
	COL8_00FFFF uint16  = 6
	COL8_FFFFFF uint16  = 7
	COL8_C6C6C6 uint16  = 8
	COL8_840000 uint16  = 9
	COL8_008400 uint16  = 10
	LIGHTBLUE   uint16  = 11 // light blue
	COL8_000084 uint16  = 12
	COL8_840084 uint16  = 13
	YELLOW      uint16  = 14 // yellow
	WHITE       uint16  = 15 // white
)

func main() {
	delay(1000)
	xsize, ysize := 320, 200
	boxFill8(xsize, 0, 0, xsize-1, ysize-29, LIGHTBLUE)
	boxFill8(xsize, 0, ysize-28, xsize-1, ysize-28, COL8_C6C6C6)
	boxFill8(xsize, 0, ysize-27, xsize-1, ysize-27, COL8_FFFFFF)
	boxFill8(xsize, 0, ysize-26, xsize-1, ysize-1, COL8_FFFF00)

	boxFill8(xsize, 3, ysize-24, 59, ysize-24, COL8_FFFFFF)
	boxFill8(xsize, 2, ysize-24, 2, ysize-4, COL8_FFFFFF)
	boxFill8(xsize, 3, ysize-4, 59, ysize-4, WHITE)
	boxFill8(xsize, 59, ysize-23, 59, ysize-5, WHITE)
	boxFill8(xsize, 2, ysize-3, 59, ysize-3, COL8_000000)
	boxFill8(xsize, 60, ysize-24, 60, ysize-3, COL8_000000)

	boxFill8(xsize, xsize-47, ysize-24, xsize-4, ysize-24, WHITE)
	boxFill8(xsize, xsize-47, ysize-23, xsize-47, ysize-4, WHITE)
	boxFill8(xsize, xsize-47, ysize-3, xsize-4, ysize-3, COL8_FFFFFF)
	boxFill8(xsize, xsize-3, ysize-24, xsize-3, ysize-3, COL8_FFFFFF)

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
