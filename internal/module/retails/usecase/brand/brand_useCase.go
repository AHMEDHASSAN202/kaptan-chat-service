package brand

import (
	"context"
	"github.com/jinzhu/copier"
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
	var brandDomain domain.Brand
	copier.Copy(&dto, &brandDomain)
	err := oRec.repo.Create(ctx, &brandDomain)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *BrandUseCase) Update(ctx *context.Context, dto *brand.UpdateBrandDto) validators.ErrorResponse {
	doc := convertDtoToCorrespondingDomain(dto)
	err := oRec.repo.Update(ctx, doc)
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

//func (oRec *BrandUseCase) ChangeStatus(ctx *context.Context, dto *cuisine.ChangeCuisineStatusDto) validators.ErrorResponse {
//	err := oRec.repo.ChangeStatus(ctx, dto)
//	if err != nil {
//		return validators.GetErrorResponseFromErr(err)
//	}
//	return validators.ErrorResponse{}
//}

func (oRec *BrandUseCase) List(ctx *context.Context, dto *brand.ListBrandDto) (*[]domain.Brand, validators.ErrorResponse) {
	cuisines, err := oRec.repo.List(ctx, dto)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}
	return cuisines, validators.ErrorResponse{}
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
