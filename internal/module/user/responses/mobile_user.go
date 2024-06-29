package responses

var (
	OtpMessage = "OTP verified successfully, but you need to complete your profile."
	NextStep   = "completeProfile"
)

type VerifyOtpResp struct {
	Message  string `json:"message"`
	Token    string `json:"token"`
	NextStep string `json:"next_step"`
}
