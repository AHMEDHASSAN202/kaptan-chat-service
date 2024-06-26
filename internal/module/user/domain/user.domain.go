package domain

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/user/dto/user"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"time"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string     `json:"name" bson:"name"`
	CountryCode      string     `json:"country_code" bson:"country_code"`
	PhoneNumber      string     `json:"phone_number" bson:"phone_number"`
	Email            string     `json:"email" bson:"email"`
	Gender           string     `json:"gender" bson:"gender"`
	Dob              string     `json:"dob" bson:"dob"`
	Otp              string     `json:"otp" bson:"otp"`
	ExpiryOtpDate    *time.Time `json:"expiry_otp_date" bson:"expiry_otp_date"`
	OtpCounter       int        `json:"otp_counter" bson:"otp_counter"`
	ImageURL         string     `json:"image_url" bson:"image_url"`
	Country          string     `json:"country" bson:"country"`
	IsActive         bool       `json:"is_active" bson:"is_active"`
	DeletedAt        *time.Time `json:"deleted_at" bson:"deleted_at"`
	Tokens           []string   `json:"tokens" bson:"tokens"`
}

type UserUseCase interface {
	StoreUser(ctx context.Context, payload *user.CreateUserDto) (err validators.ErrorResponse)
	UpdateUserProfile(ctx context.Context, payload *user.UpdateUserProfileDto) (err validators.ErrorResponse)
	FindUser(ctx context.Context, Id string) (user User, err validators.ErrorResponse)
	DeleteUser(ctx context.Context, Id string) (err validators.ErrorResponse)
	ListUser(ctx context.Context, payload *user.ListUserDto) (users []User, paginationResult utils.PaginationResult, err validators.ErrorResponse)
}

type UserRepository interface {
	StoreUser(ctx context.Context, user *User) (err error)
	UpdateUser(ctx context.Context, user *User) (err error)
	FindUser(ctx context.Context, Id primitive.ObjectID) (user *User, err error)
	DeleteUser(ctx context.Context, Id primitive.ObjectID) (err error)
	ListUser(ctx context.Context, payload *user.ListUserDto) (locations []User, paginationResult utils.PaginationResult, err error)
}
