package imagefile

import (
	"bufio"
	"encoding/binary"
	"image"
	"io"
)

// Encode writes the Image img to w in imagefile format.
func Encode(w io.Writer, img image.Image) error {
	bounds := img.Bounds()
	width := uint32(bounds.Dx())
	height := uint32(bounds.Dy())

	header := make([]byte, 9+4+4)
	copy(header, "imagefile")
	binary.BigEndian.PutUint32(header[9:], width)
	binary.BigEndian.PutUint32(header[13:], height)
	_, err := w.Write(header)
	if err != nil {
		return err
	}

	if img, ok := img.(*image.RGBA); ok {
		_, err = w.Write(img.Pix)
		return err
	}

	bw := bufio.NewWriter(w)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			_, err := bw.Write([]byte{
				byte(r / 256),
				byte(g / 256),
				byte(b / 256),
				byte(a / 256),
			})
			if err != nil {
				return err
			}
		}
	}

	return bw.Flush()
}
