package cmhp_image

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"

	"github.com/disintegration/imaging"
	"github.com/maldan/go-cmhp/cmhp_crypto"
	"github.com/maldan/go-cmhp/cmhp_file"
	"github.com/maldan/go-cmhp/cmhp_process"
)

type MagickArgs struct {
	Quality    int
	Width      int
	Height     int
	Format     string
	InputPath  string
	InputData  []byte
	OutputPath string
}

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

func Thumbnail(path string, width int, height int, format string) ([]byte, error) {
	srcImage, err := imaging.Open(path, imaging.AutoOrientation(true))
	if err != nil {
		return nil, err
	}

	thumbnail := imaging.Thumbnail(srcImage, width, height, imaging.Lanczos)
	outputFile := new(bytes.Buffer)

	if format == "png" {
		png.Encode(outputFile, thumbnail)
	} else if format == "jpg" || format == "jpeg" {
		jpeg.Encode(outputFile, thumbnail, nil)
	} else {
		return nil, errors.New("unsupported format")
	}

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

func Magick(args MagickArgs) (string, error) {
	// Create temp file
	tempFile := os.TempDir() + "/" + cmhp_crypto.UID(10)
	if args.OutputPath == "" {
		args.OutputPath = tempFile
		if args.Format == "" {
			args.OutputPath += ".jpeg"
		} else {
			args.OutputPath += "." + args.Format
		}
	}

	// Use data insted of path
	if len(args.InputData) > 0 {
		args.InputPath = os.TempDir() + "/" + cmhp_crypto.UID(10)
		cmhp_file.WriteBin(args.InputPath, args.InputData)
		defer cmhp_file.Delete(args.InputPath)
	}

	// Prepare args
	params := make([]string, 0)
	params = append(params, "magick", args.InputPath)

	// Set quality
	if args.Quality > 0 {
		params = append(params, "-quality", fmt.Sprintf("%v", args.Quality))
	}

	// Set size
	if args.Width > 0 && args.Height > 0 {
		params = append(params,
			"-thumbnail", fmt.Sprintf("%vx%v^", args.Width, args.Height),
			"-gravity", "center",
			"-extent", fmt.Sprintf("%vx%v", args.Width, args.Height),
		)
	}

	// Set output
	params = append(params, args.OutputPath)

	// Prepare dir
	os.MkdirAll(path.Dir(args.OutputPath), 0777)

	// Process image
	_, err := cmhp_process.Exec(params...)
	return args.OutputPath, err
}
