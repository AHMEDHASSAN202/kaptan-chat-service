package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/user/domain"
	"samm/internal/module/user/dto/user"
	"samm/internal/module/user/responses"
	"samm/pkg/jwt"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"time"
)

type UserUseCase struct {
	repo           domain.UserRepository
	userJwtService jwt.JwtService
	logger         logger.ILogger
}

const tag = " UserUseCase "

func NewUserUseCase(repo domain.UserRepository, jwtFactory jwt.JwtServiceFactory, logger logger.ILogger) domain.UserUseCase {
	return &UserUseCase{
		repo:           repo,
		userJwtService: jwtFactory.UserJwtService(),
		logger:         logger,
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

func (l UserUseCase) SendOtp(ctx *context.Context, payload *user.SendUserOtpDto) (err validators.ErrorResponse, tempOtp string) {
	userDomain, dbErr := l.repo.GetUserByPhoneNumber(ctx, payload.PhoneNumber, payload.CountryCode)
	// new user
	if dbErr != nil {
		userDomain.ID = primitive.NewObjectID()
		userDomain.CreatedAt = time.Now()
		userDomain.UpdatedAt = time.Now()
	}

	newOtpCounter, ctrErr := otpTrialsPerDaySetter(userDomain.OtpCounter)
	if ctrErr != nil {
		return validators.GetErrorResponse(ctx, localization.E1015, nil, nil), ""
	}

	otp, otpErr := generateOTP()
	if otpErr != nil {
		err = validators.GetErrorResponseFromErr(otpErr)
		return
	}

	newUserDomain := domainBuilderAtCreateProfile(&userDomain, payload, otp, newOtpCounter)

	// send otp sms provider in-progress

	errRe := l.repo.UpdateUser(ctx, newUserDomain)
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe), ""
	}

	return validators.ErrorResponse{}, otp
}

func (l UserUseCase) VerifyOtp(ctx *context.Context, payload *user.VerifyUserOtpDto) (res responses.VerifyOtpResp, err validators.ErrorResponse) {
	userDomain, dbErr := l.repo.GetUserByPhoneNumber(ctx, payload.PhoneNumber, payload.CountryCode)
	if dbErr != nil {
		err = validators.GetErrorResponseFromErr(dbErr)
		return
	}
	if userDomain.Otp != payload.Otp {
		err = validators.GetErrorResponse(ctx, localization.E1013, nil, nil)
		return
	}
	if userDomain.ExpiryOtpDate.Before(time.Now()) {
		err = validators.GetErrorResponse(ctx, localization.E1014, nil, nil)
		return
	}

	// remove deletion if try to log in again within allowed period
	if userDomain.DeletedAt != nil {
		userDomain.DeletedAt = nil
		deletedUser := domain.DeletedUser{User: userDomain}
		dbErr = l.repo.RemoveDeletedUser(&deletedUser)
		if dbErr != nil {
			err = validators.GetErrorResponseFromErr(dbErr)
			return
		}
	}

	userToken, tokenErr := l.userJwtService.GenerateToken(*ctx, utils.ConvertObjectIdToStringId(userDomain.ID))
	if tokenErr != nil {
		err = validators.GetErrorResponseFromErr(tokenErr)
		return
	}

	if userDomain.Name == "" || userDomain.Email == "" {
		res = responses.VerifyOtpResp{
			IsProfileCompleted: false,
			Token:              userToken, //todo use temp token
		}
	} else {
		res = responses.VerifyOtpResp{
			IsProfileCompleted: true,
			Token:              userToken,
		}
		userDomain.Tokens = append(userDomain.Tokens, userToken)
	}

	dbErr = l.repo.UpdateUser(ctx, &userDomain)
	if dbErr != nil {
		err = validators.GetErrorResponseFromErr(dbErr)
		return
	}

	return
}

func (l UserUseCase) UserSignUp(ctx *context.Context, payload *user.UserSignUpDto) (res responses.VerifyOtpResp, err validators.ErrorResponse) {
	userDomain, dbErr := l.repo.GetUserByPhoneNumber(ctx, payload.PhoneNumber, payload.CountryCode)
	if dbErr != nil {
		err = validators.GetErrorResponseFromErr(dbErr)
		return
	}

	userToken, tokenErr := l.userJwtService.GenerateToken(*ctx, utils.ConvertObjectIdToStringId(userDomain.ID))
	if tokenErr != nil {
		err = validators.GetErrorResponseFromErr(tokenErr)
		return
	}

	updatedUserDomain := domainBuilderAtSignUp(payload, userToken, &userDomain)

	dbErr = l.repo.UpdateUser(ctx, updatedUserDomain)
	if dbErr != nil {
		err = validators.GetErrorResponseFromErr(dbErr)
		return
	}

	res = responses.VerifyOtpResp{
		Token: userToken,
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
	domainUser, dbErr := l.repo.FindUser(ctx, utils.ConvertStringIdToObjectId(Id))
	if dbErr != nil {
		err = validators.GetErrorResponseFromErr(dbErr)
		return
	}
	// Add 14 days to the current time
	deletedAt := time.Now().Add(14 * 24 * time.Hour)
	domainUser.DeletedAt = &deletedAt
	dbErr = l.repo.UpdateUser(ctx, domainUser)
	if dbErr != nil {
		err = validators.GetErrorResponseFromErr(dbErr)
		return
	}

	deletedUser := domain.DeletedUser{User: *domainUser}
	// store doc in user_coll_deleted_users
	errRe := l.repo.InsertDeletedUser(ctx, &deletedUser)
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
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
