package mongodb

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"samm/internal/module/kitchen/domain"
	"samm/internal/module/kitchen/dto/kitchen"
	"samm/pkg/logger"
	"time"
)

type KitchenRepository struct {
	kitchenCollection *mgm.Collection
	logger            logger.ILogger
}

const mongoKitchenRepositoryTag = "KitchenMongoRepository"

func NewKitchenMongoRepository(dbs *mongo.Database, log logger.ILogger) domain.KitchenRepository {
	kitchenDbCollection := mgm.Coll(&domain.Kitchen{})

	return &KitchenRepository{
		kitchenCollection: kitchenDbCollection,
		logger:            log,
	}
}

func (l KitchenRepository) CreateKitchen(kitchen *domain.Kitchen) (err error) {
	err = l.kitchenCollection.Create(kitchen)
	if err != nil {
		return err
	}
	return nil

}

func (l KitchenRepository) UpdateKitchen(kitchen *domain.Kitchen) (err error) {
	upsert := true
	opts := options.UpdateOptions{Upsert: &upsert}
	err = l.kitchenCollection.Update(kitchen, &opts)
	return
}
func (l KitchenRepository) FindKitchen(ctx context.Context, Id primitive.ObjectID) (kitchen *domain.Kitchen, err error) {
	domainData := domain.Kitchen{}
	filter := bson.M{"deleted_at": nil, "_id": Id}
	err = l.kitchenCollection.FirstWithCtx(ctx, filter, &domainData)

	return &domainData, err
}

func (l KitchenRepository) DeleteKitchen(ctx context.Context, Id primitive.ObjectID) (err error) {
	kitchenData, err := l.FindKitchen(ctx, Id)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	kitchenData.DeletedAt = &now
	kitchenData.UpdatedAt = now
	return l.UpdateKitchen(kitchenData)
}

func (l *KitchenRepository) List(ctx *context.Context, dto *kitchen.ListKitchenDto) (usersRes *[]domain.Kitchen, paginationMeta *PaginationData, err error) {
	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.D{{"deleted_at", nil}},
	}}}

	if dto.Query != "" {
		pattern := ".*" + dto.Query + ".*"
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"$or": []bson.M{{"name": bson.M{"$regex": pattern, "$options": "i"}}, {"phone_number": bson.M{"$regex": pattern, "$options": "i"}}}})
	}

	data, err := New(l.kitchenCollection.Collection).Context(*ctx).Limit(dto.Limit).Page(dto.Page).Sort("created_at", -1).Aggregate(matching)

	if data == nil || data.Data == nil {
		return nil, nil, err
	}

	users := make([]domain.Kitchen, 0)
	for _, raw := range data.Data {
		model := domain.Kitchen{}
		err = bson.Unmarshal(raw, &model)
		if err != nil {
			l.logger.Error("kitchen Repo -> List -> ", err)
			break
		}
		users = append(users, model)
	}
	paginationMeta = &data.Pagination
	usersRes = &users

	return
}


