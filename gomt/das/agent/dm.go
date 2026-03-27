package agent

import (
	"gomt/core/model"

	"github.com/pkg/errors"
)

func (s *DasDeviceAgent) QueryDevices(
	schema string,
	queryString string,
	cb func(total uint8, index uint8, info *model.DeviceInfo,
	)) ([]*model.DeviceInfo, error) {
	// queryValue := "01020304060708090B0C0D0E0F"
	// if s.supportArmIni && s.info.SubID == 0 {
	// 	infos := []*model.DeviceInfo{s.info}
	// 	_infos, err := arm.ReadDeviceInfosFromArmIni(s.deviceConfigPath)
	// 	if err != nil {
	// 		return infos, errors.Wrap(err, "read device infos from arm ini file")
	// 	}
	// 	infos = append(infos, _infos...)
	// 	if cb != nil {
	// 		total := uint8(len(infos))
	// 		for i, info := range infos {
	// 			cb(total, uint8(i+1), info)
	// 		}
	// 	}
	// 	return infos, nil
	// } else
	if s.privSess != nil {
		return s.privSrv.QueryDevices(s.privSess, schema, queryString, cb)
	} else if s.supportCGI {
		return s.ServeCgiQueryDevices(schema, queryString, cb)
	}
	return nil, errors.New("not supported")
}
