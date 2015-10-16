package imagefile

import (
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"io"
)

// ErrNoMagic is returned by Decode and DecodeConfig when the image header
// doesn't start with "imagefile".
var ErrNoMagic = errors.New("no magic")

func decodeConfig(r io.Reader) (int, int, error) {
	header := make([]byte, 9+4+4)
	_, err := io.ReadFull(r, header)
	if err != nil {
		return 0, 0, err
	}

	if string(header[:9]) != "imagefile" {
		return 0, 0, ErrNoMagic
	}

	w := binary.BigEndian.Uint32(header[9:])
	h := binary.BigEndian.Uint32(header[13:])

	return int(w), int(h), nil
}

// Decode reads an imagefile image from r and returns it as an image.Image.
func Decode(r io.Reader) (image.Image, error) {
	w, h, err := decodeConfig(r)
	if err != nil {
		return nil, err
	}

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	_, err = io.ReadFull(r, img.Pix)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// DecodeConfig returns the color model and dimensions of an imagefile image without
// decoding the entire image.
func DecodeConfig(r io.Reader) (image.Config, error) {
	w, h, err := decodeConfig(r)
	if err != nil {
		return image.Config{}, err
	}

	return image.Config{
		ColorModel: color.RGBAModel,
		Width:      w,
		Height:     h,
	}, nil
}

func init() {
	image.RegisterFormat("imagefile", "imagefile", Decode, DecodeConfig)
}
