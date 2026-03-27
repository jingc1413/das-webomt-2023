package agent

import (
	"encoding/csv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type AlarmEventLog struct {
	SiteName        string
	DeviceTypeName  string
	DeviceName      string
	DeviceSubID     string
	DeviceID        string
	SerialNumber    string
	Route           string
	SoftwareVersion string
	Location        string

	AlarmName     string
	EventTime     int64
	AlarmStatus   string
	AlarmSeverity string
}

func (s *DasDeviceAgent) GetAlarmLogs() ([]AlarmEventLog, error) {
	logs := []AlarmEventLog{}
	infos, err := s.ListFilesWithPrefix("DeviceLog", "event_alarm_log", ".csv")
	if err != nil {
		return logs, errors.Wrap(err, "get alarm log files")
	}

	for _, info := range infos {
		_, f, err := s.GetFile("DeviceLog", info.FileName)
		if err != nil {
			s.log.Error(errors.Wrapf(err, "get alarm log file, %v", info.FileName))
			continue
		}
		defer f.Close()

		reader := csv.NewReader(f)
		reader.FieldsPerRecord = -1
		data, err := reader.ReadAll()
		if err != nil {
			s.log.Error(errors.Wrapf(err, "read alarm log file, %v", info.FileName))
			continue
		}
		flag := false

		base := AlarmEventLog{}
		for _, row := range data {
			if len(row) < 1 || strings.TrimSpace(row[0]) == "" {
				break
			}
			if !flag {
				if row[0] == "Serial Number" {
					flag = true
					continue
				}
				switch row[0] {
				case "Device Type Name":
					base.DeviceTypeName = strings.Trim(row[1], "`")
				case "Device ID":
					base.DeviceID = strings.Trim(row[1], "`")
				case "Device Sub ID":
					base.DeviceSubID = strings.Trim(row[1], "`")
				case "Site Name":
					base.SiteName = strings.Trim(row[1], "`")
				case "Device Name":
					base.DeviceName = strings.Trim(row[1], "`")
				case "Installed Location Label":
					base.Location = strings.Trim(row[1], "`")
				}
			} else {
				//Serial Number,Route,SW version,Location,Site name,Device name,Device sub ID,Device ID,Alarm Name,Time And Date,Status,Severity
				if len(row) == 12 {
					log := base
					log.SerialNumber = strings.Trim(row[0], "`")
					log.Route = strings.Trim(row[1], "`")
					log.SoftwareVersion = strings.Trim(row[2], "`")
					log.Location = strings.Trim(row[3], "`")
					log.SiteName = strings.Trim(row[4], "`")
					log.DeviceName = strings.Trim(row[5], "`")
					log.DeviceSubID = strings.Trim(row[6], "`")
					log.DeviceID = strings.Trim(row[7], "`")
					log.AlarmName = strings.Trim(row[8], "`")
					if t, err := time.Parse("2006/01/02 15:04:05", strings.Trim(row[9], "`")); err == nil {
						log.EventTime = t.Unix()
					} else {
						s.log.Error(errors.Wrap(err, "parse time of alarm event log"))
					}
					log.AlarmStatus = strings.Trim(row[10], "`")
					log.AlarmSeverity = strings.Trim(row[11], "`")
					logs = append(logs, log)
				}
			}
		}
	}

	return logs, nil
}
