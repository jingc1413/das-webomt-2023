package server

import (
	"gomt/core/layout"
	"gomt/core/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func (s *OMTServer) handleListDasProductModels(c echo.Context) error {
	data := model.GetDefaultProductModels()
	return c.JSON(http.StatusOK, data)
}

func (s *OMTServer) handleGetDasProductModel(c echo.Context) error {
	deviceTypeName := c.Param("name")
	data := model.GetProductModelByDeviceTypeName(deviceTypeName)
	if data == nil {
		s.log.Error(errors.Errorf("get product model %v", deviceTypeName))
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, data)
}

func (s *OMTServer) handleGetDasDeviceTypes(c echo.Context) error {
	all := model.GetAllParameterDefinesMap()
	data := make([]string, 0, len(all))
	for key := range all {
		data = append(data, key)
	}
	return c.JSON(http.StatusOK, data)
}

func (s *OMTServer) handleGetDasModelLayout(c echo.Context) error {
	deviceTypeName := c.Param("name")
	version := c.Param("version")
	var data *layout.Element

	if s.opts.Schema == "default" {
		data = layout.GetLayoutByDeviceTypeName(deviceTypeName, version)
	} else if s.opts.Schema == "corning" {
		data = layout.GetLayoutByDeviceTypeName(deviceTypeName, version)

	}
	if data == nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, data)
}

func (s *OMTServer) handleGetDasModelParamters(c echo.Context) error {
	deviceTypeName := c.Param("name")
	version := c.Param("version")

	data := model.GetParameterDefinesByDeviceTypeName(s.opts.Schema, deviceTypeName, version)
	if data == nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, data)
}
