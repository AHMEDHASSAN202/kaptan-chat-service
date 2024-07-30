package approval

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"samm/internal/module/admin/responses"
	"samm/internal/module/approval/domain"
	dto "samm/internal/module/approval/dto"
	"samm/pkg/logger"
	"samm/pkg/utils"
	utilsDto "samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"time"
)

type ApprovalUseCase struct {
	approvalRepo domain.ApprovalRepository
	logger       logger.ILogger
}

func NewApprovalUseCase(repo domain.ApprovalRepository, logger logger.ILogger) domain.ApprovalUseCase {
	return &ApprovalUseCase{
		approvalRepo: repo,
		logger:       logger,
	}
}

func (oRec *ApprovalUseCase) List(ctx context.Context, input *dto.ListApprovalDto) (interface{}, validators.ErrorResponse) {
	list, pagination, err := oRec.approvalRepo.List(ctx, input)
	listResponse := responses.SetListResponse(list, pagination)
	if err != nil {
		oRec.logger.Error("ApprovalUseCase -> List -> ", err)
		return listResponse, validators.GetErrorResponse(&ctx, localization.E1000, nil, nil)
	}
	return listResponse, validators.ErrorResponse{}
}

func (oRec *ApprovalUseCase) FindByEntity(ctx context.Context, entityId primitive.ObjectID, entityType string) (interface{}, validators.ErrorResponse) {
	approval, err := oRec.approvalRepo.FindByEntity(ctx, entityId, entityType)
	if err != nil {
		oRec.logger.Error("ApprovalUseCase -> FindByEntity -> ", err)
		return approval, validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
	}
	return approval, validators.ErrorResponse{}
}

func (oRec *ApprovalUseCase) ChangeStatus(ctx context.Context, input *dto.ChangeStatusApprovalDto) validators.ErrorResponse {
	approval, err := oRec.approvalRepo.FindById(ctx, utils.ConvertStringIdToObjectId(input.Id))
	if err != nil {
		oRec.logger.Error("ApprovalUseCase -> FindById -> ", err)
		return validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
	}
	if approval.Status != utils.APPROVAL_STATUS.WAIT_FOR_APPROVAL {
		oRec.logger.Error("ApprovalUseCase -> FindById -> WAIT_FOR_APPROVAL -> ", err)
		return validators.GetErrorResponse(&ctx, localization.CanNotChangeApprovalStatus, nil, utils.GetAsPointer(http.StatusBadRequest))
	}
	approval.AdminDetails = utilsDto.AdminDetails{Id: utils.ConvertStringIdToObjectId(input.CauserId), Name: input.CauserName, Type: input.CauserType, Operation: "Change Approval Status", UpdatedAt: time.Now()}
	approval.Status = input.Status
	now := time.Now().UTC()
	switch input.Status {
	case utils.APPROVAL_STATUS.APPROVED:
		approval.Dates.ApprovedAt = &now
		approval.Dates.RejectedAt = nil
	case utils.APPROVAL_STATUS.REJECTED:
		approval.Dates.ApprovedAt = nil
		approval.Dates.RejectedAt = &now
	}
	err = oRec.approvalRepo.ChangeStatus(ctx, approval)
	if err != nil {
		return validators.GetErrorResponse(&ctx, localization.E1000, nil, utils.GetAsPointer(http.StatusBadRequest))
	}
	return validators.ErrorResponse{}
}
