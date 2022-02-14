package cmhp_file

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Data interface {
	string | []byte
}

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

/*
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
*/

// Write bytes, text or struct as json to file
func Write(path string, data interface{}) error {
	// Create path for file
	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return err
	}

	switch data.(type) {
	case string:
		if err = ioutil.WriteFile(path, []byte(data.(string)), 0777); err != nil {
			return err
		}
	case []byte:
		if err = ioutil.WriteFile(path, data.([]byte), 0777); err != nil {
			return err
		}
	default:
		// Write as json
		data, err := json.Marshal(data)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path, data, 0777)
	}

	return nil
}

// Append bytes or text to file
func Append(path string, data interface{}) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	switch data.(type) {
	case string:
		if _, err = f.WriteString(data.(string)); err != nil {
			return err
		}
	case []byte:
		if _, err = f.Write(data.([]byte)); err != nil {
			return err
		}
	default:
		panic("Unknown type")
	}

	return nil
}

func List(path string) ([]fs.FileInfo, error) {
	return ioutil.ReadDir(path)
}

func Info(path string) (fs.FileInfo, error) {
	stat, err := os.Stat(path)
	return stat, err
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

func HashSha1(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func HashSha256(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func HashSha512(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha512.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func HashMd5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
