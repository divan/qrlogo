package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"os"

	"github.com/divan/qrlogo"
)

var (
	input  = flag.String("i", "logo.png", "Logo to be placed over QR code")
	output = flag.String("o", "qr.png", "Output filename")
	keep   = flag.Bool("k", true, "keep the color of logo")
	size   = flag.Int("size", 512, "Image size in pixels")
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

	qrlogo.DefaultEncoder.KeepColor = *keep

	qr, err := qrlogo.Encode(text, logo, *size)
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
