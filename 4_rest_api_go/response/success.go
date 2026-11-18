package response

import (
	"github.com/labstack/echo/v5"
	responsestruct "github.com/tonklaklach/swd-backend/4_rest_api_go/response/struct"
)

func Success[T any](c *echo.Context, status int, data T) error {
	return c.JSON(status, responsestruct.Success[T]{Data: data})
}
