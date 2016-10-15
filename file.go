package cachego

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type File struct {
	dir       string
	extension string
}

func (f *File) fileName(key string) string {
	extension := "cachego.data"

	if f.extension != "" {
		extension = f.extension
	}

	h := sha256.New()
	h.Write([]byte(key))
	hash := hex.EncodeToString(h.Sum(nil))

	filename := hash

	if key == "" {
		filename += "__"
	}

	filePath := fmt.Sprintf(
		"%s%s%s.%s",
		f.dir,
		string(os.PathSeparator),
		filename,
		extension,
	)

	return filePath
}

func (f *File) Contains(key string) bool {
	_, err := os.Stat(
		f.fileName(key),
	)

	if err != nil {
		return false
	}

	return true
}

func (f *File) Delete(key string) bool {
	err := os.Remove(
		f.fileName(key),
	)

	if err != nil {
		return false
	}

	return true
}

func (f *File) Fetch(key string) (string, bool) {
	value, err := ioutil.ReadFile(
		f.fileName(key),
	)

	if err != nil {
		return "", false
	}

	return string(value[:]), true
}

func (f *File) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		if value, ok := f.Fetch(key); ok {
			result[key] = value
		}
	}

	return result
}

func (f *File) Flush() bool {

	dir, err := os.Open(f.dir)

	if err != nil {
		return false
	}

	defer dir.Close()

	names, _ := dir.Readdirnames(-1)

	for _, name := range names {
		os.Remove(filepath.Join(f.dir, name))
	}

	return true
}

func (f *File) Save(key string, value string, lifeTime time.Duration) bool {

	file := f.fileName(key)

	if err := ioutil.WriteFile(file, []byte(value), 0666); err != nil {
		return false
	}

	return true
}
