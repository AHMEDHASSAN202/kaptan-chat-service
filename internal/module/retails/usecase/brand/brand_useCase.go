package brand

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/brand"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
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

func (oRec *BrandUseCase) Create(ctx *context.Context, dto *brand.CreateBrandDto) validators.ErrorResponse {
	doc := domainBuilderAtCreate(dto)
	err := oRec.repo.Create(doc)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *BrandUseCase) Update(ctx *context.Context, dto *brand.UpdateBrandDto) validators.ErrorResponse {
	findBrand, findBrandErr := oRec.Find(ctx, dto.Id)
	if findBrandErr.IsError {
		return findBrandErr
	}
	doc := domainBuilderAtUpdate(dto, findBrand)
	err := oRec.repo.Update(doc)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
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

func (oRec *BrandUseCase) List(ctx *context.Context, dto *brand.ListBrandDto) (brands *[]domain.Brand, paginationMeta *PaginationData, err validators.ErrorResponse) {
	brands, paginationMeta, resErr := oRec.repo.List(ctx, dto)
	if resErr != nil {
		err = validators.GetErrorResponseFromErr(resErr)
		return
	}
	return
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

func (oRec *BrandUseCase) ToggleSnooze(ctx *context.Context, dto *brand.BrandToggleSnoozeDto) validators.ErrorResponse {
	brand, err := oRec.repo.FindBrand(ctx, dto.Id)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	doc := domainBuilderToggleSnooze(dto, brand)
	err = oRec.repo.Update(doc)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *BrandUseCase) SoftDelete(ctx *context.Context, id string) validators.ErrorResponse {
	idDoc := utils.ConvertStringIdToObjectId(id)
	err := oRec.repo.SoftDelete(ctx, idDoc)
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
