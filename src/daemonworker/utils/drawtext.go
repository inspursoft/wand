package utils

import (
	"image"
	"image/color"
	"image/png"
	"io"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func DrawText(w io.Writer, label string) error {
	img := image.NewRGBA(image.Rect(0, 0, 150, 15))
	col := color.RGBA{0, 0, 255, 255}
	point := fixed.Point26_6{fixed.Int26_6(10 * 64), fixed.Int26_6(10 * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
	if err := png.Encode(w, img); err != nil {
		return err
	}
	return nil
}
