package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	staticPath = "/root/static"
	correctTag = "correct.png"
	wrongTag   = "wrong.png"
)

var imagePath = filepath.Join(staticPath, "images")

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

func DrawTag(w io.Writer, tagName string) error {
	tagFilePath := filepath.Join(imagePath, tagName)
	if _, err := os.Stat(tagFilePath); os.IsNotExist(err) {
		return fmt.Errorf("no tag file: %s found", tagFilePath)
	}
	f, err := os.Open(tagFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(w, f)
	return err
}
