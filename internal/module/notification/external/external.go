package external

import (
	"samm/internal/module/notification/external/user"
	"samm/internal/module/user/domain"
)

type ExtService struct {
	UserService user.IService
}

func NewExternalService() *ExtService {
	return &ExtService{}
}
func (e *ExtService) SetUserUseCase(userUseCase domain.UserUseCase) {
	e.UserService = user.IService{UserUseCase: userUseCase}
}
