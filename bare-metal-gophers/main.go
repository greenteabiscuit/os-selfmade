package main

import (
	"unsafe"
)

const (
	fbWidth            = 80
	fbHeight           = 25
	fbPhysAddr uintptr = 0xa0000
)

func main() {
	delay(1000)
	// Display a string to the top-left corner of the screen one character
	// at a time.
	for i := 0; i < 0xffff; i++ {
		// 1st byte is color, second is letter
		/*
			|  green	|    H		|
			| 00100000	| 01001000	|
		*/
		var x uint16
		x = uint16(i)
		*(*uint16)(unsafe.Pointer(fbPhysAddr + uintptr(x))) = 15
	}

	delay(10000)

	// Since both fb and logo are slices we can use the built-in copy funtion
	// to draw our logo.
	// copy(fb, logo)
	// transition(fb)

	// For the final part of the demo we will run a complex piece of code
	// that is designed so that no variables escape to the heap. This is
	// important as any calls to the runtime memory allocator will cause
	// the machine to triple-fault.

}

func boxFill(scr []uint16, color uint16, xsize, x0, y0, x1, y1 uint16) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			scr[y*xsize+x] = color
		}
	}
}

// transition implements a slide transition using the current contents of the
// supplied framebuffer.
func transition(fb []uint16) {
	delay(5000)

	for i := 0; i < fbWidth; i++ {
		for y, off := 0, 0; y < fbHeight; y, off = y+1, off+fbWidth {
			// Even rows should slide one character to the left and
			// odd rows should slide one character to the right
			if y%2 == 0 {
				copy(fb[off:off+fbWidth], fb[off+1:off+fbWidth])
				fb[off+fbWidth-1] = ' '
			} else {
				copy(fb[off+1:off+fbWidth], fb[off:off+fbWidth-1])
				fb[off] = ' '
			}
		}
		delay(50)
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
