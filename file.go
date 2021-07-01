package cmhp

import (
	"io/fs"
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

func FileList(path string) ([]fs.FileInfo, error) {
	return ioutil.ReadDir(path)
}

func DirRemove(path string) error {
	d, err := os.Open(path)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(path, name))
		if err != nil {
			return err
		}
	}
	err = os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}
