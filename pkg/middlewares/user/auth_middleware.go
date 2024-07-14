package user

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"samm/internal/module/user/domain"
	"samm/pkg/jwt"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"time"
)

func (m Middlewares) AuthenticationMiddleware(isTempToken bool) echo.MiddlewareFunc {
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

			isBearerToken, parts := utils.IsBearerToken(userToken)
			if !isBearerToken {
				m.logger.Info("AuthMiddleware -> No Bearer Token Found")
				return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
			}

			token := utils.GetValueByKey(parts, 1)
			if token == nil {
				m.logger.Info("AuthMiddleware -> No Bearer Token Found")
				return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
			}

			claims, err := m.jwtFactory.UserJwtService().ValidateToken(ctx, *token, isTempToken)
			if err != nil {
				m.logger.Info("AuthMiddleware -> ValidateToken Error -> ", err)
				return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
			}

			claim, ok := claims.(*jwt.UserJwtClaim)
			if !ok {
				m.logger.Info("AuthMiddleware -> Claims Parse Error")
				return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
			}
			c.Request().Header.Add("causer-id", claim.CauserId)
			c.Request().Header.Add("causer-type", claim.CauserType)

			return next(c)
		}
	}
}
func (m Middlewares) AuthorizationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		causerId := c.Request().Header.Get("causer-id")
		causerType := c.Request().Header.Get("causer-type")
		if causerId == "" || causerType == "" {
			m.logger.Info("AuthMiddleware -> Empty Causer ")
			return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
		}

		var user *domain.User
		userRedisKey := causerType + ":" + causerId
		err := m.redisClient.Get(userRedisKey, &user)
		if user == nil || err != nil {
			m.logger.Info("AuthMiddleware -> FindByToken MongoDB .... ")
			user, err = m.userRepository.FindUser(&ctx, utils.ConvertStringIdToObjectId(causerId))
			if err != nil {
				m.logger.Info("AuthMiddleware -> FindByToken Error -> ", err)
				return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
			}
			redisExpUserData := time.Now().AddDate(1, 0, 0).Sub(time.Now())
			setErr := m.redisClient.Set(userRedisKey, user, redisExpUserData)
			if setErr != nil {
				m.logger.Info(" REDIS -> AuthMiddleware -> Setter  Error -> ", setErr)
			}
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

func (m Middlewares) RemoveUserFromRedis(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get userId from the path parameter or form token
		userId := c.Param("id")
		if userId == "" {
			userId = c.Request().Header.Get("causer-id")
		}
		userRedisKey := "user:" + userId
		if userId != "" {
			delErr := m.redisClient.Delete(userRedisKey)
			if delErr != nil {
				m.logger.Info(" REDIS -> AuthMiddleware -> Delete Error > ", delErr)
				return next(c)
			}
			m.logger.Info(" RemoveUserFromRedis -> UserDeleted > ", userRedisKey)
		}
		return next(c)
	}
}
