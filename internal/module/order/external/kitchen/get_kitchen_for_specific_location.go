package kitchen

import (
	"context"
	"samm/pkg/validators"
)

func (i IService) GetKitchensForSpecificLocation(ctx context.Context, locId, AccountId string) (kitchenIDs []string, err validators.ErrorResponse) {
	return i.KitchenUseCase.GetKitchensForSpecificLocation(ctx, locId, AccountId)
}
