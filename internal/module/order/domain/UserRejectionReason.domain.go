package domain

type UserRejectionReason struct {
	Id   string `json:"id"`
	Name struct {
		Ar string `json:"ar"`
		En string `json:"en"`
	} `json:"name"`
	Status string `json:"status"`
}
