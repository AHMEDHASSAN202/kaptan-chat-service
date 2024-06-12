package sku

import (
	"context"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/sku"
	"samm/pkg/logger"
	"samm/pkg/validators"
	"time"

	"github.com/jinzhu/copier"
)

type SKUUseCase struct {
	repo   domain.SKURepository
	logger logger.ILogger
}

func NewSKUUseCase(repo domain.SKURepository, logger logger.ILogger) domain.SKUUseCase {
	return &SKUUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (oRec *SKUUseCase) Create(ctx context.Context, dto sku.CreateSKUDto) validators.ErrorResponse {
	var skuDoc domain.SKU
	copier.Copy(&skuDoc, &dto)
	skuDoc.CreatedAt = time.Now()
	skuDoc.UpdatedAt = time.Now()
	err := oRec.repo.Create(ctx, skuDoc)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *SKUUseCase) List(ctx context.Context, dto *sku.ListSKUDto) ([]domain.SKU, validators.ErrorResponse) {
	skus, err := oRec.repo.List(ctx, dto)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}
	return skus, validators.ErrorResponse{}
}

func (oRec *SKUUseCase) CheckExists(ctx context.Context, name string) (bool, validators.ErrorResponse) {
	isExists, err := oRec.repo.CheckExists(ctx, name)
	if err != nil {
		return isExists, validators.GetErrorResponseFromErr(err)
	}
	return isExists, validators.ErrorResponse{}
}
