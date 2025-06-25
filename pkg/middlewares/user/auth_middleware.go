package user

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"kaptan/pkg/localization"
	"kaptan/pkg/utils"
	"kaptan/pkg/validators"
	"net/http"
	"strings"
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

			userToken = strings.Replace(userToken, "Bearer ", "", 1)
			userToken = strings.Replace(userToken, "bearer ", "", 1)
			isSanctumToken, tokenParts := utils.IsSanctumToken(userToken)
			if !isSanctumToken {
				m.logger.Info("AuthMiddleware -> Token Parsing Failed -> Token Is Not Sanctum")
				return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
			}

			accessTokenId := tokenParts[0]
			if accessTokenId == "" {
				m.logger.Info("AuthMiddleware -> Token Parsing Failed -> AccessTokenId Not Found")
				return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
			}

			driver, err := m.driverRepository.FindByAccessTokenId(&ctx, uint(*utils.StringToUint(accessTokenId)))
			if err != nil {
				m.logger.Info("AuthMiddleware -> Token Parsing Failed -> Driver Not Found", "error", err)
				return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
			}

			causerId := fmt.Sprintf("%d", driver.ID)
			c.Request().Header.Add("causer-id", causerId)
			c.Request().Header.Add("causer-type", causerType)

			ctx = context.WithValue(ctx, "causer-id", causerId)
			ctx = context.WithValue(ctx, "causer-type", causerType)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
