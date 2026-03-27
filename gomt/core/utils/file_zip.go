package utils

import (
	"bytes"
	"io"

	"github.com/yeka/zip"
)

func ExtractFileFromEncryptedZip(zipFilePath string, fileName, password string) (*bytes.Buffer, error) {
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	for _, f := range r.File {
		if f.IsEncrypted() {
			f.SetPassword(password)
		}
		if f.Name != fileName {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, rc); err != nil {
			return nil, err
		}
		rc.Close()
		//b, _ := io.ReadAll(rc)
		//rt := strings.NewReader(string(b))
		//rnc := io.NopCloser(rt)
		//rc.Close()
		return buf, nil
	}
	return nil, nil
}
