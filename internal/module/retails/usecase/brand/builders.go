package brand

import (
	"github.com/jinzhu/copier"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/brand"
	"samm/pkg/utils"
)

func domainBuilderAtCreate(dto *brand.CreateBrandDto) *domain.Brand {
	brandDoc := domain.Brand{}
	copier.Copy(&brandDoc, dto)
	brandDoc.CuisineIds = utils.ConvertStringIdsToObjectIds(dto.CuisineIds)
	return &brandDoc
}

func domainBuilderAtUpdate(dto *brand.UpdateBrandDto, domainData *domain.Brand) *domain.Brand {
	brandDoc := domain.Brand{}
	copier.Copy(&brandDoc, dto)
	brandDoc.CuisineIds = utils.ConvertStringIdsToObjectIds(dto.CuisineIds)
	return &brandDoc
}

func domainBuilderChangeStatus(dto *brand.ChangeBrandStatusDto, domainData *domain.Brand) *domain.Brand {
	brandDoc := domain.Brand{}
	copier.Copy(&brandDoc, domainData)
	brandDoc.IsActive = dto.IsActive
	return &brandDoc
}
