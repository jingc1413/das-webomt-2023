package server

import (
	"net/http/cgi"
	"path/filepath"
	"sync"

	"github.com/labstack/echo/v4"
)

var cgiLock sync.Mutex

func (s *OMTServer) makeLocalCgiHandlerFunc(cgiBinPath string) echo.HandlerFunc {
	return func(c echo.Context) error {
		cgiLock.Lock()
		defer cgiLock.Unlock()

		req := c.Request()

		handler := new(cgi.Handler)
		handler.Path = filepath.FromSlash(req.URL.Path[len("/cgi-bin/"):])
		handler.Dir = cgiBinPath
		rw := c.Response().Writer
		handler.ServeHTTP(rw, req)
		return nil
	}
}
