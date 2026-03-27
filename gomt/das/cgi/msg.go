package cgi

import (
	"encoding/hex"
	"fmt"
	"gomt/core/iam"
	"gomt/core/model"
	"gomt/core/proto/priv"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

type itemCache struct {
	Param   *model.ParameterDefine
	PageKey string
	OID     string
}

func getItemCache(caches []*itemCache, pageKey string, oid string) *itemCache {
	for _, cache := range caches {
		if cache.OID == oid && cache.PageKey == pageKey {
			tmp := cache
			return tmp
		}
	}
	return nil
}

func UnmarshalQueryDevicesResponseContent(schema string, defs model.ParameterDefines, caches []*itemCache, content []byte) (uint8, uint8, *model.DeviceInfo, error) {
	total := uint8(0)
	index := uint8(0)

	params, err := UmarshalIndexResponseContent(defs, caches, content)
	if err != nil {
		return total, index, nil, err
	}
	if len(params) != 1 {
		return total, index, nil, errors.Errorf("invaild response content, %v", string(content))
	}
	// 0=date=1712273936171&param=00,Settings,Engineering CMD2>0aec,01010D01020304060708090B0C0D0E0F,00
	// 00,Settings,Engineering CMD2,00>0aec,06010D01020304060708090B0C0D0E0F^0|0000|16.7.1.1|1|240|SUNWAVE-0011|0|1|255|Primary A3|A3-001|1.2.4|121131313131,00;

	if params[0].Code != "00" {
		return total, index, nil, errors.Errorf("invaild response content, code=%v", params[0].Code)
	}

	parts := strings.Split(cast.ToString(params[0].Value), "^")
	if len(parts) != 2 {
		return total, index, nil, errors.Errorf("invaild response content, %v", string(content))
	}

	headerBytes, err := hex.DecodeString(parts[0])
	if err != nil {
		return total, index, nil, err
	}
	total = uint8(headerBytes[0])
	index = uint8(headerBytes[1])

	values := strings.Split(parts[1], "|")
	if len(values) != len(headerBytes)-3 {
		return total, index, nil, errors.Errorf("miss match values number")
	}
	// logrus.Warnf("%v", parts)
	info := model.DeviceInfo{
		Schema: schema,
	}
	for i := 0; i < len(headerBytes)-3; i++ {
		qid := headerBytes[i+3]
		v := values[i]
		switch qid {
		case priv.QUERY_DEVICE_INPUT_SUB_ID:
			info.SubID = cast.ToUint8(v)
		case priv.QUERY_DEVICE_INPUT_ROUTE:
			info.RouteAddress = []byte{}
			for _, c := range v {
				tmp, err := strconv.ParseUint(fmt.Sprintf("0%c", c), 16, 8)
				if err != nil {
					return total, index, nil, errors.Wrap(err, "parse route address")
				}
				info.RouteAddress = append(info.RouteAddress, byte(tmp))
			}
			// logrus.Warnf("%v => %v", v, info.RouteAddress)
		case priv.QUERY_DEVICE_INPUT_IPADDR:
			info.IpAddress = net.ParseIP(v).To4()
		case priv.QUERY_DEVICE_INPUT_CONNECT_STATE:
			info.ConnectState = cast.ToInt8(v)
		case priv.QUERY_DEVICE_INPUT_ALARM_STATE:
			info.AlarmState = cast.ToInt8(v)
		case priv.QUERY_DEVICE_INPUT_VERSION_STATE:
			info.VersionState = cast.ToInt8(v)
		case priv.QUERY_DEVICE_INPUT_MIXED_STATE:
			info.MixedState = cast.ToInt8(v)
		case priv.QUERY_DEVICE_INPUT_OPTICAL_STATE:
			info.OpticalState = cast.ToInt8(v)
		case priv.QUERY_DEVICE_INPUT_SUB_CODE:
			// ignore
		case priv.QUERY_DEVICE_INPUT_DEVICE_TYPE_ID:
			info.DeviceTypeID = cast.ToInt(v)
		case priv.QUERY_DEVICE_INPUT_LOCATION_INFO:
			info.InstalledLocation = v
		case priv.QUERY_DEVICE_INPUT_DEVICE_TYPE_NAME:
			info.DeviceTypeName = v
		case priv.QUERY_DEVICE_INPUT_VERSION:
			info.Version = v
		case priv.QUERY_DEVICE_INPUT_DEVICE_NAME:
			info.DeviceName = v
		case priv.QUERY_DEVICE_INPUT_ELEMENT_MODEL_NUM:
			info.ElementModelNumber = v
		}
	}
	info.Setup()
	return total, index, &info, nil
}

func MarshalQueryDeviceRequestData(defs model.ParameterDefines, total uint8, index uint8, query string) (string, []*itemCache, error) {
	size := uint8(len(query) / 2)
	inputValues := []priv.Parameter{
		{
			ID:    "T02.P0AEC",
			Value: strings.ToLower(fmt.Sprintf("%02X%02X%02X%v", total, index, size, query)),
		},
	}
	cmds, cache, err := MarshalGetParameterValuesRequestData(defs, inputValues, true)
	if err != nil || len(cmds) != 1 {
		return "", cache, err
	}
	return cmds[0], cache, err
}

func UmarshalIndexResponseContent(defs model.ParameterDefines, caches []*itemCache, content []byte) ([]priv.Parameter, error) {
	data := ""
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "Udp") {
			data = data + strings.TrimSpace(line)
		}
	}
	return getParameterValuesFromPageLine(defs, caches, data)
}

