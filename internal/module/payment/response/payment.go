package response

type PayResponse struct {
	Transaction interface{} `json:"transaction"`
	RedirectUrl *string     `json:"redirect_url"`
}
