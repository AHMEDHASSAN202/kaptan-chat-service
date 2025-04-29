package echomiddleware

import (
	"github.com/labstack/echo/v4"
	"os"
)

const HEADER_NAME = "STORAGE-URL"
const HEADER_NAME_ENV = "STORAGE_URL"

func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		browseUrl := os.Getenv(HEADER_NAME_ENV)
		c.Response().Header().Set(HEADER_NAME, browseUrl)
		return next(c)
	}
}
