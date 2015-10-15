package imagefile

import (
	"bytes"
	"image"
	"image/color"
	"testing"
)

var (
	r = color.RGBA{255, 0, 0, 255}
	g = color.RGBA{0, 255, 0, 255}
	b = color.RGBA{0, 0, 255, 255}

	testBuf = [][]color.RGBA{
		{r, g, b},
		{b, r, g},
		{g, b, r},
	}

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

func TestColorModel(t *testing.T) {
	m := Imagefile{}

	if got, want := m.ColorModel(), color.RGBAModel; got != want {
		t.Fatalf(`m.ColorModel() = %v, want %v`, got, want)
	}
}

func TestBounds(t *testing.T) {
	for i, tt := range []struct {
		w uint32
		h uint32
	}{
		{100, 100},
		{200, 100},
		{100, 200},
	} {
		m := Imagefile{Width: tt.w, Height: tt.h}

		if b := image.Rect(0, 0, int(tt.w), int(tt.h)); !m.Bounds().Eq(b) {
			t.Fatalf(`[%d] m.Bounds() = %v, want %v`, i, m.Bounds(), b.Bounds())
		}
	}
}

func TestAt(t *testing.T) {
	m := Imagefile{Buf: testBuf}

	for i, tt := range imagefileTestCases {
		got := m.At(tt.x, tt.y)

		gr, gg, gb, ga := got.RGBA()
		wr, wg, wb, wa := tt.c.RGBA()

		if gr != wr || gg != wg || gb != wb || ga != wa {
			t.Fatalf(`[%d] m.At(tt.x, tt.y).RGBA() = %v, want %v`, i, got, tt.c)
		}
	}
}

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
