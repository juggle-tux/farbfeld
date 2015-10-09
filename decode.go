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
	return color.ModelFunc(func(c color.Color) color.Color {
		return c
	})
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
	r = uint32(c.R)
	g = uint32(c.G)
	b = uint32(c.B)
	a = uint32(c.A)

	return
}

func Decode(r io.Reader) (image.Image, error) {
	var img Imagefile
	b := bufio.NewReader(r)

	io.CopyN(ioutil.Discard, b, 9)

	binary.Read(b, binary.BigEndian, &img.Width)
	binary.Read(b, binary.BigEndian, &img.Height)
	img.Buf = make([][]Color, img.Height)

	for i := range img.Buf {
		img.Buf[i] = make([]Color, img.Width)
		for j := range img.Buf[i] {
			c := make([]byte, 4)
			_, err := b.Read(c)
			if err != nil {
				return nil, err
			}
			img.Buf[i][j] = Color{c[0], c[1], c[2], c[3]}
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