package card

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"

	"samm/internal/module/payment/domain"
	"samm/internal/module/payment/dto/card"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
)

type CardUseCase struct {
	repo   domain.CardRepository
	logger logger.ILogger
}

func (c CardUseCase) StoreCard(ctx context.Context, payload *card.CreateCardDto) (err validators.ErrorResponse) {

	// Store Card Details
	var cardDomain domain.Card
	cardDomain.MFToken = "1223"
	cardDomain.Number = utils.MaskCard(payload.Number)
	cardDomain.Type = payload.Type
	cardDomain.UserId = utils.ConvertStringIdToObjectId(payload.UserId)

	errRe := c.repo.StoreCard(ctx, &cardDomain)
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return

}

func (c CardUseCase) FindCard(ctx context.Context, Id string, userId string) (card domain.Card, err validators.ErrorResponse) {
	domainCard, errRe := c.repo.FindCard(ctx, utils.ConvertStringIdToObjectId(Id), utils.ConvertStringIdToObjectId(userId))
	if errRe != nil {
		return *domainCard, validators.GetErrorResponseFromErr(errRe)
	}
	return *domainCard, validators.ErrorResponse{}
}

func (c CardUseCase) DeleteCard(ctx context.Context, Id string, userId string) (err validators.ErrorResponse) {
	domainCard, errRe := c.repo.FindCard(ctx, utils.ConvertStringIdToObjectId(Id), utils.ConvertStringIdToObjectId(userId))
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	// Delete Card From MF domainCard
	c.logger.Info(domainCard)

	delErr := c.repo.DeleteCard(ctx, utils.ConvertStringIdToObjectId(Id), utils.ConvertStringIdToObjectId(userId))
	if delErr != nil {
		return validators.GetErrorResponseFromErr(delErr)
	}
	return validators.ErrorResponse{}
}

func (c CardUseCase) ListCard(ctx context.Context, payload *card.ListCardDto) (cards []domain.Card, paginationResult PaginationData, err validators.ErrorResponse) {
	results, paginationResult, errRe := c.repo.ListCard(ctx, payload)
	if errRe != nil {
		return results, paginationResult, validators.GetErrorResponseFromErr(errRe)
	}
	return results, paginationResult, validators.ErrorResponse{}
}

func NewCardUseCase(repo domain.CardRepository, logger logger.ILogger) domain.CardUseCase {
	return &CardUseCase{
		repo:   repo,
		logger: logger,
	}
}
