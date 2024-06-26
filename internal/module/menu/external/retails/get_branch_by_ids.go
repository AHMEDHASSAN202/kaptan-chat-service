package retails

import (
	"context"
	"github.com/jinzhu/copier"
	"samm/internal/module/menu/external/retails/responses"
	"samm/internal/module/retails/dto/location"
	"samm/pkg/logger"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

func (i IService) GetBranchesByIds(ctx context.Context, ids []string) ([]responses.BranchByIdResp, validators.ErrorResponse) {
	if ids == nil || len(ids) == 0 {
		return make([]responses.BranchByIdResp, 0), validators.ErrorResponse{}
	}

	input := location.ListLocationDto{Ids: ids, Pagination: dto.Pagination{Limit: int64(len(ids))}}
	branches, _, err := i.LocationUseCase.ListLocation(ctx, &input)
	if err.IsError {
		return make([]responses.BranchByIdResp, 0), err
	}

	branchResp := []responses.BranchByIdResp{}
	if errCopy := copier.Copy(&branchResp, branches); errCopy != nil {
		logger.Logger.Error(errCopy)
		return branchResp, validators.GetErrorResponse(&ctx, localization.E1000, nil, nil)
	}

	return branchResp, validators.ErrorResponse{}
}
