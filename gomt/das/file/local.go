package file

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type LocalFileMgmt struct {
}

func (m *LocalFileMgmt) ListFiles(def *FileType) ([]FileInfo, error) {
	result := []FileInfo{}

	filenames, err := ReadLocalDirFileList(def.Dir, def.Exts)
	if err != nil {
		return nil, err
	}
	for _, filename := range filenames {
		if stat, err := os.Stat(filename); err == nil && !stat.IsDir() {
			result = append(result, FileInfo{
				FileName: stat.Name(),
				FileSize: stat.Size(),
				ModTime:  stat.ModTime().Unix(),
			})
		}
	}
	return result, nil
}

func (m *LocalFileMgmt) GetFile(def *FileType, fname string) (*FileInfo, io.ReadSeekCloser, error) {
	return ReadLocalFile(def.Dir, fname)
}

func (m *LocalFileMgmt) SaveFile(def *FileType, fname string, f io.Reader) error {
	fileExtension := filepath.Ext(fname)
	if !def.MatchExt(fileExtension) {
		return errors.New("invalid file extension")
	}

	if err := MkdirAll(def.Dir); err != nil {
		return err
	}
	fw, err := os.Create(filepath.Join(def.Dir, fname))
	if err != nil {
		return err
	}
	defer fw.Close()

	if _, err := io.Copy(fw, f); err != nil {
		return err
	}
	return nil
}

func (m *LocalFileMgmt) RemoveFile(def *FileType, fname string) error {
	filename := filepath.Join(def.Dir, fname)
	if err := os.Remove(filename); err != nil {
		return err
	}
	return nil
}
