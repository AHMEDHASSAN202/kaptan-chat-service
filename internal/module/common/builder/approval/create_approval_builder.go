package approval

import (
	"samm/internal/module/common/domain"
	"samm/internal/module/common/dto"
	"samm/pkg/utils"
)

func CreateApprovalBuilder(dto *dto.CreateApprovalDto) *domain.Approval {
	return &domain.Approval{
		EntityId:     dto.EntityId,
		EntityType:   dto.EntityType,
		CountryId:    dto.CountryId,
		Status:       utils.APPROVAL_STATUS.WAIT_FOR_APPROVAL,
		Fields:       domain.Fields{New: dto.New, Old: dto.Old},
		Dates:        domain.Dates{ApprovedAt: nil, RejectedAt: nil},
		AdminDetails: dto.AdminDetails,
	}
}
