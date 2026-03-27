package server

import (
	"gomt/das/agent"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *OMTServer) handleDasDeviceListAlarmEventLogs(c echo.Context) error {
	deviceAgent := c.Get("agent").(*agent.DasDeviceAgent)
	logs, err := deviceAgent.GetAlarmLogs()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, logs)

}
