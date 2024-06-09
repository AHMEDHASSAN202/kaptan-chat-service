package domain

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/retails/dto/location"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"time"
)

type City struct {
	Id   primitive.ObjectID `json:"_id" bson:"id"`
	Name Name               `json:"name" bson:"name"`
}
type Coordinate struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}
type BrandDetails struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Name     Name               `json:"name" bson:"name"`
	Logo     string             `json:"logo" bson:"logo"`
	IsActive bool               `json:"is_active" bson:"is_active"`
}
type WorkingHour struct {
	Day  string `json:"day" bson:"day"`
	From string `json:"from" bson:"from"`
	To   string `json:"to" bson:"to"`
}
type BankAccount struct {
	AccountNumber string `json:"account_number" bson:"account_number"`
	BankName      string `json:"bank_name" bson:"bank_name"`
	CompanyName   string `json:"company_name" bson:"company_name"`
}
type AdminDetail struct {
	AdminId   primitive.ObjectID `json:"admin_id" bson:"admin_id"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Location struct {
	mgm.DefaultModel `bson:",inline"`
	Name             Name               `json:"name" bson:"name"`
	City             City               `json:"city" bson:"city"`
	Street           Name               `json:"street" bson:"street"`
	Tags             string             `json:"tags" bson:"tags"`
	CoverImage       string             `json:"cover_image" bson:"cover_image"`
	Open             bool               `json:"open" bson:"open"`
	Phone            string             `json:"phone" bson:"phone"`
	BranchSignature  string             `json:"branch_signature" bson:"branch_signature"`
	Coordinate       Coordinate         `json:"coordinate" bson:"coordinate"`
	Index            string             `json:"index" bson:"index"`
	BrandDetails     BrandDetails       `json:"brand_details" bson:"brand_details"`
	WorkingHour      []WorkingHour      `json:"working_hour" bson:"working_hour"`
	PreparationTime  int                `json:"preparation_time" bson:"preparation_time"`
	AutoAccept       bool               `json:"auto_accept" bson:"auto_accept"`
	Status           string             `json:"status" bson:"status"`
	SnoozeTo         *time.Time         `json:"snooze_to" bson:"snooze_to"`
	BankAccount      BankAccount        `json:"bank_account" bson:"bank_account"`
	AccountId        primitive.ObjectID `json:"account_id" bson:"account_id"`
	AdminDetails     []AdminDetail      `json:"admin_details" bson:"admin_details"`
	DeletedAt        *time.Time         `json:"-" bson:"deleted_at"`
}

type LocationUseCase interface {
	StoreLocation(ctx context.Context, payload *location.StoreLocationDto) (err validators.ErrorResponse)
	UpdateLocation(ctx context.Context, id string, payload *location.StoreLocationDto) (err validators.ErrorResponse)
	ToggleLocationStatus(ctx context.Context, id string) (err validators.ErrorResponse)
	FindLocation(ctx context.Context, Id string) (location Location, err validators.ErrorResponse)
	ToggleSnooze(ctx context.Context, dto *location.LocationToggleSnoozeDto) validators.ErrorResponse

	DeleteLocation(ctx context.Context, Id string) (err validators.ErrorResponse)
	DeleteLocationByAccountId(ctx context.Context, accountId string) (err validators.ErrorResponse)
	ListLocation(ctx context.Context, payload *location.ListLocationDto) (locations []Location, paginationResult utils.PaginationResult, err validators.ErrorResponse)
}

type LocationRepository interface {
	StoreLocation(ctx context.Context, location *Location) (err error)
	UpdateLocation(ctx context.Context, location *Location) (err error)
	FindLocation(ctx context.Context, Id primitive.ObjectID) (location *Location, err error)
	DeleteLocation(ctx context.Context, Id primitive.ObjectID) (err error)
	DeleteLocationByAccountId(ctx context.Context, accountId primitive.ObjectID) (err error)
	ListLocation(ctx context.Context, payload *location.ListLocationDto) (locations []Location, paginationResult utils.PaginationResult, err error)
	UpdateBulkByBrand(ctx context.Context, brand BrandDetails) error
	SoftDeleteBulkByBrandId(ctx context.Context, brandId primitive.ObjectID) error
}
