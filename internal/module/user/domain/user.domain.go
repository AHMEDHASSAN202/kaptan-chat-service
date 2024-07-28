package domain

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/user/dto/user"
	"samm/internal/module/user/responses"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
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
	Tokens           []string           `json:"-" bson:"tokens"`
	PlayerIds        []string           `json:"-" bson:"player_ids"`
	AdminDetails     []dto.AdminDetails `json:"admin_details" bson:"admin_details,omitempty"`
}

type DeletedUser struct {
	User `bson:",inline"`
}

type UserUseCase interface {
	StoreUser(ctx *context.Context, payload *user.CreateUserDto) (err validators.ErrorResponse)
	SendOtp(ctx *context.Context, payload *user.SendUserOtpDto) (err validators.ErrorResponse, tempOtp string)
	VerifyOtp(ctx *context.Context, payload *user.VerifyUserOtpDto) (res responses.VerifyOtpResp, err validators.ErrorResponse)
	UserSignUp(ctx *context.Context, payload *user.UserSignUpDto) (res responses.VerifyOtpResp, err validators.ErrorResponse)
	UpdateUserProfile(ctx *context.Context, payload *user.UpdateUserProfileDto) (user *responses.MobileUser, err validators.ErrorResponse)
	FindUser(ctx *context.Context, Id string) (user User, err validators.ErrorResponse)
	DeleteUser(ctx *context.Context, Id string) (err validators.ErrorResponse)
	List(ctx *context.Context, dto *user.ListUserDto) (*responses.ListResponse, validators.ErrorResponse)
	ToggleUserActivation(ctx *context.Context, userId string, adminHeader *dto.AdminHeaders) (err validators.ErrorResponse)
	UserEmailExists(ctx *context.Context, email, userId string) bool
	UpdateUserPlayerId(ctx *context.Context, payload *user.UpdateUserPlayerId) (user *responses.MobileUser, err validators.ErrorResponse)
}

type UserRepository interface {
	StoreUser(ctx *context.Context, user *User) (err error)
	InsertDeletedUser(ctx *context.Context, user *DeletedUser) (err error)
	UpdateUser(ctx *context.Context, user *User) (err error)
	FindUser(ctx *context.Context, Id primitive.ObjectID) (user *User, err error)
	GetUserByPhoneNumber(ctx *context.Context, phoneNum, countryCode string) (user User, err error)
	RemoveDeletedUser(user *DeletedUser) (err error)
	FindByToken(ctx context.Context, token string) (domainData *User, err error)
	List(ctx *context.Context, dto *user.ListUserDto) (usersRes *[]User, paginationMeta *PaginationData, err error)
	UserEmailExists(ctx *context.Context, email, userId string) bool
}
