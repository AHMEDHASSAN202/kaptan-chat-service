package requests

type InitSessionRequest struct {
	CustomerIdentifier string `json:"CustomerIdentifier"`
	SaveToken          bool   `json:"SaveToken"`
}
