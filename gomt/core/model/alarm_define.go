package model

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//go:embed asserts/alarm_description.csv
var alarmDescData []byte

type AlarmDefines []AlarmDefine

func LoadDefaultAlarmDefines() (AlarmDefines, error) {
	return LoadAlaramDefinesFromData(alarmDescData)
}

func LoadAlaramDefinesFromData(data []byte) (AlarmDefines, error) {
	defs := AlarmDefines{}

	r := bytes.NewBuffer(data)
	reader := csv.NewReader(r)
	result, err := reader.ReadAll()
	if err != nil {
		return defs, errors.Wrap(err, "read csv file")
	}

	idIndex := -1
	nameIndex := -1
	descriptionIndex := -1
	probableCauseIndex := -1
	proposedActionIndex := -1
	for i, v := range result[0] {
		if v == "Alarm ID" {
			idIndex = i
		} else if v == "Alarm Name" {
			nameIndex = i
		} else if v == "Description" {
			descriptionIndex = i
		} else if v == "Probable Causes" {
			probableCauseIndex = i
		} else if v == "Proposed Action" {
			proposedActionIndex = i
		}
	}
	if nameIndex < 0 {
		return defs, nil
	}

	for i, row := range result {
		if i == 0 {
			continue
		}
		if len(row) <= idIndex {
			continue
		}
		trapAlarmIdText := row[idIndex]
		trapAlarmID, err := strconv.ParseInt(trapAlarmIdText, 10, 32)
		if err != nil {
			continue
		}

		name := row[nameIndex]
		name = strings.ReplaceAll(name, "-power", " Power")
		name = strings.ReplaceAll(name, "-slave", " Slave")
		name = strings.ReplaceAll(name, "-master", " Master")
		name = strings.ReplaceAll(name, "-temperature", " Temperature")
		//name = strings.ReplaceAll(name, "OP-Master", "OP Master Transceiver")
		//name = strings.ReplaceAll(name, "OP-Slave", "OP Slave Transceiver")
		//name = strings.ReplaceAll(name, "OP Master", "OP Master Transceiver")
		//name = strings.ReplaceAll(name, "OP Slave", "OP Slave Transceiver")
		name = strings.ReplaceAll(name, "-", " ")

		description := ""
		probableCause := ""
		proposedAction := ""
		id := strings.ReplaceAll(name, " ", "")
		if id == "" {
			continue
		}
		if len(row) > descriptionIndex {
			description = row[descriptionIndex]
		}
		if len(row) > proposedActionIndex {
			proposedAction = row[proposedActionIndex]
			proposedAction = strings.ReplaceAll(proposedAction, "\n", "; ")
		}
		if len(row) > probableCauseIndex {
			probableCause = row[probableCauseIndex]
			probableCause = strings.ReplaceAll(probableCause, "\n", "; ")
		}
		if description == "" {
			description = name
		}

		def := AlarmDefine{
			ID:             id,
			TrapAlarmID:    int(trapAlarmID),
			Name:           name,
			Description:    description,
			ProbableCause:  probableCause,
			ProposedAction: proposedAction,
		}
		if err := defs.AddAlarmDefine(def); err != nil {
			logrus.Error(err, "add alarm define")
			continue
		}
		if id == "FanDeviceLossAlarm" {
			def2 := AlarmDefine{
				ID:             "Fan1DeviceLossAlarm",
				Name:           name,
				Description:    description,
				ProbableCause:  probableCause,
				ProposedAction: proposedAction,
			}
			if err := defs.AddAlarmDefine(def2); err != nil {
				logrus.Error(err, "add alarm define")
			}
			def3 := AlarmDefine{
				ID:             "Fan2DeviceLossAlarm",
				Name:           name,
				Description:    description,
				ProbableCause:  probableCause,
				ProposedAction: proposedAction,
			}
			if err := defs.AddAlarmDefine(def3); err != nil {
				logrus.Error(err, "add alarm define")
			}
		}
		if id == "PowerInterruptionAlarm" {
			def2 := AlarmDefine{
				ID:             "Power1InterruptionAlarm",
				Name:           name,
				Description:    description,
				ProbableCause:  probableCause,
				ProposedAction: proposedAction,
			}
			if err := defs.AddAlarmDefine(def2); err != nil {
				logrus.Error(err, "add alarm define")
			}
			def3 := AlarmDefine{
				ID:             "Power2InterruptionAlarm",
				Name:           name,
				Description:    description,
				ProbableCause:  probableCause,
				ProposedAction: proposedAction,
			}
			if err := defs.AddAlarmDefine(def3); err != nil {
				logrus.Error(err, "add alarm define")
			}
		}
	}
	return defs, nil
}

func (m AlarmDefines) GetAlarmDefine(name string) *AlarmDefine {
	for _, v := range m {
		if v.Name == name {
			tmp := v
			return &tmp
		}
	}

	if re := regexp.MustCompile(`(FAN|Fan)\s*(\d) Device Loss Alarm`); re != nil {
		name = re.ReplaceAllString(name, "FAN $2 Alarm")
	}

	for _, v := range m {
		if v.Name == name {
			tmp := v
			return &tmp
		}
	}
	return nil
}

func (m *AlarmDefines) AddAlarmDefine(def AlarmDefine) error {
	// for _, v := range *m {
	// 	if v.ID == def.ID {
	// 		return errors.Errorf("alarm already exists, id=%v", def.ID)
	// 	}
	// }
	*m = append(*m, def)
	return nil
}

func (m *AlarmDefines) SetAlarmDefine(def AlarmDefine) {
	for i, v := range *m {
		if v.ID == def.ID {
			(*m)[i] = def
			return
		}
	}
	*m = append(*m, def)
}

type AlarmDefine struct {
	ID             string `json:"ID"`
	TrapAlarmID    int    `json:"TrapAlarmID"`
	Name           string `json:"Name"`
	Description    string `json:"Description"`
	ProbableCause  string `json:"ProbableCause,omitempty"`
	ProposedAction string `json:"ProposedAction,omitempty"`
}
