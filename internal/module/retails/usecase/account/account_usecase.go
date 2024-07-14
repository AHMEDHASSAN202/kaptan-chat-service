package account

import (
	"context"
	"errors"
	mongopagination "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/admin/consts"
	domain2 "samm/internal/module/admin/domain"
	"samm/internal/module/admin/dto/admin"
	domain3 "samm/internal/module/kitchen/domain"
	"samm/internal/module/kitchen/dto/kitchen"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/account"
	"samm/internal/module/retails/dto/brand"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"time"
)

type AccountUseCase struct {
	repo            domain.AccountRepository
	locationUseCase domain.LocationUseCase
	brandUseCase    domain.BrandUseCase
	adminUseCase    domain2.AdminUseCase
	logger          logger.ILogger
	kitchenUseCase  domain3.KitchenUseCase
}

func (l AccountUseCase) StoreAccount(ctx context.Context, payload *account.StoreAccountDto) (err validators.ErrorResponse) {
	accountDomain := CreateAccountBuilder(payload)
	accountDomain.ID = primitive.NewObjectID()

	erre := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {

		errRe := l.repo.StoreAccount(sc, &accountDomain)
		if errRe != nil {
			return errRe
		}
		storeAdminDto := admin.CreateAdminDTO{
			ID:              primitive.NewObjectID(),
			Name:            payload.Name.En,
			Email:           payload.Email,
			Password:        payload.Password,
			ConfirmPassword: payload.ConfirmPassword,
			Status:          "active",
			Type:            consts.PORTAL_TYPE,
			RoleId:          consts.SUPER_PORTAL_ROLE,
			CountryIds:      []string{payload.Country.Id},
			Account:         &admin.Account{Id: utils.ConvertObjectIdToStringId(accountDomain.ID), Name: admin.Name{Ar: accountDomain.Name.Ar, En: accountDomain.Name.En}},
		}
		_, errR := l.adminUseCase.Create(ctx, &storeAdminDto)
		if errR.IsError {
			return errors.New(errR.ErrorMessageObject.Text)
		}
		return session.CommitTransaction(sc)
	})

	if erre != nil {
		return validators.GetErrorResponseFromErr(erre)
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
	accountDomain.Country.Id = payload.Country.Id
	accountDomain.Country.PhonePrefix = payload.Country.PhonePrefix
	accountDomain.Country.Currency = payload.Country.Currency
	accountDomain.Country.Timezone = payload.Country.Timezone
	accountDomain.Country.Name.Ar = payload.Country.Name.Ar
	accountDomain.Country.Name.En = payload.Country.Name.En
	accountDomain.AllowedBrandIds = utils.ConvertStringIdsToObjectIds(payload.AllowedBrandIds)
	accountDomain.UpdatedAt = time.Now()
	accountDomain.BankAccount.AccountNumber = payload.BankAccount.AccountNumber
	accountDomain.BankAccount.BankName = payload.BankAccount.BankName
	accountDomain.BankAccount.CompanyName = payload.BankAccount.CompanyName

	errRe = l.repo.UpdateAccount(ctx, accountDomain)
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	l.adminUseCase.SyncAccount(ctx, admin.Account{Id: utils.ConvertObjectIdToStringId(accountDomain.ID), Name: admin.Name{Ar: accountDomain.Name.Ar, En: accountDomain.Name.En}})
	return
}
func (l AccountUseCase) FindAccount(ctx context.Context, Id string) (account domain.Account, err validators.ErrorResponse) {
	domainAccount, errRe := l.repo.FindAccount(ctx, utils.ConvertStringIdToObjectId(Id))
	if errRe != nil {
		return *domainAccount, validators.GetErrorResponseFromErr(errRe)
	}
	if len(domainAccount.AllowedBrandIds) > 0 {
		dto := brand.ListBrandDto{
			Ids: utils.ConvertObjectIdsToStringIds(domainAccount.AllowedBrandIds),
		}
		brandsRes, _ := l.brandUseCase.List(&ctx, &dto)
		brands, ok := brandsRes.Docs.(*[]domain.Brand)
		domainAccount.Brands = make([]domain.Brand, 0)
		if ok {
			domainAccount.Brands = *brands
		}
	}
	return *domainAccount, validators.ErrorResponse{}
}
func (l AccountUseCase) CheckAccountEmail(ctx context.Context, email string, accountId string) bool {
	return l.repo.CheckAccountEmail(ctx, email, accountId)
}
func (l AccountUseCase) DeleteAccount(ctx context.Context, Id string) (err validators.ErrorResponse) {

	// Check if it has any kitchen
	listKitchen := kitchen.ListKitchenDto{
		AccountId: Id,
	}
	kitchens := l.kitchenUseCase.KitchenExists(&ctx, &listKitchen)
	if kitchens {
		return validators.GetErrorResponse(&ctx, localization.Account_Used_in_kitchen, nil, nil)
	}

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

func (l AccountUseCase) ListAccount(ctx context.Context, payload *account.ListAccountDto) (accounts []domain.Account, paginationResult mongopagination.PaginationData, err validators.ErrorResponse) {
	results, paginationResult, errRe := l.repo.ListAccount(ctx, payload)
	if errRe != nil {
		return results, paginationResult, validators.GetErrorResponseFromErr(errRe)
	}
	return results, paginationResult, validators.ErrorResponse{}

}

const tag = " AccountUseCase "

func NewAccountUseCase(repo domain.AccountRepository, adminUseCase domain2.AdminUseCase, brandUseCase domain.BrandUseCase, locationUseCase domain.LocationUseCase, logger logger.ILogger, kitchenUseCase domain3.KitchenUseCase) domain.AccountUseCase {
	return &AccountUseCase{
		repo:            repo,
		brandUseCase:    brandUseCase,
		locationUseCase: locationUseCase,
		adminUseCase:    adminUseCase,
		logger:          logger,
		kitchenUseCase:  kitchenUseCase,
	}
}
