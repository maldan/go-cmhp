package cmhp_file

import (
	"encoding/json"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ReadBin(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	return data, err
}

func ReadText(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	return string(data), err
}

func ReadJSON(path string, v interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, v)
	return err
}

func WriteBin(path string, data []byte) error {
	os.MkdirAll(filepath.Dir(path), 0777)
	err := ioutil.WriteFile(path, data, 0777)
	return err
}

func WriteText(path string, data string) error {
	os.MkdirAll(filepath.Dir(path), 0777)
	err := ioutil.WriteFile(path, []byte(data), 0777)
	return err
}

func WriteJSON(path string, v interface{}) error {
	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return err
	}

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, []byte(data), 0777)
	return err
}

func List(path string) ([]fs.FileInfo, error) {
	return ioutil.ReadDir(path)
}

func Exists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return true
}

func Delete(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func Size(path string) int64 {
	stat, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return stat.Size()
}

func DeleteDir(path string) error {
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