func MarshalIndexRequestData(prefix string, defs model.ParameterDefines, values []priv.Parameter) ([]string, []*itemCache, error) {
	out := []string{}
	pages, caches, err := getPageLinesFromParameterValues(defs, values)
	if err != nil {
		return out, caches, err
	}
	for pageKey, items := range pages {
		size := 0
		items2 := []string{}
		for _, item := range items {
			size += len(item)
			items2 = append(items2, item)
			if size >= 512 {
				data := fmt.Sprintf("%v,%v>%v", prefix, pageKey, strings.Join(items2, ";"))
				out = append(out, data)
				size = 0
				items2 = []string{}
			}
		}
		if len(items2) > 0 {
			data := fmt.Sprintf("%v,%v>%v", prefix, pageKey, strings.Join(items2, ";"))
			out = append(out, data)
		}
	}
	return out, caches, nil
}

func MarshalGetParameterValuesRequestData(defs model.ParameterDefines, values []priv.Parameter, statusEnable bool) ([]string, []*itemCache, error) {
	if statusEnable {
		return MarshalIndexRequestData("parastatus=00", defs, values)
	}
	return MarshalIndexRequestData("param=00", defs, values)
}

func MarshalSetParameterValuesRequestData(defs model.ParameterDefines, values []priv.Parameter, statusEnable bool) ([]string, []*itemCache, error) {
	if statusEnable {
		return MarshalIndexRequestData("parastatus=01", defs, values)
	}
	return MarshalIndexRequestData("param=01", defs, values)
}

func getParameterValuesFromPageLine(defs model.ParameterDefines, caches []*itemCache, content string) ([]priv.Parameter, error) {
	//00,Settings,Engineering CMD1,00>0b31,1,00;0b32,1,00;0b33,admin,00;0b34,1,00;

	values := []priv.Parameter{}
	parts := strings.Split(content, ">")
	if len(parts) < 2 {
		return values, errors.Errorf("invalid response: %v", content)
	}

	headers := strings.Split(parts[0], ",")
	if len(headers) < 4 {
		return values, errors.Errorf("invalid response: %v", content)
	}
	pageKey := fmt.Sprintf("%v,%v", headers[1], headers[2])

	items := strings.Split(parts[1], ";")
	for _, item := range items {
		if item == "" {
			continue
		}
		itemParts := strings.Split(item, ",")
		if len(itemParts) != 3 {
			logrus.Warnf("invalid object content: %v", item)
			continue
		}
		cache := getItemCache(caches, pageKey, itemParts[0])
		if cache == nil {
			logrus.Warnf("invalid response oid: %v", item)
			continue
		}
		param := cache.Param
		v := itemParts[1]
		c := itemParts[2]

		if strings.HasPrefix(string(param.DataType), "string") {
			v = strings.ReplaceAll(v, "#a#", ",")
			v = strings.ReplaceAll(v, "#b#", ";")
			v = strings.ReplaceAll(v, "#c#", ">")
		}

		vv, err := param.UnmarshalCgiStringValue(v)
		if err != nil {
			logrus.Error(errors.Wrapf(err, "ummarshal parameter value: %v", v))
			continue
		}
		values = append(values, priv.Parameter{
			ID:    string(param.PrivOid),
			Value: vv,
			Code:  c,
		})
	}
	return values, nil
}

