package retails

import (
	"context"
	"samm/internal/module/admin/external/retails/responses"
	"samm/pkg/validators"
)

type Interface interface {
	GetAccountById(ctx context.Context, id string) (*responses.AccountByIdResp, validators.ErrorResponse)
}
