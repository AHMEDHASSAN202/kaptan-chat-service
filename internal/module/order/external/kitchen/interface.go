package kitchen

import (
	"context"
	"samm/pkg/validators"
)

type Interface interface {
	GetKitchensForSpecificLocation(ctx context.Context, locId, AccountId string) (kitchenIDs []string, err validators.ErrorResponse)
}
