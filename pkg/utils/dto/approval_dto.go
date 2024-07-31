package dto

type ApprovalData struct {
	ApprovalStatus string `json:"approval_status" bson:"approval_status"`
	HasOriginal    bool   `json:"has_original" bson:"has_original"`
}
