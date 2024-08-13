package bitmap

import (
	"os"
)

type color_t uint32

type bitmap struct {
	// File Header
	bfType       [2]byte
	bfSize       uint32
	bfReserved1  [4]byte
	bfReserved2  [2]byte
	bfOffsetBits [4]byte

	// Info Header
	biSize          [4]byte
	biWidth         uint32
	biHeight        uint32
	biPlane         [2]byte
	biBitCount      [2]byte
	biCompression   [4]byte
	biSizeImage     [4]byte
	biXPelsPerMeter [4]byte
	biYPelsPerMeter [4]byte
	biClrUsed       [4]byte
	biClrImportant  [4]byte

	// 4 bytes each [FF000000, 00FFFF00, 00FFFF00, 00FFFF00, BB000000]
	pixels []color_t
}

func NewBitmap(width, height int) *bitmap {
	return &bitmap{
		biHeight: uint32(height),
		biWidth:  uint32(width),
		pixels:   make([]color_t, height*width),
	}
}

func OpenBitmap(filename string) (*bitmap, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return nil, nil
}
