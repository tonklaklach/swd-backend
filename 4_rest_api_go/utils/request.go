package utils

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

func ParseID(c *echo.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("id must be positive")
	}
	return id, nil
}
func Bind(c *echo.Context, dst any) error {
	r := c.Request()
	media, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || media != "application/json" {
		return echo.NewHTTPError(415, "Content-Type must be application/json")
	}
	r.Body = http.MaxBytesReader(c.Response(), r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(dst); err != nil {
		return echo.NewHTTPError(400, "invalid JSON body")
	}
	if err = decoder.Decode(new(any)); err != io.EOF {
		return echo.NewHTTPError(400, "body must contain one JSON object")
	}
	return nil
}
