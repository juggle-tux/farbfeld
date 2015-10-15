package imagefile

import (
	"bytes"
	"image"
	"image/color"
	"io"
	"testing"
)

var (
	r = color.RGBA{255, 0, 0, 255}
	g = color.RGBA{0, 255, 0, 255}
	b = color.RGBA{0, 0, 255, 255}

	imagefileTestCases = []struct {
		x int
		y int
		c color.RGBA
	}{
		{0, 0, r},
		{0, 1, b},
		{0, 2, g},
		{1, 0, g},
		{1, 1, r},
		{1, 2, b},
		{2, 0, b},
		{2, 1, g},
		{2, 2, r},
	}

	imagefileData = []byte("imagefile" +
		"\x00\x00\x00\x03\x00\x00\x00\x03\xff\x00\x00\xff\x00\xff\x00" +
		"\xff\x00\x00\xff\xff\x00\x00\xff\xff\xff\x00\x00\xff\x00\xff" +
		"\x00\xff\x00\xff\x00\xff\x00\x00\xff\xff\xff\x00\x00\xff")
)

func TestDecode(t *testing.T) {
	m, err := Decode(bytes.NewReader(imagefileData))
	if err != nil {
		t.Fatalf(`unexpected error: %v`, err)
	}

	for i, tt := range imagefileTestCases {
		got := m.At(tt.x, tt.y)

		gr, gg, gb, ga := got.RGBA()
		wr, wg, wb, wa := tt.c.RGBA()

		if gr != wr || gg != wg || gb != wb || ga != wa {
			t.Fatalf(`[%d] m.At(tt.x, tt.y).RGBA() = %v, want %v`, i, got, tt.c)
		}
	}
}

func TestDecode_withBlankImagefile(t *testing.T) {
	_, err := Decode(bytes.NewReader([]byte{}))
	if err == nil {
		t.Fatalf(`expected error`)
	}

	if err != io.ErrUnexpectedEOF {
		t.Fatalf(`err = %v, want io.ErrUnexpectedEOF`, err)
	}
}

func TestDecode_withInvalidImagefile(t *testing.T) {
	_, err := Decode(bytes.NewReader([]byte("invalid-imagefile")))
	if err == nil {
		t.Fatalf(`expected error`)
	}

	if err.Error() != "invalid Imagefile format: unexpected magic number" {
		t.Fatalf(`unexpected error: %v`, err)
	}
}

func TestDecodeConfig(t *testing.T) {
	dc, err := DecodeConfig(nil)
	if err != nil {
		t.Fatalf(`unexpected error: %v`, err)
	}

	if got, want := dc.Width, 0; got != want {
		t.Fatalf(`dc.Width = %d, want %d`, got, want)
	}

	if got, want := dc.Height, 0; got != want {
		t.Fatalf(`dc.Height = %d, want %d`, got, want)
	}
}

func TestRegisteredFormat(t *testing.T) {
	m, format, err := image.Decode(bytes.NewReader(imagefileData))
	if err != nil {
		t.Fatalf(`unexpected error: %v`, err)
	}

	if got, want := format, "imagefile"; got != want {
		t.Fatalf(`format = %q, want %q`, got, want)
	}

	for i, tt := range imagefileTestCases {
		got := m.At(tt.x, tt.y)

		gr, gg, gb, ga := got.RGBA()
		wr, wg, wb, wa := tt.c.RGBA()

		if gr != wr || gg != wg || gb != wb || ga != wa {
			t.Fatalf(`[%d] m.At(tt.x, tt.y).RGBA() = %v, want %v`, i, got, tt.c)
		}
	}
}
