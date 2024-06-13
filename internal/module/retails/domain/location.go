package domain

import (
	"context"
	mongopagination "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/retails/dto/location"
	"samm/pkg/validators"
	"time"
)

type Country struct {
	Id   string `json:"_id" bson:"_id"`
	Name struct {
		Ar string `json:"ar" bson:"ar"`
		En string `json:"en" bson:"en"`
	} `json:"name" bson:"name"`
	Timezone    string `json:"timezone" bson:"timezone"`
	Currency    string `json:"currency" bson:"currency"`
	PhonePrefix string `json:"phone_prefix" bson:"phone_prefix"`
}

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
	Day       string `json:"day" bson:"day"`
	From      string `json:"from" bson:"from"`
	IsFullDay bool   `json:"is_full_day" bson:"is_full_day"`
	To        string `json:"to" bson:"to"`
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
type PercentsDate struct {
	From    time.Time `json:"from" bson:"from"`
	To      time.Time `json:"to" bson:"to"`
	Percent float64   ` json:"percent" bson:"percent"`
}
type Location struct {
	mgm.DefaultModel   `bson:",inline"`
	Name               Name          `json:"name" bson:"name"`
	City               City          `json:"city" bson:"city"`
	Street             Name          `json:"street" bson:"street"`
	Tags               []string      `json:"tags" bson:"tags"`
	CoverImage         string        `json:"cover_image" bson:"cover_image"`
	Logo               string        `json:"logo" bson:"logo"`
	Open               bool          `json:"open" bson:"open"`
	Phone              string        `json:"phone" bson:"phone"`
	BranchSignature    string        `json:"branch_signature" bson:"branch_signature"`
	Coordinate         Coordinate    `json:"coordinate" bson:"coordinate"`
	Index              string        `json:"index" bson:"index"`
	BrandDetails       BrandDetails  `json:"brand_details" bson:"brand_details"`
	WorkingHour        []WorkingHour `json:"working_hour" bson:"working_hour"`
	WorkingHourRamadan []WorkingHour `json:"working_hour_ramadan" bson:"working_hour_ramadan"`
	WorkingHourEid     []WorkingHour `json:"working_hour_eid" bson:"working_hour_eid"`
	//ActiveWorkingHour  string             `json:"active_working_hour" bson:"active_working_hour"`
	PreparationTime int                `json:"preparation_time" bson:"preparation_time"`
	AutoAccept      bool               `json:"auto_accept" bson:"auto_accept"`
	Status          string             `json:"status" bson:"status"`
	Percent         float64            ` json:"percent" bson:"percent"`
	PercentsDate    []PercentsDate     `json:"percents_date" bson:"percents_date"`
	SnoozeTo        *time.Time         `json:"snooze_to" bson:"snooze_to"`
	BankAccount     BankAccount        `json:"bank_account" bson:"bank_account"`
	Country         Country            `json:"country" bson:"country"`
	AccountId       primitive.ObjectID `json:"account_id" bson:"account_id"`
	AdminDetails    []AdminDetail      `json:"admin_details" bson:"admin_details"`
	DeletedAt       *time.Time         `json:"-" bson:"deleted_at"`
}

type LocationUseCase interface {
	StoreLocation(ctx context.Context, payload *location.StoreLocationDto) (err validators.ErrorResponse)
	BulkStoreLocation(ctx context.Context, payload location.StoreBulkLocationDto) (err validators.ErrorResponse)
	UpdateLocation(ctx context.Context, id string, payload *location.StoreLocationDto) (err validators.ErrorResponse)
	ToggleLocationStatus(ctx context.Context, id string) (err validators.ErrorResponse)
	FindLocation(ctx context.Context, Id string) (location Location, err validators.ErrorResponse)
	ToggleSnooze(ctx context.Context, dto *location.LocationToggleSnoozeDto) validators.ErrorResponse

	DeleteLocation(ctx context.Context, Id string) (err validators.ErrorResponse)
	DeleteLocationByAccountId(ctx context.Context, accountId string) (err validators.ErrorResponse)
	ListLocation(ctx context.Context, payload *location.ListLocationDto) (locations []Location, paginationResult *mongopagination.PaginationData, err validators.ErrorResponse)

	ListMobileLocation(ctx context.Context, payload *location.ListLocationMobileDto) (locations []LocationMobile, paginationResult *mongopagination.PaginationData, err validators.ErrorResponse)
	FindMobileLocation(ctx context.Context, Id string, payload *location.FindLocationMobileDto) (location LocationMobile, err validators.ErrorResponse)
}

type LocationRepository interface {
	StoreLocation(ctx context.Context, location *Location) (err error)
	BulkStoreLocation(ctx context.Context, data []Location) (err error)
	UpdateLocation(ctx context.Context, location *Location) (err error)
	FindLocation(ctx context.Context, Id primitive.ObjectID) (location *Location, err error)
	DeleteLocation(ctx context.Context, Id primitive.ObjectID) (err error)
	DeleteLocationByAccountId(ctx context.Context, accountId primitive.ObjectID) (err error)
	ListLocation(ctx context.Context, payload *location.ListLocationDto) (locations []Location, paginationResult *mongopagination.PaginationData, err error)
	UpdateBulkByBrand(ctx context.Context, brand BrandDetails) error
	SoftDeleteBulkByBrandId(ctx context.Context, brandId primitive.ObjectID) error

	ListMobileLocation(ctx context.Context, payload *location.ListLocationMobileDto) (locations []LocationMobile, paginationResult *mongopagination.PaginationData, err error)
	FindMobileLocation(ctx context.Context, Id primitive.ObjectID, payload *location.FindLocationMobileDto) (location *LocationMobile, err error)
}
