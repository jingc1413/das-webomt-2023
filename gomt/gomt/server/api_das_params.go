package server

import (
	"encoding/json"
	"gomt/core/proto/priv"
	"gomt/das/agent"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func (s *OMTServer) handleDasDeviceCheckAvailable(c echo.Context) error {
	deviceAgent := c.Get("agent").(*agent.DasDeviceAgent)
	if deviceAgent.IsServiceAvailable(true) {
		return c.NoContent(http.StatusOK)
	}
	return c.NoContent(http.StatusServiceUnavailable)
}

func (s *OMTServer) handleDasDeviceGetParameterValues(c echo.Context) error {
	deviceAgent := c.Get("agent").(*agent.DasDeviceAgent)
	var inputValues []priv.Parameter
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := json.Unmarshal(body, &inputValues); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	values, err := deviceAgent.GetParameterValues(inputValues)
	if err != nil {
		s.log.Error(errors.Wrap(err, "query parameter values"))
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, values)
}

func (s *OMTServer) handleDasDeviceSetParameterValues(c echo.Context) error {
	deviceAgent := c.Get("agent").(*agent.DasDeviceAgent)

	var inputValues []priv.Parameter
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := json.Unmarshal(body, &inputValues); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	values, err := deviceAgent.SetParameterValues(inputValues)
	if err != nil {
		s.log.Error(errors.Wrap(err, "set parameter values"))
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, values)
}

type RegisterData struct {
	Module uint8
	Offset uint16
	Size   uint8
	Buffer string
}

func (s *OMTServer) handleDasDeviceReadRegister(c echo.Context) error {
	deviceAgent := c.Get("agent").(*agent.DasDeviceAgent)

	var reqData RegisterData
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := json.Unmarshal(body, &reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if reqData.Size%4 != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid buffer size")
	}

	buffer := ""
	size := uint8(0)
	num := int(reqData.Size / 4)
	for i := 0; i < num; i++ {
		tmp, err := deviceAgent.ReadRegister(reqData.Module, reqData.Offset+uint16(i*4))
		if err != nil {
			s.log.Error(errors.Wrap(err, "read register"))
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		size += 4
		buffer += tmp
	}

	reqData.Size = size
	reqData.Buffer = buffer
	return c.JSON(http.StatusOK, reqData)
}

func (s *OMTServer) handleDasDeviceWriteRegister(c echo.Context) error {
	deviceAgent := c.Get("agent").(*agent.DasDeviceAgent)

	var reqData RegisterData
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := json.Unmarshal(body, &reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if reqData.Size%4 != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid buffer size")
	}
	if len(reqData.Buffer) != int(reqData.Size)*2 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid buffer")
	}

	buf := ""
	num := int(reqData.Size / 4)
	for i := 0; i < num; i++ {
		tmp, err := deviceAgent.WriteRegister(reqData.Module, reqData.Offset, reqData.Buffer[i*8:i*8+8])
		if err != nil {
			s.log.Error(errors.Wrap(err, "write register"))
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		buf += tmp
	}
	reqData.Buffer = buf
	return c.JSON(http.StatusOK, reqData)
}
