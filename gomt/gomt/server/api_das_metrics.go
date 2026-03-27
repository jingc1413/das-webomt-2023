package server

import (
	"encoding/json"
	"gomt/das/agent"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type queryStatsRequestData struct {
	BeginTime int64    `json:"beginTime"`
	EndTime   int64    `json:"endTime"`
	Keys      []string `json:"keys"`
}

func (s *OMTServer) handleDasDeviceQueryMetricsData(c echo.Context) error {
	deviceAgent := c.Get("agent").(*agent.DasDeviceAgent)
	var reqData queryStatsRequestData
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := json.Unmarshal(body, &reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	data, err := deviceAgent.QueryStats(reqData.BeginTime, reqData.EndTime, reqData.Keys)
	if err != nil {
		s.log.Error(errors.Wrap(err, "query stats"))
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, data)
}

func (s *OMTServer) handleeDasDeviceGetMetricsCurrent(c echo.Context) error {
	deviceAgent := c.Get("agent").(*agent.DasDeviceAgent)
	metris := deviceAgent.GetMetricsData()
	if metris == nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, metris.GaugeMap)
}
