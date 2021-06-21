package cmhp

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func FileReadAsBin(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	return data, err
}

func FileReadAsText(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	return string(data), err
}

func FileWriteAsBin(path string, data []byte) error {
	os.MkdirAll(filepath.Dir(path), 0777)
	err := ioutil.WriteFile(path, data, 0777)
	return err
}

func FileWriteAsText(path string, data string) error {
	os.MkdirAll(filepath.Dir(path), 0777)
	err := ioutil.WriteFile(path, []byte(data), 0777)
	return err
}
