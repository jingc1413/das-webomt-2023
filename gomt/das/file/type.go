package file

import "strings"

type FileType struct {
	Name          string
	Dir           string
	Path          string
	Exts          string
	SupportUpload bool
	FieldName     string
}

func (m FileType) MatchExt(ext string) bool {
	if m.Exts == "" || m.Exts == "*" {
		return true
	}
	exts := strings.Split(m.Exts, ",")
	for _, v := range exts {
		if v == ext {
			return true
		}
	}
	return false
}

type FileTypes []*FileType

func (m FileTypes) Get(name string) *FileType {
	for _, v := range m {
		if v.Name == name {
			return v
		}
	}
	return nil
}
