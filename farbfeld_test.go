package farbfeld

import (
	"bytes"
	"crypto/rand"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"testing/quick"
)

// TODO: add more tests, in updated format.
var (
	imageData   = []byte("farbfeld")
	red         = color.RGBA{255, 0, 0, 255}
	green       = color.RGBA{0, 255, 0, 255}
	blue        = color.RGBA{0, 0, 255, 255}
	grayTr      = color.RGBA{128, 128, 128, 128}
	imagePixels = [3][3]color.RGBA{
		{red, green, blue},
		{blue, grayTr, green},
		{green, blue, red},
	}
	red64         = color.RGBA64{65535, 0, 0, 65535}
	green64       = color.RGBA64{0, 65535, 0, 65535}
	blue64        = color.RGBA64{0, 0, 65535, 65535}
	grayTr64      = color.RGBA64{32768, 32768, 32768, 32768}
	imagePixels64 = [3][3]color.RGBA64{
		{red64, green64, blue64},
		{blue64, grayTr64, green64},
		{green64, blue64, red64},
	}
)

func Test(t *testing.T) {
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}

	for _, fi := range files {
		name := fi.Name()
		if !strings.HasSuffix(name, ".png") {
			continue
		}

		f, err := os.Open("testdata/" + name)
		if err != nil {
			t.Errorf("failed to open file %q", "testdata/"+name)
			continue
		}

		want, err := ioutil.ReadFile("testdata/" + name + ".ff")
		if err != nil {
			t.Errorf("failed to open file %q", "testdata/"+name+".ff")
			continue
		}

		img, err := png.Decode(f)
		if err != nil {
			t.Errorf("failed to decode png %q: %v", "testdata/"+name, err)
			continue
		}

		var buf bytes.Buffer
		err = Encode(&buf, img)
		if err != nil {
			t.Errorf("failed to encode farbfeld: %v", err)
			continue
		}

		if !bytes.Equal(buf.Bytes(), want) {
			t.Errorf("invalid output for input file %q", "testdata/"+name)
			continue
		}
	}
}

func TestQuickCheck(t *testing.T) {
	f := func(w, h uint8, pix [1 << 14]byte) bool {
		if w == 0 || h == 0 {
			return true
		}
		img1 := image.NewRGBA64(image.Rect(0, 0, int(w>>2), int(h>>2)))
		copy(img1.Pix, pix[:])
		var buf bytes.Buffer
		if Encode(&buf, img1) != nil {
			return false
		}
		img2raw, err := Decode(&buf)
		if err != nil {
			return false
		}
		img2, _ := img2raw.(*image.RGBA64)
		return bytes.Equal(img1.Pix, img2.Pix)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func BenchmarkEncode(b *testing.B) {
	img := image.NewRGBA64(image.Rect(0, 0, 256, 256))
	io.ReadFull(rand.Reader, img.Pix)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		Encode(&buf, img)
	}
}

func BenchmarkDecode(b *testing.B) {
	img := image.NewRGBA64(image.Rect(0, 0, 256, 256))
	io.ReadFull(rand.Reader, img.Pix)
	var buf bytes.Buffer
	Encode(&buf, img)
	b.ResetTimer()
	var r io.Reader
	for i := 0; i < b.N; i++ {
		r = bytes.NewReader(buf.Bytes())
		Decode(r)
	}
}
