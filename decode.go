package farbfeld

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"io"
)

// A FormatError reports that the input is not a valid farbfeld.
// It is returned by Decode and DecodeConfig when the image header
// doesn't start with "farbfeld".
type FormatError string

func (e FormatError) Error() string {
	return "invalid farbfeld format: " + string(e)
}

func decodeConfig(r io.Reader) (int, int, error) {
	header := make([]byte, len("farbfeld")+4+4)
	_, err := io.ReadFull(r, header)
	if err != nil {
		return 0, 0, err
	}

	if !bytes.HasPrefix(header, []byte("farbfeld")) {
		return 0, 0, FormatError("unexpected magic number")
	}

	w := binary.BigEndian.Uint32(header[len("farbfeld"):])
	h := binary.BigEndian.Uint32(header[len("farbfeld")+4:])

	return int(w), int(h), nil
}

// Decode reads an farbfeld image from r and returns it as an image.Image.
func Decode(r io.Reader) (image.Image, error) {
	w, h, err := decodeConfig(r)
	if err != nil {
		return nil, err
	}

	img := image.NewRGBA64(image.Rect(0, 0, w, h))
	_, err = io.ReadFull(r, img.Pix)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// DecodeConfig returns the color model and dimensions of an farbfeld image without
// decoding the entire image.
func DecodeConfig(r io.Reader) (image.Config, error) {
	w, h, err := decodeConfig(r)
	if err != nil {
		return image.Config{}, err
	}

	return image.Config{
		ColorModel: color.RGBA64Model,
		Width:      w,
		Height:     h,
	}, nil
}

func init() {
	image.RegisterFormat("farbfeld", "farbfeld", Decode, DecodeConfig)
}
