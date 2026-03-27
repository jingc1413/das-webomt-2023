package server

import (
	"encoding/json"
	"gomt/das/agent"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type StartUpgradeRequestData struct {
	Filename string
	Force    bool
	ByArm    bool
}

func (s *OMTServer) handleDasDeviceStartUpgrade(c echo.Context) error {
	deviceAgent := c.Get("agent").(*agent.DasDeviceAgent)

	var reqData StartUpgradeRequestData
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := json.Unmarshal(body, &reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	result, err := deviceAgent.StartUpgrade(reqData.Filename, reqData.Force, reqData.ByArm)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

func (s *OMTServer) handleDasDeviceSetUpgradeReboot(c echo.Context) error {
	deviceAgent := c.Get("agent").(*agent.DasDeviceAgent)

	ok, err := deviceAgent.ServeCgiSetUpradeReboot()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "set upgrade reboot failed")
	}
	return c.NoContent(http.StatusOK)
}
