package bitmap

import (
	"bytes"
	"encoding/binary"
	"os"
)

func (b *bitmap) setPixel(x, y int, c color_t) {
	if (x < 0 || x >= int(b.biWidth)) || (y < 0 || y >= int(b.biHeight)) {
		panic("error: set_pixel, out of bounds")
	}

	b.pixels[y*int(b.biWidth)+x] = c
}

func (b *bitmap) Fill(x1, x2, y1, y2 int, c color_t) {
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			b.setPixel(x, y, c)
		}
	}
}

func (b *bitmap) WriteFileHeader(f *os.File) error {
	b.bfType = [2]byte{'B', 'M'}
	file_header_size := 14
	info_header_size := 40

	_, err := f.Write(b.bfType[:])
	if err != nil {
		return err
	}

	b.bfSize = 14 + 40 + calc_imagesize(b.biWidth, b.biHeight)

	err = binary.Write(f, binary.LittleEndian, b.bfSize)
	if err != nil {
		return err
	}

	reserved := uint32(0)
	err = binary.Write(f, binary.LittleEndian, reserved)
	if err != nil {
		return err
	}

	offset := uint32(file_header_size + info_header_size)
	return binary.Write(f, binary.LittleEndian, offset)
}

func (b *bitmap) WriteInfoHeader(f *os.File) error {
	info_header := make([]byte, 0)
	var header_size uint32 = 40
	var planes uint16 = 1
	var bits_per_pixel uint16 = 24
	var compression uint32 = 0
	var pixel_data_section_size = uint32(calc_imagesize(b.biWidth, b.biHeight))
	var pixelX_per_meter uint32 = 0
	var pixelY_per_meter uint32 = 0
	var colors_used uint32 = 0
	var important_colors uint32 = 0

	info_header = binary.LittleEndian.AppendUint32(info_header, header_size)
	info_header = binary.LittleEndian.AppendUint32(info_header, uint32(b.biWidth))
	info_header = binary.LittleEndian.AppendUint32(info_header, uint32(b.biHeight))
	info_header = binary.LittleEndian.AppendUint16(info_header, planes)
	info_header = binary.LittleEndian.AppendUint16(info_header, bits_per_pixel)
	info_header = binary.LittleEndian.AppendUint32(info_header, compression)
	info_header = binary.LittleEndian.AppendUint32(info_header, pixel_data_section_size)
	info_header = binary.LittleEndian.AppendUint32(info_header, pixelX_per_meter)
	info_header = binary.LittleEndian.AppendUint32(info_header, pixelY_per_meter)
	info_header = binary.LittleEndian.AppendUint32(info_header, colors_used)
	info_header = binary.LittleEndian.AppendUint32(info_header, important_colors)

	return binary.Write(f, binary.LittleEndian, info_header)
}

func (b *bitmap) WritePixelData(f *os.File) error {
	row_length := int(b.biWidth) * 3

	for y := int(b.biHeight - 1); y >= 0; y-- {
		for x := 0; x < int(b.biWidth); x++ {
			color := b.pixels[y*int(b.biWidth)+x]
			red := byte(color)
			green := byte(color >> 8)
			blue := byte(color >> 16)

			f.Write([]byte{blue, green, red})
		}

		if row_length%4 != 0 {
			var padding bytes.Buffer
			for j := 0; j < 4-(row_length%4); j++ {
				if err := padding.WriteByte(0); err != nil {
					return err
				}
			}
			f.Write(padding.Bytes())
		}
	}

	return nil
}

func (b *bitmap) WriteBitmap(filename string) error {
	f, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = b.WriteFileHeader(f)
	if err != nil {
		return err
	}

	err = b.WriteInfoHeader(f)
	if err != nil {
		return err
	}

	err = b.WritePixelData(f)
	if err != nil {
		return err
	}

	return nil
}
