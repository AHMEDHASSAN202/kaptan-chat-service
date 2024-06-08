package item

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/item"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type ItemUseCase struct {
	repo   domain.ItemRepository
	logger logger.ILogger
}

func NewItemUseCase(repo domain.ItemRepository, logger logger.ILogger) domain.ItemUseCase {
	return &ItemUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (oRec *ItemUseCase) Create(ctx context.Context, dto []item.CreateItemDto) validators.ErrorResponse {
	err := oRec.repo.Create(ctx, convertDtoArrToCorrespondingDomain(dto))
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *ItemUseCase) Update(ctx context.Context, dto item.UpdateItemDto) validators.ErrorResponse {
	id := utils.ConvertStringIdToObjectId(dto.Id)
	doc := convertDtoToCorrespondingDomain(dto)
	err := oRec.repo.Update(ctx, &id, &doc)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}
func (oRec *ItemUseCase) SoftDelete(ctx context.Context, id string) validators.ErrorResponse {
	idDoc := utils.ConvertStringIdToObjectId(id)
	err := oRec.repo.SoftDelete(ctx, &idDoc)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *ItemUseCase) ChangeStatus(ctx context.Context, id string, dto *item.ChangeItemStatusDto) validators.ErrorResponse {
	idDoc := utils.ConvertStringIdToObjectId(id)
	err := oRec.repo.ChangeStatus(ctx, &idDoc, dto)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *ItemUseCase) List(ctx context.Context, dto *item.ListItemsDto) ([]domain.Item, validators.ErrorResponse) {
	items, err := oRec.repo.List(ctx, dto)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}
	return items, validators.ErrorResponse{}
}

func (oRec *ItemUseCase) GetById(ctx context.Context, id string) (domain.Item, validators.ErrorResponse) {
	items, err := oRec.repo.GetByIds(ctx, []primitive.ObjectID{utils.ConvertStringIdToObjectId(id)})
	if err != nil {
		return domain.Item{}, validators.GetErrorResponseFromErr(err)
	}
	if len(items) <= 0 {
		return domain.Item{}, validators.GetErrorResponse(&ctx, localization.E1002, nil)
	}
	return items[0], validators.ErrorResponse{}
}
