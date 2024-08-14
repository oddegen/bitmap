package main

import (
	"log"

	bitmap "github.com/oddegen/bitmap/pkg"
)

func main() {
	b := bitmap.NewBitmap(5, 5)
	b.SetPixel(0, 0, 0xFF0000)
	b.SetPixel(1, 0, 0x00FFFF)
	b.SetPixel(2, 0, 0x00FFFF)
	b.SetPixel(3, 0, 0x00FFFF)
	b.SetPixel(4, 0, 0xBB0000)

	if err := b.WriteBitmap("out.bmp"); err != nil {
		log.Fatal(err)
	}

	b, err := bitmap.OpenBitmap("out.bmp")
	if err != nil {
		log.Fatal(err)
	}

	b.PrintTerminal()

}
