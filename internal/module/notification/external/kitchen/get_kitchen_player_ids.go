package kitchen

import (
	"context"
	"samm/pkg/validators"
)

func (i IService) GetKitchenPlayerIds(ctx context.Context, locationIds []string, accountIds []string) (playerIDs []string, err validators.ErrorResponse) {
	return i.KitchenUseCase.GetKitchensPlayerId(&ctx, locationIds, accountIds)
}
