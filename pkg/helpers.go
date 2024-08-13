package bitmap

import "fmt"

func calc_imagesize(width, height uint32) uint32 {
	row_length := width * 3
	row_length += 4 - (row_length % 4)
	if row_length%4 != 0 {
		panic("row_length % 4 != 0")
	}

	return height * row_length
}

func cprintf(r, g, b uint8) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm \033[0m", r, g, b)
}
