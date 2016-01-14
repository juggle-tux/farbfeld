package farbfeld

import (
	"bufio"
	"encoding/binary"
	"image"
	"io"
)

// Encode writes the Image img to w in Farbfeld format.
func Encode(w io.Writer, img image.Image) error {
	bounds := img.Bounds()
	width := uint32(bounds.Dx())
	height := uint32(bounds.Dy())

	header := make([]byte, len(Magic)+4+4)
	copy(header, Magic)
	binary.BigEndian.PutUint32(header[len(Magic):], width)
	binary.BigEndian.PutUint32(header[len(Magic)+4:], height)
	_, err := w.Write(header)
	if err != nil {
		return err
	}

	if img, ok := img.(*image.RGBA64); ok {
		_, err = w.Write(img.Pix)
		return err
	}

	bw := bufio.NewWriter(w)

	var r, g, b, a uint32
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a = img.At(x, y).RGBA()
			err = binary.Write(bw, binary.BigEndian, uint16(r))
			if err != nil {
				return err
			}
			err = binary.Write(bw, binary.BigEndian, uint16(g))
			if err != nil {
				return err
			}
			err = binary.Write(bw, binary.BigEndian, uint16(b))
			if err != nil {
				return err
			}
			err = binary.Write(bw, binary.BigEndian, uint16(a))
			if err != nil {
				return err
			}
		}
	}

	return bw.Flush()
}
