package delivery

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"samm/internal/module/admin/custom_validators"
	"samm/internal/module/admin/domain"
	dto "samm/internal/module/admin/dto/role"
	"samm/pkg/logger"
	"samm/pkg/middlewares/admin"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type RoleHandler struct {
	roleUseCase         domain.RoleUseCase
	roleCustomValidator custom_validators.RoleCustomValidator
	validator           *validator.Validate
	logger              logger.ILogger
}

// InitMenuGroupController will initialize the article's HTTP controller
func InitRoleController(e *echo.Echo, roleUseCase domain.RoleUseCase, roleCustomValidator custom_validators.RoleCustomValidator, validator *validator.Validate, logger logger.ILogger, adminMiddlewares *admin.ProviderMiddlewares) {
	handler := &RoleHandler{
		roleUseCase:         roleUseCase,
		validator:           validator,
		logger:              logger,
		roleCustomValidator: roleCustomValidator,
	}
	role := e.Group("api/v1/admin/role")
	role.Use(adminMiddlewares.AuthMiddleware)
	{
		role.GET("", handler.List)
		role.GET("/:id", handler.Find)
		role.POST("", handler.Create)
		role.PUT("/:id", handler.Update)
		role.DELETE("/:id", handler.Delete)
	}

	permissions := e.Group("api/v1/admin/permissions")
	permissions.Use(adminMiddlewares.AuthMiddleware)
	{
		permissions.GET("", handler.ListPermissions)
	}
}

func (a *RoleHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input dto.ListRoleDTO
	binder := &echo.DefaultBinder{}
	err := binder.BindHeaders(c, &input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	err = c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	data, errResp := a.roleUseCase.List(ctx, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"roles": data})
}

func (a *RoleHandler) Find(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, localization.E1002, nil, nil))
	}

	data, errResp := a.roleUseCase.Find(ctx, utils.ConvertStringIdToObjectId(id))
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"role": data})
}

func (a *RoleHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input dto.CreateRoleDTO
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	id, errResp := a.roleUseCase.Create(ctx, &input)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"id": id})
}

func (a *RoleHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, localization.E1002, nil, nil))
	}

	var input dto.CreateRoleDTO
	input.ID = utils.ConvertStringIdToObjectId(id)

	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	id, errResp := a.roleUseCase.Update(ctx, &input)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"id": id})
}

func (a *RoleHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input dto.DeleteRoleDTO
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusBadRequest(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator, a.roleCustomValidator.ValidateRoleHasAdmins(), a.roleCustomValidator.ValidateStaticRoles())
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusBadRequest(c, validationErr)
	}

	errResp := a.roleUseCase.Delete(ctx, utils.ConvertStringIdToObjectId(input.ID))
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *RoleHandler) ListPermissions(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// Read the permissions JSON file
	dir, err := os.Getwd()
	if err != nil {
		return validators.ErrorStatusBadRequest(c, validators.GetErrorResponseFromErr(err))
	}
	permissionsPath := filepath.Join(dir, "internal", "module", "admin", "assets", "permissions.json")
	data, err := os.ReadFile(permissionsPath)
	if err != nil {
		a.logger.Error("ListPermissions -> Error -> ", err)
		return validators.ErrorStatusBadRequest(c, validators.GetErrorResponseFromErr(errors.New("Unable to read permissions file")))
	}

	var permissions []map[string]interface{}
	if err := json.Unmarshal(data, &permissions); err != nil {
		a.logger.Error("ListPermissions -> Error -> ", err)
		return validators.ErrorStatusBadRequest(c, validators.GetErrorResponseFromErr(errors.New("Unable to parse permissions file")))
	}

	return validators.SuccessResponse(c, map[string]interface{}{"permissions": permissions})
}
