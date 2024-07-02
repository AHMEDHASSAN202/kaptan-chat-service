package responses

import (
	"github.com/kamva/mgm/v3"
)

type VerifyOtpResp struct {
	Token              string `json:"token"`
	IsProfileCompleted bool   `json:"is_profile_completed"`
}

type MobileUser struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name" bson:"name"`
	CountryCode      string `json:"country_code" bson:"country_code"`
	PhoneNumber      string `json:"phone_number" bson:"phone_number"`
	Email            string `json:"email" bson:"email"`
	Gender           string `json:"gender" bson:"gender"`
	Dob              string `json:"dob" bson:"dob"`
	ImageURL         string `json:"image_url" bson:"image_url"`
	Country          string `json:"country" bson:"country"`
	IsActive         bool   `json:"is_active" bson:"is_active"`
}
