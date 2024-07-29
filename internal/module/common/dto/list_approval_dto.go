package dto

import "samm/pkg/utils/dto"

type ListApprovalDto struct {
	dto.AdminHeaders
	dto.Pagination
	Type string `json:"type" validate:"required"`
}
