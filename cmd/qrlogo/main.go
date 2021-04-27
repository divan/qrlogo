package main

// #include <stdlib.h>
import "C"
import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"unsafe"

	"os"

	"github.com/divan/qrlogo"
)

// the "export"-declaration are required

//export CreateQrCode
func CreateQrCode(qrCodeString string, qrCodePath string, qrCodeSize int) {
	qr, err := qrlogo.Encode(qrCodeString, nil, qrCodeSize)
	errcheck(err, "Failed to encode QR:")
	writeFileToFilesystem(*qr, qrCodePath)
}

//export CreateQrCodeAsBase64String
func CreateQrCodeAsBase64String(qrCodeString string, qrCodeSize int) *C.char {
	qr, err := qrlogo.Encode(qrCodeString, nil, qrCodeSize)
	errcheck(err, "Failed to encode QR:")
	base64String := base64.StdEncoding.EncodeToString(qr.Bytes())
	return C.CString(base64String)
}

//export CreateQrCodeWithLogo
func CreateQrCodeWithLogo(qrCodeString string, qrCodePath string, overlayLogoPath string, qrCodeSize int) {
	qr := qrCodeWithLogo(qrCodeString, overlayLogoPath, qrCodeSize)
	writeFileToFilesystem(*qr, qrCodePath)
}

//export CreateQrCodeWithLogoAsBase64String
func CreateQrCodeWithLogoAsBase64String(qrCodeString string, overlayLogoPath string, qrCodeSize int) *C.char {
	qr := qrCodeWithLogo(qrCodeString, overlayLogoPath, qrCodeSize)
	base64String := base64.StdEncoding.EncodeToString(qr.Bytes())
	return C.CString(base64String)
}

//export FreeUnsafePointer
func FreeUnsafePointer(cPointer *C.char) {
	C.free(unsafe.Pointer(cPointer))
}

func qrCodeWithLogo(qrCodeString string, overlayLogoPath string, qrCodeSize int) *bytes.Buffer {
	file, err := os.Open(overlayLogoPath)
	errcheck(err, "Failed to open logo:")
	defer file.Close()

	logo, _, err := image.Decode(file)
	errcheck(err, "Failed to decode PNG with logo:")

	qr, err := qrlogo.Encode(qrCodeString, logo, qrCodeSize)
	errcheck(err, "Failed to encode QR:")

	return qr
}

func writeFileToFilesystem(qrCode bytes.Buffer, qrCodePath string) {
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

func main() {
}
