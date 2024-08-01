package retails

import (
	"context"
	"samm/internal/module/menu/external/retails/responses"
	"samm/pkg/validators"
)

type Interface interface {
	GetBranchesByIds(ctx context.Context, ids []string) ([]responses.BranchByIdResp, validators.ErrorResponse)
	GetAccountById(ctx context.Context, id string) (*responses.AccountByIdResp, validators.ErrorResponse)
}
