package user

import (
	"context"
	"samm/internal/module/user/domain"
	"samm/internal/module/user/dto/user"
	"samm/internal/module/user/responses"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"time"
)

type UserUseCase struct {
	repo   domain.UserRepository
	logger logger.ILogger
}

const tag = " UserUseCase "

func NewUserUseCase(repo domain.UserRepository, logger logger.ILogger) domain.UserUseCase {
	return &UserUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (l UserUseCase) StoreUser(ctx *context.Context, payload *user.CreateUserDto) (err validators.ErrorResponse) {
	userDomain := domain.User{}
	errRe := l.repo.StoreUser(ctx, &userDomain)
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return
}

func (l UserUseCase) SendOtp(ctx *context.Context, payload *user.SendUserOtpDto) (err validators.ErrorResponse) {
	userDomain, _ := l.repo.GetUserByPhoneNumber(ctx, payload.PhoneNumber, payload.CountryCode)
	if userDomain.OtpCounter == 0 {
		return validators.GetErrorResponse(ctx, localization.E1015, nil, nil)
	}
	otp, otpErr := generateOTP()
	if otpErr != nil {
		err = validators.GetErrorResponseFromErr(otpErr)
		return
	}

	expiry := time.Now().Add(5 * time.Minute)
	userDomain.Otp = otp
	userDomain.ExpiryOtpDate = &expiry
	userDomain.PhoneNumber = payload.PhoneNumber
	userDomain.CountryCode = payload.CountryCode

	// send otp sms provider

	errRe := l.repo.UpdateUser(ctx, &userDomain)
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return
}

func (l UserUseCase) VerifyOtp(ctx *context.Context, payload *user.VerifyUserOtpDto) (err validators.ErrorResponse) {
	userDomain, dbErr := l.repo.GetUserByPhoneNumber(ctx, payload.PhoneNumber, payload.CountryCode)
	if dbErr != nil {
		return validators.GetErrorResponseFromErr(dbErr)
	}
	if userDomain.Otp != payload.Otp {
		return validators.GetErrorResponse(ctx, localization.E1013, nil, nil)
	}
	if userDomain.ExpiryOtpDate.Before(time.Now()) {
		return validators.GetErrorResponse(ctx, localization.E1014, nil, nil)
	}

	userDomain.PhoneNumber = payload.PhoneNumber
	userDomain.CountryCode = payload.CountryCode

	dbErr = l.repo.UpdateUser(ctx, &userDomain)
	if dbErr != nil {
		return validators.GetErrorResponseFromErr(dbErr)
	}
	return
}

func (l UserUseCase) UpdateUserProfile(ctx *context.Context, payload *user.UpdateUserProfileDto) (err validators.ErrorResponse) {
	userDomain, errRe := l.repo.FindUser(ctx, utils.ConvertStringIdToObjectId(payload.ID))
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	updatedUserDomain := domainBuilderAtUpdateProfile(payload, userDomain)
	errRe = l.repo.UpdateUser(ctx, updatedUserDomain)
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return
}
func (l UserUseCase) FindUser(ctx *context.Context, Id string) (user domain.User, err validators.ErrorResponse) {
	domainUser, errRe := l.repo.FindUser(ctx, utils.ConvertStringIdToObjectId(Id))
	if errRe != nil {
		return *domainUser, validators.GetErrorResponseFromErr(errRe)
	}
	return *domainUser, validators.ErrorResponse{}
}

func (l UserUseCase) DeleteUser(ctx *context.Context, Id string) (err validators.ErrorResponse) {

	delErr := l.repo.DeleteUser(ctx, utils.ConvertStringIdToObjectId(Id))
	if delErr != nil {
		return validators.GetErrorResponseFromErr(delErr)
	}
	return validators.ErrorResponse{}
}

func (oRec *UserUseCase) List(ctx *context.Context, dto *user.ListUserDto) (*responses.ListResponse, validators.ErrorResponse) {
	brands, paginationMeta, resErr := oRec.repo.List(ctx, dto)
	if resErr != nil {
		return nil, validators.GetErrorResponseFromErr(resErr)
	}
	return responses.SetListResponse(brands, paginationMeta), validators.ErrorResponse{}
}

func (l UserUseCase) ToggleUserActivation(ctx *context.Context, userId string) (err validators.ErrorResponse) {
	userDomain, errRe := l.repo.FindUser(ctx, utils.ConvertStringIdToObjectId(userId))
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	if userDomain.IsActive {
		userDomain.IsActive = false
	} else {
		userDomain.IsActive = true
	}
	errRe = l.repo.UpdateUser(ctx, userDomain)
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return
}

func (l UserUseCase) UserEmailExists(ctx *context.Context, email, userId string) bool {
	return l.repo.UserEmailExists(ctx, email, userId)
}
