package server

import (
	"gomt/das/agent"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func (s *OMTServer) handleGetDasQueryDevicesProgress(c echo.Context) error {
	s.queryDevicesLock.RLock()
	data := s.queryDevicesProgress
	s.queryDevicesLock.RUnlock()
	return c.JSON(http.StatusOK, data)
}

func (s *OMTServer) handleGetDasDeviceInfos(c echo.Context) error {
	queryParams := c.QueryParams()
	_, force := queryParams["force"]
	go s.updateDasDevices(force)

	data := s.dassys.GetAllDeviceInfos()
	return c.JSON(http.StatusOK, data)
}

func (s *OMTServer) handleGetDasDeviceInfo(c echo.Context) error {
	deviceAgent := c.Get("agent").(*agent.DasDeviceAgent)
	if err := deviceAgent.UpdateDeviceInfo(); err != nil {
		s.log.Error(errors.Wrapf(err, "update device info sub=%v", deviceAgent.GetDeviceInfo().SubID))
	}
	info := deviceAgent.GetDeviceInfo()
	return c.JSON(http.StatusOK, info)
}
