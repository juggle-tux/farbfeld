package farbfeld

import (
	"bytes"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// TODO: re-add tests, in updated format.
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
