package commmon

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"kaptan/pkg/localization"
	"kaptan/pkg/utils"
	"kaptan/pkg/validators"
	"net/http"
)

func (m ProviderMiddlewares) PermissionMiddleware(permissions ...string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			if c.Request().Method == http.MethodOptions {
				return next(c)
			}

			if len(permissions) == 0 {
				return next(c)
			}

			myPermissionsJSON := c.Request().Header.Get("causer-permissions")
			if myPermissionsJSON == "" {
				m.logger.Info("PermissionMiddleware -> Missing permissions")
				return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1403, nil, utils.GetAsPointer(http.StatusForbidden)))
			}

			var myPermissions []string
			if err := json.Unmarshal([]byte(myPermissionsJSON), &myPermissions); err != nil {
				m.logger.Info("PermissionMiddleware -> Failed to parse permissions")
				return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1403, nil, utils.GetAsPointer(http.StatusForbidden)))
			}

			if !hasAtLeastOnePermission(myPermissions, permissions) {
				m.logger.Info("PermissionMiddleware -> Insufficient permissions")
				return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.E1403, nil, utils.GetAsPointer(http.StatusForbidden)))
			}

			return next(c)
		}
	}
}

func hasAtLeastOnePermission(myPermissions []string, requiredPermissions []string) bool {
	for _, reqPerm := range requiredPermissions {
		for _, userPerm := range myPermissions {
			if reqPerm == userPerm {
				return true
			}
		}
	}
	return false
}
