package mongodb

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"samm/internal/module/notification/consts"
	"samm/internal/module/notification/domain"
	"samm/internal/module/notification/dto/notification"
	"samm/internal/module/notification/responses/Notification"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"time"
)

type NotificationRepository struct {
	notificationCollection *mgm.Collection
	logger                 logger.ILogger
}

const mongoNotificationRepositoryTag = "NotificationMongoRepository"

func NewNotificationMongoRepository(dbs *mongo.Database, log logger.ILogger) domain.NotificationRepository {
	notificationDbCollection := mgm.Coll(&domain.Notification{})

	return &NotificationRepository{
		notificationCollection: notificationDbCollection,
		logger:                 log,
	}
}

func (l NotificationRepository) CreateNotification(notification *domain.Notification) (err error) {
	err = l.notificationCollection.Create(notification)
	if err != nil {
		return err
	}
	return nil

}

func (l NotificationRepository) UpdateNotification(notification *domain.Notification) (err error) {
	upsert := true
	opts := options.UpdateOptions{Upsert: &upsert}
	err = l.notificationCollection.Update(notification, &opts)
	return
}
func (l NotificationRepository) FindNotification(ctx context.Context, Id primitive.ObjectID) (notification *domain.Notification, err error) {
	domainData := domain.Notification{}
	filter := bson.M{"deleted_at": nil, "_id": Id}
	err = l.notificationCollection.FirstWithCtx(ctx, filter, &domainData)

	return &domainData, err
}

func (l NotificationRepository) DeleteNotification(ctx context.Context, Id primitive.ObjectID) (err error) {
	notificationData, err := l.FindNotification(ctx, Id)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	notificationData.DeletedAt = &now
	notificationData.UpdatedAt = now
	return l.UpdateNotification(notificationData)
}

func (l *NotificationRepository) List(ctx *context.Context, dto *notification.ListNotificationDto) (usersRes *[]domain.Notification, paginationMeta *PaginationData, err error) {
	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.D{{"deleted_at", nil}},
		bson.D{{"country_id", dto.CountryId}},
	}}}
	var pipeline []interface{}
	if dto.Query != "" {
		pattern := ".*" + dto.Query + ".*"
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"$or": []bson.M{{"title.ar": bson.M{"$regex": pattern, "$options": "i"}}, {"title.en": bson.M{"$regex": pattern, "$options": "i"}}, {"text.ar": bson.M{"$regex": pattern, "$options": "i"}}, {"text.en": bson.M{"$regex": pattern, "$options": "i"}}}})
	}
	pipeline = append(pipeline, matching)
	pipeline = append(pipeline,
		bson.M{
			"$addFields": bson.M{
				"users_count": bson.M{"$size": "$user_ids"},
			},
		},
	)

	data, err := New(l.notificationCollection.Collection).Context(*ctx).Limit(dto.Limit).Page(dto.Page).Sort("created_at", -1).Aggregate(pipeline...)

	if data == nil || data.Data == nil {
		return nil, nil, err
	}

	items := make([]domain.Notification, 0)
	for _, raw := range data.Data {
		model := domain.Notification{}
		err = bson.Unmarshal(raw, &model)
		if err != nil {
			l.logger.Error("notification Repo -> List -> ", err)
			break
		}
		items = append(items, model)
	}
	paginationMeta = &data.Pagination
	usersRes = &items

	return
}
func (l *NotificationRepository) ListMobile(ctx *context.Context, dto *notification.ListNotificationMobileDto) (usersRes *[]Notification.NotificationMobile, paginationMeta *PaginationData, err error) {
	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.D{{"deleted_at", nil}},
		bson.D{{"country_id", dto.CountryId}},
		bson.D{{"$or", []bson.M{
			{"type": consts.TYPE_PUBLIC},
			{"user_ids": utils.ConvertStringIdToObjectId(dto.CauserId)},
		},
		}},
	}}}
	var pipeline []interface{}
	pipeline = append(pipeline, matching)

	data, err := New(l.notificationCollection.Collection).Context(*ctx).Limit(dto.Limit).Page(dto.Page).Sort("created_at", -1).Aggregate(pipeline...)

	if data == nil || data.Data == nil {
		return nil, nil, err
	}

	items := make([]Notification.NotificationMobile, 0)
	for _, raw := range data.Data {
		model := Notification.NotificationMobile{}
		err = bson.Unmarshal(raw, &model)
		if err != nil {
			l.logger.Error("notification Repo -> List -> ", err)
			break
		}
		items = append(items, model)
	}
	paginationMeta = &data.Pagination
	usersRes = &items

	return
}
