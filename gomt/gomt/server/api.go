package server

import (
	"gomt/core"
	"gomt/das/agent"
	"net/http"

	"github.com/labstack/echo/v4"
)

type appInfoData struct {
	AppName        string
	AppVersion     string
	AppBuild       string
	Schema         string
	DeviceTypeName string
}

func (s *OMTServer) handleGetAppInfo(c echo.Context) error {
	data := appInfoData{
		AppName:    core.NAME,
		AppVersion: core.VERSION,
		AppBuild:   core.BUILD,
		Schema:     s.opts.Schema,
	}
	if localAgent := s.getDasDeviceAgent("local"); localAgent != nil {
		data.DeviceTypeName = localAgent.GetDeviceTypeName()
	}
	return c.JSON(http.StatusOK, data)
}

func (s *OMTServer) getDasDeviceAgentFromContext(c echo.Context) *agent.DasDeviceAgent {
	deviceSub := c.Param("device_sub")
	agent := s.dassys.GetAgent(deviceSub)
	if agent == nil {
		s.log.Warnf("cant find device agent for sub=%v", deviceSub)
	}
	return agent
}

func (s *OMTServer) getDasDeviceAgent(deviceSub string) *agent.DasDeviceAgent {
	return s.dassys.GetAgent(deviceSub)
}

func (s *OMTServer) agentAvailableMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		agent := s.getDasDeviceAgentFromContext(c)
		if agent == nil {
			return echo.NewHTTPError(http.StatusNotFound, "invalid device")
		}
		if !agent.IsServiceAvailable(false) {
			return echo.NewHTTPError(http.StatusServiceUnavailable, "device not available")
		}
		c.Set("agent", agent)
		return next(c)
	}
}
