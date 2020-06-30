package main

import (
	"flag"
	"fmt"
	"github.com/divan/qrlogo"
	"image"
	_ "image/png"
	"os"
)

var (
	input     = flag.String("i", "logo.png", "Logo to be placed over QR code")
	output    = flag.String("o", "qr.png", "Output filename")
	size      = flag.Int("size", 512, "Image size in pixels")
	withColor = flag.Bool("t", false, "\nUse logo with color")
)

func main() {
	flag.Usage = Usage
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	text := flag.Arg(0)

	file, err := os.Open(*input)
	errcheck(err, "Failed to open logo:")
	defer file.Close()

	logo, _, err := image.Decode(file)
	errcheck(err, "Failed to decode PNG with logo:")

	qr, err := Encode(text, logo, *size, *withColor)
	errcheck(err, "Failed to encode QR:")

	out, err := os.Create(*output)
	errcheck(err, "Failed to open output file:")
	out.Write(qr.Bytes())
	out.Close()

	fmt.Println("Done! Written QR image to", *output)
}

// Usage overloads flag.Usage.
func Usage() {
	fmt.Fprintln(os.Stderr, "Usage: qrlogo [options] text")
	flag.PrintDefaults()
}

func errcheck(err error, str string) {
	if err != nil {
		fmt.Println(str, err)
		os.Exit(1)
	}
}
