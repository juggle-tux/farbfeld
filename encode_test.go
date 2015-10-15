package imagefile

import (
	"bytes"
	"image"
	"testing"
)

func TestEncode(t *testing.T) {
	w := new(bytes.Buffer)
	m := image.NewRGBA(image.Rect(0, 0, 3, 3))

	m.Set(0, 0, r)
	m.Set(0, 1, b)
	m.Set(0, 2, g)
	m.Set(1, 0, g)
	m.Set(1, 1, r)
	m.Set(1, 2, b)
	m.Set(2, 0, b)
	m.Set(2, 1, g)
	m.Set(2, 2, r)

	if err := Encode(w, m); err != nil {
		t.Fatalf(`unexpected error: %v`, err)
	}

	if !bytes.Equal(w.Bytes(), imagefileData) {
		t.Fatalf(`w.Bytes() != imagefileData`)
	}
}
