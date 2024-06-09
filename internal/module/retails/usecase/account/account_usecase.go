package account

import (
	"context"
	"errors"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/account"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"time"
)

type AccountUseCase struct {
	repo            domain.AccountRepository
	locationUseCase domain.LocationUseCase
	logger          logger.ILogger
}

func (l AccountUseCase) StoreAccount(ctx context.Context, payload *account.StoreAccountDto) (err validators.ErrorResponse) {
	accountDomain := domain.Account{}
	accountDomain.Name.Ar = payload.Name.Ar
	accountDomain.Name.En = payload.Name.En
	accountDomain.Email = payload.Email
	password, er := utils.HashPassword(payload.Password)
	if er != nil {
		return validators.GetErrorResponseFromErr(er)
	}
	accountDomain.Password = password
	accountDomain.CreatedAt = time.Now()
	accountDomain.UpdatedAt = time.Now()

	errRe := l.repo.StoreAccount(ctx, &accountDomain)
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return
}

func (l AccountUseCase) UpdateAccount(ctx context.Context, id string, payload *account.UpdateAccountDto) (err validators.ErrorResponse) {
	accountDomain, errRe := l.repo.FindAccount(ctx, utils.ConvertStringIdToObjectId(id))
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	accountDomain.Name.Ar = payload.Name.Ar
	accountDomain.Name.En = payload.Name.En
	accountDomain.Email = payload.Email

	if payload.Password != "" {
		password, er := utils.HashPassword(payload.Password)
		if er != nil {
			return validators.GetErrorResponseFromErr(er)
		}
		accountDomain.Password = password
	}
	accountDomain.UpdatedAt = time.Now()

	errRe = l.repo.UpdateAccount(ctx, accountDomain)
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return
}
func (l AccountUseCase) FindAccount(ctx context.Context, Id string) (account domain.Account, err validators.ErrorResponse) {
	domainAccount, errRe := l.repo.FindAccount(ctx, utils.ConvertStringIdToObjectId(Id))
	if errRe != nil {
		return *domainAccount, validators.GetErrorResponseFromErr(errRe)
	}
	return *domainAccount, validators.ErrorResponse{}
}

func (l AccountUseCase) DeleteAccount(ctx context.Context, Id string) (err validators.ErrorResponse) {

	erre := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		// Delete Account

		errRe := l.repo.DeleteAccount(sc, utils.ConvertStringIdToObjectId(Id))
		if errRe != nil {
			return errRe
		}
		// Delete Locations
		errRee := l.locationUseCase.DeleteLocationByAccountId(sc, Id)
		if errRee.IsError {
			return errors.New(errRee.ErrorMessageObject.Text)
		}
		return session.CommitTransaction(sc)
	})

	if erre != nil {
		return validators.GetErrorResponseFromErr(erre)
	}
	return validators.ErrorResponse{}
}

func (l AccountUseCase) ListAccount(ctx context.Context, payload *account.ListAccountDto) (accounts []domain.Account, paginationResult utils.PaginationResult, err validators.ErrorResponse) {
	results, paginationResult, errRe := l.repo.ListAccount(ctx, payload)
	if errRe != nil {
		return results, paginationResult, validators.GetErrorResponseFromErr(errRe)
	}
	return results, paginationResult, validators.ErrorResponse{}

}

const tag = " AccountUseCase "

func NewAccountUseCase(repo domain.AccountRepository, locationUseCase domain.LocationUseCase, logger logger.ILogger) domain.AccountUseCase {
	return &AccountUseCase{
		repo:            repo,
		locationUseCase: locationUseCase,
		logger:          logger,
	}
}
