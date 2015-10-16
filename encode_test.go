package imagefile

import (
	"bytes"
	"image"
	"image/draw"
	"io/ioutil"
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

func BenchmarkEncode(b *testing.B) {
	m := image.NewRGBA(image.Rect(0, 0, 1000, 1000))

	draw.Draw(m, m.Bounds(), &image.Uniform{green}, image.ZP, draw.Src)

	m.Set(0, 0, red)
	m.Set(0, 10, blue)
	m.Set(100, 10, red)
	m.Set(600, 100, blue)
	m.Set(201, 20, blue)
	m.Set(12, 25, red)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Encode(ioutil.Discard, m)
	}
}
