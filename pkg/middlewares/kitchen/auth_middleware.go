package kitchen

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"samm/internal/module/admin/consts"
	"samm/internal/module/admin/domain"
	"samm/pkg/jwt"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"time"
)

func (m ProviderMiddlewares) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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

		claims, err := m.jwtFactory.KitchenJwtService().ValidateToken(ctx, *token)
		if err != nil {
			m.logger.Info("AuthMiddleware -> ValidateToken Error -> ", err)
			return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
		}

		data, ok := claims.(*jwt.PortalJwtClaim)
		if !ok {
			m.logger.Info("AuthMiddleware -> Claims Parse Error")
			return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
		}

		var admin *domain.Admin
		key := data.CauserType + ":" + data.CauserId
		err = m.redisClient.Get(key, &admin)
		if admin == nil || err != nil {
			m.logger.Info("kitchen -> AuthMiddleware -> FindByToken MongoDB .... ")
			admin, err = m.adminRepository.FindByToken(ctx, *token, []string{consts.KITCHEN_TYPE})
			if err != nil {
				m.logger.Info("kitchen -> AuthMiddleware -> FindByToken Error -> ", err)
				return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
			}
			setErr := m.redisClient.Set(key, admin, data.ExpiresAt.Sub(time.Now()))
			if setErr != nil {
				m.logger.Info("kitchen -> REDIS -> AuthMiddleware -> Setter > ", setErr)
			}
		}

		if !admin.IsActive() {
			m.logger.Info("AuthMiddleware -> Admin Is Not Active")
			return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
		}

		jsonPermissionsByte, err := json.Marshal(admin.Role.Permissions)
		if err != nil {
			m.logger.Info("AuthMiddleware -> Marshal Error -> ", err)
			return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
		}

		c.Request().Header.Add("causer-id", utils.ConvertObjectIdToStringId(admin.ID))
		c.Request().Header.Add("causer-type", admin.Type)
		c.Request().Header.Add("causer-name", admin.Name)
		c.Request().Header.Add("causer-permissions", string(jsonPermissionsByte))

		//accountId := ""
		//if admin.Account != nil {
		//	accountId = utils.ConvertObjectIdToStringId(admin.Account.Id)
		//} else if utils.SafeMapGet(data.CauserDetails, "id", "") != "" {
		//	accountId = utils.SafeMapGet(data.CauserDetails, "id", "").(string)
		//}
		//
		//c.Request().Header.Add("causer-account-id", accountId)
		//c.Request().Header.Set("account-id", accountId)

		if data.CauserDetails != nil {
			jsonDetailsByte, err := json.Marshal(data.CauserDetails)
			if err != nil {
				m.logger.Info("AuthMiddleware -> Marshal Error -> jsonDetailsByte -> ", err)
				return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1401, nil, utils.GetAsPointer(http.StatusUnauthorized)))
			}
			c.Request().Header.Add("causer-details", string(jsonDetailsByte))
		}

		ctx = context.WithValue(ctx, "causer-id", utils.ConvertObjectIdToStringId(admin.ID))
		ctx = context.WithValue(ctx, "causer-type", admin.Type)
		ctx = context.WithValue(ctx, "causer-details", admin)
		//ctx = context.WithValue(ctx, "causer-account-id", accountId)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
