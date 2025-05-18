package user

import (
	"context"
	"github.com/labstack/echo/v4"
	"kaptan/pkg/localization"
	"kaptan/pkg/utils"
	"kaptan/pkg/validators"
	"net/http"
)

func (m Middlewares) AuthenticationMiddleware(causerType string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			if c.Request().Method == http.MethodOptions {
				return next(c)
			}

			userToken := c.Request().Header.Get("Authorization")
			if len(userToken) == 0 {
				m.logger.Info("AuthMiddleware -> UserToken Not found")
				return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
			}

			isSanctumToken, tokenParts := utils.IsSanctumToken(userToken)
			if !isSanctumToken {
				m.logger.Info("AuthMiddleware -> Token Parsing Failed -> Token Is Not Sanctum")
				return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
			}

			causerId := tokenParts[0]
			c.Request().Header.Add("causer-id", causerId)
			c.Request().Header.Add("causer-type", causerType)

			ctx = context.WithValue(ctx, "causer-id", causerId)
			ctx = context.WithValue(ctx, "causer-type", causerType)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
