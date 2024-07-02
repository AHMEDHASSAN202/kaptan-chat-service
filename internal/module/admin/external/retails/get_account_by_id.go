package retails

import (
	"context"
	"samm/internal/module/admin/external/retails/responses"
	"samm/pkg/validators"
)

func (i IService) GetAccountById(ctx context.Context, id string) (*responses.AccountByIdResp, validators.ErrorResponse) {
	//account, err := i.AccountUseCase.OnlyFindAccount(ctx, id)
	//if err.IsError {
	//	return &responses.AccountByIdResp{}, err
	//}
	//accountResp := responses.AccountByIdResp{}
	//if errCopy := copier.Copy(&accountResp, account); errCopy != nil {
	//	logger.Logger.Error(errCopy)
	//	return &accountResp, validators.GetErrorResponse(&ctx, localization.E1000, nil, nil)
	//}
	//return &accountResp, validators.ErrorResponse{}
	return &responses.AccountByIdResp{}, validators.ErrorResponse{}
}
