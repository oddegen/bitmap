package main

import (
	"log"

	bitmap "github.com/oddegen/bitmap/pkg"
)

func main() {
	b := bitmap.NewBitmap(9, 9)
	b.Fill(0, 8, 0, 2, 0x00FF00)
	b.Fill(0, 8, 3, 5, 0x00FFFF)
	b.Fill(0, 8, 6, 8, 0x0000FF)

	if err := b.WriteBitmap("out.bmp"); err != nil {
		log.Fatal(err)
	}
}
