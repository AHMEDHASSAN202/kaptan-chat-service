package collection_method

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	commonDomain "samm/internal/module/common/domain"
	"samm/internal/module/user/domain"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type CollectionMethodUseCase struct {
	repo          domain.CollectionMethodRepository
	commonUseCase commonDomain.CommonUseCase
	logger        logger.ILogger
}

const tag = " CollectionMethodUseCase "

func NewCollectionMethodUseCase(repo domain.CollectionMethodRepository, logger logger.ILogger, commonUseCase commonDomain.CommonUseCase) domain.CollectionMethodUseCase {
	return &CollectionMethodUseCase{
		repo:          repo,
		commonUseCase: commonUseCase,
		logger:        logger,
	}
}

func (l CollectionMethodUseCase) StoreCollectionMethod(ctx context.Context, collectMethodDomain *domain.CollectionMethods) (err validators.ErrorResponse) {
	errRe := l.repo.StoreCollectionMethod(ctx, collectMethodDomain)
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return
}

func (l CollectionMethodUseCase) UpdateCollectionMethod(ctx context.Context, id string, collectMethodDomain *domain.CollectionMethods) (err validators.ErrorResponse) {
	doc, errRe := l.repo.FindCollectionMethod(ctx, utils.ConvertStringIdToObjectId(id), collectMethodDomain.UserId)
	if errRe != nil {
		if errRe == mongo.ErrNoDocuments {
			return validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
		}
		return validators.GetErrorResponseFromErr(errRe)
	}
	collectMethodDomain.ID = utils.ConvertStringIdToObjectId(id)
	collectMethodDomain.CreatedAt = doc.CreatedAt
	errRe = l.repo.UpdateCollectionMethod(ctx, collectMethodDomain)
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return
}
func (l CollectionMethodUseCase) FindCollectionMethod(ctx context.Context, Id string, userId string) (user domain.CollectionMethods, err validators.ErrorResponse) {
	//find the default location id
	collectionMethod, _ := l.commonUseCase.FindCollectionMethodByDefaultId(ctx, Id)
	if collectionMethod != nil {
		if val, ok := collectionMethod["type"].(string); ok {
			return domain.CollectionMethods{Type: val, UserId: utils.ConvertStringIdToObjectId(userId)}, err
		}
	}
	////find the user collection method
	domainCollectionMethod, errRe := l.repo.FindCollectionMethod(ctx, utils.ConvertStringIdToObjectId(Id), utils.ConvertStringIdToObjectId(userId))
	if errRe != nil {
		if errRe == mongo.ErrNoDocuments {
			return *domainCollectionMethod, validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
		}
		return *domainCollectionMethod, validators.GetErrorResponseFromErr(errRe)
	}
	errResp := buildResponse(ctx, l, &[]domain.CollectionMethods{*domainCollectionMethod})
	if errResp.IsError {
		return *domainCollectionMethod, errResp
	}
	return *domainCollectionMethod, validators.ErrorResponse{}
}

func (l CollectionMethodUseCase) DeleteCollectionMethod(ctx context.Context, Id string, userId string) (err validators.ErrorResponse) {
	_, errRe := l.repo.FindCollectionMethod(ctx, utils.ConvertStringIdToObjectId(Id), utils.ConvertStringIdToObjectId(Id))
	if errRe != nil {
		if errRe == mongo.ErrNoDocuments {
			return validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
		}
		return validators.GetErrorResponseFromErr(errRe)
	}
	delErr := l.repo.DeleteCollectionMethod(ctx, utils.ConvertStringIdToObjectId(Id), utils.ConvertStringIdToObjectId(userId))
	if delErr != nil {
		return validators.GetErrorResponseFromErr(delErr)
	}
	return validators.ErrorResponse{}
}

func (l CollectionMethodUseCase) ListCollectionMethod(ctx context.Context, collectionMethodType string, userId string) ([]domain.CollectionMethods, validators.ErrorResponse) {
	results, errRe := l.repo.ListCollectionMethod(ctx, collectionMethodType, utils.ConvertStringIdToObjectId(userId))
	if errRe != nil {
		return results, validators.GetErrorResponseFromErr(errRe)
	}
	errResp := buildResponse(ctx, l, &results)
	if errResp.IsError {
		return results, errResp
	}
	return results, validators.ErrorResponse{}

}

func buildResponse(ctx context.Context, l CollectionMethodUseCase, results *[]domain.CollectionMethods) validators.ErrorResponse {
	assets, errRe := l.commonUseCase.ListAssets(ctx, true, true)
	if errRe.IsError {
		return errRe
	}
	carColorsResult := assets.(map[string]any)["car_colors"]
	carBrandsResult := assets.(map[string]any)["car_brands"]
	for i, collection := range *results {
		if collection.Type == "drive_through" {
			for _, color := range carColorsResult.([]map[string]interface{}) {
				if color["id"] == collection.Values["car_color"] {
					collection.Values["car_color"] = color
				}
			}
			for _, brand := range carBrandsResult.([]map[string]interface{}) {
				if brand["id"] == collection.Values["car_brand"] {
					collection.Values["car_brand"] = brand
				}
			}
			(*results)[i].Values = collection.Values
		}
	}
	return validators.ErrorResponse{}
}
