package main

import (
	"unsafe"
)

const (
	fbPhysAddr  uintptr = 0xa0000
	IDTAddr     uintptr = 0x0026f800
	GDTAddr     uintptr = 0x101100
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

type SegmentDescriptor struct {
	LimitLow, BaseLow    uint16
	BaseMid, AccessRight uint8
	LimitHigh, BaseHigh  uint8
}

type GateDescriptor struct {
	OffsetLow, Selector  uint16
	DWCount, AccessRight uint8
	OffsetHigh           uint16
}

func asmIntHandler21()

func asmIntHandler2c()

func asmIntHandler27()

func load_idtr(uint32, uint32) uint32

func load_gdtr(uint32, uint32) uint32

func add(i int16, j int16) (int16, int16)

func GetGDTRAddress() int32

func GetGDTRSize() int16

func io_sti()

// func HandleInterrupt(f func())

func main() {
	delay(1000)
	idtAddr := GetGDTRAddress()
	f21 := asmIntHandler21
	f2c := asmIntHandler2c
	f27 := asmIntHandler27

	for i := 0; i < 8192; i++ {
		switch i {
		case 1:
			ar := 0x4092 | 0x8000
			limit := 0xffffffff / 0x1000
			base := 0
			*(*SegmentDescriptor)(unsafe.Pointer(GDTAddr + uintptr(i*8))) = *(&SegmentDescriptor{
				LimitLow:    uint16(limit & 0xffff),
				BaseLow:     uint16(base & 0xffff),
				BaseMid:     uint8((base >> 16) & 0xff),
				AccessRight: uint8(ar & 0xff),
				LimitHigh:   uint8(((limit >> 16) & 0x0f) | ((ar >> 8) & 0xf0)),
				BaseHigh:    uint8((base >> 24) & 0xff),
			})
		case 2:
			limit := 0x0007ffff
			// base := 0x00280000
			base := 0x00142e30
			ar := 0x409a
			*(*SegmentDescriptor)(unsafe.Pointer(GDTAddr + uintptr(i*8))) = *(&SegmentDescriptor{
				LimitLow:    uint16(limit & 0xffff),
				BaseLow:     uint16(base & 0xffff),
				BaseMid:     uint8((base >> 16) & 0xff),
				AccessRight: uint8(ar & 0xff),
				LimitHigh:   uint8(((limit >> 16) & 0x0f) | ((ar >> 8) & 0xf0)),
				BaseHigh:    uint8((base >> 24) & 0xff),
			})
		default:
			*(*SegmentDescriptor)(unsafe.Pointer(IDTAddr + uintptr(i*8))) = *(&SegmentDescriptor{
				LimitLow:    0,
				BaseLow:     0,
				BaseMid:     0,
				AccessRight: 0,
				LimitHigh:   0,
				BaseHigh:    0,
			})
		}
	}

	// gdtrsize := load_gdtr(0xFFFF, uint32(GDTAddr))

	gdtrsize := GetGDTRSize()

	for i := 0; i < 256; i++ {
		*(*GateDescriptor)(unsafe.Pointer(IDTAddr + uintptr(i*8))) = *(&GateDescriptor{
			OffsetLow:   0,
			Selector:    0,
			DWCount:     0,
			AccessRight: 0,
			OffsetHigh:  0,
		})
	}

	_ = load_idtr(0x7FF, uint32(IDTAddr))
	// _ = InitIDT()
	InitPIC()
	io_sti()

	for i := 0; i < 256; i++ {
		if i == 0x21 {
			*(*GateDescriptor)(unsafe.Pointer(IDTAddr + uintptr(i*8))) = *(&GateDescriptor{
				OffsetLow:   uint16(uintptr(unsafe.Pointer(&f21)) & 0xffff),
				Selector:    2 * 8,
				DWCount:     uint8((0x8e >> 8) & 0xff),
				AccessRight: uint8(0x8e & 0xff),
				OffsetHigh:  uint16((uintptr(unsafe.Pointer(&f21)) >> 16) & 0xffff),
			})
		}
		if i == 0x2c {
			*(*GateDescriptor)(unsafe.Pointer(IDTAddr + uintptr(i*8))) = *(&GateDescriptor{
				OffsetLow:   uint16(uintptr(unsafe.Pointer(&f2c)) & 0xffff),
				Selector:    2 * 8,
				DWCount:     uint8((0x8e >> 8) & 0xff),
				AccessRight: uint8(0x8e & 0xff),
				OffsetHigh:  uint16((uintptr(unsafe.Pointer(&f2c)) >> 16) & 0xffff),
			})
		}
		if i == 0x27 {
			*(*GateDescriptor)(unsafe.Pointer(IDTAddr + uintptr(i*8))) = *(&GateDescriptor{
				OffsetLow:   uint16(uintptr(unsafe.Pointer(&f27)) & 0xffff),
				Selector:    2 * 8,
				DWCount:     uint8((0x8e >> 8) & 0xff),
				AccessRight: uint8(0x8e & 0xff),
				OffsetHigh:  uint16((uintptr(unsafe.Pointer(&f27)) >> 16) & 0xffff),
			})
		}
	}
	// setGatedesc(*(*GateDescriptor)(unsafe.Pointer(IDTAddr + uintptr(0x21*8))), unsafe.Pointer(&f), 2*8, 0x8e)
	// HandleInterrupt(IntHandler21)

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

	variable := "Golang OS"

	putfont8Asc(xsize, 11, 11, WHITE, []byte(variable))
	putfont8Asc(xsize, 10, 10, BLACK, []byte(variable))

	putfont8Asc(xsize, 10, 31, WHITE, []byte("scrnx = "))
	putfont8Asc(xsize, 10, 30, BLACK, []byte("scrnx = "))
	bs := convertIntToByteArray(xsize)
	putfont8Asc(xsize, 101, 31, WHITE, bs[:])
	putfont8Asc(xsize, 100, 30, BLACK, bs[:])

	res1, res2 := add(int16(102), int16(102))

	resByte1 := convertIntToByteArray(int(res1))

	putfont8Asc(xsize, 101, 51, WHITE, resByte1[:])
	putfont8Asc(xsize, 100, 50, BLACK, resByte1[:])

	resByte2 := convertIntToByteArray(int(res2))

	putfont8Asc(xsize, 200, 51, WHITE, resByte2[:])
	putfont8Asc(xsize, 200, 50, BLACK, resByte2[:])

	idtbyte := convertIntToByteArray(int(InitIDT()) + 8)

	putfont8Asc(xsize, 250, 51, WHITE, idtbyte[:])
	putfont8Asc(xsize, 250, 50, BLACK, idtbyte[:])

	PortWriteByte(PIC0_IMR, 0xf9) // Allow PIC1&keyboard (11111001), if this is commented out, nothing will happen
	PortWriteByte(PIC1_IMR, 0xef) // Allow mouse (11101111)

	idtAddrByte := convertIntToByteArray(int(idtAddr))

	putfont8Asc(xsize, 180, 71, WHITE, []byte("gdtraddr"))
	putfont8Asc(xsize, 180, 70, BLACK, []byte("gdtraddr:"))
	putfont8Asc(xsize, 250, 71, WHITE, idtAddrByte[:])
	putfont8Asc(xsize, 250, 70, BLACK, idtAddrByte[:])

	// size++
	gdtr := (*SegmentDescriptor)(unsafe.Pointer(GDTAddr + 2*8))

	sizeByte := convertIntToByteArray(int(gdtr.AccessRight))

	putfont8Asc(xsize, 180, 91, WHITE, []byte("ar: "))
	putfont8Asc(xsize, 180, 90, BLACK, []byte("ar: "))
	putfont8Asc(xsize, 250, 91, WHITE, sizeByte[:])
	putfont8Asc(xsize, 250, 90, BLACK, sizeByte[:])

	gdtrsizebyte := convertIntToByteArray(int(gdtrsize))

	putfont8Asc(xsize, 180, 121, WHITE, []byte("gdtrsize: "))
	putfont8Asc(xsize, 180, 120, BLACK, []byte("gdtrsize: "))
	putfont8Asc(xsize, 250, 121, WHITE, gdtrsizebyte[:])
	putfont8Asc(xsize, 250, 120, BLACK, gdtrsizebyte[:])

	mouse := [256]uint16{}
	cursor := "**************.." +
		"*OOOOOOOOOOO*..." +
		"*OOOOOOOOOO*...." +
		"*OOOOOOOOO*....." +
		"*OOOOOOOO*......" +
		"*OOOOOOO*......." +
		"*OOOOOOO*......." +
		"*OOOOOOOO*......" +
		"*OOOO**OOO*....." +
		"*OOO*..*OOO*...." +
		"*OO*....*OOO*..." +
		"*O*......*OOO*.." +
		"**........*OOO*." +
		"*..........*OOO*" +
		"............*OO*" +
		".............***"
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			if cursor[y*16+x] == '*' {
				mouse[y*16+x] = BLACK
			}
			if cursor[y*16+x] == 'O' {
				mouse[y*16+x] = WHITE
			}
			if cursor[y*16+x] == '.' {
				mouse[y*16+x] = LIGHTBLUE
			}
		}
	}

	putBlock8_8(xsize, 16, 16, 100, 100, 16, mouse[:])

	delay(10000)
}

