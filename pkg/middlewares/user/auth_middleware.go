package user

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"samm/internal/module/user/domain"
	"samm/pkg/jwt"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"time"
)

func (m Middlewares) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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

		claims, err := m.jwtFactory.UserJwtService().ValidateToken(ctx, *token)
		if err != nil {
			m.logger.Info("AuthMiddleware -> ValidateToken Error -> ", err)
			return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
		}

		claim, ok := claims.(*jwt.UserJwtClaim)
		if !ok {
			m.logger.Info("AuthMiddleware -> Claims Parse Error")
			return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
		}

		var user *domain.User
		err = m.redisClient.Get(claim.CauserId, &user)
		if user == nil || err != nil {
			m.logger.Info("AuthMiddleware -> FindByToken MongoDB .... ")
			user, err = m.userRepository.FindByToken(ctx, *token)
			if err != nil {
				m.logger.Info("AuthMiddleware -> FindByToken Error -> ", err)
				return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
			}
			err := m.redisClient.Set(claim.CauserId, user, time.Hour*60)
			if err != nil {
				fmt.Println(" REDIS -> AuthMiddleware -> Setter > ", err)
			}
		}

		if !user.IsActive {
			m.logger.Info("AuthMiddleware -> Admin Is Not Active")
			return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
		}

		jsonByte, err := json.Marshal(user)
		if err != nil {
			m.logger.Info("AuthMiddleware -> Marshal Error -> ", err)
			return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
		}

		c.Request().Header.Add("causer-id", utils.ConvertObjectIdToStringId(user.ID))
		c.Request().Header.Add("causer-type", "user")
		c.Request().Header.Add("causer-details", string(jsonByte))

		return next(c)
	}
}

func (m Middlewares) TempAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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

		claims, err := m.jwtFactory.UserJwtService().ValidateToken(ctx, *token, true)
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
		c.Request().Header.Add("causer-type", "user")

		return next(c)
	}
}
