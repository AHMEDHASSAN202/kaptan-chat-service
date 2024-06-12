package responses

import mongopagination "github.com/gobeam/mongo-go-pagination"

type ListResponse struct {
	Docs interface{}                     `json:"docs"`
	Meta *mongopagination.PaginationData `json:"meta"`
}

func SetListResponse(docs interface{}, meta *mongopagination.PaginationData) *ListResponse {
	listResponse := ListResponse{
		Docs: docs,
		Meta: meta,
	}
	if listResponse.Meta == nil {
		listResponse.Meta = &mongopagination.PaginationData{}
	}
	if docs == nil {
		listResponse.Docs = make([]interface{}, 0)
	}
	return &listResponse
}
