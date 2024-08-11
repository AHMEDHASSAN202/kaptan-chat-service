package requests

type UpdateSessionRequest struct {
	SessionId    string `json:"SessionId"`
	Token        string `json:"Token"`
	TokenType    string `json:"TokenType"`
	SecurityCode string `json:"SecurityCode,omitempty"`
}
