package echomiddleware

import (
	"context"
	"github.com/labstack/echo/v4"
)

const ACCEPT_LANGUAGE = "Accept-Language"

func AppendLangMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		req := c.Request()

		lang := req.Header.Get(ACCEPT_LANGUAGE)
		if lang == "" {
			lang = "en"
		}

		c.Response().Header().Set(ACCEPT_LANGUAGE, lang)
		newReq := req.WithContext(context.WithValue(req.Context(), "lang", lang))
		c.SetRequest(newReq)

		return next(c)
	}
}
