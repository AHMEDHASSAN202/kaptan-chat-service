package domain

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/notification/dto/notification"
	"samm/internal/module/notification/responses"
	Notification2 "samm/internal/module/notification/responses/Notification"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"time"
)

type Notification struct {
	mgm.DefaultModel `bson:",inline"`
	Title            Name                 `json:"title" bson:"title"`
	Image            string               `json:"image" bson:"image"`
	Text             Name                 `json:"text" bson:"text"`
	Type             string               `json:"type" bson:"type"`
	UserIds          []primitive.ObjectID `json:"-" bson:"user_ids"`
	RedirectType     string               `json:"redirect_type" bson:"redirect_type"`
	RedirectData     *RedirectData        `json:"redirect_data" bson:"redirect_data"`
	CountryId        string               `json:"country_id" bson:"country_id"`
	AdminDetails     dto.AdminDetails     `json:"admin_details" bson:"admin_details"`
	DeletedAt        *time.Time           `json:"deleted_at" bson:"deleted_at"`
	UsersCount       int                  `json:"users_count" bson:"users_count,omitempty"`
}

type Name struct {
	Ar string `json:"ar" bson:"ar"`
	En string `json:"en" bson:"en"`
}

type RedirectData struct {
	LocationId primitive.ObjectID `json:"location_id" bson:"location_id"`
}

type NotificationUseCase interface {
	CreateNotification(ctx context.Context, payload *notification.StoreNotificationDto) (err validators.ErrorResponse)
	FindNotification(ctx context.Context, Id string) (notification Notification, err validators.ErrorResponse)
	DeleteNotification(ctx context.Context, Id string) (err validators.ErrorResponse)
	List(ctx *context.Context, dto *notification.ListNotificationDto) (*responses.ListResponse, validators.ErrorResponse)
	ListMobile(ctx *context.Context, dto *notification.ListNotificationMobileDto) (*responses.ListResponse, validators.ErrorResponse)
}

type NotificationRepository interface {
	CreateNotification(notification *Notification) (err error)
	FindNotification(ctx context.Context, Id primitive.ObjectID) (notification *Notification, err error)
	DeleteNotification(ctx context.Context, Id primitive.ObjectID) (err error)
	List(ctx *context.Context, dto *notification.ListNotificationDto) (usersRes *[]Notification, paginationMeta *PaginationData, err error)
	ListMobile(ctx *context.Context, dto *notification.ListNotificationMobileDto) (usersRes *[]Notification2.NotificationMobile, paginationMeta *PaginationData, err error)
}
