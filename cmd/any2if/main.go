package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/mehlon/imagefile"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Fprintf(os.Stderr, "usage: %s\n", os.Args[0])
		os.Exit(1)
	}

	img, _, err := image.Decode(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode input: %v\n", err)
		os.Exit(1)
	}

	err = imagefile.Encode(os.Stdout, img)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
