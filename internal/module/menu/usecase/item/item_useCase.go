package item

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/item"
	"samm/internal/module/menu/responses"
	responseItem "samm/internal/module/menu/responses/item"
	"samm/pkg/gate"
	"samm/pkg/logger"
	"samm/pkg/utils"
	utilsDto "samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"time"
)

type ItemUseCase struct {
	repo       domain.ItemRepository
	logger     logger.ILogger
	skuUsecase domain.SKUUseCase
	gate       *gate.Gate
}

func NewItemUseCase(repo domain.ItemRepository, logger logger.ILogger, skuUsecase domain.SKUUseCase, gate *gate.Gate) domain.ItemUseCase {
	return &ItemUseCase{
		repo:       repo,
		logger:     logger,
		skuUsecase: skuUsecase,
		gate:       gate,
	}
}

func (oRec *ItemUseCase) Create(ctx context.Context, dto []item.CreateItemDto) validators.ErrorResponse {
	err := oRec.repo.Create(ctx, convertDtoArrToCorrespondingDomain(dto))
	if err != nil {
		oRec.logger.Error("ItemUseCase", "Create", err)
		return validators.GetErrorResponseFromErr(err)
	}

	//create sku
	skus := make([]string, 0)
	for _, i := range dto {
		skus = append(skus, i.SKU)
	}
	errResp := oRec.skuUsecase.CreateBulk(ctx, skus)
	if errResp.IsError {
		oRec.logger.Error("itemuseCase", "createSku", errResp.ErrorMessageObject)
		return errResp
	}
	return validators.ErrorResponse{}
}

func (oRec *ItemUseCase) Update(ctx context.Context, dto item.UpdateItemDto) validators.ErrorResponse {
	id := utils.ConvertStringIdToObjectId(dto.Id)
	item, err := oRec.repo.GetByIds(ctx, []primitive.ObjectID{id})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return validators.GetErrorResponse(&ctx, localization.E1002, nil, nil)
		}
		return validators.GetErrorResponseFromErr(err)
	}

	convertDtoToCorrespondingDomain(dto, &item[0])
	doc := &item[0]
	if !oRec.gate.Authorize(doc, gate.MethodNames.Update, ctx) {
		oRec.logger.Error("AuthorizeMenuGroup -> UnAuthorized Update Admin -> ", doc.ID)
		return validators.GetErrorResponse(&ctx, localization.E1006, nil, utils.GetAsPointer(http.StatusForbidden))
	}
	err = oRec.repo.Update(ctx, &id, doc)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	//create sku
	errResp := oRec.skuUsecase.CreateBulk(ctx, []string{dto.SKU})
	if errResp.IsError {
		oRec.logger.Error("itemuseCase", "createSku", errResp.ErrorMessageObject)
	}
	return validators.ErrorResponse{}
}
func (oRec *ItemUseCase) SoftDelete(ctx context.Context, id string, input item.DeleteItemDto) validators.ErrorResponse {
	idDoc := utils.ConvertStringIdToObjectId(id)
	item, err := oRec.repo.GetByIds(ctx, []primitive.ObjectID{idDoc})
	if err != nil || len(item) <= 0 {
		if err == mongo.ErrNoDocuments {
			return validators.GetErrorResponse(&ctx, localization.E1002, nil, nil)
		}
		return validators.GetErrorResponseFromErr(err)
	}
	if !oRec.gate.Authorize(&item[0], gate.MethodNames.Delete, ctx) {
		oRec.logger.Error("AuthorizeMenuGroup -> UnAuthorized Delete Admin -> ", item[0].ID)
		return validators.GetErrorResponse(&ctx, localization.E1006, nil, utils.GetAsPointer(http.StatusForbidden))
	}

	t := time.Now()
	item[0].DeletedAt = &t
	item[0].AdminDetails = append(item[0].AdminDetails, utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: input.CauserName, Operation: "Delete", UpdatedAt: time.Now()})
	err = oRec.repo.SoftDelete(ctx, &item[0])
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *ItemUseCase) ChangeStatus(ctx context.Context, id string, dto *item.ChangeItemStatusDto) validators.ErrorResponse {
	idDoc := utils.ConvertStringIdToObjectId(id)
	item, err := oRec.repo.GetByIds(ctx, []primitive.ObjectID{idDoc})
	if err != nil || len(item) <= 0 {
		if err == mongo.ErrNoDocuments {
			return validators.GetErrorResponse(&ctx, localization.E1002, nil, nil)
		}
		return validators.GetErrorResponseFromErr(err)
	}

	if !oRec.gate.Authorize(&item[0], gate.MethodNames.Update, ctx) {
		oRec.logger.Error("AuthorizeMenuGroup -> UnAuthorized Update Admin -> ", item[0].ID)
		return validators.GetErrorResponse(&ctx, localization.E1006, nil, utils.GetAsPointer(http.StatusForbidden))
	}
	item[0].Status = dto.Status
	item[0].AdminDetails = append(item[0].AdminDetails, utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: dto.CauserName, Operation: "Change Status", UpdatedAt: time.Now()})
	err = oRec.repo.ChangeStatus(ctx, &item[0])
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *ItemUseCase) List(ctx context.Context, dto *item.ListItemsDto) (*responses.ListResponse, validators.ErrorResponse) {
	items, pgination, err := oRec.repo.List(ctx, dto)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}

	return responses.SetListResponse(items, pgination), validators.ErrorResponse{}
}

func (oRec *ItemUseCase) GetById(ctx context.Context, id string) (responseItem.ItemResponse, validators.ErrorResponse) {
	items, err := oRec.repo.Find(ctx, utils.ConvertStringIdToObjectId(id))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return responseItem.ItemResponse{}, validators.GetErrorResponse(&ctx, localization.E1002, nil, nil)
		}
		return responseItem.ItemResponse{}, validators.GetErrorResponseFromErr(err)
	}
	oRec.logger.Info(items)
	if !oRec.gate.Authorize(&domain.Item{AccountId: items.AccountId}, gate.MethodNames.Find, ctx) {
		oRec.logger.Error("AuthorizeMenuGroup -> UnAuthorized Find Admin -> ", items.ID)
		return responseItem.ItemResponse{}, validators.GetErrorResponse(&ctx, localization.E1006, nil, utils.GetAsPointer(http.StatusForbidden))
	}
	return items, validators.ErrorResponse{}
}

func (oRec *ItemUseCase) CheckExists(ctx context.Context, accountId, name string, exceptProductIds ...string) (bool, validators.ErrorResponse) {
	isExists, err := oRec.repo.CheckExists(ctx, accountId, name, exceptProductIds...)
	if err != nil {
		return isExists, validators.GetErrorResponseFromErr(err)
	}
	return isExists, validators.ErrorResponse{}
}
