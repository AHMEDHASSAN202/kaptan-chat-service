package external

import (
	domain2 "samm/internal/module/kitchen/domain"
	"samm/internal/module/notification/external/kitchen"
	"samm/internal/module/notification/external/user"
	"samm/internal/module/user/domain"
)

type ExtService struct {
	UserService    user.IService
	KitchenService kitchen.IService
}

func NewExternalService() *ExtService {
	return &ExtService{}
}
func (e *ExtService) SetUserUseCase(userUseCase domain.UserUseCase) {
	e.UserService = user.IService{UserUseCase: userUseCase}
}
func (e *ExtService) SetKitchenUseCase(kitchenUseCase domain2.KitchenUseCase) {
	e.KitchenService = kitchen.IService{KitchenUseCase: kitchenUseCase}
}
