package brand

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/brand"
	"samm/internal/module/retails/responses"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"time"
)

type BrandUseCase struct {
	repo   domain.BrandRepository
	logger logger.ILogger
}

func NewBrandUseCase(repo domain.BrandRepository, logger logger.ILogger) domain.BrandUseCase {
	return &BrandUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (oRec *BrandUseCase) Create(ctx context.Context, dto *brand.CreateBrandDto) (*domain.Brand, validators.ErrorResponse) {
	doc := domainBuilderAtCreate(dto)
	err := oRec.repo.Create(doc)
	if err != nil {
		return doc, validators.GetErrorResponseFromErr(err)
	}
	return doc, validators.ErrorResponse{}
}

func (oRec *BrandUseCase) Update(ctx *context.Context, dto *brand.UpdateBrandDto) validators.ErrorResponse {
	findBrand, findBrandErr := oRec.Find(ctx, dto.Id)
	if findBrandErr.IsError {
		return findBrandErr
	}
	doc := domainBuilderAtUpdate(dto, findBrand)
	if isAllowedToCascadeUpdates(findBrand, doc) {
		err := oRec.repo.UpdateBrandAndLocations(doc)
		if err != nil {
			return validators.GetErrorResponseFromErr(err)
		}
	} else {
		err := oRec.repo.Update(doc)
		if err != nil {
			return validators.GetErrorResponseFromErr(err)
		}
	}
	return validators.ErrorResponse{}
}

func (oRec *BrandUseCase) Find(ctx *context.Context, id string) (*domain.Brand, validators.ErrorResponse) {
	brand, err := oRec.repo.FindBrand(ctx, utils.ConvertStringIdToObjectId(id))
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}
	if brand == nil {
		return nil, validators.GetErrorResponseFromErr(errors.New(localization.E1002))
	}

	return brand, validators.ErrorResponse{}
}

func (oRec *BrandUseCase) List(ctx *context.Context, dto *brand.ListBrandDto) (*responses.ListResponse, validators.ErrorResponse) {
	brands, paginationMeta, resErr := oRec.repo.List(ctx, dto)
	if resErr != nil {
		return nil, validators.GetErrorResponseFromErr(resErr)
	}
	return responses.SetListResponse(brands, paginationMeta), validators.ErrorResponse{}
}

func (oRec *BrandUseCase) ChangeStatus(ctx *context.Context, dto *brand.ChangeBrandStatusDto) validators.ErrorResponse {
	brand, err := oRec.repo.FindBrand(ctx, dto.Id)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	doc := domainBuilderChangeStatus(dto, brand)
	err = oRec.repo.Update(doc)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *BrandUseCase) SoftDelete(ctx *context.Context, id string) validators.ErrorResponse {
	brand, err := oRec.repo.FindBrand(ctx, utils.ConvertStringIdToObjectId(id))
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	currentTime := time.Now()
	brand.DeletedAt = &currentTime
	err = oRec.repo.SoftDelete(brand)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *BrandUseCase) GetById(ctx *context.Context, id string) (*domain.Brand, validators.ErrorResponse) {
	brands, err := oRec.repo.GetByIds(ctx, &[]primitive.ObjectID{utils.ConvertStringIdToObjectId(id)})
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}
	brandsVal := *brands
	if len(brandsVal) == 0 {
		return nil, validators.GetErrorResponseFromErr(errors.New(localization.E1002))
	}

	return &brandsVal[0], validators.ErrorResponse{}
}
func (oRec *BrandUseCase) FindWithCuisines(ctx *context.Context, id string) (*domain.Brand, validators.ErrorResponse) {
	brandItem, err := oRec.repo.FindWithCuisines(ctx, utils.ConvertStringIdToObjectId(id))
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}
	if brandItem == nil {
		return nil, validators.GetErrorResponseFromErr(errors.New(localization.E1002))
	}

	return brandItem, validators.ErrorResponse{}
}
