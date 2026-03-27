package file

import (
	"gomt/das/cgi"
	"io"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type FileInfo struct {
	FileName string
	FileSize int64
	ModTime  int64
}

type FileMgmtHandler interface {
	ListFiles(def *FileType) ([]FileInfo, error)
	GetFile(def *FileType, fname string) (*FileInfo, io.ReadSeekCloser, error)
	SaveFile(def *FileType, fname string, f io.Reader) error
	RemoveFile(def *FileType, fname string) error
}

type FileMgmt struct {
	types FileTypes
	h     FileMgmtHandler
	log   *logrus.Entry
}

func NewRemoteFileMgmt(ftypes FileTypes, handler *cgi.CGIHandler, urlBase string) *FileMgmt {
	h := &RemoteFileMgmt{
		handler: handler,
		urlBase: urlBase,
	}
	s := &FileMgmt{
		types: ftypes,
		h:     h,
	}
	s.log = logrus.WithFields(logrus.Fields{"filemgmt": "remote"})
	return s
}

func NewLocalFileMgmt(ftypes FileTypes) *FileMgmt {
	h := &LocalFileMgmt{}
	s := &FileMgmt{
		types: ftypes,
		h:     h,
	}
	s.log = logrus.WithFields(logrus.Fields{"filemgmt": "local"})
	return s
}

func (s *FileMgmt) ListFiles(ftype string) ([]FileInfo, error) {
	def := s.types.Get(ftype)
	if def == nil {
		return nil, errors.New("invalid file type")
	}
	return s.h.ListFiles(def)
}

func (s *FileMgmt) GetFile(ftype string, fname string) (*FileInfo, io.ReadSeekCloser, error) {
	def := s.types.Get(ftype)
	if def == nil {
		return nil, nil, errors.New("invalid file type")
	}
	return s.h.GetFile(def, fname)
}

func (s *FileMgmt) SaveFile(ftype string, fname string, f io.Reader) error {
	def := s.types.Get(ftype)
	if def == nil {
		return errors.New("invalid file type")
	}
	if !def.SupportUpload {
		return errors.New("not support")
	}
	return s.h.SaveFile(def, fname, f)
}

func (s *FileMgmt) RemoveFile(ftype string, fname string) error {
	def := s.types.Get(ftype)
	if def == nil {
		return errors.New("invalid file type")
	}
	return s.h.RemoveFile(def, fname)
}
