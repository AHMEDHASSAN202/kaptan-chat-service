package approval

import (
	"samm/internal/module/approval/domain"
	"samm/internal/module/approval/dto"
	"samm/pkg/utils"
	"time"
)

func CreateApprovalBuilder(dto *dto.CreateApprovalDto) *domain.Approval {
	approval := &domain.Approval{
		EntityId:     dto.EntityId,
		EntityType:   dto.EntityType,
		CountryId:    dto.CountryId,
		Status:       utils.APPROVAL_STATUS.WAIT_FOR_APPROVAL,
		Fields:       domain.Fields{New: dto.New, Old: dto.Old},
		Dates:        domain.Dates{ApprovedAt: nil, RejectedAt: nil},
		AdminDetails: dto.AdminDetails,
	}
	approval.UpdatedAt = time.Now().UTC()
	return approval
}
