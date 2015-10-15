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

// Decode reads a Imagefile image from r and returns it as an image.Image.
func Decode(r io.Reader) (image.Image, error) {
	bb := bufio.NewReader(r)

	if _, err := io.CopyN(ioutil.Discard, bb, 9); err != nil {
		return nil, err
	}

	var (
		width  uint32
		height uint32
	)

	if err := binary.Read(bb, binary.BigEndian, &width); err != nil {
		return nil, err
	}
	if err := binary.Read(bb, binary.BigEndian, &height); err != nil {
		return nil, err
	}

	w := int(width)
	h := int(height)

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			rgba := make([]byte, 4)

			if _, err := io.ReadFull(bb, rgba); err == io.EOF {
				return nil, io.ErrUnexpectedEOF
			} else if err != nil {
				return nil, err
			}

			img.Set(x, y, color.RGBA{rgba[0], rgba[1], rgba[2], rgba[3]})
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
