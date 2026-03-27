package server

import (
	"gomt/das/system"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type deviceProxyBalancer struct {
	ds *system.DasSystem
}

func NewDeviceProxyBalancer(ds *system.DasSystem) *deviceProxyBalancer {
	b := deviceProxyBalancer{
		ds: ds,
	}
	return &b
}

func (b *deviceProxyBalancer) AddTarget(target *middleware.ProxyTarget) bool {
	return false
}

func (b *deviceProxyBalancer) RemoveTarget(name string) bool {
	return false
}

func (b *deviceProxyBalancer) Next(c echo.Context) *middleware.ProxyTarget {
	deviceSub := c.Param("device_sub")
	deviceAgent := b.ds.GetAgent(deviceSub)
	if deviceAgent == nil {
		logrus.Errorf("invalid device sub, sub=%v", deviceSub)
		return nil
	}
	urlString := deviceAgent.MakeHttpUrl("")
	u, err := url.Parse(urlString)
	if err != nil {
		logrus.Error(errors.Wrap(err, "invalid device address"))
		return nil
	}
	if req := c.Request(); req != nil && req.URL != nil {
		logrus.Debugf("proxy: %v, path=%v", u.String(), req.URL.Path)
	} else {
		logrus.Debugf("proxy: %v", u.String())
	}
	return &middleware.ProxyTarget{
		URL: u,
	}
}
