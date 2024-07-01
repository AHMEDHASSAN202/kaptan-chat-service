package responses

type VerifyOtpResp struct {
	Token              string `json:"token"`
	IsProfileCompleted bool   `json:"is_profile_completed"`
}
