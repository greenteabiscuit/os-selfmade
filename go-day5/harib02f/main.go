package main

import (
	"unsafe"
)

const (
	fbPhysAddr  uintptr = 0xa0000
	BLACK       uint16  = 0
	BLUE        uint16  = 1
	GREEN       uint16  = 2
	COL8_FFFF00 uint16  = 3
	COL8_0000FF uint16  = 4
	COL8_FF00FF uint16  = 5
	COL8_00FFFF uint16  = 6
	LIGHTGRAY   uint16  = 7
	DARKGRAY    uint16  = 8
	COL8_840000 uint16  = 9
	COL8_008400 uint16  = 10
	LIGHTBLUE   uint16  = 11
	RED         uint16  = 12
	PINK        uint16  = 13
	YELLOW      uint16  = 14
	WHITE       uint16  = 15
)

func main() {
	delay(1000)
	xsize, ysize := 320, 200
	boxFill8(xsize, 0, 0, xsize-1, ysize-29, LIGHTBLUE)
	boxFill8(xsize, 0, ysize-28, xsize-1, ysize-28, LIGHTGRAY)
	boxFill8(xsize, 0, ysize-27, xsize-1, ysize-27, WHITE)
	boxFill8(xsize, 0, ysize-26, xsize-1, ysize-1, LIGHTGRAY)

	boxFill8(xsize, 3, ysize-24, 59, ysize-24, WHITE)
	boxFill8(xsize, 2, ysize-24, 2, ysize-4, WHITE)
	boxFill8(xsize, 3, ysize-4, 59, ysize-4, DARKGRAY)
	boxFill8(xsize, 59, ysize-23, 59, ysize-5, DARKGRAY)
	boxFill8(xsize, 2, ysize-3, 59, ysize-3, BLACK)
	boxFill8(xsize, 60, ysize-24, 60, ysize-3, BLACK)

	boxFill8(xsize, xsize-47, ysize-24, xsize-4, ysize-24, DARKGRAY)
	boxFill8(xsize, xsize-47, ysize-23, xsize-47, ysize-4, DARKGRAY)
	boxFill8(xsize, xsize-47, ysize-3, xsize-4, ysize-3, WHITE)
	boxFill8(xsize, xsize-3, ysize-24, xsize-3, ysize-3, WHITE)

	variable := "Golang OS 1"

	putfont8Asc(xsize, 11, 11, WHITE, []byte(variable))
	putfont8Asc(xsize, 10, 10, BLACK, []byte(variable))

	putfont8Asc(xsize, 10, 31, WHITE, []byte("scrnx = "))
	putfont8Asc(xsize, 10, 30, BLACK, []byte("scrnx = "))
	bs := convertIntToByteArray(xsize)
	putfont8Asc(xsize, 101, 31, WHITE, bs[:])
	putfont8Asc(xsize, 100, 30, BLACK, bs[:])
	delay(10000)
}

// can only show til 10 digits for now.
func convertIntToByteArray(n int) [10]byte {
	t := n
	count := 0
	for n > 0 {
		n = n / 10
		count++
	}
	bs := [10]byte{}

	i := count - 1
	for t > 0 {
		bs[i] = byte(t%10 + 48)
		t = t / 10
		i--
	}
	return bs
}

func putfont8Asc(xsize, x, y int, color uint16, s []byte) {
	for _, b := range s {
		putfont8(xsize, x, y, color, Letters[int(b)*16:])
		x += 8
	}
}

func putfont8(xsize, x, y int, color uint16, font []byte) {
	for i := 0; i < 16; i++ {
		d := font[i]
		if d&0x80 != 0 {
			*(*uint16)(unsafe.Pointer(fbPhysAddr + uintptr(y+i)*uintptr(xsize) + uintptr(x) + 0)) = color
		}
		if d&0x40 != 0 {
			*(*uint16)(unsafe.Pointer(fbPhysAddr + uintptr(y+i)*uintptr(xsize) + uintptr(x) + 1)) = color
		}
		if d&0x20 != 0 {
			*(*uint16)(unsafe.Pointer(fbPhysAddr + uintptr(y+i)*uintptr(xsize) + uintptr(x) + 2)) = color
		}
		if d&0x10 != 0 {
			*(*uint16)(unsafe.Pointer(fbPhysAddr + uintptr(y+i)*uintptr(xsize) + uintptr(x) + 3)) = color
		}
		if d&0x08 != 0 {
			*(*uint16)(unsafe.Pointer(fbPhysAddr + uintptr(y+i)*uintptr(xsize) + uintptr(x) + 4)) = color
		}
		if d&0x04 != 0 {
			*(*uint16)(unsafe.Pointer(fbPhysAddr + uintptr(y+i)*uintptr(xsize) + uintptr(x) + 5)) = color
		}
		if d&0x02 != 0 {
			*(*uint16)(unsafe.Pointer(fbPhysAddr + uintptr(y+i)*uintptr(xsize) + uintptr(x) + 6)) = color
		}
		if d&0x01 != 0 {
			*(*uint16)(unsafe.Pointer(fbPhysAddr + uintptr(y+i)*uintptr(xsize) + uintptr(x) + 7)) = color
		}
	}
}

func boxFill8(xsize, x0, y0, x1, y1 int, color uint16) {
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
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
