package cuisine

import (
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/cuisine"
	"samm/pkg/utils"
	utilsDto "samm/pkg/utils/dto"
	"time"
)

func convertDtoArrToCorrespondingDomain(dto *cuisine.CreateCuisineDto) *domain.Cuisine {
	var cuisineDoc domain.Cuisine
	copier.Copy(&cuisineDoc, dto)
	cuisineDoc.ID = primitive.NewObjectID()
	cuisineDoc.AdminDetails = append(cuisineDoc.AdminDetails, utilsDto.AdminDetails{Id: utils.ConvertStringIdToObjectId(dto.CauserId), Name: dto.CauserName, Type: dto.CauserType, Operation: "Create Cuisine", UpdatedAt: time.Now()})
	return &cuisineDoc
}
func domainBuilderAtUpdate(dto *cuisine.UpdateCuisineDto, domainData *domain.Cuisine) *domain.Cuisine {
	var cuisineDoc domain.Cuisine
	copier.Copy(&cuisineDoc, dto)
	cuisineDoc.CreatedAt = domainData.CreatedAt
	cuisineDoc.ID = utils.ConvertStringIdToObjectId(dto.Id)
	cuisineDoc.AdminDetails = append(cuisineDoc.AdminDetails, utilsDto.AdminDetails{Id: utils.ConvertStringIdToObjectId(dto.CauserId), Name: dto.CauserName, Type: dto.CauserType, Operation: "Create Cuisine", UpdatedAt: time.Now()})
	return &cuisineDoc
}

func getCuisinesIds(cuisines *[]domain.Cuisine) (ids []string) {
	ids = make([]string, 0)
	for _, val := range *cuisines {
		ids = append(ids, utils.ConvertObjectIdToStringId(val.ID))
	}
	return
}
