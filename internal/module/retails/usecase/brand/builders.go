package brand

import (
	"github.com/jinzhu/copier"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/brand"
	"samm/pkg/utils"
	"time"
)

func domainBuilderAtCreate(dto *brand.CreateBrandDto) *domain.Brand {
	brandDoc := domain.Brand{}
	copier.Copy(&brandDoc, dto)
	brandDoc.CuisineIds = utils.ConvertStringIdsToObjectIds(dto.CuisineIds)
	brandDoc.CreatedAt = time.Now()
	brandDoc.UpdatedAt = time.Now()
	return &brandDoc
}

func domainBuilderAtUpdate(dto *brand.UpdateBrandDto, domainData *domain.Brand) *domain.Brand {
	brandDoc := domain.Brand{}
	copier.Copy(&brandDoc, dto)
	brandDoc.CuisineIds = utils.ConvertStringIdsToObjectIds(dto.CuisineIds)
	brandDoc.UpdatedAt = time.Now()
	brandDoc.CreatedAt = domainData.CreatedAt
	brandDoc.ID = domainData.ID
	return &brandDoc
}

func domainBuilderChangeStatus(dto *brand.ChangeBrandStatusDto, domainData *domain.Brand) *domain.Brand {
	brandDoc := domain.Brand{}
	copier.Copy(&brandDoc, domainData)
	brandDoc.UpdatedAt = time.Now()
	brandDoc.ID = domainData.ID
	brandDoc.IsActive = dto.IsActive
	return &brandDoc
}
