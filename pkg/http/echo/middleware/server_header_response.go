package echomiddleware

import (
	"github.com/labstack/echo/v4"
	"os"
)

func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		HEADER_NAME := "STORAGE-URL"
		HEADER_NAME_ENV := "STORAGE_URL"
		browseUrl := os.Getenv(HEADER_NAME_ENV)
		c.Response().Header().Set(HEADER_NAME, browseUrl)
		return next(c)
	}
}
