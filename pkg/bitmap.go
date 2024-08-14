package bitmap

import (
	"fmt"
	"os"
)

type color_t uint32

type bitmap struct {
	// File Header
	bfType       [2]byte //0x4D42
	bfSize       uint32
	bfReserved1  uint16
	bfReserved2  uint16
	bfOffsetBits uint32

	// Info Header
	biSize          uint32
	biWidth         uint32
	biHeight        uint32
	biPlane         uint16
	biBitCount      uint16
	biCompression   uint32
	biSizeImage     uint32
	biXPelsPerMeter uint32
	biYPelsPerMeter uint32
	biClrUsed       uint32
	biClrImportant  uint32

	// 4 bytes each [FF000000, 00FFFF00, 00FFFF00, 00FFFF00, BB000000]
	pixels []color_t
}

func NewBitmap(width, height int) *bitmap {
	return &bitmap{
		biHeight:    uint32(height),
		biWidth:     uint32(width),
		bfType:      [2]byte{'B', 'M'},
		bfReserved1: 0,
		bfReserved2: 0,
		biPlane:     1,
		biBitCount:  24,
		pixels:      make([]color_t, height*width),
	}
}

func (b *bitmap) SetBitCount(bpp int) {
	b.biBitCount = uint16(bpp)
}

func OpenBitmap(filename string) (*bitmap, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	b, err := read(data)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (b *bitmap) PrintTerminal() {
	for y := 0; y < int(b.biHeight); y++ {
		for x := 0; x < int(b.biWidth); x++ {
			color := b.pixels[y*int(b.biWidth)+x]
			fmt.Print(cprintf(uint8(color>>16), uint8(color>>8), uint8(color)))
		}
		fmt.Println()
	}
}
