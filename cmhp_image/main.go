package cmhp_image

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/disintegration/imaging"
)

func Compress(path string, quality int) ([]byte, error) {
	srcImage, err := imaging.Open(path, imaging.AutoOrientation(true))
	if err != nil {
		return nil, err
	}

	outputFile := new(bytes.Buffer)
	jpeg.Encode(outputFile, srcImage, &jpeg.Options{
		Quality: quality,
	})
	return outputFile.Bytes(), nil
}

func Resolution(path string) (int, int, error) {
	srcImage, err := imaging.Open(path, imaging.AutoOrientation(true))
	if err != nil {
		return 0, 0, err
	}

	return srcImage.Bounds().Dx(), srcImage.Bounds().Dy(), nil
}

func Thumbnail(path string, width int, height int) ([]byte, error) {
	srcImage, err := imaging.Open(path, imaging.AutoOrientation(true))
	if err != nil {
		return nil, err
	}

	thumbnail := imaging.Thumbnail(srcImage, width, height, imaging.Lanczos)
	outputFile := new(bytes.Buffer)
	jpeg.Encode(outputFile, thumbnail, nil)
	return outputFile.Bytes(), nil
}

func Resize(path string, width int, height int) ([]byte, error) {
	srcImage, err := imaging.Open(path, imaging.AutoOrientation(true))
	if err != nil {
		return nil, err
	}

	thumbnail := imaging.Resize(srcImage, width, height, imaging.Lanczos)
	outputFile := new(bytes.Buffer)
	jpeg.Encode(outputFile, thumbnail, nil)
	return outputFile.Bytes(), nil
}

func Crop(path string, area [4]float64) ([]byte, error) {
	srcImage, err := imaging.Open(path, imaging.AutoOrientation(true))
	if err != nil {
		return nil, err
	}

	width, height := srcImage.Bounds().Dx(), srcImage.Bounds().Dy()
	x := area[0] * float64(width)
	y := area[1] * float64(height)
	w := area[2] * float64(width)
	h := area[3] * float64(height)
	thumbnail := imaging.Crop(srcImage, image.Rect(int(x), int(y), int(w), int(h)))
	outputFile := new(bytes.Buffer)
	jpeg.Encode(outputFile, thumbnail, nil)
	return outputFile.Bytes(), nil
}
