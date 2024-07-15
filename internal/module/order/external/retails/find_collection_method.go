package retails

import (
	"context"
	"github.com/jinzhu/copier"
	"samm/internal/module/order/external/retails/responses"
	"samm/pkg/validators"
)

func (i IService) FindCollectionMethod(ctx context.Context, id string, userId string) (responses.CollectionMethod, validators.ErrorResponse) {
	collectionMethod, err := i.CollectionUseCase.FindCollectionMethod(ctx, id, userId)
	if err.IsError {
		return responses.CollectionMethod{}, err
	}
	resp := responses.CollectionMethod{}
	copier.Copy(&resp, &collectionMethod)
	return resp, validators.ErrorResponse{}
}
