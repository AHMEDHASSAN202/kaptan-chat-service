package kitchen

import (
	"context"
	"samm/pkg/validators"
)

func (i IService) GetKitchenPlayerIds(ctx context.Context, kitchenIds []string) (playerIDs []string, err validators.ErrorResponse) {
	return i.KitchenUseCase.GetKitchensPlayerId(&ctx, kitchenIds)
}
