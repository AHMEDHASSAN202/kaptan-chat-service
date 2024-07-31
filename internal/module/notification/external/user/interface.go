package user

import (
	"context"
	"samm/pkg/validators"
)

type Interface interface {
	GetUsersPlayerIds(ctx context.Context, userIds []string) (playerIDs []string, err validators.ErrorResponse)
}