func setGatedesc(gd GateDescriptor, offset unsafe.Pointer, selector, ar int) {
	gd.OffsetLow = uint16(uintptr(offset) & 0xffff)
	gd.Selector = uint16(selector)
	gd.DWCount = uint8((ar >> 8) & 0xff)
	gd.AccessRight = uint8(ar & 0xff)
	gd.OffsetHigh = uint16((uintptr(offset) >> 16) & 0xffff)
	return
}

func putBlock8_8(vxsize, pxsize, pysize, px0, py0, bxsize int, buf []uint16) {
	for y := 0; y < pysize; y++ {
		for x := 0; x < pxsize; x++ {
			*(*uint16)(unsafe.Pointer(fbPhysAddr + uintptr(py0+y)*uintptr(vxsize) + uintptr(px0+x))) = buf[y*bxsize+x]
		}
	}
}

// can only show til 10 digits for now.
func convertIntToByteArray(n int) [20]byte {
	t := n
	count := 0
	for n > 0 {
		n = n / 10
		count++
	}
	bs := [20]byte{}

	i := count - 1
	if t < 0 {
		i = count
		bs[i] = '-'
	}

	for t > 0 {
		bs[i] = byte(t%10 + 48)
		t = t / 10
		i--
	}
	return bs
}

// can only show til 10 digits for now.
func convertIntToHexByteArray(n int) [20]byte {
	t := n
	count := 0
	for n > 0 {
		n = n / 16
		count++
	}
	bs := [20]byte{}

	i := count + 2
	if t < 0 {
		i = count
		bs[i] = '-'
	}

	for t > 0 {
		bs[i] = byte(t%10 + 48)
		t = t / 16
		i--
	}
	bs[i] = 'x'
	i--
	bs[i] = '0'
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

func sub(i int16) int16 {
	return i - 20
}
