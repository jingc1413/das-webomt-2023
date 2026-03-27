package agent

import (
	"encoding/json"
	"fmt"
	"gomt/core/model"
	"gomt/core/proto/priv"
	"io"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

func (s *DasDeviceAgent) GetParameterDefine(privOid model.PrivObjectId) *model.ParameterDefine {
	return s.defs.GetParameterDefine(privOid)
}
func (s *DasDeviceAgent) GetParameterValues(inputValues []priv.Parameter) ([]priv.Parameter, error) {
	// s.log.Tracef("get parameter values")
	values, err := s.doGetParameterValues(inputValues)
	if err == nil {
		s.setServiceAvailable(true)
	}
	return values, err
}

func (s *DasDeviceAgent) SetParameterValues(inputValues []priv.Parameter) ([]priv.Parameter, error) {
	// s.log.Tracef("set parameter values")
	values, err := s.doSetParameterValues(inputValues)
	if err == nil {
		s.setServiceAvailable(true)
	}
	return values, err
}

func (s *DasDeviceAgent) GetParameterValue(id string, value any) (*priv.Parameter, error) {
	values := []priv.Parameter{
		{ID: id, Value: value},
	}
	returnValues, err := s.GetParameterValues(values)
	if err != nil {
		return nil, err
	}
	if len(returnValues) >= 1 {
		return &returnValues[0], nil
	}
	return nil, errors.New("get parameter value failed")
}
func (s *DasDeviceAgent) SetParameterValue(id string, value any) (*priv.Parameter, error) {
	values := []priv.Parameter{
		{ID: id, Value: value},
	}
	returnValues, err := s.SetParameterValues(values)
	if err != nil {
		return nil, err
	}
	if len(returnValues) >= 1 {
		return &returnValues[0], nil
	}
	return nil, errors.New("set parameter value failed")
}

func (s *DasDeviceAgent) SetParameterValueExists(id string, value any) (*priv.Parameter, error) {
	if def := s.defs.GetParameterDefine(model.PrivObjectId(id)); def != nil {
		return s.SetParameterValue(id, value)
	}
	return nil, nil
}

func (s *DasDeviceAgent) doGetParameterValues(inputValues []priv.Parameter) ([]priv.Parameter, error) {
	if s.supportPriv && s.privSess != nil {
		return s.privSrv.GetParameterValues(s.privSess, inputValues)
	} else if s.supportCGI {
		return s.ServeCgiGetParameterValues(inputValues)
	} else if s.supportAPI {
		values := []priv.Parameter{}
		resp, err := s.ServeHttpApiPostJson("/parameters/get", inputValues)
		if err != nil {
			return values, errors.Wrap(err, "post api request")
		}
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return values, errors.Wrap(err, "read api response")
		}
		if err := json.Unmarshal(data, &values); err != nil {
			return values, errors.Wrap(err, "read api response")
		}
		return values, nil
	}
	return nil, errors.New("not supported")
}

func (s *DasDeviceAgent) doSetParameterValues(inputValues []priv.Parameter) ([]priv.Parameter, error) {
	if s.supportPriv && s.privSess != nil {
		return s.privSrv.SetParameterValues(s.privSess, inputValues)
	} else if s.supportCGI {
		return s.ServeCgiSetParameterValues(inputValues)
	} else if s.supportAPI {
		values := []priv.Parameter{}
		resp, err := s.ServeHttpApiPostJson("/parameters/set", inputValues)
		if err != nil {
			return values, errors.Wrap(err, "post request")
		}
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return values, errors.Wrap(err, "read response body")
		}
		if err := json.Unmarshal(data, &values); err != nil {
			return values, errors.Wrap(err, "read response content")
		}
		return values, nil
	}
	return nil, errors.New("not supported")
}

