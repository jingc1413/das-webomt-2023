package server

import (
	"gomt/das/agent"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func (s *OMTServer) handleDasDeviceListFiles(c echo.Context) error {
	deviceAgent := c.Get("agent").(*agent.DasDeviceAgent)
	ftypeName := c.Param("ftype")
	def := s.opts.FileTypes.Get(ftypeName)
	if def == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid file type")
	}
	result, err := deviceAgent.ListFiles(ftypeName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

func (s *OMTServer) handleDasDeviceGetFile(c echo.Context) error {
	deviceAgent := c.Get("agent").(*agent.DasDeviceAgent)
	fname := c.Param("fname")
	ftypeName := c.Param("ftype")
	def := s.opts.FileTypes.Get(ftypeName)
	if def == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid file type")
	}
	_, f, err := deviceAgent.GetFile(ftypeName, fname)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer f.Close()
	// contentType := fmt.Sprintf("attachment; filename=%v", fname)
	contentType := mime.TypeByExtension(filepath.Ext(fname))
	return c.Stream(http.StatusOK, contentType, f)
}

func (s *OMTServer) handleDasDeviceDeleteFile(c echo.Context) error {
	deviceAgent := c.Get("agent").(*agent.DasDeviceAgent)
	fname := c.Param("fname")
	ftypeName := c.Param("ftype")
	def := s.opts.FileTypes.Get(ftypeName)
	if def == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid file type")
	}
	if err := deviceAgent.RemoveFile(ftypeName, fname); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func (s *OMTServer) handleDasDeviceUploadFile(c echo.Context) error {
	deviceAgent := c.Get("agent").(*agent.DasDeviceAgent)
	ftypeName := c.Param("ftype")
	def := s.opts.FileTypes.Get(ftypeName)
	if def == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid file type")
	}

	formFile, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	fileSize := formFile.Size
	if fileSize == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "File is empty")
	}
	if s.opts.MaxFileUploadSize > 0 && fileSize > s.opts.MaxFileUploadSize {
		return echo.NewHTTPError(http.StatusBadRequest, "file is too large")
	}

	fr, err := formFile.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to read file")
	}
	defer fr.Close()

	if err := deviceAgent.SaveFile(ftypeName, formFile.Filename, fr); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "file uploaded successfully")
}

func (s *OMTServer) handleDasDeviceGetUpgradeFilePacketInfo(c echo.Context) error {
	deviceAgent := c.Get("agent").(*agent.DasDeviceAgent)
	fname := c.Param("fname")
	ftypeName := c.Param("ftype")
	def := s.opts.FileTypes.Get(ftypeName)
	if def == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid file type")
	}
	if def.Name != "UpgradeFile" {
		return c.NoContent(http.StatusNotFound)
	}

	if err := deviceAgent.ServeCgiGetUpgradeFilePacketInfo(def.Dir, fname); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	_, f, err := deviceAgent.GetFile("PacketUpdateFile", "fileList.txt")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer f.Close()
	// contentType := fmt.Sprintf("attachment; filename=%v", fname)
	contentType := mime.TypeByExtension(".txt")
	return c.Stream(http.StatusOK, contentType, f)
}
