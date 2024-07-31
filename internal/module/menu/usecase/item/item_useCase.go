package item

import (
	"bytes"
	"context"
	"fmt"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"samm/internal/module/menu/approval_helper"
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
	repo           domain.ItemRepository
	logger         logger.ILogger
	skuUsecase     domain.SKUUseCase
	gate           *gate.Gate
	approvalHelper *approval_helper.ApprovalItemHelper
}

func NewItemUseCase(repo domain.ItemRepository, logger logger.ILogger, skuUsecase domain.SKUUseCase, gate *gate.Gate, approvalHelper *approval_helper.ApprovalItemHelper) domain.ItemUseCase {
	return &ItemUseCase{
		repo:           repo,
		logger:         logger,
		skuUsecase:     skuUsecase,
		gate:           gate,
		approvalHelper: approvalHelper,
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

	oldDoc := item[0]
	convertDtoToCorrespondingDomain(dto, &item[0])
	doc := &item[0]
	if !oRec.gate.Authorize(doc, gate.MethodNames.Update, ctx) {
		oRec.logger.Error("AuthorizeMenuGroup -> UnAuthorized Update Admin -> ", doc.ID)
		return validators.GetErrorResponse(&ctx, localization.E1006, nil, utils.GetAsPointer(http.StatusForbidden))
	}
	err = oRec.repo.Update(ctx, &id, doc, &oldDoc)
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

func (oRec *ItemUseCase) GetByIdAndHandleApproval(ctx context.Context, id string) (responseItem.ItemResponse, validators.ErrorResponse) {
	item, err := oRec.GetById(ctx, id)
	if err.IsError {
		return item, err
	}
	errApproval := oRec.approvalHelper.UpdateItemByApproval(ctx, &item)
	return item, errApproval
}

func (oRec *ItemUseCase) CheckExists(ctx context.Context, accountId, name string, exceptProductIds ...string) (bool, validators.ErrorResponse) {
	isExists, err := oRec.repo.CheckExists(ctx, accountId, name, exceptProductIds...)
	if err != nil {
		return isExists, validators.GetErrorResponseFromErr(err)
	}
	return isExists, validators.ErrorResponse{}
}

func (oRec *ItemUseCase) ExportItems(ctx context.Context, dto utilsDto.PortalHeaders) (*excelize.File, *bytes.Buffer, validators.ErrorResponse) {
	items, err := oRec.repo.GetAllActiveItems(ctx, dto.AccountId)
	if err != nil {
		oRec.logger.Error("ExportItems", "GetAllActiveItems", err)
		return nil, nil, validators.GetErrorResponseFromErr(err)
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			oRec.logger.Error(ctx, err)
		}
	}()

	index, errS := f.NewSheet("Sheet1")
	if errS != nil {
		oRec.logger.Error(ctx, errS)
		return nil, nil, validators.GetErrorResponseFromErr(errS)
	}

	f.SetCellValue("Sheet1", "A1", "Item Id")
	f.SetCellValue("Sheet1", "B1", "Name En")
	f.SetCellValue("Sheet1", "C1", "Name Ar")
	f.SetCellValue("Sheet1", "D1", "Price")
	f.SetCellValue("Sheet1", "E1", "Category Name En")
	f.SetCellValue("Sheet1", "F1", "Category Name Ar")

	if items != nil && len(items) >= 1 {
		for key, itemDomain := range items {
			cellNumber := fmt.Sprintf("%v", key+2)
			fmt.Println(utils.ConvertObjectIdToStringId(itemDomain.ID))
			f.SetCellValue("Sheet1", "A"+cellNumber, utils.ConvertObjectIdToStringId(itemDomain.ID))
			f.SetCellValue("Sheet1", "B"+cellNumber, itemDomain.Name.En)
			f.SetCellValue("Sheet1", "C"+cellNumber, itemDomain.Name.Ar)
			f.SetCellValue("Sheet1", "D"+cellNumber, itemDomain.Price)
		}
	}

	f.SetActiveSheet(index)

	buf, err := f.WriteToBuffer()
	if err != nil {
		oRec.logger.Error(ctx, err)
		return nil, nil, validators.GetErrorResponseFromErr(err)
	}

	return f, buf, validators.ErrorResponse{}
}
