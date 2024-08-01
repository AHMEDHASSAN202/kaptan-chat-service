package admin

type AuthKitchenResponse struct {
	Profile       interface{} `json:"profile"`
	Token         string      `json:"token"`
	FirebaseToken string      `json:"firebase_token"`
}
