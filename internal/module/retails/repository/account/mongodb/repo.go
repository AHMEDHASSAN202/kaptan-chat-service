package mongodb

import (
	"context"
	"fmt"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/account"
	"samm/pkg/utils"
	"time"
)

type AccountRepository struct {
	accountCollection *mgm.Collection
}

const mongoAccountRepositoryTag = "AccountMongoRepository"

func NewAccountMongoRepository(dbs *mongo.Database) domain.AccountRepository {
	accountDbCollection := mgm.Coll(&domain.Account{})

	return &AccountRepository{
		accountCollection: accountDbCollection,
	}
}

func (l AccountRepository) StoreAccount(ctx context.Context, account *domain.Account) (err error) {
	_, err = mgm.Coll(&domain.Account{}).InsertOne(ctx, account)
	if err != nil {
		return err
	}
	return nil

}

func (l AccountRepository) UpdateAccount(ctx context.Context, account *domain.Account) (err error) {
	update := bson.M{"$set": account}
	_, err = mgm.Coll(&domain.Account{}).UpdateByID(ctx, account.ID, update)
	return
}
func (l AccountRepository) FindAccount(ctx context.Context, Id primitive.ObjectID) (account *domain.Account, err error) {
	domainData := domain.Account{}
	filter := bson.M{"deleted_at": nil, "_id": Id}
	err = l.accountCollection.FirstWithCtx(ctx, filter, &domainData)

	return &domainData, err
}
func (l AccountRepository) CheckAccountEmail(ctx context.Context, email string, accountId string) bool {
	filter := bson.M{"deleted_at": nil, "email": email}
	if accountId != "" {
		filter = bson.M{"deleted_at": nil, "email": email, "_id": bson.M{"$ne": utils.ConvertStringIdToObjectId(accountId)}}
	}
	count, _ := l.accountCollection.CountDocuments(ctx, filter)
	return count > 0
}
func (l AccountRepository) DeleteAccount(ctx context.Context, Id primitive.ObjectID) (err error) {
	accountData, err := l.FindAccount(ctx, Id)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	accountData.DeletedAt = &now
	accountData.UpdatedAt = now
	return l.UpdateAccount(ctx, accountData)
}

func (l AccountRepository) ListAccount(ctx context.Context, payload *account.ListAccountDto) (accounts []domain.Account, paginationResult PaginationData, err error) {
	models := make([]domain.Account, 0)

	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.M{"deleted_at": nil},
	}}}

	if payload.Query != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{
			"$or": []bson.M{
				{"name.ar": bson.M{"$regex": payload.Query, "$options": "i"}},
				{"name.en": bson.M{"$regex": payload.Query, "$options": "i"}},
				{"email": bson.M{"$regex": payload.Query, "$options": "i"}},
			},
		})

	}
	var pipeline []interface{}

	pipeline = append(pipeline, matching)
	pipeline = append(pipeline,
		bson.M{
			"$lookup": bson.M{
				"from":         "locations",
				"localField":   "_id",
				"foreignField": "account_id",
				"as":           "locations",
				"pipeline":     []interface{}{bson.M{"$match": bson.M{"deleted_at": nil}}},
			},
		},
	)
	pipeline = append(pipeline,
		bson.M{
			"$addFields": bson.M{
				"locations_count": bson.M{"$size": "$locations"},
			},
		},
	)
	data, err := New(l.accountCollection.Collection).Context(ctx).Limit(payload.Limit).Page(payload.Page).Sort("created_at", -1).Aggregate(pipeline...)
	if data == nil || data.Data == nil {
		return models, paginationResult, err
	}

	for _, raw := range data.Data {
		model := domain.Account{}
		errUnmarshal := bson.Unmarshal(raw, &model)
		if errUnmarshal != nil {
			fmt.Println(errUnmarshal)
			break
		}
		models = append(models, model)
	}
	return models, data.Pagination, err
}
