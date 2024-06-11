package mongodb

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
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

func (l AccountRepository) ListAccount(ctx context.Context, payload *account.ListAccountDto) (accounts []domain.Account, paginationResult utils.PaginationResult, err error) {

	offset := (payload.Page - 1) * payload.Limit
	findOptions := options.Find().SetLimit(payload.Limit).SetSkip(offset)

	filter := bson.M{}
	match := []bson.M{}
	match = append(match, bson.M{"deleted_at": nil})
	if payload.Query != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"name.ar": bson.M{"$regex": payload.Query, "$options": "i"}},
				{"name.en": bson.M{"$regex": payload.Query, "$options": "i"}},
				{"email": bson.M{"$regex": payload.Query, "$options": "i"}},
			},
		}
	}
	filter["$and"] = match

	// Query the collection for the total count of documents
	collection := mgm.Coll(&domain.Account{})
	totalItems, err := collection.CountDocuments(ctx, filter)

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalItems) / float64(payload.Limit)))

	var data []domain.Account
	err = l.accountCollection.SimpleFind(&data, filter, findOptions)

	return data, utils.PaginationResult{Page: payload.Page, TotalPages: int64(totalPages), TotalItems: totalItems}, err

}
