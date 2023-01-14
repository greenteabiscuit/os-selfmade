package main

func installIDT() int32

// InitIDT runs the appropriate CPU-specific initialization code for enabling
// support for interrupt handling.
func InitIDT() int32 {
	return installIDT()
}

func IntHandler21() {
	boxFill8(320, 0, 0, 32*8-1, 15, BLACK)
	delay(1000)
}

func IntHandler2c() {
	boxFill8(320, 0, 0, 32*8-1, 15, BLACK)
	delay(1000)
}

func IntHandler27() {
	boxFill8(320, 0, 0, 32*8-1, 15, BLACK)
	delay(1000)
}
