package external

import (
	"samm/internal/module/admin/external/retails"
	"samm/internal/module/retails/usecase/account"
)

type ExtService struct {
	RetailsIService retails.IService
}

func NewExternalService(accountUseCase account.AccountUseCase) ExtService {
	return ExtService{
		RetailsIService: retails.IService{AccountUseCase: accountUseCase},
	}
}
