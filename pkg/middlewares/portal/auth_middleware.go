package portal

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"samm/internal/module/admin/consts"
	"samm/pkg/jwt"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
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

		claims, err := m.jwtFactory.PortalJwtService().ValidateToken(ctx, *token)
		if err != nil {
			m.logger.Info("AuthMiddleware -> ValidateToken Error -> ", err)
			return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
		}

		_, ok := claims.(*jwt.PortalJwtClaim)
		if !ok {
			m.logger.Info("AuthMiddleware -> Claims Parse Error")
			return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
		}

		admin, err := m.adminRepository.FindByToken(ctx, *token, []string{consts.ADMIN_TYPE, consts.PORTAL_TYPE})
		if err != nil {
			m.logger.Info("AuthMiddleware -> FindByToken Error -> ", err)
			return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
		}

		if !admin.IsActive() {
			m.logger.Info("AuthMiddleware -> Admin Is Not Active")
			return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
		}

		jsonByte, err := json.Marshal(admin)
		if err != nil {
			m.logger.Info("AuthMiddleware -> Marshal Error -> ", err)
			return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
		}

		c.Request().Header.Add("causer-id", utils.ConvertObjectIdToStringId(admin.ID))
		c.Request().Header.Add("causer-type", admin.Type)
		c.Request().Header.Add("causer-details", string(jsonByte))

		return next(c)
	}
}
