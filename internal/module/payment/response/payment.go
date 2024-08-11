package response

type PayResponse struct {
	Transaction interface{} `json:"transaction"`
	RedirectUrl *string     `json:"redirect_url"`
}
type UpdateSessionResponse struct {
	RedirectUrl *string `json:"payment_url"`
}