func getPageLinesFromParameterValues(defs model.ParameterDefines, values []priv.Parameter) (map[string][]string, []*itemCache, error) {
	pages := map[string][]string{}
	caches := []*itemCache{}

	for _, value := range values {
		privOid, err := model.PrivObjectIdFromString(value.ID)
		if err != nil {
			logrus.Warnf(errors.Wrapf(err, "parse private object id %v", value.ID).Error())
			continue
		}
		_, _, short := privOid.Values()
		oid := fmt.Sprintf("%04x", short)
		param := defs.GetParameterDefine(privOid)
		if param == nil {
			logrus.Warnf("unknow parameter: %v", value.ID)
			continue
		}
		if len(param.Groups) < 1 {
			logrus.Warnf("invalid parameter group: %v", value.ID)
			continue
		}

		args := strings.Split(param.Groups[0], ",")
		if len(args) < 2 {
			logrus.Warnf("invalid parameter group: %v", value.ID)
			continue
		}
		key := fmt.Sprintf("%v,%v", args[0], args[1])
		caches = append(caches, &itemCache{
			Param:   param,
			PageKey: key,
			OID:     oid,
		})

		page, ok := pages[key]
		if !ok {
			page = []string{}
		}
		v, err := param.MarshalCgiStringValue(value.Value)
		if err != nil {
			logrus.Error(errors.Wrapf(err, "marshal parameter value: %v", value.Value))
			continue
		}
		if strings.HasPrefix(string(param.DataType), "string") {
			v = strings.ReplaceAll(v, ",", "#a#")
			v = strings.ReplaceAll(v, ";", "#b#")
			v = strings.ReplaceAll(v, ">", "#c#")
		}

		page = append(page, fmt.Sprintf("%v,%v,00", oid, v))
		pages[key] = page
	}

	return pages, caches, nil
}

func parseUserList(data []byte) []*iam.User {
	result := []*iam.User{}
	if parts := strings.Split(string(data), ";"); len(parts) > 0 {
		for _, part := range parts {
			name := ""
			if parts2 := strings.Split(part, ","); len(parts2) > 0 {
				for _, part2 := range parts2 {
					if parts3 := strings.Split(part2, "="); len(parts3) == 2 {
						k := strings.TrimSpace(parts3[0])
						v := strings.TrimSpace(parts3[1])
						switch k {
						case "username":
							name = v
						case "permission":
							//ignore
						}
					}
				}
			}
			if name != "" {
				user := iam.MakeUser(name, "", nil)
				result = append(result, user)
			}
		}
	}
	// sort.SliceStable(result, func(i, j int) bool {
	// 	return true
	// })
	return result
}

type UpgradeResponseData struct {
	Status                string
	Message               string
	Code                  int
	CRC                   int
	Debug                 string
	Text                  string
	Text2                 string
	NeedReboot            bool
	NeedUpdateCRC         bool
	NeedUpdateHostPackage bool
	NeedUpdateSubPackage  bool
}

