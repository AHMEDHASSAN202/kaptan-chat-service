package cuisine

import (
	"github.com/jinzhu/copier"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/cuisine"
	"time"
)

func convertDtoArrToCorrespondingDomain(dto *[]cuisine.CreateCuisineDto) *[]domain.Cuisine {
	cuisineDocs := make([]domain.Cuisine, 0)
	copier.Copy(&cuisineDocs, dto)
	for i, _ := range cuisineDocs {
		cuisineDocs[i].CreatedAt = time.Now()
		cuisineDocs[i].UpdatedAt = time.Now()
	}
	return &cuisineDocs
}
func domainBuilderAtUpdate(dto *cuisine.UpdateCuisineDto, domainData *domain.Cuisine) *domain.Cuisine {
	cuisineDoc := domain.Cuisine{}
	copier.Copy(&cuisineDoc, domainData)
	cuisineDoc.UpdatedAt = time.Now()
	return &cuisineDoc
}
