package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	_ "github.com/mehlon/imagefile"
)

func main() {
	f, err := os.Open("apple.if")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	i, s, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
	w, err := os.Create("apple.png")
	if err != nil {
		panic(err)
	}
	defer w.Close()
	png.Encode(w, i)

}
