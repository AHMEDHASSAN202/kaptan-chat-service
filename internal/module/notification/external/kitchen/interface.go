package kitchen

import (
	"context"
	"samm/pkg/validators"
)

type Interface interface {
	GetKitchenPlayerIds(ctx context.Context, kitchenIds []string) (playerIDs []string, err validators.ErrorResponse)
}
