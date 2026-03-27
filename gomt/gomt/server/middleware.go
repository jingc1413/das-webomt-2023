package server

import (
	"bytes"
	"io"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

func (s *OMTServer) loggerMiddlewareFunc() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogRoutePath: true,
		LogRemoteIP:  true,
		LogStatus:    true,
		LogMethod:    true,
		LogError:     true,
		BeforeNextFunc: func(c echo.Context) {
			if req := c.Request(); req != nil && req.ContentLength > 0 {
				if body := req.Body; body != nil {
					defer body.Close()
					all, _ := io.ReadAll(body)
					c.Set("requestContent", string(all))
					r := bytes.NewReader(all)
					req.Body = io.NopCloser(r)
				}
			}
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Status >= 200 && v.Status < 300 {
				s.log.WithFields(logrus.Fields{
					"method": v.Method,
					"path":   v.URI,
					"status": v.Status,
				}).Tracef("request")
			} else {
				s.log.WithFields(logrus.Fields{
					"method": v.Method,
					"path":   v.URI,
					"status": v.Status,
				}).Error(errors.Wrap(v.Error, "request"))
			}
			uri := v.RoutePath
			if strings.HasPrefix(uri, "/api") {
				uri = uri[4:]
			}
			if id, rule := getApiRule(v.Method, uri); rule != nil {
				if strings.HasSuffix(id, ".login") ||
					strings.HasSuffix(id, ".logout") ||
					strings.HasSuffix(id, ".get") ||
					strings.HasSuffix(id, ".list") ||
					strings.HasSuffix(id, ".read") {
					// ignore
				} else {
					args := strings.Split(id, ".")
					event := strings.Join(args[1:], ".")
					name := cast.ToString(c.Get("login"))
					content := c.Get("requestContent")
					if id == "api.current.change-password" {
						content = name
					}
					s.auditLogger.WriteApiLog(v.StartTime, "info", event, name, v.RemoteIP, v.Status, content)
				}

			}
			return nil
		},
	})
}
