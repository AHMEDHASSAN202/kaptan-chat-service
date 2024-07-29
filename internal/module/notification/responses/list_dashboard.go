package responses

import (
	mongopagination "github.com/gobeam/mongo-go-pagination"
	"samm/pkg/utils"
)

type ListResponse struct {
	Docs interface{}                     `json:"docs" bson:"docs"`
	Meta *mongopagination.PaginationData `json:"meta" bson:"meta"`
}

func SetListResponse(docs interface{}, meta *mongopagination.PaginationData) *ListResponse {
	listResponse := ListResponse{
		Docs: docs,
		Meta: meta,
	}
	if listResponse.Meta == nil {
		listResponse.Meta = &mongopagination.PaginationData{}
	}
	if utils.IsNil(docs) {
		listResponse.Docs = make([]interface{}, 0)
	}
	return &listResponse
}
