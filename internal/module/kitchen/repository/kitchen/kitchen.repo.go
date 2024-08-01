package mongodb

import (
	"context"
	"errors"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"samm/internal/module/kitchen/domain"
	"samm/internal/module/kitchen/dto/kitchen"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
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
func (l KitchenRepository) FindKitchen(ctx context.Context, Id primitive.ObjectID, withLookups bool) (kitchen *domain.Kitchen, err error) {
	domainData := domain.Kitchen{}
	var filters []interface{}

	filters = append(filters, bson.M{"$match": bson.M{"deleted_at": nil, "_id": Id}})
	if withLookups {
		filters = append(filters, bson.M{"$lookup": bson.M{"from": "locations", "localField": "location_ids", "foreignField": "_id", "as": "locations"}})
		filters = append(filters, bson.M{"$lookup": bson.M{"from": "accounts", "localField": "account_ids", "foreignField": "_id", "as": "accounts"}})
	}

	exists, err := l.kitchenCollection.SimpleAggregateFirstWithCtx(ctx, &domainData, filters...)
	if !exists {
		return &domainData, errors.New("Not Found")
	}
	return &domainData, err
}

func (l KitchenRepository) DeleteKitchen(ctx context.Context, Id primitive.ObjectID, causer *dto.AdminDetails) (err error) {
	kitchenData, err := l.FindKitchen(ctx, Id, false)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	kitchenData.DeletedAt = &now
	kitchenData.UpdatedAt = now
	kitchenData.AdminDetails = append(kitchenData.AdminDetails, *causer)
	return l.UpdateKitchen(kitchenData)
}

func (l *KitchenRepository) List(ctx *context.Context, dto *kitchen.ListKitchenDto) (usersRes *[]domain.Kitchen, paginationMeta *PaginationData, err error) {
	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.D{{"deleted_at", nil}},
	}}}

	if dto.Query != "" {
		pattern := ".*" + dto.Query + ".*"
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"$or": []bson.M{{"name.ar": bson.M{"$regex": pattern, "$options": "i"}}, {"name.en": bson.M{"$regex": pattern, "$options": "i"}}}})
	}
	if dto.LocationId != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{
			"location_ids": utils.ConvertStringIdToObjectId(dto.LocationId),
		})
	}
	if dto.AccountId != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{
			"account_ids": utils.ConvertStringIdToObjectId(dto.AccountId),
		})
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
func (l KitchenRepository) GetKitchensPlayerId(ctx *context.Context, locationIds []string, accountIds []string) (playerIds []string, err error) {
	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.D{{"deleted_at", nil}},
		bson.M{"$or": []bson.M{
			{"location_ids": bson.M{"$in": utils.ConvertStringIdsToObjectIds(locationIds)}},
			{"account_ids": bson.M{"$in": utils.ConvertStringIdsToObjectIds(accountIds)}},
		}},
	}}}
	var users []domain.Kitchen
	err = l.kitchenCollection.SimpleAggregate(&users, matching)
	playerIds = make([]string, 0)

	if err != nil {
		return playerIds, err
	}

	for _, model := range users {
		playerIds = append(playerIds, model.PlayerIds...)
	}

	return playerIds, err

}

func (l *KitchenRepository) GetKitchensForSpecificLocation(ctx context.Context, locId, AccountId primitive.ObjectID) (kitchensResult []domain.Kitchen, errResp validators.ErrorResponse) {
	filter := bson.M{"$or": bson.A{bson.M{"account_ids": bson.M{"$in": bson.A{AccountId}}}, bson.M{"location_ids": bson.M{"$in": bson.A{locId}}}}}
	kitchensResult = make([]domain.Kitchen, 0)

	err := l.kitchenCollection.SimpleFindWithCtx(ctx, &kitchensResult, filter)
	if err != nil {
		l.logger.Error("GetKitchensForSpecificLocation=>Repo", err)
		return nil, validators.GetErrorResponseFromErr(err)
	}
	return kitchensResult, validators.ErrorResponse{}
}
