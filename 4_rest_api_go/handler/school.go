package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/model"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/response"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/service"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/utils"
)

type SchoolHandler struct{ service service.SchoolService }

func NewSchool(svc service.SchoolService) *SchoolHandler {
	return &SchoolHandler{
		service: svc,
	}
}

func (h *SchoolHandler) List(c *echo.Context) error {
	filter, err := utils.ParseFilter("schools", c.QueryParams())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	res, err := h.service.List(c.Request().Context(), filter)
	if err != nil {
		return err
	}
	return response.Success(c, http.StatusOK, res)
}

func (h *SchoolHandler) Get(c *echo.Context) error {
	id, err := utils.ParseID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	res, err := h.service.Get(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return response.Success(c, http.StatusOK, res)
}

func (h *SchoolHandler) Create(c *echo.Context) error {
	var req model.SchoolRequest
	if err := utils.Bind(c, &req); err != nil {
		return err
	}
	if err := utils.Validate(req); err != nil {
		return err
	}
	res, err := h.service.Create(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return response.Success(c, http.StatusCreated, res)
}

func (h *SchoolHandler) Update(c *echo.Context) error {
	id, err := utils.ParseID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var req model.SchoolRequest
	if err = utils.Bind(c, &req); err != nil {
		return err
	}
	if err = utils.Validate(req); err != nil {
		return err
	}
	res, err := h.service.Update(c.Request().Context(), id, req)
	if err != nil {
		return err
	}
	return response.Success(c, http.StatusOK, res)
}

func (h *SchoolHandler) Delete(c *echo.Context) error {
	id, err := utils.ParseID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = h.service.Delete(c.Request().Context(), id); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
