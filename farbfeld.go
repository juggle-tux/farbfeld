// Package farbfeld implements a Farbfeld decoder and encoder.
//
// The Farbfeld specification is at http://git.suckless.org/farbfeld/tree/FORMAT.
package farbfeld

import (
	"bufio"
	"encoding/binary"
	"image"
	"image/color"
	"io"
)

const (
	// Magic is the first 16 (8 + 4 + 4) bytes of a farbfeld file.
	Magic = "farbfeld" + "????" + "????"
)

// Encode writes the Image img to w in Farbfeld format.
func Encode(w io.Writer, img image.Image) error {
	bounds := img.Bounds()
	width := uint32(bounds.Dx())
	height := uint32(bounds.Dy())

	header := []byte(Magic)
	binary.BigEndian.PutUint32(header[8:12], width)
	binary.BigEndian.PutUint32(header[12:16], height)
	_, err := w.Write(header)
	if err != nil {
		return err
	}

	if img, ok := img.(*image.NRGBA64); ok {
		_, err = w.Write(img.Pix)
		return err
	}

	bw := bufio.NewWriter(w)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.NRGBA64Model.Convert(img.At(x, y)).(color.NRGBA64)
			err = binary.Write(bw, binary.BigEndian, c.R)
			if err != nil {
				return err
			}
			err = binary.Write(bw, binary.BigEndian, c.G)
			if err != nil {
				return err
			}
			err = binary.Write(bw, binary.BigEndian, c.B)
			if err != nil {
				return err
			}
			err = binary.Write(bw, binary.BigEndian, c.A)
			if err != nil {
				return err
			}
		}
	}

	return bw.Flush()
}

// Decode reads a Farbfeld image from r and returns it as an image.Image.
func Decode(r io.Reader) (image.Image, error) {
	cfg, err := DecodeConfig(r)
	if err != nil {
		return nil, err
	}

	img := image.NewNRGBA64(image.Rect(0, 0, cfg.Width, cfg.Height))
	_, err = io.ReadFull(r, img.Pix)
	return img, err
}

// DecodeConfig returns the color model and dimensions of a Farbfeld image without
// decoding the entire image.
func DecodeConfig(r io.Reader) (image.Config, error) {
	header := make([]byte, len(Magic))
	_, err := io.ReadFull(r, header)

	w := int(binary.BigEndian.Uint32(header[8:12]))
	h := int(binary.BigEndian.Uint32(header[12:16]))

	return image.Config{
		ColorModel: color.NRGBA64Model,
		Width:      w,
		Height:     h,
	}, err
}

func init() {
	image.RegisterFormat("farbfeld", Magic, Decode, DecodeConfig)
}
