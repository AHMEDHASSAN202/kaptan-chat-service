package domain

import (
	"context"
	"github.com/kamva/mgm/v3"
	"kaptan/pkg/utils/dto"
	"time"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string             `json:"name" bson:"name"`
	CountryCode      string             `json:"country_code" bson:"country_code"`
	PhoneNumber      string             `json:"phone_number" bson:"phone_number"`
	Email            string             `json:"email" bson:"email"`
	Gender           string             `json:"gender" bson:"gender"`
	Dob              string             `json:"dob" bson:"dob"`
	Otp              string             `json:"otp" bson:"otp"`
	ExpiryOtpDate    *time.Time         `json:"expiry_otp_date" bson:"expiry_otp_date"`
	OtpCounter       string             `json:"otp_counter" bson:"otp_counter"`
	ImageURL         string             `json:"image_url" bson:"image_url"`
	Country          string             `json:"country" bson:"country"`
	IsActive         bool               `json:"is_active" bson:"is_active"`
	VerifiedAt       *time.Time         `json:"verified_at" bson:"verified_at"`
	DeletedAt        *time.Time         `json:"deleted_at" bson:"deleted_at"`
	DeleteReason     UserDeletionReason `json:"delete_reason" bson:"delete_reason"`
	Tokens           []string           `json:"-" bson:"tokens"`
	PlayerIds        []string           `json:"-" bson:"player_ids"`
	AdminDetails     []dto.AdminDetails `json:"admin_details" bson:"admin_details,omitempty"`
}

type UserDeletionReason struct {
	Id   string `json:"id"`
	Name struct {
		Ar string `json:"ar"`
		En string `json:"en"`
	} `json:"name"`
}

type DeletedUser struct {
	User `bson:",inline"`
}

type DriverRepository interface {
	FindByToken(ctx *context.Context, token string) (domainData *User, err error)
}
