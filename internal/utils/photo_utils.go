package utils

import (
	"fmt"
	"os"

	"github.com/h2non/bimg"
)

func WaterMark(img []byte, watermark []byte) ([]byte, error) {
	newImg, err := bimg.NewImage(img).Watermark(bimg.Watermark{
		Text:       "© 2021 Public Surf",
		Opacity:    0.1,
		Width:      200,
		DPI:        100,
		Background: bimg.Color{R: 255, G: 255, B: 255},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return newImg, nil
}

func ResizeImg(img []byte, width int, height int) ([]byte, error) {
	newImg, err := bimg.NewImage(img).Resize(width, height)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return newImg, nil
}

func LoadImg(dir string, imageName string) ([]byte, error) {
	buffer, err := bimg.Read(dir + imageName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return buffer, nil
}

func SaveImg(img []byte, dest string) (err error) {
	err = bimg.Write(dest, img)
	if err != nil {
		return err
	}
	return nil
}
