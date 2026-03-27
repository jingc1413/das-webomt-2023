package agent

import (
	"encoding/json"
	"gomt/core/model"
	"gomt/core/proto/priv"
	"gomt/das/cgi"
	"strings"

	"github.com/pkg/errors"
)

// func (s *DasDeviceAgent) ServeCgiGetRequest(script string, query url.Values, timeout time.Duration) (*cgi.CGIResponse, error) {
// 	if !s.supportCGI || s.cgiHandler == nil {
// 		return nil, errors.New("not supported")
// 	}
// 	return s.cgiHandler.ServeGetRequest(script, query, timeout)
// }

// func (s *DasDeviceAgent) ServeCgiPostFile(script string, fieldName string, filename string, f io.Reader, timeout time.Duration) error {
// 	if !s.supportCGI || s.cgiHandler == nil {
// 		return errors.New("not supported")
// 	}
// 	return s.cgiHandler.ServePostFile(script, fieldName, fieldName, f, timeout)
// }

// func (s *DasDeviceAgent) ServeCgiPostJsonRequest(script string, query url.Values, data any, timeout time.Duration) (*cgi.CGIResponse, error) {
// 	if !s.supportCGI || s.cgiHandler == nil {
// 		return nil, errors.New("not supported")
// 	}
// 	return s.cgiHandler.ServePostJsonRequest(script, query, data, timeout)
// }

// func (s *DasDeviceAgent) ServeCgiRequest(method string, script string, query url.Values, body io.Reader, contentType string, timeout time.Duration) (*cgi.CGIResponse, error) {
// 	if !s.supportCGI || s.cgiHandler == nil {
// 		return nil, errors.New("not supported")
// 	}
// 	return s.cgiHandler.ServeRequest(method, script, query, body, contentType, timeout)
// }

func (s *DasDeviceAgent) ServeCgiQueryDevices(schema string, query string, cb func(total uint8, index uint8, info *model.DeviceInfo)) ([]*model.DeviceInfo, error) {
	if !s.supportCGI || s.cgiHandler == nil {
		return nil, errors.New("not supported")
	}
	return s.cgiHandler.ServeQueryDevices(schema, query, cb)
}

func (s *DasDeviceAgent) ServeCgiGetParameterValues(inputValues []priv.Parameter) ([]priv.Parameter, error) {
	if !s.supportCGI || s.cgiHandler == nil {
		return nil, errors.New("not supported")
	}
	return s.cgiHandler.ServeGetParameterValues(inputValues)
}

func (s *DasDeviceAgent) ServeCgiSetParameterValues(inputValues []priv.Parameter) ([]priv.Parameter, error) {
	if !s.supportCGI || s.cgiHandler == nil {
		return nil, errors.New("not supported")
	}
	return s.cgiHandler.ServeSetParameterValues(inputValues)
}

func (s *DasDeviceAgent) ServeCgiSetUpradeReboot() (bool, error) {
	if !s.supportCGI || s.cgiHandler == nil {
		return false, errors.New("not supported")
	}
	return s.cgiHandler.ServeSetUpgradeReboot()
}

func (s *DasDeviceAgent) ServeCgiStartUpgrade(filename string, force bool, byArm bool) (*cgi.UpgradeResponseData, error) {
	if s.supportCGI == false || s.cgiHandler == nil {
		return nil, errors.New("not supported")
	}

	resp, err := s.cgiHandler.ServeStartUpgrade(filename, force, byArm)
	if err != nil {
		return nil, errors.Wrap(err, "start upgrade")
	}
	if tmp, err := json.Marshal(resp); err == nil {
		s.log.Tracef("%v", tmp)
	}

	isSnmp := strings.Contains(filename, "SNMP")
	if byArm {
		oid := "TB4.P0B19"
		if force {
			oid = "TB4.P0B22"
		} else if isSnmp {
			oid = "TB4.P0B25"
		}
		if value, err := s.SetParameterValue(oid, "00"); err != nil {
			return nil, errors.Wrap(err, "set parameter value")
		} else if value.Code != "00" {
			return nil, errors.Errorf("the ugprade package is not found")
		}
	}

	return resp, nil
}

func (s *DasDeviceAgent) ServeCgiGetUpgradeFilePacketInfo(dir string, filename string) error {
	if s.supportCGI == false || s.cgiHandler == nil {
		return errors.New("not supported")
	}
	return s.cgiHandler.ServeGetUpgradeFilePacketInfo(dir, filename)
}

func (s *DasDeviceAgent) ServeCgiDeleteKeyAndLogs() error {
	if s.supportCGI == false || s.cgiHandler == nil {
		return errors.New("not supported")
	}
	return s.cgiHandler.ServeDeleteKeyAndLogs()
}
