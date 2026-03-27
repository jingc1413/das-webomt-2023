package arm

import (
	"bufio"
	"encoding/hex"
	"gomt/core/model"
	"gomt/core/utils"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

func IsSupportArmIni(configPath string) bool {
	filename := filepath.Join(configPath, "Arm.ini")
	return utils.ExistsFile(filename)
}

func ReadDeviceTypeNameFromXmlFile(data []byte) (string, error) {
	// repeatertype="Primary A3"
	re := regexp.MustCompile(`repeatertype\s*=\s*"([0-9a-zA-Z_\- ]+)"`)
	if subs := re.FindStringSubmatch(string(data)); len(subs) == 2 {
		productModel := model.GetProductModelByDeviceTypeName(subs[1])
		if productModel == nil {
			return "", errors.Errorf("get product model %v", subs[1])
		}
		return subs[1], nil
	}
	return "", errors.New("parse type file failed")
}

func ReadDeviceTypeNameFromTypeFile(role uint8, data []byte) (string, error) {
	//machine=iDAS_A402
	re := regexp.MustCompile(`machine=(\S+)`)
	if subs := re.FindStringSubmatch(string(data)); len(subs) == 2 {
		productModel := model.GetProductModelByModelID(subs[1])
		if productModel == nil {
			return "", errors.Errorf("get product model %v", productModel)
		}
		parts := strings.Split(productModel.DeviceTypeName, "/")
		if len(parts) == 1 {
			return parts[0], nil
		}
		if role == 2 && len(parts) > 1 {
			return parts[1], nil
		}
		return parts[0], nil
	}
	return "", errors.New("parse type file failed")
}

func ReadElementRoleFromArmIni(configPath string) (uint8, error) {
	// N0AE7=1
	// 1 = AU
	// 2 = SAU
	filename := filepath.Join(configPath, "Arm.ini")
	f, err := os.Open(filename)
	if err != nil {
		return 0, errors.Wrap(err, "read arm ini file")
	}
	defer f.Close()

	scaner := bufio.NewScanner(f)
	scaner.Split(bufio.ScanLines)

	re := regexp.MustCompile(`^N0AE7=(.*)$`)
	for scaner.Scan() {
		line := scaner.Text()
		if subs := re.FindStringSubmatch(line); len(subs) == 2 {
			return cast.ToUint8(subs[1]), nil
		}
	}
	return 0, nil
}

func ReadDeviceInfosFromArmIni(configPath string) ([]*model.DeviceInfo, error) {
	lines := map[string]string{}
	filename := filepath.Join(configPath, "Arm.ini")
	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "read arm ini file")
	}
	defer f.Close()

	scaner := bufio.NewScanner(f)
	scaner.Split(bufio.ScanLines)

	re := regexp.MustCompile(`^Slv([0-9a-fA-F]+)=(.*)$`)
	for scaner.Scan() {
		line := scaner.Text()
		if subs := re.FindStringSubmatch(line); len(subs) == 3 {
			k := subs[1]
			v := subs[2]

			line := lines[k]
			line = line + v
			lines[k] = line
		}
	}
	infos := []*model.DeviceInfo{}
	for _, line := range lines {
		//01020000000000010000000B070A014E332D52552F2D33300000000000000000000000000000000000000000000000 123N3333N

		size := len(line)
		serialNubmer := ""
		if size < 47*2 {
			logrus.Errorf("parse line %v", line)
			continue
		}
		if size > 47*2 {
			serialNubmer = line[94:]
		}

		buf, err := hex.DecodeString(line[0:94])
		if err != nil {
			logrus.Error(errors.Wrapf(err, "parse line %v", line))
			continue
		}

		//06 03 0000 0000 00 02000000 0B071401 45332D4F2F540000 0000000000000000 0000000000000000 0000000000000000

		connectState := int8(buf[0])
		subID := uint8(buf[1])

		// deviceTypeID := cast.ToInt8(buf[6])
		routeAddress := buf[7:11]
		ipAddress := buf[11:15]
		deviceTypeString := strings.TrimSpace(string(buf[15:47]))
		deviceTypeParts := strings.Split(deviceTypeString, "/")

		info := &model.DeviceInfo{}
		info.ConnectState = connectState
		info.SubID = subID

		info.RouteAddress = routeAddress
		info.IpAddress = ipAddress
		info.DeviceTypeName = deviceTypeParts[0]
		info.ElementSerialNumber = serialNubmer

		info.Setup()
		// logrus.Warnf("%v: %v, %v", k, info.SubID, info.DeviceTypeName)
		infos = append(infos, info)
	}

	return infos, nil
}
