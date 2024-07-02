package usecase

import (
	"context"
	"github.com/kamva/mgm/v3"
	"samm/internal/module/example/domain"
	"samm/internal/module/example/dto"
	"samm/pkg/logger"
	"samm/pkg/validators"
)

type MenuUseCase struct {
	repo   domain.MenuRepository
	logger logger.ILogger
}

const tag = " LocationUseCase "

func NewMenuUseCase(repo domain.MenuRepository, logger logger.ILogger) domain.MenuUseCase {
	return &MenuUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (oRec *MenuUseCase) UpdateTokens(ctx context.Context) validators.ErrorResponse {
	//TODO implement me
	panic("implement me")
}

func (oRec *MenuUseCase) MenuWebhook(ctx context.Context, webhook *dto.MenusWebhook) validators.ErrorResponse {
	//TODO implement me
	panic("implement me")
}

func (oRec *MenuUseCase) LocationActiveStatusWebhook(ctx context.Context, webhook *dto.BusyModeLocationWebhook) validators.ErrorResponse {
	//TODO implement me
	panic("implement me")
}

func (oRec *MenuUseCase) MenuActiveStatusWebhook(ctx context.Context, webhook *dto.SnoozeMenuWebhook) validators.ErrorResponse {
	//TODO implement me
	panic("implement me")
}

func (oRec *MenuUseCase) LocationRegisterWebhook(ctx context.Context, payload *dto.LocationRegisterWebhook) (errResponse validators.ErrorResponse) {
	oRec.repo.UpdateTokens(ctx, &domain.Token{
		DefaultModel: mgm.DefaultModel{},
		AccessToken:  "aqw		q",
		ExpiresAt:    1111,
		ExpiresIn:    22222,
		TokenType:    "231123",
		Scope:        "123123312",
	})
	return validators.ErrorResponse{}
}

func (oRec *MenuUseCase) FindLocation(ctx context.Context, locationId string) (response domain.Location, errResponse validators.ErrorResponse) {
	token, err := oRec.repo.FindToken(ctx)
	if err != nil {
		return domain.Location{}, validators.ErrorResponse{}
	}
	oRec.logger.Info("fawaaaaaaz", token)
	return domain.Location{}, validators.ErrorResponse{}
}
