package user

import (
	"context"
	"samm/internal/module/user/domain"
	"samm/internal/module/user/dto/user"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
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

func (l UserUseCase) StoreUser(ctx context.Context, payload *user.CreateUserDto) (err validators.ErrorResponse) {
	userDomain := domain.User{}
	errRe := l.repo.StoreUser(ctx, &userDomain)
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return
}

func (l UserUseCase) UpdateUserProfile(ctx context.Context, payload *user.UpdateUserProfileDto) (err validators.ErrorResponse) {
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
func (l UserUseCase) FindUser(ctx context.Context, Id string) (user domain.User, err validators.ErrorResponse) {
	domainUser, errRe := l.repo.FindUser(ctx, utils.ConvertStringIdToObjectId(Id))
	if errRe != nil {
		return *domainUser, validators.GetErrorResponseFromErr(errRe)
	}
	return *domainUser, validators.ErrorResponse{}
}

func (l UserUseCase) DeleteUser(ctx context.Context, Id string) (err validators.ErrorResponse) {

	delErr := l.repo.DeleteUser(ctx, utils.ConvertStringIdToObjectId(Id))
	if delErr != nil {
		return validators.GetErrorResponseFromErr(delErr)
	}
	return validators.ErrorResponse{}
}

func (l UserUseCase) ListUser(ctx context.Context, payload *user.ListUserDto) (users []domain.User, paginationResult utils.PaginationResult, err validators.ErrorResponse) {
	results, paginationResult, errRe := l.repo.ListUser(ctx, payload)
	if errRe != nil {
		return results, paginationResult, validators.GetErrorResponseFromErr(errRe)
	}
	return results, paginationResult, validators.ErrorResponse{}

}
