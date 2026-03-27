package server

import (
	"gomt/das/agent"
	"mime"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *OMTServer) handleDasDeviceListFirmwares(c echo.Context) error {
	deviceAgent := c.Get("agent").(*agent.DasDeviceAgent)
	data, err := deviceAgent.GetFirmwareInfos()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, data)
}

func (s *OMTServer) handleDasDeviceDeleteFirmware(c echo.Context) error {
	deviceAgent := c.Get("agent").(*agent.DasDeviceAgent)
	name := c.Param("name")

	if err := deviceAgent.DeleteFirmwareInfo(name); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func (s *OMTServer) handleDasDeviceGetTypeInfo(c echo.Context) error {
	deviceAgent := c.Get("agent").(*agent.DasDeviceAgent)
	_, f, err := deviceAgent.GetFile("Version", "type.txt")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer f.Close()
	// contentType := fmt.Sprintf("attachment; filename=%v", fname)
	contentType := mime.TypeByExtension(".txt")
	return c.Stream(http.StatusOK, contentType, f)
}

func (s *OMTServer) handleDasDeviceGetVersionPacketInfo(c echo.Context) error {
	deviceAgent := c.Get("agent").(*agent.DasDeviceAgent)
	_, f, err := deviceAgent.GetFile("Version", "fileList.txt")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer f.Close()
	// contentType := fmt.Sprintf("attachment; filename=%v", fname)
	contentType := mime.TypeByExtension(".txt")
	return c.Stream(http.StatusOK, contentType, f)
}

func (s *OMTServer) handleDasDeviceDeleteKeyAndLogs(c echo.Context) error {
	deviceAgent := c.Get("agent").(*agent.DasDeviceAgent)
	if err := deviceAgent.ServeCgiDeleteKeyAndLogs(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
