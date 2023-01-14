package main

const PIC0_ICW1 = 0x0020
const PIC0_OCW2 = 0x0020
const PIC0_IMR = 0x0021
const PIC0_ICW2 = 0x0021
const PIC0_ICW3 = 0x0021
const PIC0_ICW4 = 0x0021
const PIC1_ICW1 = 0x00a0
const PIC1_OCW2 = 0x00a0
const PIC1_IMR = 0x00a1
const PIC1_ICW2 = 0x00a1
const PIC1_ICW3 = 0x00a1
const PIC1_ICW4 = 0x00a1

func PortWriteByte(uint16, uint16)

func InitPIC() {
	PortWriteByte(PIC0_IMR, 0xFF) // disable all interrupts
	PortWriteByte(PIC1_IMR, 0xFF) // disable all interrupts

	PortWriteByte(PIC0_ICW1, 0x11) // edge trigger mode
	PortWriteByte(PIC0_ICW2, 0x20) // receive IRQ0-7 on INT20-27
	PortWriteByte(PIC0_ICW3, 1<<2) // PIC1 is connected with IRQ2
	PortWriteByte(PIC0_ICW4, 0x01) // non buffer mode

	PortWriteByte(PIC0_ICW1, 0x11) // edge trigger mode
	PortWriteByte(PIC0_ICW2, 0x28) // receive IRQ8-15 on INT28-2f
	PortWriteByte(PIC0_ICW3, 2)    // PIC1 is connected with IRQ2
	PortWriteByte(PIC0_ICW4, 0x01) // non buffer mode

	PortWriteByte(PIC0_IMR, 0xFB) // disable all interrupts except for PIC1
	PortWriteByte(PIC1_IMR, 0xFF) // disable all interrupts
}
