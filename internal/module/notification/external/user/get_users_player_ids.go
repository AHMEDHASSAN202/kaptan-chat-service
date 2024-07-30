package user

import (
	"context"
	"samm/pkg/validators"
)

func (i IService) GetUsersPlayerIds(ctx context.Context, userIDs []string) (playerIDs []string, err validators.ErrorResponse) {
	return i.UserUseCase.GetUsersPlayerId(&ctx, userIDs)
}
