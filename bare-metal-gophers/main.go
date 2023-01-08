package main

import (
	"unsafe"
)

const (
	fbPhysAddr uintptr = 0xa0000
)

func main() {
	delay(1000)
	for i := 0; i < 0xffff; i++ {
		var x uint16
		x = uint16(i)
		// white screen
		*(*uint16)(unsafe.Pointer(fbPhysAddr + uintptr(x))) = 15
	}

	delay(10000)
}

// delay implements a simple loop-based delay. The outer loop value is selected
// so that a reasonable delay is generated when running on virtualbox.
func delay(v int) {
	for i := 0; i < 684000; i++ {
		for j := 0; j < v; j++ {
		}
	}
}
