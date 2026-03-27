package utils

import (
	"io"
	"os"
	"path"
	"strings"

	"github.com/alexmullins/zip"
	"github.com/pkg/errors"
)

func ExistsFile(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}

func ExistsDir(dir string) bool {
	stat, err := os.Stat(dir)
	if err == nil && stat.IsDir() {
		return true
	}
	return false
}

func GetFileList(dir string, ext string) ([]string, error) {
	out := []string{}
	stat, err := os.Stat(dir)
	if err != nil {
		return out, err
	}
	if !stat.IsDir() {
		return out, errors.New("invalid file path")
	}

	entrise, err := os.ReadDir(dir)
	if err != nil {
		return out, errors.Wrap(err, "read dir")
	}
	for _, entry := range entrise {
		filename := path.Join(dir, entry.Name())
		if entry.IsDir() {
			sub, err := GetFileList(filename, ext)
			if err != nil {
				return out, errors.Wrap(err, "read dir "+entry.Name())
			}
			out = append(out, sub...)
		} else {
			if ext == "" || ext == "*" || strings.HasSuffix(filename, ext) {
				out = append(out, filename)
			}
		}
	}
	return out, nil
}

func GetDirList(dir string) ([]string, error) {
	out := []string{}
	stat, err := os.Stat(dir)
	if err != nil {
		return out, err
	}
	if !stat.IsDir() {
		return out, errors.New("invalid file path")
	}
	entrise, err := os.ReadDir(dir)
	if err != nil {
		return out, errors.Wrap(err, "read dir")
	}
	for _, entry := range entrise {
		filename := path.Join(dir, entry.Name())
		if entry.IsDir() {
			out = append(out, filename)
		}
	}
	return out, nil
}

func ReadFileFromZipFile(zipfile string, pass string, filename string) ([]byte, error) {
	r, err := zip.OpenReader(zipfile)
	if err != nil {
		return nil, errors.Wrap(err, "open zip")
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name != filename {
			continue
		}
		f.SetPassword(pass)
		rc, err := f.Open()
		if err != nil {
			return nil, errors.Wrap(err, "open file")
		}
		defer rc.Close()

		data, err := io.ReadAll(rc)
		if err != nil {
			return nil, errors.Wrap(err, "read file")
		}
		return data, nil
	}
	return nil, errors.New("cant find file")
}
