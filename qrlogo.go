package qrlogo

import (
	"bytes"
	"image"
	"image/color"
	"image/png"

	qr "github.com/skip2/go-qrcode"
)

// Encoder defines settings for QR/Overlay encoder.
type Encoder struct {
	AlphaThreshold int
	GreyThreshold  int
	QRLevel        qr.RecoveryLevel
	KeepColor      bool
}

// DefaultEncoder is the encoder with default settings.
var DefaultEncoder = Encoder{
	AlphaThreshold: 2000,       // FIXME: don't remember where this came from
	GreyThreshold:  30,         // in percent
	QRLevel:        qr.Highest, // recommended, as logo steals some redundant space
	KeepColor:      true,       // keep the color of logo
}

// Encode encodes QR image, adds logo overlay and renders result as PNG.
func Encode(str string, logo image.Image, size int) (*bytes.Buffer, error) {
	return DefaultEncoder.Encode(str, logo, size)
}

// EncodeToBg encodes QR image, adds to bg image overlay and renders result as PNG.
func EncodeToBg(str string, bg image.Image, size, offsetX, offsetY, sx, sy, swidth, sheight int) (*bytes.Buffer, error) {
	return DefaultEncoder.EncodeToBg(str, bg, size, offsetX, offsetY, sx, sy, swidth, sheight)
}

// EncodeToBg encodes QR image, adds to bg image overlay and renders result as PNG.
func (e Encoder) EncodeToBg(str string, bg image.Image, size, offsetX, offsetY, sx, sy, swidth, sheight int) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	code, err := qr.New(str, e.QRLevel)
	if err != nil {
		return nil, err
	}

	img := code.Image(size)
	if swidth <= 0 {
		swidth = img.Bounds().Max.X
	}
	if sheight <= 0 {
		sheight = img.Bounds().Max.Y
	}

	e.overlayWithOffset(bg, img, sx, sy, swidth, sheight, offsetX, offsetY)

	err = png.Encode(&buf, bg)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

// Encode encodes QR image, adds logo overlay and renders result as PNG.
func (e Encoder) Encode(str string, logo image.Image, size int) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	code, err := qr.New(str, e.QRLevel)
	if err != nil {
		return nil, err
	}

	img := code.Image(size)
	//e.overlayLogo(img, logo)
	src := PalettedToRGBA(img.(*image.Paletted))
	e.overlayOrinalLogo(src, logo)

	err = png.Encode(&buf, src)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

// overlayLogo blends logo to the center of the QR code,
// changing all colors to black.
func (e Encoder) overlayLogo(dst, src image.Image) {
	grey := uint32(^uint16(0)) * uint32(e.GreyThreshold) / 100
	alphaOffset := uint32(e.AlphaThreshold)
	offset := dst.Bounds().Max.X/2 - src.Bounds().Max.X/2
	for x := 0; x < src.Bounds().Max.X; x++ {
		for y := 0; y < src.Bounds().Max.Y; y++ {
			if r, g, b, alpha := src.At(x, y).RGBA(); alpha > alphaOffset {
				col := color.Black
				if r > grey && g > grey && b > grey {
					col = color.White
				}
				dst.(*image.Paletted).Set(x+offset, y+offset, col)
			}
		}
	}
}

// overlayLogo blends logo to the center of the QR code,
// keep colors.
func (e Encoder) overlayOrinalLogo(dst, src image.Image) {
	offset := dst.Bounds().Max.X/2 - src.Bounds().Max.X/2
	for x := 0; x < src.Bounds().Max.X; x++ {
		for y := 0; y < src.Bounds().Max.Y; y++ {
			dst.(*image.RGBA).Set(x+offset, y+offset, src.At(x, y))
		}
	}
}

func (e Encoder) overlayWithOffset(dst, src image.Image, sx, sy, swidth, sheight, offsetX, offsetY int) {
	switch dst.(type) {
	case *image.NRGBA:
		dstRGBA := dst.(*image.NRGBA)
		for x := 0; x < swidth; x++ {
			for y := 0; y < sheight; y++ {
				dstRGBA.Set(x+offsetX, y+offsetY, src.At(sx+x, sy+y))
			}
		}
	case *image.RGBA:
		dstRGBA := dst.(*image.RGBA)
		for x := 0; x < swidth; x++ {
			for y := 0; y < sheight; y++ {
				dstRGBA.Set(x+offsetX, y+offsetY, src.At(sx+x, sy+y))
			}
		}
	}
}

// PalettedToRGBA use unsafe.Pointer for optimize
func PalettedToRGBA(src *image.Paletted) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	for x := 0; x < src.Bounds().Max.X; x++ {
		for y := 0; y < src.Bounds().Max.Y; y++ {
			dst.Set(x, y, src.At(x, y))
		}
	}

	return dst
}
