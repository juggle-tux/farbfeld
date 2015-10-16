# imagefile

[![Build Status](https://travis-ci.org/mehlon/imagefile.svg?branch=master)](https://travis-ci.org/mehlon/imagefile)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/mehlon/imagefile)

Imagefile is a simple image encoding format from suckless. See [SPECIFICATION](http://git.2f30.org/imagefile/tree/SPECIFICATION).

## Installation

    go get -u github.com/mehlon/imagefile

## Public API

- **[func Decode(r io.Reader) (image.Image, error)](https://godoc.org/github.com/mehlon/imagefile#Decode)**
- **[func DecodeConfig(r io.Reader) (image.Config, error)](https://godoc.org/github.com/mehlon/imagefile#DecodeConfig)**
- **[func Encode(w io.Writer, img image.Image) error](https://godoc.org/github.com/mehlon/imagefile#Encode)**

## Usage

**imagefile.Decode (`if2png`)**

```go
package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/mehlon/imagefile"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Fprintf(os.Stderr, "usage: %s\n", os.Args[0])
		os.Exit(1)
	}

	img, err := imagefile.Decode(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode imagefile: %v\n", err)
		os.Exit(1)
	}

	err = png.Encode(os.Stdout, img)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
```

**imagefile.Encode (`any2if`)**

```go
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
```

## License

This package is free and unemcumbered software released into the public domain.
For more information, see the included [UNLICENSE](https://github.com/mehlon/imagefile/blob/master/UNLICENSE) file.
