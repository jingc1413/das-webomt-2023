package parser

import (
	"github.com/pkg/errors"
)

type ReleaseInfo struct {
	Schema   string `yaml:"schema"`
	Revision string `yaml:"revision"`
	Devices  []struct {
		DeviceTypeName   string `yaml:"deviceTypeName"`
		SoftwareVersion  string `yaml:"softwareVersion"`
		SoftwareRevision string `yaml:"softwareRevision"`
	} `yaml:"devices"`
}

func (m ReleaseInfo) GetVersion(deviceTypeName string) (string, string, error) {
	for _, v := range m.Devices {
		if v.DeviceTypeName == deviceTypeName {
			return v.SoftwareVersion, v.SoftwareRevision, nil
		}
	}
	return "", "", errors.Errorf("invalid device type name, %v", deviceTypeName)
}
