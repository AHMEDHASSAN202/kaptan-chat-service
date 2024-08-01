package retails

import (
	"context"
	"github.com/jinzhu/copier"
	"samm/internal/module/menu/external/retails/responses"
	"samm/pkg/logger"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

func (i IService) GetAccountById(ctx context.Context, id string) (*responses.AccountByIdResp, validators.ErrorResponse) {
	account, err := i.AccountUseCase.OnlyFindAccount(ctx, id)
	if err.IsError {
		return &responses.AccountByIdResp{}, err
	}
	accountResp := responses.AccountByIdResp{}
	if errCopy := copier.Copy(&accountResp, account); errCopy != nil {
		logger.Logger.Error(errCopy)
		return &accountResp, validators.GetErrorResponse(&ctx, localization.E1000, nil, nil)
	}
	return &accountResp, validators.ErrorResponse{}
}
