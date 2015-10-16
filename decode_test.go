package imagefile

import (
	"bytes"
	"image"
	"image/color"
	"io"
	"strings"
	"testing"
)

func TestDecodeConfig(t *testing.T) {
	config, err := DecodeConfig(bytes.NewReader(imageData))
	if err != nil {
		t.Fatal(err)
	}

	if config.ColorModel != color.RGBAModel {
		t.Error("expected RGBA color model")
	}

	if config.Width != 3 || config.Height != 3 {
		t.Errorf("want size %dx%d, got %dx%d", 3, 3, config.Width, config.Height)
	}
}

func TestDecodeConfigRejectsInvalidMagic(t *testing.T) {
	_, err := DecodeConfig(strings.NewReader("imagefil\xff\x00\x00\x00\x03\x00\x00\x00\x03"))
	if err != ErrNoMagic {
		t.Errorf("expected ErrNoMagic error, got %v", err)
	}
}

func TestDecodeConfigRejectsTruncatedHeader(t *testing.T) {
	_, err := DecodeConfig(strings.NewReader("imagefile\x00\x00\x00\x03\x00\x00\x00\x03"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = DecodeConfig(strings.NewReader("imagefile\x00\x00\x00\x03\x00\x00\x00"))
	if err != io.ErrUnexpectedEOF {
		t.Errorf("expected ErrUnexpectedEOF error, got %v", err)
	}

	_, err = DecodeConfig(strings.NewReader("imagefile\x00\x00\x00"))
	if err != io.ErrUnexpectedEOF {
		t.Errorf("expected ErrUnexpectedEOF error, got %v", err)
	}
}

func colorsEqual(c, c2 color.Color) bool {
	r, g, b, a := c.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r == r2 && g == g2 && b == b2 && a == a2
}

func TestDecode(t *testing.T) {
	img, err := Decode(bytes.NewReader(imageData))
	if err != nil {
		t.Fatal(err)
	}

	if img.ColorModel() != color.RGBAModel {
		t.Error("expected RGBA color model")
	}

	bounds := img.Bounds()
	want := image.Rect(0, 0, 3, 3)
	if bounds != want {
		t.Fatalf("want bounds %v, got %v", want, bounds)
	}

	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			got := img.At(x, y)
			want := imagePixels[y][x]
			if !colorsEqual(got, want) {
				t.Errorf("pixel at (%d,%d): want %v, got %v", x, y, want, got)
			}
		}
	}
}

func TestDecodeRejectsInvalidMagic(t *testing.T) {
	_, err := Decode(strings.NewReader("imagefil\xff\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00"))
	if err != ErrNoMagic {
		t.Errorf("expected ErrNoMagic error, got %v", err)
	}
}

func TestDecodeRejectsTruncatedData(t *testing.T) {
	_, err := Decode(strings.NewReader("imagefile\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = Decode(strings.NewReader("imagefile\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00"))
	if err != io.ErrUnexpectedEOF {
		t.Errorf("expected ErrUnexpectedEOF error, got %v", err)
	}
}
