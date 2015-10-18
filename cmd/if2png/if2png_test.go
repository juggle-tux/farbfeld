package main

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func colorsEqual(c, c2 color.Color) bool {
	r, g, b, a := c.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r/256 == r2/256 && g/256 == g2/256 && b/256 == b2/256 && a/256 == a2/256
}

func color2Str(c color.Color) string {
	r, g, b, a := c.RGBA()
	return fmt.Sprintf("(%d,%d,%d,%d)", r/256, g/256, b/256, a/256)
}

func Test(t *testing.T) {
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}

outer:
	for _, fi := range files {
		name := fi.Name()
		if !strings.HasSuffix(name, ".if") {
			continue
		}

		f, err := os.Open("testdata/" + name)
		if err != nil {
			t.Errorf("failed to open file %q", "testdata/"+name)
			continue
		}

		pngfile, err := os.Open("testdata/" + name + ".png")
		if err != nil {
			t.Errorf("failed to open file %q", "testdata/"+name+".png")
			continue
		}

		img, err := png.Decode(pngfile)
		pngfile.Close()
		if err != nil {
			t.Errorf("failed to decode %q", "testdata/"+name+".png")
			continue
		}

		var buf bytes.Buffer
		cmd := exec.Command("go", "run", "if2png.go")
		cmd.Stdin = f
		cmd.Stdout = &buf
		err = cmd.Run()
		f.Close()
		if err != nil {
			if e, ok := err.(*exec.ExitError); ok && !e.Success() {
				t.Error("if2png exited with non-zero exit status")
			} else {
				t.Errorf("failed to start if2png: %v", err)
			}
			continue
		}

		imgGot, err := png.Decode(&buf)
		if err != nil {
			t.Errorf("failed to decode output: %v", err)
			continue
		}

		if img.Bounds() != imgGot.Bounds() {
			t.Errorf("want bounds to be %v, got %v", img.Bounds(), imgGot.Bounds())
			continue
		}

		bounds := img.Bounds()
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				if !colorsEqual(img.At(x, y), imgGot.At(x, y)) {
					t.Errorf("invalid output for input file %q; want color at (%d,%d) to be %s, got %s", "testdata/"+name, x, y, color2Str(img.At(x, y)), color2Str(imgGot.At(x, y)))
					continue outer
				}
			}
		}
	}
}
