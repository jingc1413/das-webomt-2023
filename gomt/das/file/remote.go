package file

import (
	"fmt"
	"gomt/das/cgi"
	"io"

	"github.com/pkg/errors"
)

type RemoteFileMgmt struct {
	handler *cgi.CGIHandler
	urlBase string
}

func (m *RemoteFileMgmt) ListFiles(def *FileType) ([]FileInfo, error) {
	result := []FileInfo{}
	files, err := m.handler.ServeListFiles(def.Dir)
	if err != nil {
		return nil, errors.Wrap(err, "list files by cgi")
	}
	for _, v := range files {
		result = append(result, FileInfo{
			FileName: v.FileName,
			FileSize: v.FileSize,
			ModTime:  v.ModTime,
		})
	}
	return result, nil
}

func (m *RemoteFileMgmt) GetFile(def *FileType, fname string) (*FileInfo, io.ReadSeekCloser, error) {
	if def.Path != "" {
		if def.Name == "Config" || def.Name == "Version" || def.Name == "PacketUpdateFile" {
			return ReadHttpGetFile(fmt.Sprintf("%v/%v", m.urlBase, def.Path), fname)
		}
	}

	f, err := m.handler.ServeGetFile(def.Dir, fname)
	if err != nil {
		return nil, nil, errors.Wrap(err, "get file by cgi")
	}

	info := FileInfo{
		FileName: fname,
		FileSize: 0,
		ModTime:  0,
	}
	return &info, ReadSeekCloser{f}, err
}

func (m *RemoteFileMgmt) SaveFile(def *FileType, fname string, f io.Reader) error {
	if err := m.handler.ServeSaveFile(def.Dir, fname, def.FieldName, f); err != nil {
		return errors.Wrap(err, "set file by cgi")
	}
	return nil
}

func (m *RemoteFileMgmt) RemoveFile(def *FileType, fname string) error {
	if err := m.handler.ServeRemoveFile(def.Dir, fname); err != nil {
		return errors.Wrap(err, "remove file by cgi")
	}
	return nil
}
