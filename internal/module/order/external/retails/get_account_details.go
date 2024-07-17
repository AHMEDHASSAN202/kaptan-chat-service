package retails

import (
	"context"
	"github.com/jinzhu/copier"
	"samm/internal/module/order/external/retails/responses"
	"samm/pkg/validators"
)

func (i IService) GetAccountDetails(ctx context.Context, id string) (accountDetails responses.AccountDetails, err validators.ErrorResponse) {

	account, err := i.AccountUseCase.FindAccount(ctx, id)
	if err.IsError {
		return accountDetails, err
	}
	resp := responses.AccountDetails{}
	copier.Copy(&resp, &account)

	return resp, validators.ErrorResponse{}
}
