package domain

import (
	"context"
	"github.com/kamva/mgm/v3"
	"samm/internal/module/menu/dto"
	"samm/pkg/validators"
)

type Token struct {
	// DefaultModel adds _id, created_at and updated_at fields to the Model.
	mgm.DefaultModel `bson:",inline"`
	AccessToken      string `json:"access_token" bson:"access_token"`
	ExpiresAt        int    `json:"expires_at" bson:"expires_at"`
	ExpiresIn        int    `json:"expires_in" bson:"expires_in"`
	TokenType        string `json:"token_type" bson:"token_type"`
	Scope            string `json:"scope" bson:"scope"`
}

type MenuUseCase interface {
	UpdateTokens(ctx context.Context) validators.ErrorResponse
	MenuWebhook(ctx context.Context, webhook *dto.MenusWebhook) validators.ErrorResponse
	LocationActiveStatusWebhook(ctx context.Context, webhook *dto.BusyModeLocationWebhook) validators.ErrorResponse
	MenuActiveStatusWebhook(ctx context.Context, webhook *dto.SnoozeMenuWebhook) validators.ErrorResponse
	LocationRegisterWebhook(ctx context.Context, payload *dto.LocationRegisterWebhook) (errResponse validators.ErrorResponse)
	FindLocation(ctx context.Context, locationId string) (response Location, errResponse validators.ErrorResponse)

	//ListMenus(ctx context.Context, storeId string) (menuValidator.ErrorResponse, dto.MenuData)

}

type MenuRepository interface {
	UpdateTokens(ctx context.Context, domainData *Token) (err error)
	FindToken(ctx context.Context) (result Token, err error)
}
