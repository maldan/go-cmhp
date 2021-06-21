package cmhp

import (
	"bytes"
	"compress/flate"
	"io"
)

func DataCompress(data []byte) []byte {
	inputFile := new(bytes.Buffer)
	inputFile.Write(data)

	outputFile := new(bytes.Buffer)
	flateWriter, _ := flate.NewWriter(outputFile, flate.BestCompression)
	defer flateWriter.Close()
	io.Copy(flateWriter, inputFile)
	flateWriter.Flush()
	return outputFile.Bytes()
}

func DataDecompress(data []byte) []byte {
	inputFile := new(bytes.Buffer)
	inputFile.Write(data)

	outputFile := new(bytes.Buffer)

	flateReader := flate.NewReader(inputFile)
	flateReader.Close()

	io.Copy(outputFile, flateReader)
	return outputFile.Bytes()
}
