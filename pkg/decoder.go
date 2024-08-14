package bitmap

import (
	"bytes"
	"encoding/binary"
	"io"
)

func read(data []byte) (*bitmap, error) {
	b := new(bitmap)
	buf := bytes.NewReader(data)

	err := b.readFileHeader(buf)
	if err != nil {
		return nil, err
	}

	err = b.readInfoHeader(buf)
	if err != nil {
		return nil, err
	}

	b.pixels = make([]color_t, b.biHeight*b.biWidth)
	row_length := int(((b.biWidth*uint32(b.biBitCount) + 31) &^ 31) >> 3)
	for y := int(b.biHeight) - 1; y >= 0; y-- {
		for x := 0; x < int(b.biWidth); x++ {
			pix, err := readNBytes(buf, 3)
			if err != nil {
				return nil, err
			}
			b.pixels[y*int(b.biWidth)+x] = color_t(pix[0]) | color_t(pix[1])<<8 | color_t(pix[2])<<16
		}

		padding := row_length % int(b.biWidth)
		_, err := io.CopyN(io.Discard, buf, int64(padding))
		if err != nil {
			return nil, err
		}
	}

	return b, err
}

func (bitmap *bitmap) readFileHeader(buf *bytes.Reader) error {

	bfType, err := readNBytes(buf, 2)
	if err != nil {
		return err
	}
	bitmap.bfType = [2]byte(bfType)

	bfSize, err := readNBytes(buf, 4)
	if err != nil {
		return err
	}
	bitmap.bfSize = uint32(bfSize[0]) | uint32(bfSize[1])<<8 | uint32(bfSize[2])<<16 | uint32(bfSize[3])<<24

	bfReserved1, err := readNBytes(buf, 2)
	if err != nil {
		return err
	}
	bitmap.bfReserved1 = binary.LittleEndian.Uint16(bfReserved1[:])

	bfReserved2, err := readNBytes(buf, 2)
	if err != nil {
		return err
	}
	bitmap.bfReserved2 = binary.LittleEndian.Uint16(bfReserved2[:])

	bfOffBits, err := readNBytes(buf, 4)
	if err != nil {
		return err
	}
	bitmap.bfOffsetBits = binary.LittleEndian.Uint32(bfOffBits[:])

	return nil
}

func (b *bitmap) readInfoHeader(buf *bytes.Reader) error {
	biSize, err := readNBytes(buf, 4)
	if err != nil {
		return err
	}
	b.biSize = binary.LittleEndian.Uint32(biSize)

	biWidth, err := readNBytes(buf, 4)
	if err != nil {
		return err
	}
	b.biWidth = binary.LittleEndian.Uint32(biWidth)

	biHeight, err := readNBytes(buf, 4)
	if err != nil {
		return err
	}
	b.biHeight = binary.LittleEndian.Uint32(biHeight)

	biPlane, err := readNBytes(buf, 2)
	if err != nil {
		return err
	}
	b.biPlane = binary.LittleEndian.Uint16(biPlane)

	biBitCount, err := readNBytes(buf, 2)
	if err != nil {
		return err
	}
	b.biBitCount = binary.LittleEndian.Uint16(biBitCount)

	biCompression, err := readNBytes(buf, 4)
	if err != nil {
		return err
	}
	b.biCompression = binary.LittleEndian.Uint32(biCompression)

	biSizeImage, err := readNBytes(buf, 4)
	if err != nil {
		return err
	}
	b.biSizeImage = binary.LittleEndian.Uint32(biSizeImage)

	biXPelsPerMeter, err := readNBytes(buf, 4)
	if err != nil {
		return err
	}
	b.biXPelsPerMeter = binary.LittleEndian.Uint32(biXPelsPerMeter)

	biYPelsPerMeter, err := readNBytes(buf, 4)
	if err != nil {
		return err
	}
	b.biYPelsPerMeter = binary.LittleEndian.Uint32(biYPelsPerMeter)

	biClrUsed, err := readNBytes(buf, 4)
	if err != nil {
		return err
	}
	b.biClrUsed = binary.LittleEndian.Uint32(biClrUsed)

	biClrImportant, err := readNBytes(buf, 4)
	if err != nil {
		return err
	}
	b.biClrImportant = binary.LittleEndian.Uint32(biClrImportant)

	return nil
}
