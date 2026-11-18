package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/model"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/response"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/service"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/utils"
)

type TeacherHandler struct{ service service.TeacherService }

func NewTeacher(svc service.TeacherService) *TeacherHandler {
	return &TeacherHandler{
		service: svc,
	}
}

func (h *TeacherHandler) List(c *echo.Context) error {
	filter, err := utils.ParseFilter("teachers", c.QueryParams())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	res, err := h.service.List(c.Request().Context(), filter)
	if err != nil {
		return err
	}
	return response.Success(c, http.StatusOK, res)
}

func (h *TeacherHandler) Get(c *echo.Context) error {
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

func (h *TeacherHandler) Create(c *echo.Context) error {
	var req model.TeacherRequest
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

func (h *TeacherHandler) Update(c *echo.Context) error {
	id, err := utils.ParseID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var req model.TeacherRequest
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

func (h *TeacherHandler) Delete(c *echo.Context) error {
	id, err := utils.ParseID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = h.service.Delete(c.Request().Context(), id); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
