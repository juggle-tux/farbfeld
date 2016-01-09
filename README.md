# farbfeld

[![Build Status](https://travis-ci.org/mehlon/farbfeld.svg?branch=master)](https://travis-ci.org/mehlon/farbfeld)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/mehlon/farbfeld)

Farbfeld is a simple image encoding format from suckless. See [FORMAT](http://git.suckless.org/farbfeld/tree/FORMAT).

## Installation

    go get -u github.com/mehlon/farbfeld

See [godoc](https://godoc.org/github.com/mehlon/farbfeld) for documentation.

## Usage

**farbfeld.Decode (`ff2png`)**

```go
package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/mehlon/farbfeld"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Fprintf(os.Stderr, "usage: %s\n", os.Args[0])
		os.Exit(1)
	}

	img, err := farbfeld.Decode(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode farbfeld: %v\n", err)
		os.Exit(1)
	}

	err = png.Encode(os.Stdout, img)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
```

**farbfeld.Encode (`any2ff`)**

```go
package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/mehlon/farbfeld"
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

	err = farbfeld.Encode(os.Stdout, img)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
```

## License

This package is free and unemcumbered software released into the public domain.
For more information, see the included [UNLICENSE](blob/master/UNLICENSE) file.
