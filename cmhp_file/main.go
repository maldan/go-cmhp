package cmhp_file

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_compress"
	"github.com/maldan/go-cmhp/cmhp_crypto"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Data interface {
	string | []byte
}

type FileInfo struct {
	FullPath string
	Name     string
	Dir      string
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

func ReadCompressedJSON(path string, v interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	cdata, err := cmhp_compress.Inflate(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(cdata, v)
	return err
}

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

// WriteCompressed bytes, text or struct as json to file and compress it
func WriteCompressed(path string, data interface{}) error {
	// Create path for file
	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return err
	}

	switch data.(type) {
	case string:
		cdata, err := cmhp_compress.Deflate([]byte(data.(string)))
		if err != nil {
			return err
		}
		if err = ioutil.WriteFile(path, cdata, 0777); err != nil {
			return err
		}
	case []byte:
		cdata, err := cmhp_compress.Deflate(data.([]byte))
		if err != nil {
			return err
		}
		if err = ioutil.WriteFile(path, cdata, 0777); err != nil {
			return err
		}
	default:
		// Write as json
		data, err := json.Marshal(data)
		if err != nil {
			return err
		}
		cdata, err := cmhp_compress.Deflate(data)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path, cdata, 0777)
	}

	return nil
}

func WriteTemp(prefix string, data interface{}) (string, error) {
	tmpPath := fmt.Sprintf("%v/%v/%v",
		strings.ReplaceAll(os.TempDir(), "\\", "/"),
		prefix,
		cmhp_crypto.UID(24),
	)
	tmpPath = strings.ReplaceAll(tmpPath, "//", "/")
	err := Write(tmpPath, data)
	if err != nil {
		return "", err
	}
	return tmpPath, nil
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

func ReadJSONList[T any](path string) ([]T, error) {
	list, err := List(path)
	if err != nil {
		return nil, err
	}

	out := make([]T, 0)
	for _, f := range list {
		t := new(T)
		ReadJSON(f.FullPath, &t)
		out = append(out, *t)
	}
	return out, nil
}

func List(path string) ([]FileInfo, error) {
	list, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	out := make([]FileInfo, 0)
	for _, f := range list {
		absPath, _ := filepath.Abs(path + "/" + f.Name())
		absPath = strings.ReplaceAll(absPath, "\\", "/")

		out = append(out, FileInfo{
			FullPath: absPath,
			Dir:      strings.ReplaceAll(filepath.Dir(absPath), "\\", "/"),
			Name:     f.Name(),
		})
	}
	return out, nil
}

func ListAll(path string) ([]FileInfo, error) {
	list := make([]FileInfo, 0)

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip dir
			if info.IsDir() {
				return nil
			}

			absPath, _ := filepath.Abs(path)
			absPath = strings.ReplaceAll(absPath, "\\", "/")
			list = append(list, FileInfo{
				FullPath: absPath,
				Dir:      strings.ReplaceAll(filepath.Dir(absPath), "\\", "/"),
				Name:     info.Name(),
			})

			return nil
		})
	if err != nil {
		return list, err
	}

	return list, nil
}

func Info(path string) (fs.FileInfo, error) {
	stat, err := os.Stat(path)
	return stat, err
}

func Copy(from string, to string) error {
	source, err := os.Open(from)
	if err != nil {
		return err
	}
	defer source.Close()

	// Prepare dir
	os.MkdirAll(filepath.Dir(to), 0777)

	destination, err := os.Create(to)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
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

func BuildJSONIndex[T any](path string, outPath string) error {
	list, err := ReadJSONList[T](path)
	if err != nil {
		return err
	}
	err = Write(outPath, &list)
	return err
}
