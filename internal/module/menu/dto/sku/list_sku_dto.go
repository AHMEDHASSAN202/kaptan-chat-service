package sku

import "samm/pkg/utils/dto"

type ListSKUDto struct {
	dto.Pagination
	Query string `json:"query" form:"query" query:"query"`
}
