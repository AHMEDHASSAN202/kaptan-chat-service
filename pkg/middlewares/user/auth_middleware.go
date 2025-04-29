package user

import (
	"context"
	"encoding/json"
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

func (m Middlewares) AuthorizationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		causerId := c.Request().Header.Get("causer-id")
		causerType := c.Request().Header.Get("causer-type")
		causerToken := c.Request().Header.Get("Authorization")
		if causerId == "" || causerType == "" || causerToken == "" {
			m.logger.Info("AuthMiddleware -> Empty Causer ")
			return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
		}

		user, err := m.driverRepository.FindByToken(&ctx, utils.ExtractToken(causerToken))
		if err != nil {
			m.logger.Info("AuthMiddleware -> FindByToken Error -> ", err)
			return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
		}

		if !user.IsActive {
			m.logger.Info("AuthMiddleware -> User Is Not Active")
			return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
		}

		jsonByte, err := json.Marshal(user)
		if err != nil {
			m.logger.Info("AuthMiddleware -> Marshal Error -> ", err)
			return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
		}

		c.Request().Header.Add("causer-details", string(jsonByte))

		ctx = context.WithValue(ctx, "causer-id", causerId)
		ctx = context.WithValue(ctx, "causer-type", causerType)
		ctx = context.WithValue(ctx, "causer-details", user)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
