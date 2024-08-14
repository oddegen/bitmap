package bitmap

import "fmt"

func calc_imagesize(width, height uint32, bpp int) uint32 {
	row_length := ((width*uint32(bpp) + 31) &^ 31) >> 3
	if row_length%4 != 0 {
		panic("row_length % 4 != 0")
	}

	return height * row_length
}

func cprintf(r, g, b uint8) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm \033[0m", r, g, b)
}
