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

func domainBuilderToggleSnooze(dto *brand.BrandToggleSnoozeDto, domainData *domain.Brand) *domain.Brand {
	var snoozedTill string
	if dto.IsSnooze && dto.SnoozeMinutesInterval > 0 {
		snoozedTill = time.Now().Add(time.Duration(dto.SnoozeMinutesInterval) * time.Minute).Format("2006-01-02 15:04:05")
	}
	brandDoc := domain.Brand{}
	copier.Copy(&brandDoc, domainData)
	brandDoc.SnoozedTill = snoozedTill
	return &brandDoc
}
