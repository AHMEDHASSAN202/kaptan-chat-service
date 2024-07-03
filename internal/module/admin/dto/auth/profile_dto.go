package auth

type ProfileDTO struct {
	AdminId       string
	AccountId     string
	CauserDetails map[string]interface{}
}
