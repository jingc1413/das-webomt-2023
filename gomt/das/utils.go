package das

import (
	"fmt"
	"gomt/core/utils"
	"time"
)

func IsSupportCGI(baseURL string) bool {
	u := baseURL + "/cgi-bin/index.cgi"
	resp, err := utils.HttpGetWithTimeout(u, time.Second)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true
}

func IsSupportAPI(baseURL string) bool {
	u := baseURL + "/api/info"
	resp, err := utils.HttpGetWithTimeout(u, time.Second)
	if err != nil {
		// logrus.Error(err)
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func MakeDeviceUrlBase(deviceAddr string, port uint16, tlsEnable bool) string {
	base := ""
	proto := "http"
	if tlsEnable {
		proto = "https"
	}
	if port <= 0 {
		base = fmt.Sprintf("%v://%s", proto, deviceAddr)
	} else {
		base = fmt.Sprintf("%v://%s:%v", proto, deviceAddr, port)
	}
	return base
}