func (s *DasDeviceAgent) GetParameterValuesOfDeviceInfo() (*model.DeviceInfo, error) {
	info := &model.DeviceInfo{
		DeviceTypeName: s.deviceTypeName,
	}
	if s.privSess != nil {
		info.SubID = s.privSess.SubID
	} else {
		info.SubID = cast.ToUint8(s.deviceSub)
	}

	inputs := []priv.Parameter{
		{ID: "T02.P0102"}, // Device Sub ID
		{ID: "T02.P000A"}, // Software Version
		{ID: "T02.P0030"}, // Device Name
		{ID: "T02.P0024"}, // Site Name
		{ID: "T02.P0023"}, // Installed Location Label
		{ID: "T02.P0150"}, // System time
	}

	if def := s.defs.GetParameterDefine("T02.P0004"); def != nil {
		inputs = append(inputs, priv.Parameter{ID: "T02.P0004"}) // Element Model Number
	}
	if def := s.defs.GetParameterDefine("T02.P0005"); def != nil {
		inputs = append(inputs, priv.Parameter{ID: "T02.P0005"}) // Element Serial Number
	}
	if def := s.defs.GetParameterDefine("T02.P0006"); def != nil {
		inputs = append(inputs, priv.Parameter{ID: "T02.P0006"}) // Element Model Number
	}
	if def := s.defs.GetParameterDefine("T02.P0007"); def != nil {
		inputs = append(inputs, priv.Parameter{ID: "T02.P0007"}) // Element Serial Number
	}

	if def := s.defs.GetParameterDefine("T02.P002C"); def != nil {
		inputs = append(inputs, priv.Parameter{ID: "T02.P002C"}) // Up time
	}
	if def := s.defs.GetParameterDefine("T02.P002B"); def != nil {
		inputs = append(inputs, priv.Parameter{ID: "T02.P002B"}) // Life uptime
	}
	values, err := s.GetParameterValues(inputs)
	if err != nil {
		return info, errors.Wrapf(err, "get parameter values")
	}
	for _, value := range values {
		if value.Code != "00" {
			return info, errors.Errorf("get parameter value of %v with fault code %v", value.ID, value.Code)
		}
		switch value.ID {
		case "T02.P0102":
			info.SubID = cast.ToUint8(value.Value)
		case "T02.P0004", "T02.P0006":
			info.ElementModelNumber = cast.ToString(value.Value)
		case "T02.P0005", "T02.P0007":
			info.ElementSerialNumber = cast.ToString(value.Value)
		case "T02.P000A":
			info.Version = cast.ToString(value.Value)
		case "T02.P0023":
			info.InstalledLocation = cast.ToString(value.Value)
		case "T02.P0024":
			info.SiteName = cast.ToString(value.Value)
		case "T02.P0030":
			info.DeviceName = cast.ToString(value.Value)
		case "T02.P0150":
			info.SystemTime = cast.ToInt64(value.Value)
		case "T02.P002C":
			info.UpTime = 0
			if args := strings.Split(cast.ToString(value.Value), ":"); len(args) == 3 {
				uptime := int64(0)
				for i, arg := range args {
					v2 := cast.ToInt64(arg)
					switch i {
					case 0:
						uptime += v2 * 86400
					case 1:
						uptime += v2 * 3600
					case 2:
						uptime += v2 * 60
					}
				}
				info.UpTime = uptime
			}
		case "T02.P002B":
			info.LifeTime = cast.ToInt64(value.Value) * 3600
		}
	}

	return info, nil
}

func (s *DasDeviceAgent) GetParameterValueOfSessionExpireSeconds() (int64, error) {
	param, err := s.GetParameterValue("T02.P0026", nil)
	if err != nil {
		return 0, errors.Wrap(err, "get parameter of session timeout")
	}
	if v := cast.ToInt64(param.Value); v > 0 {
		return v, nil
	}
	return 0, nil
}

func (s *DasDeviceAgent) SetParameterValueOfJumpEnable(enable bool) (bool, error) {
	input := uint8(0)
	if enable {
		input = 1
	}

	result, err := s.SetParameterValueExists("TB2-DEBUG.P01D3", input)
	if err != nil {
		return false, err
	}
	if result == nil {
		return true, nil
	}
	v := cast.ToUint8(result.Value)
	if v > 0 {
		s.log.Tracef("jump enabled")
	} else {
		s.log.Tracef("jump disabled")
	}
	return false, nil
}

func (s *DasDeviceAgent) ReadRegister(module uint8, offset uint16) (string, error) {
	inputValue := []any{
		fmt.Sprintf("%02X", module),
		offset,
		"04",
		"00000000",
	}

	result, err := s.GetParameterValue("TB2.P0CCC", inputValue)
	if err != nil {
		return "", errors.Wrap(err, "get parameter value")
	}
	logrus.Warnf("%v", result.Value)

	out := ""
	if v, ok := result.Value.([]any); ok {
		if len(v) != 4 {
			return "", errors.New("incorrect value length")
		}
		out = cast.ToString(v[3])
	}
	if result.Code == "00" {
		return out, nil
	}
	return out, errors.New("response with error code")
}

func (s *DasDeviceAgent) WriteRegister(module uint8, offset uint16, value string) (string, error) {
	inputValue := []any{
		fmt.Sprintf("%02X", module),
		offset,
		"04",
		value,
	}

	result, err := s.GetParameterValue("TB2.P0CCC", inputValue)
	if err != nil {
		return "", errors.Wrap(err, "get parameter value")
	}
	logrus.Warnf("%v", result.Value)

	out := ""
	if v, ok := result.Value.([]any); ok {
		if len(v) != 4 {
			return "", errors.New("incorrect value length")
		}
		out = cast.ToString(v[3])
	}
	if result.Code == "00" {
		return out, nil
	}
	return out, errors.New("response with error code")
}
