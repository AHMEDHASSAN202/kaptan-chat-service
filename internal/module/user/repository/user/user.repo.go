package mongodb

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"samm/internal/module/user/domain"
	"samm/internal/module/user/dto/user"
	"samm/pkg/utils"
	"time"
)

type UserRepository struct {
	userCollection *mgm.Collection
}

const mongoUserRepositoryTag = "UserMongoRepository"

func NewUserMongoRepository(dbs *mongo.Database) domain.UserRepository {
	userDbCollection := mgm.Coll(&domain.User{})

	return &UserRepository{
		userCollection: userDbCollection,
	}
}

func (l UserRepository) StoreUser(ctx context.Context, user *domain.User) (err error) {
	_, err = mgm.Coll(&domain.User{}).InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil

}

func (l UserRepository) UpdateUser(ctx context.Context, user *domain.User) (err error) {
	update := bson.M{"$set": user}
	_, err = mgm.Coll(&domain.User{}).UpdateByID(ctx, user.ID, update)
	return
}
func (l UserRepository) FindUser(ctx context.Context, Id primitive.ObjectID) (user *domain.User, err error) {
	domainData := domain.User{}
	filter := bson.M{"deleted_at": nil, "_id": Id}
	err = l.userCollection.FirstWithCtx(ctx, filter, &domainData)

	return &domainData, err
}

func (l UserRepository) DeleteUser(ctx context.Context, Id primitive.ObjectID) (err error) {
	userData, err := l.FindUser(ctx, Id)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	userData.DeletedAt = &now
	userData.UpdatedAt = now
	return l.UpdateUser(ctx, userData)
}

func (l UserRepository) ListUser(ctx context.Context, payload *user.ListUserDto) (users []domain.User, paginationResult utils.PaginationResult, err error) {

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
	collection := mgm.Coll(&domain.User{})
	totalItems, err := collection.CountDocuments(ctx, filter)

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalItems) / float64(payload.Limit)))

	var data []domain.User
	err = l.userCollection.SimpleFind(&data, filter, findOptions)

	return data, utils.PaginationResult{Page: payload.Page, TotalPages: int64(totalPages), TotalItems: totalItems}, err

}
