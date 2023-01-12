package main

func InitMouseCursor8(bc uint16) [256]uint16 {
	mouse := [256]uint16{}
	cursor := []byte(
		"**************.." +
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
			".............***")
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			if cursor[y*16+x] == '*' {
				mouse[y*16+x] = BLACK
			}
			if cursor[y*16+x] == 'O' {
				mouse[y*16+x] = WHITE
			}
			if cursor[y*16+x] == '.' {
				mouse[y*16+x] = bc
			}
		}
	}
	return mouse
}
