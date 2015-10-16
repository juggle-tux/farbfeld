package imagefile

import (
	"bufio"
	"encoding/binary"
	"image"
	"io"
)

// Encode writes the Image m to w in Imagefile format.
func Encode(w io.Writer, m image.Image) error {
	bb := bufio.NewWriter(w)
	defer bb.Flush()

	if _, err := bb.WriteString(imagefileHeader); err != nil {
		return err
	}

	b := m.Bounds()

	if err := binary.Write(bb, binary.BigEndian, uint32(b.Dx())); err != nil {
		return err
	}
	if err := binary.Write(bb, binary.BigEndian, uint32(b.Dy())); err != nil {
		return err
	}

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := m.At(x, y).RGBA()

			if err := bb.WriteByte(byte(r)); err != nil {
				return err
			}
			if err := bb.WriteByte(byte(g)); err != nil {
				return err
			}
			if err := bb.WriteByte(byte(b)); err != nil {
				return err
			}
			if err := bb.WriteByte(byte(a)); err != nil {
				return err
			}
		}
	}

	return nil
}
