package imagefile

import (
	"bufio"
	"encoding/binary"
	"image"
	"io"
)

func Encode(w io.Writer, m image.Image) error {
	bb := bufio.NewWriter(w)
	defer bb.Flush()

	_, err := bb.WriteString("imagefile")
	if err != nil {
		return err
	}

	b := m.Bounds()
	width := uint32(b.Max.X - b.Min.X)
	height := uint32(b.Max.Y - b.Min.Y)
	binary.Write(bb, binary.BigEndian, width)
	binary.Write(bb, binary.BigEndian, height)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := m.At(x, y).RGBA()
			r1, g1, b1, a1 := byte(r), byte(g), byte(b), byte(a)

			binary.Write(bb, binary.BigEndian, r1)
			binary.Write(bb, binary.BigEndian, g1)
			binary.Write(bb, binary.BigEndian, b1)
			binary.Write(bb, binary.BigEndian, a1)
		}
	}

	return nil
}