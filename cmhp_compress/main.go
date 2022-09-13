package cmhp_compress

import (
	"bytes"
	"compress/flate"
	"io"
)

// Deflate Compress data
func Deflate(data []byte) ([]byte, error) {
	inputFile := new(bytes.Buffer)
	_, err := inputFile.Write(data)
	if err != nil {
		return nil, err
	}

	outputFile := new(bytes.Buffer)
	flateWriter, err := flate.NewWriter(outputFile, flate.BestCompression)
	if err != nil {
		return nil, err
	}

	defer flateWriter.Close()
	_, err = io.Copy(flateWriter, inputFile)
	if err != nil {
		return nil, err
	}

	err = flateWriter.Flush()
	if err != nil {
		return nil, err
	}

	return outputFile.Bytes(), nil
}

// Inflate Decompress data
func Inflate(data []byte) ([]byte, error) {
	inputFile := new(bytes.Buffer)
	_, err := inputFile.Write(data)
	if err != nil {
		return nil, err
	}
	flateReader := flate.NewReader(inputFile)
	defer flateReader.Close()

	outputFile := new(bytes.Buffer)

	io.Copy(outputFile, flateReader)

	return outputFile.Bytes(), nil
}
