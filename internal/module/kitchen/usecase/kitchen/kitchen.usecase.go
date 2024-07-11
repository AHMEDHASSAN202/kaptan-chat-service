package kitchen

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/admin/consts"
	domain2 "samm/internal/module/admin/domain"
	"samm/internal/module/admin/dto/admin"
	"samm/internal/module/kitchen/domain"
	"samm/internal/module/kitchen/dto/kitchen"
	"samm/internal/module/kitchen/responses"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"time"
)

type KitchenUseCase struct {
	repo         domain.KitchenRepository
	logger       logger.ILogger
	adminUseCase domain2.AdminUseCase
}

const tag = " KitchenUseCase "

func NewKitchenUseCase(repo domain.KitchenRepository, logger logger.ILogger, adminUseCase domain2.AdminUseCase) domain.KitchenUseCase {
	return &KitchenUseCase{
		repo:         repo,
		logger:       logger,
		adminUseCase: adminUseCase,
	}
}

func (l KitchenUseCase) CreateKitchen(ctx context.Context, payload *kitchen.StoreKitchenDto) (err validators.ErrorResponse) {

	erre := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {

		kitchenDomain := domain.Kitchen{}
		kitchenDomain.Name.Ar = payload.Name.Ar
		kitchenDomain.Name.En = payload.Name.En
		kitchenDomain.AccountIds = utils.ConvertStringIdsToObjectIds(payload.AccountIds)
		kitchenDomain.LocationIds = utils.ConvertStringIdsToObjectIds(payload.LocationIds)
		kitchenDomain.CreatedAt = time.Now().UTC()
		kitchenDomain.UpdatedAt = time.Now().UTC()
		kitchenDomain.Country.Id = payload.Country.Id
		kitchenDomain.Country.PhonePrefix = payload.Country.PhonePrefix
		kitchenDomain.Country.Currency = payload.Country.Currency
		kitchenDomain.Country.Timezone = payload.Country.Timezone
		kitchenDomain.Country.Name.Ar = payload.Country.Name.Ar
		kitchenDomain.Country.Name.En = payload.Country.Name.En
		kitchenDomain.ID = primitive.NewObjectID()

		dbErr := l.repo.CreateKitchen(&kitchenDomain)
		if dbErr != nil {
			return dbErr
		}
		storeAdminDto := admin.CreateAdminDTO{
			ID:              primitive.NewObjectID(),
			Name:            payload.Name.En,
			Email:           payload.Email,
			Password:        payload.Password,
			ConfirmPassword: payload.ConfirmPassword,
			Status:          "active",
			Type:            consts.KITCHEN_TYPE,
			RoleId:          consts.SUPER_KITCHEN_ROLE,
			CountryIds:      []string{payload.Country.Id},
			Kitchen:         &admin.Kitchen{Id: utils.ConvertObjectIdToStringId(kitchenDomain.ID), Name: admin.Name{Ar: kitchenDomain.Name.Ar, En: kitchenDomain.Name.En}, AllowedStatus: payload.AllowedStatus},
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

func (l KitchenUseCase) UpdateKitchen(ctx context.Context, id string, payload *kitchen.UpdateKitchenDto) (err validators.ErrorResponse) {
	kitchenDomain, dbErr := l.repo.FindKitchen(ctx, utils.ConvertStringIdToObjectId(id))
	if dbErr != nil {
		return validators.GetErrorResponseFromErr(dbErr)
	}
	kitchenDomain.Name.Ar = payload.Name.Ar
	kitchenDomain.Name.En = payload.Name.En

	kitchenDomain.AccountIds = utils.ConvertStringIdsToObjectIds(payload.AccountIds)
	kitchenDomain.LocationIds = utils.ConvertStringIdsToObjectIds(payload.LocationIds)
	kitchenDomain.UpdatedAt = time.Now().UTC()
	kitchenDomain.Country.Id = payload.Country.Id
	kitchenDomain.Country.PhonePrefix = payload.Country.PhonePrefix
	kitchenDomain.Country.Currency = payload.Country.Currency
	kitchenDomain.Country.Timezone = payload.Country.Timezone
	kitchenDomain.Country.Name.Ar = payload.Country.Name.Ar
	kitchenDomain.Country.Name.En = payload.Country.Name.En

	dbErr = l.repo.UpdateKitchen(kitchenDomain)
	if dbErr != nil {
		return validators.GetErrorResponseFromErr(dbErr)
	}
	return
}
func (l KitchenUseCase) FindKitchen(ctx context.Context, Id string) (kitchen domain.Kitchen, err validators.ErrorResponse) {
	domainKitchen, dbErr := l.repo.FindKitchen(ctx, utils.ConvertStringIdToObjectId(Id))
	if dbErr != nil {
		return *domainKitchen, validators.GetErrorResponseFromErr(dbErr)
	}
	return *domainKitchen, validators.ErrorResponse{}
}

func (l KitchenUseCase) DeleteKitchen(ctx context.Context, Id string) (err validators.ErrorResponse) {

	delErr := l.repo.DeleteKitchen(ctx, utils.ConvertStringIdToObjectId(Id))
	if delErr != nil {
		return validators.GetErrorResponseFromErr(delErr)
	}
	return validators.ErrorResponse{}
}

func (l KitchenUseCase) List(ctx *context.Context, dto *kitchen.ListKitchenDto) (*responses.ListResponse, validators.ErrorResponse) {
	users, paginationMeta, resErr := l.repo.List(ctx, dto)
	if resErr != nil {
		return nil, validators.GetErrorResponseFromErr(resErr)
	}
	return responses.SetListResponse(users, paginationMeta), validators.ErrorResponse{}
}
