package layout

import (
	"fmt"
	"gomt/core/model"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func MustMakeElementFromParam(obj *model.ParameterTreeNode, info DeviceInfo, path []string) *Element {
	elem, err := MakeElementFromParam(obj, info, path)
	if err != nil {
		logrus.Fatal(errors.Wrapf(err, "make element from parameter %v", obj.Name))
	}
	return elem
}

var ignorePrivObjectIds = []string{
	"T02.P0572",
	"TB2-PA.P0060",
	"TB2-DEBUG.P01D3",
	"T02.P0831",
}

func ignorePrivObjectId(oid model.PrivObjectId) bool {
	for _, v := range ignorePrivObjectIds {
		if v == string(oid) {
			return true
		}
	}
	return false
}

func MakeElementFromParam(obj *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	if len(obj.Params) == 1 {
		def := obj.Params[0]
		param := makeParameter(def, info)
		def.SetPath(path)
		return param, nil
	} else if len(obj.Params) == 2 {
		def := obj.Params[0]
		def2 := obj.Params[1]
		param := makeParameter(def, info)
		param2 := makeParameter(def2, info)
		obj := NewParamGroupComponent(obj.Name, 2, param, param2)
		def.SetPath(path)
		def2.SetPath(path)
		return obj, nil
	}
	return nil, errors.Errorf("invalid parameter number for tree node, %v, %v", obj.Name, len(obj.Params))
}

func makeParameter(def *model.ParameterDefine, info DeviceInfo) *Element {
	if ignorePrivObjectId((def.PrivOid)) {
		return nil
	}
	param := NewParam(def.Name, string(def.PrivOid), string(def.Access))

	switch def.InputType {
	case "number":
		param.SetStyle("input", "number")
	case "binary":
		param.SetStyle("input", "binary")
	case "password":
		param.SetStyle("input", "password")
	case "datetime":
		param.SetStyle("input", "datetime")
	case "ipv4":
		param.SetStyle("input", "ipv4")
	case "ipv6":
		param.SetStyle("input", "ipv6")
	case "button":
		param.Access = "wo"
		param.SetStyle("input", "button")
	case "buttonGroup":
		param.Access = "wo"
		param.SetStyle("input", "buttonGroup")
		param.Items = []*Element{}
		keys := []string{}
		for k, _ := range def.Options {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] > keys[j]
		})
		for _, k := range keys {
			tmp := NewParam(param.Name, param.OID, param.Access)
			tmp.SetValue((k))
			action := NewSetParameterValuesAction(tmp)
			action.SetStyle("showMessage", true)
			param2 := NewParam(param.Name, param.OID, param.Access)
			param2.SetValue(k)
			param2.SetAction("click", action)
			param2.SetStyle("input", "button")
			param2.SetStyle("plain", true)
			param2.SetStyle("type", "primary")
			param.Items = append(param.Items, param2)
		}
	case "switch":
		activeValue := ""
		activeText := ""
		inactiveValue := ""
		inactiveText := ""
		if len(def.Options) == 2 {
			for k, v := range def.Options {
				text := strings.ToLower(fmt.Sprintf("%v", v))
				switch text {
				case "on":
					activeValue = k
					activeText = "ON"
				case "enable":
					activeValue = k
					activeText = "Enable"
				case "off":
					inactiveValue = k
					inactiveText = "OFF"
				case "disable":
					inactiveValue = k
					inactiveText = "Disable"
				case "factory pattern":
					activeValue = k
					activeText = ""
				case "quit":
					inactiveValue = k
					inactiveText = ""
				}
			}
		}
		if activeValue != "" && inactiveValue != "" {
			param.SetStyle("input", "switch")
			param.SetStyle("activeValue", activeValue)
			param.SetStyle("activeText", activeText)
			param.SetStyle("inactiveValue", inactiveValue)
			param.SetStyle("inactiveText", inactiveText)
		} else {
			param.SetStyle("input", "select")
		}
	case "radio":
		param.SetStyle("input", "radio")
	case "select":
		param.SetStyle("input", "select")
	case "treeSelect":
		param.SetStyle("input", "treeSelect")
	case "status:sync":
		param.SetStyle("input", "status")
		param.SetStyle("status", "sync")
	case "status:alarm":
		param.SetStyle("input", "status")
		param.SetStyle("status", "alarm")
	case "default":
		fallthrough
	default:
		param.SetStyle("input", "default")
		if def.Access == model.AccessReadOnly {
			param.Access = "ro"
			param.SetStyle("readonly", true)
		}
	}

	if param.Key == "element_role" {
		param.SetStyle("confirmTitle", "WARNING")
		param.SetStyle("confirmMessage", fmt.Sprintf("This command will change %v work mode, do you want to continue?", info.Device.DeviceTypeName))
		param.SetStyle("confirmType", "warning")
	}

	if param.Access == "wo" && len(def.Options) == 1 {
		if def.Name != "Edit User Confirm" && def.Name != "Update" {
			tmp := *param
			ids := def.Options.IDs()
			tmp.SetValue(ids[0])
			action := NewSetParameterValuesAction(&tmp)
			action.SetStyle("showMessage", true)

			param.SetValue(ids[0])
			param.SetAction("click", action)
			param.SetStyle("input", "button")
			param.SetStyle("plain", true)
			param.SetStyle("type", "primary")

			if param.Key == "delete_history_alarm" {
				param.SetStyle("confirmTitle", "WARNING")
				param.SetStyle("confirmMessage", "This command will delete the history alarms, do you want to continue?")
				param.SetStyle("confirmType", "warning")
			} else if param.Key == "hardware_reset" {
				if clickAction := param.Actions["click"]; clickAction != nil {
					param.SetAction("click", NewMultipleActionsAction(clickAction, NewAction("SetDeviceUnavailable")))
				}
				param.SetStyle("confirmTitle", "WARNING")
				param.SetStyle("confirmMessage", "This command will reboot device, do you want to continue?")
				param.SetStyle("confirmType", "warning")
				param.SetStyle("willReboot", true)
			} else if param.Key == "triggered_delay_activation" {
				param.SetStyle("enableParam", "T02.P0BE4")
				param.SetStyle("enableValue", "01")
			}
		}
	}
	if param.OID == "T02.P0006" {
		param.SetStyle("enableParam", "TB0.P0AFF")
		param.SetStyle("enableValue", "5A")
	} else if param.OID == "T02.P0166" {
		param.Access = "wo"
		param.SetStyle("input", "buttonGroup")
		param.Items = []*Element{}
		keys := []string{}
		for k, _ := range def.Options {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] > keys[j]
		})
		for _, k := range keys {
			tmp := NewParam(param.Name, param.OID, param.Access)
			tmp.SetValue((k))
			action := NewSetParameterValuesAction(tmp)
			param2 := NewParam(param.Name, param.OID, param.Access)
			param2.SetValue(k)
			param2.SetAction("click", action)
			param2.SetStyle("input", "button")
			param2.SetStyle("plain", true)
			param2.SetStyle("type", "primary")
			param.Items = append(param.Items, param2)
		}
	}

	return param
}
