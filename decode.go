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
	Buf    [][]color.RGBA // Buf[y][x]
}

func (i Imagefile) ColorModel() color.Model {
	return color.RGBAModel
}
func (i Imagefile) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(i.Width), int(i.Height))
}
func (i Imagefile) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(i.Bounds())) {
		return color.RGBA{}
	}
	return i.Buf[y][x]
}

func Decode(r io.Reader) (image.Image, error) {
	var img Imagefile
	bb := bufio.NewReader(r)

	io.CopyN(ioutil.Discard, bb, 9)

	if err := binary.Read(bb, binary.BigEndian, &img.Width); err != nil {
		return nil, err
	}
	if err := binary.Read(bb, binary.BigEndian, &img.Height); err != nil {
		return nil, err
	}

	img.Buf = make([][]color.RGBA, img.Height)
	for y := range img.Buf {
		img.Buf[y] = make([]color.RGBA, img.Width)
		for x := range img.Buf[y] {
			rgba := make([]byte, 4)
			if _, err := io.ReadFull(bb, rgba); err == io.EOF {
				return nil, io.ErrUnexpectedEOF
			} else if err != nil {
				return nil, err
			}
			img.Buf[y][x] = color.RGBA{rgba[0], rgba[1], rgba[2], rgba[3]}
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
