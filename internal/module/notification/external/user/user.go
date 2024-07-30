package user

import "samm/internal/module/user/domain"

type IService struct {
	UserUseCase domain.UserUseCase
}
