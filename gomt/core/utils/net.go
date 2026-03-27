package utils

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func GetIpByInterfaceName(name string) (net.IP, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		return ip, nil
	}
	return nil, nil
}
func MustGetAddressByInterfaceName(name string) string {
	ip, err := GetIpByInterfaceName(name)
	if err != nil {
		logrus.Error(errors.Wrap(err, "get address by interface name"))
	}
	if ip == nil {
		return ""
	}
	return ip.String()
}

func IsLocalAddress(v string) bool {
	if v == "127.0.0.1" || v == "localhost" {
		return true
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return false
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, addr := range addrs {
			if addr.String() == v {
				return true
			}
		}
	}
	return false
}
