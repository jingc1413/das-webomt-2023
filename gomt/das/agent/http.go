package agent

import (
	"fmt"
	"gomt/core/utils"
	"io"
	"net/http"
)

func (s *DasDeviceAgent) MakeHttpUrl(path string) string {
	url := fmt.Sprintf("%v%v", s.deviceUrlBase, path)
	return url
}

func (s *DasDeviceAgent) ServeHttpApiPostJson(path string, data any) (*http.Response, error) {
	u := s.deviceApiUrlBase + path
	return utils.HttpPostJson(u, data)
}

func (s *DasDeviceAgent) ServeHttpGet(path string) (*http.Response, error) {
	u := s.deviceUrlBase + path
	return utils.HttpGet(u)
}

func (s *DasDeviceAgent) ServeHttpGet2(path string) ([]byte, error) {
	u := s.deviceUrlBase + path
	return utils.HttpGet2(u)
}

func (s *DasDeviceAgent) ServeHttpPostJson(path string, data any) (*http.Response, error) {
	u := s.deviceUrlBase + path
	return utils.HttpPostJson(u, data)
}

func (s *DasDeviceAgent) ServeHttpPostFile(path string, fieldname string, filename string, f io.Reader) (*http.Response, error) {
	u := s.deviceUrlBase + path
	return utils.HttpPostFile(u, fieldname, filename, f)
}
