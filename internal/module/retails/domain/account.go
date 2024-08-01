package domain

import (
	"context"
	mongopagination "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/retails/dto/account"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"time"
)

type Account struct {
	mgm.DefaultModel `bson:",inline"`
	Name             Name `json:"name" bson:"name"`
	//Email            string               `json:"email" bson:"email"`
	//Password         string               `json:"-" bson:"password"`
	AllowedBrandIds []primitive.ObjectID `json:"allowed_brand_ids" bson:"allowed_brand_ids"`
	Country         Country              `json:"country" bson:"country"`
	Brands          []Brand              `json:"brands" bson:"-"`
	Percent         float64              ` json:"percent" bson:"percent"`
	BankAccount     BankAccount          `json:"bank_account" bson:"bank_account"`
	LocationsCount  int                  `json:"locations_count" bson:"locations_count,omitempty"`
	DeletedAt       *time.Time           `json:"-" bson:"deleted_at"`
	AdminDetails    []dto.AdminDetails   `json:"admin_details" bson:"admin_details,omitempty"`
}

type AccountUseCase interface {
	StoreAccount(ctx context.Context, payload *account.StoreAccountDto) (err validators.ErrorResponse)
	UpdateAccount(ctx context.Context, id string, payload *account.UpdateAccountDto) (err validators.ErrorResponse)
	FindAccount(ctx context.Context, Id string) (account Account, err validators.ErrorResponse)
	OnlyFindAccount(ctx context.Context, Id string) (account Account, err validators.ErrorResponse)
	CheckAccountEmail(ctx context.Context, email string, accountId string) bool
	DeleteAccount(ctx context.Context, Id string, adminDetails *dto.AdminHeaders) (err validators.ErrorResponse)
	ListAccount(ctx context.Context, payload *account.ListAccountDto) (accounts []Account, paginationResult mongopagination.PaginationData, err validators.ErrorResponse)
	CheckAccountExists(ctx context.Context, accountId string) (bool, validators.ErrorResponse)
}

type AccountRepository interface {
	StoreAccount(ctx context.Context, account *Account) (err error)
	UpdateAccount(ctx context.Context, account *Account) (err error)
	FindAccount(ctx context.Context, Id primitive.ObjectID) (account *Account, err error)
	CheckAccountEmail(ctx context.Context, email string, accountId string) bool
	DeleteAccount(ctx context.Context, Id primitive.ObjectID) (err error)
	ListAccount(ctx context.Context, payload *account.ListAccountDto) (locations []Account, paginationResult mongopagination.PaginationData, err error)
	CheckExists(ctx context.Context, id string) (bool, error)
}
