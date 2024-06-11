package domain

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/retails/dto/account"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"time"
)

type Account struct {
	mgm.DefaultModel `bson:",inline"`
	Name             Name       `json:"name" bson:"name"`
	Email            string     `json:"email" bson:"email"`
	Password         string     `json:"-" bson:"password"`
	Country          Country    `json:"country" bson:"country"`
	DeletedAt        *time.Time `json:"-" bson:"deleted_at"`
}

type AccountUseCase interface {
	StoreAccount(ctx context.Context, payload *account.StoreAccountDto) (err validators.ErrorResponse)
	UpdateAccount(ctx context.Context, id string, payload *account.UpdateAccountDto) (err validators.ErrorResponse)
	FindAccount(ctx context.Context, Id string) (account Account, err validators.ErrorResponse)
	CheckAccountEmail(ctx context.Context, email string, accountId string) bool
	DeleteAccount(ctx context.Context, Id string) (err validators.ErrorResponse)
	ListAccount(ctx context.Context, payload *account.ListAccountDto) (accounts []Account, paginationResult utils.PaginationResult, err validators.ErrorResponse)
}

type AccountRepository interface {
	StoreAccount(ctx context.Context, account *Account) (err error)
	UpdateAccount(ctx context.Context, account *Account) (err error)
	FindAccount(ctx context.Context, Id primitive.ObjectID) (account *Account, err error)
	CheckAccountEmail(ctx context.Context, email string, accountId string) bool
	DeleteAccount(ctx context.Context, Id primitive.ObjectID) (err error)
	ListAccount(ctx context.Context, payload *account.ListAccountDto) (locations []Account, paginationResult utils.PaginationResult, err error)
}
