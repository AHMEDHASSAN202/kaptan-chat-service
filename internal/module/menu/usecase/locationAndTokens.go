package usecase

import (
	"context"
	"example.com/fxdemo/internal/module/menu/domain"
	"example.com/fxdemo/internal/module/menu/dto"
	"example.com/fxdemo/pkg/logger"
	"github.com/kamva/mgm/v3"
	menu_validator "github.com/n-goo/ngo-menu-service/pkg/validators"
)

type MenuUseCase struct {
	repo   domain.MenuRepository
	logger logger.ILogger
}

const tag = " MenuUseCase "

func NewMenuUseCase(repo domain.MenuRepository, logger logger.ILogger) domain.MenuUseCase {
	return &MenuUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (oRec *MenuUseCase) UpdateTokens(ctx context.Context) menu_validator.ErrorResponse {
	//TODO implement me
	panic("implement me")
}

func (oRec *MenuUseCase) MenuWebhook(ctx context.Context, webhook *dto.MenusWebhook) menu_validator.ErrorResponse {
	//TODO implement me
	panic("implement me")
}

func (oRec *MenuUseCase) LocationActiveStatusWebhook(ctx context.Context, webhook *dto.BusyModeLocationWebhook) menu_validator.ErrorResponse {
	//TODO implement me
	panic("implement me")
}

func (oRec *MenuUseCase) MenuActiveStatusWebhook(ctx context.Context, webhook *dto.SnoozeMenuWebhook) menu_validator.ErrorResponse {
	//TODO implement me
	panic("implement me")
}

func (oRec *MenuUseCase) LocationRegisterWebhook(ctx context.Context, payload *dto.LocationRegisterWebhook) (errResponse menu_validator.ErrorResponse) {
	oRec.repo.UpdateTokens(ctx, &domain.Token{
		DefaultModel: mgm.DefaultModel{},
		AccessToken:  "aqw		q",
		ExpiresAt:    1111,
		ExpiresIn:    22222,
		TokenType:    "231123",
		Scope:        "123123312",
	})
	return menu_validator.ErrorResponse{}
}

func (oRec *MenuUseCase) FindLocation(ctx context.Context, locationId string) (response domain.Location, errResponse menu_validator.ErrorResponse) {
	token, err := oRec.repo.FindToken(ctx)
	if err != nil {
		return domain.Location{}, menu_validator.ErrorResponse{}
	}
	oRec.logger.Info("fawaaaaaaz", token)
	return domain.Location{}, menu_validator.ErrorResponse{}
}
