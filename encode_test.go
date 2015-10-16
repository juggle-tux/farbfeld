package imagefile

import (
	"bytes"
	"image"
	"testing"
)

func TestEncode(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 3, 3))
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			img.Set(x, y, imagePixels[y][x])
		}
	}

	var buf bytes.Buffer
	err := Encode(&buf, img)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(imageData, buf.Bytes()) {
		t.Fatal("encoding error")
	}
}

func TestEncode64(t *testing.T) {
	img := image.NewRGBA64(image.Rect(0, 0, 3, 3))
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			img.SetRGBA64(x, y, imagePixels64[y][x])
		}
	}

	var buf bytes.Buffer
	err := Encode(&buf, img)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(imageData, buf.Bytes()) {
		t.Fatal("encoding error")
	}
}
