package qrlogo

import (
	"bytes"
	"image"
	"image/color"
	"image/png"

	qr "github.com/skip2/go-qrcode"
)

type ColorImageMap struct {
	image.Image
	custom map[image.Point]color.Color
}

func NewColorImageMap(img image.Image) *ColorImageMap {
	return &ColorImageMap{img, map[image.Point]color.Color{}}
}

func (m *ColorImageMap) Set(x, y int, c color.Color) {
	m.custom[image.Point{x, y}] = c
}

func (m *ColorImageMap) At(x, y int) color.Color {
	// Explicitly changed part: custom colors of the changed pixels:
	if c := m.custom[image.Point{x, y}]; c != nil {
		return c
	}
	// Unchanged part: colors of the original image:
	return m.Image.At(x, y)
}

// Encode encodes QR image, adds logo overlay and renders result as PNG.
func (e Encoder) EncodeColor(str string, logo image.Image, size int) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	code, err := qr.New(str, e.QRLevel)
	if err != nil {
		return nil, err
	}

	img := code.Image(size)
	my := NewColorImageMap(img)
	e.overlayLogoWithColor(my, logo)

	err = png.Encode(&buf, my)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

// overlayLogo blends logo to the center of the QR code
func (e Encoder) overlayLogoWithColor(dst *ColorImageMap, src image.Image) {
	offset := dst.Bounds().Max.X/2 - src.Bounds().Max.X/2
	for x := 0; x < src.Bounds().Max.X; x++ {
		for y := 0; y < src.Bounds().Max.Y; y++ {
			r, g, b, _ := src.At(x, y).RGBA()
			mycolor := color.RGBA{uint8(r), uint8(g), uint8(b), 255}
			dst.Set(x+offset, y+offset, mycolor)
		}
	}
}
