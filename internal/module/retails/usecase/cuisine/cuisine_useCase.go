package cuisine

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/cuisine"
	"samm/internal/module/retails/responses"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"strings"
	"time"
)

type CuisineUseCase struct {
	repo      domain.CuisineRepository
	brandRepo domain.BrandRepository
	logger    logger.ILogger
}

func NewCuisineUseCase(repo domain.CuisineRepository, brandRepo domain.BrandRepository, logger logger.ILogger) domain.CuisineUseCase {
	return &CuisineUseCase{
		repo:      repo,
		logger:    logger,
		brandRepo: brandRepo,
	}
}

func (oRec *CuisineUseCase) Create(ctx *context.Context, dto *cuisine.CreateCuisineDto) (*domain.Cuisine, validators.ErrorResponse) {
	doc := convertDtoArrToCorrespondingDomain(dto)
	err := oRec.repo.Create(doc)
	if err != nil {
		return doc, validators.GetErrorResponseFromErr(err)
	}
	return doc, validators.ErrorResponse{}
}

func (oRec *CuisineUseCase) Update(ctx *context.Context, dto *cuisine.UpdateCuisineDto) validators.ErrorResponse {
	findCuisine, findCuisineErr := oRec.GetById(ctx, dto.Id)
	if findCuisineErr.IsError {
		return findCuisineErr
	}
	doc := domainBuilderAtUpdate(dto, findCuisine)
	err := oRec.repo.UpdateCuisineAndLocations(doc)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}
func (oRec *CuisineUseCase) SoftDelete(ctx *context.Context, id string, adminDetails *dto.AdminHeaders) validators.ErrorResponse {
	idDoc := utils.ConvertStringIdToObjectId(id)
	causerDetails := dto.AdminDetails{Id: utils.ConvertStringIdToObjectId(adminDetails.CauserId), Name: adminDetails.CauserName, Type: adminDetails.CauserType, Operation: "Delete Brand", UpdatedAt: time.Now()}
	transactionErr := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		brandErr := oRec.brandRepo.DeleteCuisinesFromBrand(sc, idDoc)
		if brandErr != nil {
			return brandErr
		}
		cuisineErr := oRec.repo.SoftDelete(sc, idDoc, &causerDetails)
		if cuisineErr != nil {
			return cuisineErr
		}
		return session.CommitTransaction(sc)
	})

	if transactionErr != nil {
		return validators.GetErrorResponseFromErr(transactionErr)
	}
	return validators.ErrorResponse{}
}

func (oRec *CuisineUseCase) ChangeStatus(ctx *context.Context, dto *cuisine.ChangeCuisineStatusDto) validators.ErrorResponse {
	err := oRec.repo.ChangeStatus(ctx, dto)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *CuisineUseCase) ListCuisinesForDashboard(ctx *context.Context, dto *cuisine.ListCuisinesDto) (*responses.ListResponse, validators.ErrorResponse) {
	cuisines, paginationMeta, resErr := oRec.repo.List(ctx, false, dto)
	if resErr != nil {
		return nil, validators.GetErrorResponseFromErr(resErr)
	}
	return responses.SetListResponse(cuisines, paginationMeta), validators.ErrorResponse{}
}

func (oRec *CuisineUseCase) ListCuisinesForMobile(ctx *context.Context, dto *cuisine.ListCuisinesDto) (*responses.ListResponse, validators.ErrorResponse) {
	cuisines, paginationMeta, resErr := oRec.repo.List(ctx, true, dto)
	if resErr != nil {
		return nil, validators.GetErrorResponseFromErr(resErr)
	}
	return responses.SetListResponse(mobileListCuisineBuilder(cuisines), paginationMeta), validators.ErrorResponse{}
}

func (oRec *CuisineUseCase) Find(ctx *context.Context, id string) (*domain.Cuisine, validators.ErrorResponse) {
	cuisine, err := oRec.repo.Find(ctx, utils.ConvertStringIdToObjectId(id))
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}
	if cuisine == nil {
		return nil, validators.GetErrorResponseFromErr(errors.New(localization.E1002))
	}

	return cuisine, validators.ErrorResponse{}
}

func (oRec *CuisineUseCase) GetById(ctx *context.Context, id string) (*domain.Cuisine, validators.ErrorResponse) {
	cuisines, err := oRec.repo.GetByIds(ctx, &[]primitive.ObjectID{utils.ConvertStringIdToObjectId(id)})
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}
	cuisinesVal := *cuisines
	if len(cuisinesVal) == 0 {
		return nil, validators.GetErrorResponseFromErr(errors.New(localization.E1002))
	}

	return &cuisinesVal[0], validators.ErrorResponse{}
}

func (oRec *CuisineUseCase) CheckExists(ctx *context.Context, ids []string) validators.ErrorResponse {
	objIds := utils.ConvertStringIdsToObjectIds(ids)
	cuisines, err := oRec.repo.GetByIds(ctx, &objIds)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	if len(*cuisines) == 0 {
		return validators.GetErrorResponseFromErr(errors.New(localization.E1002))
	}
	cuisineIds := getCuisinesIds(cuisines)
	diffIds := utils.ElementsDiff(ids, cuisineIds)
	if diffIds != nil && len(diffIds) > 0 {
		oRec.logger.Error(strings.Join(diffIds, ", ") + " cuisine not exist")
		return validators.GetErrorResponseFromErr(errors.New(localization.E1002))
	}
	return validators.ErrorResponse{}
}

func (oRec *CuisineUseCase) CheckNameExists(ctx context.Context, name string) (bool, validators.ErrorResponse) {
	isExists, err := oRec.repo.CheckNameExists(ctx, name)
	if err != nil {
		return isExists, validators.GetErrorResponseFromErr(err)
	}
	return isExists, validators.ErrorResponse{}
}
