package main

import "C"
import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"os"

	"github.com/divan/qrlogo"
)

func main() {
}

// the "export"-declaration are required

//export CreateQrCode
func CreateQrCode(qrCodeString string, qrCodePath string, qrCodeSize int) {
	qr, err := qrlogo.Encode(qrCodeString, nil, qrCodeSize)
	errcheck(err, "Failed to encode QR:")

	writeFile(*qr, qrCodePath)
}

//export CreateQrCodeWithLogo
func CreateQrCodeWithLogo(qrCodeString string, qrCodePath string, overlayLogoPath string, qrCodeSize int) {
	file, err := os.Open(overlayLogoPath)
	errcheck(err, "Failed to open logo:")
	defer file.Close()

	logo, _, err := image.Decode(file)
	errcheck(err, "Failed to decode PNG with logo:")

	qr, err := qrlogo.Encode(qrCodeString, logo, qrCodeSize)
	errcheck(err, "Failed to encode QR:")

	writeFile(*qr, qrCodePath)
}

func writeFile(qrCode bytes.Buffer, qrCodePath string) {
	out, err := os.Create(qrCodePath)
	errcheck(err, "Failed to open output file:")
	out.Write(qrCode.Bytes())
	out.Close()
}

func errcheck(err error, str string) {
	if err != nil {
		fmt.Println(str, err)
		os.Exit(1)
	}
}
