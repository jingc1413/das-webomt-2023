package file

import (
	"bytes"
	"fmt"
	"gomt/core/utils"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ReadSeekCloser struct {
	*bytes.Reader
}

func (rsc ReadSeekCloser) Close() error {
	return nil
}

func ReadHttpGetFile(base string, fname string) (*FileInfo, io.ReadSeekCloser, error) {
	url := fmt.Sprintf("%v/%v", base, fname)
	resp, err := utils.HttpGet(url)
	if err != nil {
		return nil, nil, errors.Wrap(err, "http get")
	}
	defer resp.Body.Close()
	buf, _ := io.ReadAll(resp.Body)

	info := FileInfo{
		FileName: fname,
		FileSize: int64(len(buf)),
		ModTime:  0,
	}
	return &info, ReadSeekCloser{bytes.NewReader(buf)}, err
}

func ReadLocalFile(dir string, name string) (*FileInfo, io.ReadSeekCloser, error) {
	filename := filepath.Join(dir, name)
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, nil, err
	} else if fi.IsDir() {
		return nil, nil, nil
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	return &FileInfo{
		FileName: fi.Name(),
		FileSize: fi.Size(),
		ModTime:  fi.ModTime().Unix(),
	}, f, nil
}

func ReadLocalDirFileList(dir string, exts string) ([]string, error) {
	_exts := strings.Split(exts, ",")
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
			sub, err := ReadLocalDirFileList(filename, exts)
			if err != nil {
				return out, errors.Wrap(err, "read dir "+entry.Name())
			}
			out = append(out, sub...)
		} else {
			matchExt := false
			for _, v := range _exts {
				if v == "*" || strings.HasSuffix(filename, v) {
					matchExt = true
					break
				}
			}
			if matchExt {
				out = append(out, filename)
			}
		}
	}
	return out, nil
}

func MkdirAll(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func MustDir(dir string) {
	if ok, err := IsDir(dir); err != nil {
		logrus.Fatal(err)
	} else if !ok {
		logrus.Fatal("invalid dir, " + dir)
	}
}

func IsDir(dir string) (bool, error) {
	stat, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return stat.IsDir(), nil
}
