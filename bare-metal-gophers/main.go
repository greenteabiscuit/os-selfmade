package main

import (
	"unsafe"
)

const (
	fbPhysAddr uintptr = 0xa0000
)

func main() {
	delay(1000)
	boxFill8(320, 20, 20, 120, 120, 15)
	boxFill8(320, 70, 50, 170, 150, 14)

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
