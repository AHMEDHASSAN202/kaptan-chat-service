package retails

import (
	"context"
	"github.com/jinzhu/copier"
	"samm/internal/module/order/external/retails/responses"
	"samm/internal/module/retails/dto/location"
	"samm/pkg/validators"
)

func (i IService) GetLocationDetails(ctx context.Context, id string) (responses.LocationDetails, validators.ErrorResponse) {

	input := location.FindLocationMobileDto{}
	branch, err := i.LocationUseCase.FindMobileLocation(ctx, id, &input)
	if err.IsError {
		return responses.LocationDetails{}, err
	}
	resp := responses.LocationDetails{}
	copier.Copy(&resp, &branch)

	return resp, validators.ErrorResponse{}
}
