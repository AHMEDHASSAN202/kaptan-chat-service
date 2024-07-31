package kitchen

import "samm/internal/module/kitchen/domain"

type IService struct {
	KitchenUseCase domain.KitchenUseCase
}
