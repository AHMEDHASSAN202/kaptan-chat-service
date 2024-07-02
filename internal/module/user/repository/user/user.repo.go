package mongodb

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"samm/internal/module/user/domain"
	"samm/internal/module/user/dto/user"
	"samm/pkg/logger"
	"samm/pkg/utils"
)

type UserRepository struct {
	userCollection        *mgm.Collection
	deletedUserCollection *mgm.Collection
	logger                logger.ILogger
}

const mongoUserRepositoryTag = "UserMongoRepository"

func NewUserMongoRepository(dbs *mongo.Database, log logger.ILogger) domain.UserRepository {
	userDbCollection := mgm.Coll(&domain.User{})
	deletedUserCollection := mgm.Coll(&domain.DeletedUser{})

	return &UserRepository{
		userCollection:        userDbCollection,
		deletedUserCollection: deletedUserCollection,
		logger:                log,
	}
}

func (l UserRepository) StoreUser(ctx *context.Context, user *domain.User) (err error) {
	_, err = l.userCollection.InsertOne(*ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (l UserRepository) UpdateUser(ctx *context.Context, user *domain.User) (err error) {
	update := bson.M{"$set": user}
	upsert := true
	opts := options.UpdateOptions{Upsert: &upsert}
	_, err = mgm.Coll(&domain.User{}).UpdateByID(*ctx, user.ID, update, &opts)
	return
}

func (l UserRepository) FindUser(ctx *context.Context, Id primitive.ObjectID) (user *domain.User, err error) {
	domainData := domain.User{}
	filter := bson.M{"deleted_at": nil, "_id": Id}
	err = l.userCollection.FirstWithCtx(*ctx, filter, &domainData)

	return &domainData, err
}

func (l UserRepository) GetUserByPhoneNumber(ctx *context.Context, phoneNum, countryCode string) (user domain.User, err error) {
	domainData := domain.User{}
	filter := bson.M{"phone_number": phoneNum, "country_code": countryCode}
	err = l.userCollection.FirstWithCtx(*ctx, filter, &domainData)

	return domainData, err
}

func (i *UserRepository) List(ctx *context.Context, dto *user.ListUserDto) (usersRes *[]domain.User, paginationMeta *PaginationData, err error) {
	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.D{{"deleted_at", nil}},
	}}}

	if dto.Query != "" {
		pattern := ".*" + dto.Query + ".*"
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"$or": []bson.M{{"name": bson.M{"$regex": pattern, "$options": "i"}}, {"phone_number": bson.M{"$regex": pattern, "$options": "i"}}}})
	}

	if dto.Status != "" {
		if dto.Status == "active" {
			matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"is_active": true})
		} else if dto.Status == "inactive" {
			matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"is_active": false})
		}
	}
	if dto.Dob != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"dob": dto.Dob})
	}

	data, err := New(i.userCollection.Collection).Context(*ctx).Limit(dto.Limit).Page(dto.Page).Sort("created_at", -1).Aggregate(matching)

	if data == nil || data.Data == nil {
		return nil, nil, err
	}

	users := make([]domain.User, 0)
	for _, raw := range data.Data {
		model := domain.User{}
		err = bson.Unmarshal(raw, &model)
		if err != nil {
			i.logger.Error("user Repo -> List -> ", err)
			break
		}
		users = append(users, model)
	}
	paginationMeta = &data.Pagination
	usersRes = &users

	return
}

func (l UserRepository) UserEmailExists(ctx *context.Context, email, userId string) bool {
	filter := bson.M{"deleted_at": nil, "email": email}
	if userId != "" {
		filter = bson.M{"deleted_at": nil, "email": email, "_id": bson.M{"$ne": utils.ConvertStringIdToObjectId(userId)}}
	}
	count, _ := l.userCollection.CountDocuments(*ctx, filter)
	return count > 0
}

func (l UserRepository) InsertDeletedUser(ctx *context.Context, user *domain.DeletedUser) (err error) {
	_, err = l.deletedUserCollection.InsertOne(*ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (l UserRepository) RemoveDeletedUser(user *domain.DeletedUser) (err error) {
	err = l.deletedUserCollection.Delete(user)
	if err != nil {
		return err
	}
	return
}

func (r *UserRepository) FindByToken(ctx context.Context, token string) (*domain.User, error) {
	var domainData domain.User
	err := r.userCollection.FirstWithCtx(ctx, bson.M{"tokens": token, "deleted_at": nil}, &domainData, nil)
	return &domainData, err
}
