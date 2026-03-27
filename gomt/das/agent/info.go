package agent

import (
	"gomt/core/model"
	"gomt/das/arm"
	"gomt/das/file"
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

func (s *DasDeviceAgent) GetDeviceInfo() *model.DeviceInfo {
	return s.info
}

func (s *DasDeviceAgent) UpdateDeviceInfoByQuery(queryInfo *model.DeviceInfo, needUpdate bool) error {
	if queryInfo != nil {
		s.info.RouteAddress = queryInfo.RouteAddress
		s.info.IpAddress = queryInfo.IpAddress
		s.info.ConnectState = queryInfo.ConnectState
		s.info.AlarmState = queryInfo.AlarmState
		s.info.MixedState = queryInfo.MixedState
		s.info.OpticalPort = queryInfo.OpticalPort
		s.info.VersionState = queryInfo.VersionState
		s.info.Setup()
	}
	if needUpdate {
		if err := s.UpdateDeviceInfo(); err != nil {
			return err
		}
	}
	return nil
}

func (s *DasDeviceAgent) UpdateDeviceInfo() error {
	updateInfo, err := s.GetParameterValuesOfDeviceInfo()
	if err != nil {
		return errors.Wrap(err, "get paramter values of device info")
	}
	if updateInfo != nil {
		s.info.SubID = updateInfo.SubID
		s.info.Version = updateInfo.Version
		s.info.ElementModelNumber = updateInfo.ElementModelNumber
		s.info.ElementSerialNumber = updateInfo.ElementSerialNumber
		s.info.DeviceName = updateInfo.DeviceName
		s.info.SiteName = updateInfo.SiteName
		s.info.InstalledLocation = updateInfo.InstalledLocation
		s.info.SystemTime = updateInfo.SystemTime
		s.info.UpTime = updateInfo.UpTime
		s.info.LifeTime = updateInfo.LifeTime
	}

	if s.info.SubID == 0 {
		s.info.RouteAddress = []byte{0, 0, 0, 0}
	}
	s.info.Setup()
	return nil
}

func (s *DasDeviceAgent) getDeviceTypeName() (string, error) {
	if s.isLocalDevice && s.supportArmIni {
		return s.getDeviceTypeNameFromLocal()
	} else {
		return s.getDeviceTypeNameFromRemote()
		//return s.getDeviceTypeNameFromRemote2()
	}
}

func (s *DasDeviceAgent) getDeviceTypeNameFromLocal() (string, error) {
	def := s.fileTypes.Get("XmlConfig")
	if def == nil {
		return "", errors.New("invalid file type")
	}
	_, f, err := file.ReadLocalFile(def.Dir, "MachineXml.xml")
	if err != nil {
		return "", errors.Wrap(err, "read MachineXml.xml")
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return "", errors.Wrap(err, "read content")
	}

	return arm.ReadDeviceTypeNameFromXmlFile(data)
}

func (s *DasDeviceAgent) getDeviceTypeNameFromRemote() (string, error) {
	resp, err := s.ServeHttpGet("/config/MachineXml.xml")
	if err != nil {
		return "", errors.Wrap(err, "http get")
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "read content")
	}
	return arm.ReadDeviceTypeNameFromXmlFile(data)
}

func (s *DasDeviceAgent) getDeviceTypeNameFromRemote2() (string, error) {
	data, err := s.ServeHttpGet2("/config/MachineXml.xml")
	if err != nil {
		return "", errors.Wrap(err, "http get")
	}
	return arm.ReadDeviceTypeNameFromXmlFile(data)
}

func (s *DasDeviceAgent) getSoftwareVersion() (string, error) {
	value, err := s.GetParameterValue("T02.P000A", nil)
	if err != nil {
		return "", errors.Wrap(err, "serve get parameter value")
	}
	if value == nil || value.Code != "00" {
		return "", errors.New("get parameter value error")
	}
	return cast.ToString(value.Value), nil
}
