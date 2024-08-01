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
		Doc:          domain.Doc{ID: dto.Doc.ID, Name: domain.Name{Ar: dto.Doc.Name.Ar, En: dto.Doc.Name.En}, Image: dto.Doc.Image},
		Account:      domain.Account{ID: dto.Account.ID, Name: domain.Name{Ar: dto.Account.Name.Ar, En: dto.Account.Name.En}},
	}
	approval.UpdatedAt = time.Now().UTC()
	return approval
}
