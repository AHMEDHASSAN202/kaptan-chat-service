package brand

import (
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/brand"
	"samm/pkg/utils"
	utilsDto "samm/pkg/utils/dto"
	"time"
)

var LocationBrandAtt = []string{"name.ar", "name.en", "logo", "is_active"}

func domainBuilderAtCreate(dto *brand.CreateBrandDto) *domain.Brand {
	brandDoc := domain.Brand{}
	copier.Copy(&brandDoc, dto)
	brandDoc.ID = primitive.NewObjectID()
	brandDoc.CuisineIds = utils.ConvertStringIdsToObjectIds(dto.CuisineIds)
	brandDoc.AdminDetails = append(brandDoc.AdminDetails, utilsDto.AdminDetails{Id: utils.ConvertStringIdToObjectId(dto.CauserId), Name: dto.CauserName, Type: dto.CauserType, Operation: "Create Brand", UpdatedAt: time.Now()})
	return &brandDoc
}

func domainBuilderAtUpdate(dto *brand.UpdateBrandDto, domainData *domain.Brand) *domain.Brand {
	brandDoc := domain.Brand{}
	copier.Copy(&brandDoc, dto)
	brandDoc.ID = utils.ConvertStringIdToObjectId(dto.Id)
	brandDoc.CuisineIds = utils.ConvertStringIdsToObjectIds(dto.CuisineIds)
	brandDoc.AdminDetails = domainData.AdminDetails
	brandDoc.AdminDetails = append(brandDoc.AdminDetails, utilsDto.AdminDetails{Id: utils.ConvertStringIdToObjectId(dto.CauserId), Name: dto.CauserName, Type: dto.CauserType, Operation: "Update Brand", UpdatedAt: time.Now()})
	return &brandDoc
}

func domainBuilderChangeStatus(dto *brand.ChangeBrandStatusDto, domainData *domain.Brand) *domain.Brand {
	brandDoc := domain.Brand{}
	copier.Copy(&brandDoc, domainData)
	brandDoc.IsActive = dto.IsActive
	return &brandDoc
}

func isAllowedToCascadeUpdates(old *domain.Brand, new *domain.Brand) bool {

	return true

	//if old.Name.Ar != new.Name.Ar || old.Name.En != new.Name.En || old.Logo != new.Logo || old.IsActive != new.IsActive {
	//	return true
	//}
	//return false
}
