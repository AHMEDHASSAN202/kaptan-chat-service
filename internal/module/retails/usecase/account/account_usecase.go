package account

import (
	"context"
	"errors"
	mongopagination "github.com/gobeam/mongo-go-pagination"
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
	brandUseCase    domain.BrandUseCase
	logger          logger.ILogger
}

func (l AccountUseCase) StoreAccount(ctx context.Context, payload *account.StoreAccountDto) (err validators.ErrorResponse) {
	accountDomain := CreateAccountBuilder(payload)
	errRe := l.repo.StoreAccount(ctx, &accountDomain)
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return validators.ErrorResponse{}
}

func (l AccountUseCase) UpdateAccount(ctx context.Context, id string, payload *account.UpdateAccountDto) (err validators.ErrorResponse) {
	accountDomain, errRe := l.repo.FindAccount(ctx, utils.ConvertStringIdToObjectId(id))
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	accountDomain.Name.Ar = payload.Name.Ar
	accountDomain.Name.En = payload.Name.En
	accountDomain.Email = payload.Email
	accountDomain.Country.Id = payload.Country.Id
	accountDomain.Country.PhonePrefix = payload.Country.PhonePrefix
	accountDomain.Country.Currency = payload.Country.Currency
	accountDomain.Country.Timezone = payload.Country.Timezone
	accountDomain.Country.Name.Ar = payload.Country.Name.Ar
	accountDomain.Country.Name.En = payload.Country.Name.En
	accountDomain.AllowedBrandIds = utils.ConvertStringIdsToObjectIds(payload.AllowedBrandIds)

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
	//if len(domainAccount.AllowedBrandIds) > 0 {
	//	dto := brand.ListBrandDto{
	//		Ids: utils.ConvertObjectIdsToStringIds(domainAccount.AllowedBrandIds),
	//	}
	//	brandsDomain, _, _ := l.brandUseCase.List(&ctx, &dto)
	//	domainAccount.Brands = *brandsDomain
	//}
	return *domainAccount, validators.ErrorResponse{}
}
func (l AccountUseCase) CheckAccountEmail(ctx context.Context, email string, accountId string) bool {
	return l.repo.CheckAccountEmail(ctx, email, accountId)
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

func (l AccountUseCase) ListAccount(ctx context.Context, payload *account.ListAccountDto) (accounts []domain.Account, paginationResult *mongopagination.PaginationData, err validators.ErrorResponse) {
	results, paginationResult, errRe := l.repo.ListAccount(ctx, payload)
	if errRe != nil {
		return results, paginationResult, validators.GetErrorResponseFromErr(errRe)
	}
	return results, paginationResult, validators.ErrorResponse{}

}

const tag = " AccountUseCase "

func NewAccountUseCase(repo domain.AccountRepository, brandUseCase domain.BrandUseCase, locationUseCase domain.LocationUseCase, logger logger.ILogger) domain.AccountUseCase {
	return &AccountUseCase{
		repo:            repo,
		brandUseCase:    brandUseCase,
		locationUseCase: locationUseCase,
		logger:          logger,
	}
}
