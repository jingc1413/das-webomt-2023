package agent

import (
	"gomt/das/file"
	"io"
	"strings"

	"github.com/pkg/errors"
)

func (s *DasDeviceAgent) ListFiles(ftype string) ([]file.FileInfo, error) {
	if s.fileMgmt == nil {
		return nil, errors.New("not supported")
	}
	return s.fileMgmt.ListFiles(ftype)
}

func (s *DasDeviceAgent) ListFilesWithPrefix(ftype string, prefix string, ext string) ([]file.FileInfo, error) {
	if s.fileMgmt == nil {
		return nil, errors.New("not supported")
	}
	infos, err := s.fileMgmt.ListFiles(ftype)
	if err != nil {
		return nil, err
	}
	out := []file.FileInfo{}
	for _, info := range infos {
		if prefix != "" && strings.HasPrefix(info.FileName, prefix) {
			if ext == "" || strings.HasSuffix(info.FileName, ext) {
				out = append(out, info)
			}
		}
	}
	return out, nil
}

func (s *DasDeviceAgent) GetFile(ftype string, fname string) (*file.FileInfo, io.ReadSeekCloser, error) {
	if s.fileMgmt == nil {
		return nil, nil, errors.New("not supported")
	}
	return s.fileMgmt.GetFile(ftype, fname)
}

func (s *DasDeviceAgent) RemoveFile(ftype string, fname string) error {
	if s.fileMgmt == nil {
		return errors.New("not supported")
	}
	return s.fileMgmt.RemoveFile(ftype, fname)
}

func (s *DasDeviceAgent) SaveFile(ftype string, fname string, f io.Reader) error {
	if s.fileMgmt == nil {
		return errors.New("not supported")
	}
	return s.fileMgmt.SaveFile(ftype, fname, f)
}
