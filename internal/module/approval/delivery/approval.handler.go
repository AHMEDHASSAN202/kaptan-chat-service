package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/approval/domain"
	dto "samm/internal/module/approval/dto"
	"samm/pkg/logger"
	"samm/pkg/middlewares/admin"
	commmon "samm/pkg/middlewares/common"
	"samm/pkg/validators"
)

type ApprovalHandler struct {
	approvalUseCase domain.ApprovalUseCase
	validator       *validator.Validate
	logger          logger.ILogger
}

func InitApprovalController(e *echo.Echo, us domain.ApprovalUseCase, validator *validator.Validate, logger logger.ILogger, adminMiddlewares *admin.ProviderMiddlewares, commonMiddlewares *commmon.ProviderMiddlewares) {
	handler := &ApprovalHandler{
		approvalUseCase: us,
		validator:       validator,
		logger:          logger,
	}
	approvalRoutes := e.Group("api/v1/admin/approval")
	approvalRoutes.Use(adminMiddlewares.AuthMiddleware)
	approvalRoutes.GET("", handler.ListApprovals, commonMiddlewares.PermissionMiddleware("list-approval"))
	approvalRoutes.PUT("/:id", handler.ChangeApprovalStatus, commonMiddlewares.PermissionMiddleware("update-status-approval"))
}

func (a *ApprovalHandler) ListApprovals(c echo.Context) error {
	ctx := c.Request().Context()
	var input dto.ListApprovalDto
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
	data, errResp := a.approvalUseCase.List(ctx, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"approvals": data})
}

func (a *ApprovalHandler) ChangeApprovalStatus(c echo.Context) error {
	ctx := c.Request().Context()
	var input dto.ChangeStatusApprovalDto
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &input); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}
	errResp := a.approvalUseCase.ChangeStatus(ctx, &input)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
