package retails

import (
	"context"
	"samm/internal/module/order/external/retails/responses"
	"samm/pkg/validators"
)

type Interface interface {
	GetLocationDetails(ctx context.Context, id string) (responses.LocationDetails, validators.ErrorResponse)
}
