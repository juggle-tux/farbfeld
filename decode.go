// Package imagefile enables reading and writing files
// in the Suckless imagefile format:
// http://git.2f30.org/imagefile/tree/SPECIFICATION.
package imagefile

import (
	"bufio"
	"encoding/binary"
	"image"
	"image/color"
	"io"
	"io/ioutil"
)

type Imagefile struct {
	Width  uint32
	Height uint32
	Buf    [][]Color // Buf[y][x]
}

func (i Imagefile) ColorModel() color.Model {
	return color.RGBAModel
}
func (i Imagefile) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(i.Width), int(i.Height))
}
func (i Imagefile) At(x, y int) color.Color {
	return i.Buf[y][x]
}

type Color struct {
	R byte // Red
	G byte // Green
	B byte // Blue
	A byte // Alpha
}

func (c Color) RGBA() (r, g, b, a uint32) {
	return uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
}

func Decode(r io.Reader) (image.Image, error) {
	var img Imagefile
	bb := bufio.NewReader(r)

	io.CopyN(ioutil.Discard, bb, 9)

	binary.Read(bb, binary.BigEndian, &img.Width)
	binary.Read(bb, binary.BigEndian, &img.Height)

	img.Buf = make([][]Color, img.Height)
	for y := range img.Buf {
		img.Buf[y] = make([]Color, img.Width)
		for x := range img.Buf[y] {
			var r, g, b, a byte
			binary.Read(bb, binary.BigEndian, &r)
			binary.Read(bb, binary.BigEndian, &g)
			binary.Read(bb, binary.BigEndian, &b)
			binary.Read(bb, binary.BigEndian, &a)

			img.Buf[y][x] = Color{r, g, b, a}
		}
	}

	return img, nil

}

// DecodeConfig returns an empty image.Config:
// imagefile has no configuration.
func DecodeConfig(r io.Reader) (image.Config, error) {
	return image.Config{}, nil
}

func init() {
	image.RegisterFormat("imagefile", "imagefile", Decode, DecodeConfig)
}
