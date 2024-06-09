package cuisine

import (
	"github.com/jinzhu/copier"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/cuisine"
)

func convertDtoArrToCorrespondingDomain(dto *cuisine.CreateCuisineDto) *domain.Cuisine {
	var cuisineDoc domain.Cuisine
	copier.Copy(&cuisineDoc, dto)
	return &cuisineDoc
}
func domainBuilderAtUpdate(dto *cuisine.UpdateCuisineDto, domainData *domain.Cuisine) *domain.Cuisine {
	var cuisineDoc domain.Cuisine
	copier.Copy(&cuisineDoc, dto)
	return &cuisineDoc
}
