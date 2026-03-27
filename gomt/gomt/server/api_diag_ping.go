package server

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"gomt/core/task"
	"net/http"
)

func (s *OMTServer) handleCreatePingJob(c echo.Context) error {
	opts := task.PingOptions{}
	if err := c.Bind(&opts); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	job, err := task.CreateTaskJob(opts)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	task.SetTaskJob(job)
	return c.JSON(http.StatusOK, job)
}

//func (s *OMTServer) handleGetPingJob(c echo.Context) error {
//	token := cast.ToString(c.Param("token"))
//	job := task.GetTaskJob(token)
//	if job == nil {
//		return echo.NewHTTPError(http.StatusNotFound, "Task not found")
//	}
//	return c.JSON(http.StatusOK, job)
//}

func (s *OMTServer) handlePingJobWebSocket(c echo.Context) error {
	token := cast.ToString(c.Param("token"))
	job := task.GetTaskJob(token)
	if job == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Task not found")
	}
	ws, err := s.wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer ws.Close()
	job.SetWebSocket(ws)

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			logrus.Errorf("read ws, %v", err)
			break
		}
	}
	job.Cancel()
	return nil
}

func (s *OMTServer) handleRunPingJob(c echo.Context) error {
	token := cast.ToString(c.Param("token"))
	job := task.GetTaskJob(token)
	if job == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Task not found")
	}
	if err := job.Run(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, job)
}
func (s *OMTServer) handleCancelPingJob(c echo.Context) error {
	token := cast.ToString(c.Param("token"))
	job := task.GetTaskJob(token)
	if job == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Task not found")
	}
	if err := job.Cancel(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	task.RemoveTaskJob(job.Token)
	return c.NoContent(http.StatusOK)
}