func parseUpgradeResult(data []byte) (UpgradeResponseData, error) {
	resp := UpgradeResponseData{
		Status:  "error",
		Message: "",

		Code:  -1,
		CRC:   -1,
		Debug: "",
		Text:  "",
		Text2: "",
	}

	re := regexp.MustCompile(`<!(.*)>>(.*)-->`)
	re2 := regexp.MustCompile(`\(CRC:(\S*)\)`)

	allSubs := re.FindAllStringSubmatch(string(data), -1)

	if subs := re2.FindStringSubmatch(string(data)); len(subs) == 2 {
		resp.CRC = cast.ToInt(subs[1])
	}
	for _, subs := range allSubs {
		if len(subs) != 3 {
			return resp, errors.New("parse response content failed")
		}
		key := subs[1]
		switch key {
		case "debug":
			resp.Debug = subs[2]
		case "rstmain":
			resp.Text = subs[2]
		case "rsttext":
			resp.Text2 = subs[2]
		case "resultcode":
			resp.Code = cast.ToInt(subs[2])
			// case "img":
			// 	img = subs[2]
		}
	}

	const (
		UPGRADE_FAILED             = -1
		UPGRADE_SUCCESS            = 1
		UPGRADE_SUCCESS_DAS        = 2
		UPGRADE_SUCCESS_DAS_REBOOT = 3
		UPGRADE_SUCCESS_WEBOMT     = 4
		UPGRADE_SUCCESS_FPGA       = 5
		UPGRADE_SUCCESS_ARM        = 6
		UPGRADE_SUCCESS_PA         = 7
		UPGRADE_SUCCESS_DAS_CRC    = 8
		UPGRADE_NORMAL_CANT        = 9
		UPGRADE_SUCCESS_AUTO       = 10
		UPGRADE_SUCCESS_SNMP       = 11
		UPGRADE_SUCCESS_LINUX      = 12
	)

	switch resp.Code {
	case UPGRADE_NORMAL_CANT:
		resp.Message = "Upgrade not allow."
	case UPGRADE_FAILED:
		resp.Message = "Upgrade failed, the device will restart, please wait."
		resp.NeedReboot = true
	case UPGRADE_SUCCESS:
		resp.Status = "success"
		resp.Message = "Upgrade succesfully, the device will restart, please wait."
		resp.NeedReboot = true
	case UPGRADE_SUCCESS_DAS: // au upgrade slave
		resp.Status = "success"
		resp.Message = "Upgrade sub device succesfully, the device will restart, please wait."
		resp.NeedUpdateSubPackage = true // need update sub package 0b21
	case UPGRADE_SUCCESS_DAS_CRC:
		resp.Status = "success"
		resp.Message = "Upgrade succesfully, the device will restart, please wait."
		resp.NeedUpdateCRC = true
		resp.NeedReboot = true
	case UPGRADE_SUCCESS_DAS_REBOOT: // au upgrade au
		resp.Status = "success"
		resp.Message = "Upgrade local device succesfully, the device will restart, please wait."
		resp.NeedReboot = true
		resp.NeedUpdateHostPackage = true // need update host package 0b19
	case UPGRADE_SUCCESS_LINUX:
		resp.Status = "success"
		resp.Message = "Upgrade linux succesfully, the device will restart, please wait."
		resp.NeedReboot = true
		resp.NeedUpdateHostPackage = true // need update host package 0b19
	case UPGRADE_SUCCESS_WEBOMT:
		resp.Status = "success"
		resp.Message = "Upgrade omt succesfully, the device will restart, please wait."
		resp.NeedReboot = true
	case UPGRADE_SUCCESS_FPGA:
		resp.Status = "success"
		resp.Message = "Upgrade fpga succesfully, the device will restart, please wait."
		resp.NeedReboot = true
	case UPGRADE_SUCCESS_ARM:
		resp.Status = "success"
		resp.Message = "Upgrade arm succesfully, the device will restart, please wait."
		resp.NeedReboot = true
	case UPGRADE_SUCCESS_SNMP:
		resp.Status = "success"
		resp.Message = "Upgrade snmp succesfully, the device will restart, please wait."
		resp.NeedReboot = true
	case UPGRADE_SUCCESS_PA:
		resp.Status = "success"
		resp.Message = "Upgrade 485 sub module succesfully, the device will restart, please wait."
		resp.NeedReboot = true
	case UPGRADE_SUCCESS_AUTO:
		resp.Status = "success"
		resp.Message = "Upgrade succesfully, the device will restart, please wait."
		resp.NeedReboot = true
	}
	return resp, nil
}
