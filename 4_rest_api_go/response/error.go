package response

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/domain"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/errorcode"
	responsestruct "github.com/tonklaklach/swd-backend/4_rest_api_go/response/struct"
)

func ErrorHandler(c *echo.Context, err error) {
	if response, _ := echo.UnwrapResponse(c.Response()); response != nil && response.Committed {
		return
	}
	status, message, code := http.StatusInternalServerError, "internal server error", errorcode.InternalError
	var invalid *domain.ValidationError
	var httpErr *echo.HTTPError
	var statusErr echo.HTTPStatusCoder
	switch {
	case errors.As(err, &httpErr):
		status, message = httpErr.StatusCode(), httpErr.Message
		code = codeForStatus(status)
	case errors.As(err, &statusErr):
		status = statusErr.StatusCode()
		message = http.StatusText(status)
		code = codeForStatus(status)
	case errors.As(err, &invalid):
		status, message, code = http.StatusUnprocessableEntity, invalid.Message, errorcode.ValidationFailed
	case errors.Is(err, domain.ErrNotFound):
		status, message, code = http.StatusNotFound, err.Error(), errorcode.NotFound
	case errors.Is(err, domain.ErrConflict):
		status, message, code = http.StatusConflict, err.Error(), errorcode.Conflict
	case errors.Is(err, domain.ErrInvalidReference):
		status, message, code = http.StatusUnprocessableEntity, err.Error(), errorcode.InvalidReference
	default:
		slog.Error("request failed", "error", err)
	}
	if err = c.JSON(status, responsestruct.Error{Code: code, Error: message}); err != nil {
		slog.Error("write error response failed", "error", err)
	}
}

func codeForStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return errorcode.BadRequest
	case http.StatusNotFound:
		return errorcode.NotFound
	case http.StatusUnsupportedMediaType:
		return errorcode.UnsupportedMediaType
	case http.StatusServiceUnavailable:
		return errorcode.Unavailable
	default:
		return errorcode.HTTPError
	}
}
