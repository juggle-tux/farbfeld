package farbfeld

import (
	"bytes"
	"crypto/rand"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"testing/quick"
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
		img1 := image.NewNRGBA64(image.Rect(0, 0, int(w>>2), int(h>>2)))
		copy(img1.Pix, pix[:])
		var buf bytes.Buffer
		if Encode(&buf, img1) != nil {
			return false
		}
		img2raw, err := Decode(&buf)
		if err != nil {
			return false
		}
		img2 := img2raw.(*image.NRGBA64)
		return bytes.Equal(img1.Pix, img2.Pix)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func BenchmarkEncodeNRGBA64(b *testing.B) {
	img := image.NewNRGBA64(image.Rect(0, 0, 1<<10, 1<<10))
	io.ReadFull(rand.Reader, img.Pix)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode(ioutil.Discard, img)
	}
}

func BenchmarkEncodeRGBA64(b *testing.B) {
	img := image.NewRGBA64(image.Rect(0, 0, 1<<10, 1<<10))
	io.ReadFull(rand.Reader, img.Pix)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode(ioutil.Discard, img)
	}
}

func BenchmarkDecode(b *testing.B) {
	img := image.NewNRGBA64(image.Rect(0, 0, 256, 256))
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
