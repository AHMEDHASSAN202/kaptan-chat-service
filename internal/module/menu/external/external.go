package external

import (
	"samm/internal/module/menu/external/retails"
	"samm/internal/module/retails/domain"
)

type ExtService struct {
	RetailsIService retails.IService
}

func NewExternalService(locationUseCase domain.LocationUseCase) ExtService {
	return ExtService{
		RetailsIService: retails.IService{LocationUseCase: locationUseCase},
	}
}
