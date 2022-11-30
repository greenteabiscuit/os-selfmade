package main

//go:nosplit
func io_hlt();

//go:nosplit
func main() {
fin:
	io_hlt()
	goto fin
}
