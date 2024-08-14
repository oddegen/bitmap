package bitmap

import (
	"bytes"
	"encoding/binary"
	"os"
)

func (b *bitmap) SetPixel(x, y int, c color_t) {
	if (x < 0 || x >= int(b.biWidth)) || (y < 0 || y >= int(b.biHeight)) {
		panic("error: set_pixel, out of bounds")
	}

	b.pixels[y*int(b.biWidth)+x] = c
}

func (b *bitmap) Fill(x1, x2, y1, y2 int, c color_t) {
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			b.SetPixel(x, y, c)
		}
	}
}

func (b *bitmap) writeFileHeader(f *os.File) error {
	file_header_size := 14
	info_header_size := 40

	_, err := f.Write(b.bfType[:])
	if err != nil {
		return err
	}

	b.bfSize = 14 + 40 + calc_imagesize(b.biWidth, b.biHeight, int(b.biBitCount))

	err = binary.Write(f, binary.LittleEndian, b.bfSize)
	if err != nil {
		return err
	}

	err = binary.Write(f, binary.LittleEndian, b.bfReserved1)
	if err != nil {
		return err
	}
	err = binary.Write(f, binary.LittleEndian, b.bfReserved2)
	if err != nil {
		return err
	}

	b.bfOffsetBits = uint32(file_header_size + info_header_size)
	return binary.Write(f, binary.LittleEndian, b.bfOffsetBits)
}

func (b *bitmap) writeInfoHeader(f *os.File) error {
	info_header := make([]byte, 0)
	b.biSize = 40
	b.biCompression = 0
	b.biSizeImage = calc_imagesize(b.biWidth, b.biHeight, int(b.biBitCount))
	b.biXPelsPerMeter = 0
	b.biYPelsPerMeter = 0
	b.biClrUsed = 0
	b.biClrImportant = 0

	info_header = binary.LittleEndian.AppendUint32(info_header, b.biSize)
	info_header = binary.LittleEndian.AppendUint32(info_header, b.biWidth)
	info_header = binary.LittleEndian.AppendUint32(info_header, b.biHeight)
	info_header = binary.LittleEndian.AppendUint16(info_header, b.biPlane)
	info_header = binary.LittleEndian.AppendUint16(info_header, b.biBitCount)
	info_header = binary.LittleEndian.AppendUint32(info_header, b.biCompression)
	info_header = binary.LittleEndian.AppendUint32(info_header, b.biSizeImage)
	info_header = binary.LittleEndian.AppendUint32(info_header, b.biXPelsPerMeter)
	info_header = binary.LittleEndian.AppendUint32(info_header, b.biYPelsPerMeter)
	info_header = binary.LittleEndian.AppendUint32(info_header, b.biClrUsed)
	info_header = binary.LittleEndian.AppendUint32(info_header, b.biClrImportant)

	return binary.Write(f, binary.LittleEndian, info_header)
}

func (b *bitmap) writePixelData(f *os.File) error {
	var buffer bytes.Buffer
	buffer.Grow(int(calc_imagesize(b.biWidth, b.biHeight, int(b.biBitCount))))

	row_length := int(((b.biWidth*uint32(b.biBitCount) + 31) &^ 31) >> 3)

	for y := int(b.biHeight - 1); y >= 0; y-- {
		for x := 0; x < int(b.biWidth); x++ {
			color := b.pixels[y*int(b.biWidth)+x]
			red := byte(color)
			green := byte(color >> 8)
			blue := byte(color >> 16)

			buffer.Write([]byte{blue, green, red})
		}

		if padding := row_length % int(b.biWidth); padding != 0 {
			for j := 0; j < padding; j++ {
				if err := buffer.WriteByte(0); err != nil {
					return err
				}
			}
		}
	}
	f.Write(buffer.Bytes())

	return nil
}

func (b *bitmap) WriteBitmap(filename string) error {
	f, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = b.writeFileHeader(f)
	if err != nil {
		return err
	}

	err = b.writeInfoHeader(f)
	if err != nil {
		return err
	}

	err = b.writePixelData(f)
	if err != nil {
		return err
	}

	return nil
}
