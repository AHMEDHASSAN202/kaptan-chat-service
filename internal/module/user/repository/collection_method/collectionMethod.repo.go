package mongodb

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/user/domain"
	"time"
)

type CollectionMethodRepository struct {
	collectionMethodCollection *mgm.Collection
}

const mongoCollectionMethodRepositoryTag = "CollectionMethodMongoRepository"

func NewCollectionMethodMongoRepository(dbs *mongo.Database) domain.CollectionMethodRepository {
	collectionMethodCollection := mgm.Coll(&domain.CollectionMethods{})
	return &CollectionMethodRepository{
		collectionMethodCollection: collectionMethodCollection,
	}
}

func (l CollectionMethodRepository) StoreCollectionMethod(ctx context.Context, user *domain.CollectionMethods) (err error) {
	err = mgm.Coll(&domain.CollectionMethods{}).CreateWithCtx(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (l CollectionMethodRepository) UpdateCollectionMethod(ctx context.Context, collectionMethod *domain.CollectionMethods) (err error) {
	collectionMethod.UpdatedAt = time.Now().UTC()
	update := bson.M{"$set": collectionMethod}
	_, err = mgm.Coll(&domain.CollectionMethods{}).UpdateByID(ctx, collectionMethod.ID, update)
	return
}
func (l CollectionMethodRepository) FindCollectionMethod(ctx context.Context, Id primitive.ObjectID, userId primitive.ObjectID) (user *domain.CollectionMethods, err error) {
	domainData := domain.CollectionMethods{}
	filter := bson.M{"deleted_at": nil, "_id": Id}
	err = l.collectionMethodCollection.FirstWithCtx(ctx, filter, &domainData)
	return &domainData, err
}

func (l CollectionMethodRepository) DeleteCollectionMethod(ctx context.Context, Id primitive.ObjectID, userId primitive.ObjectID) (err error) {
	now := time.Now().UTC()
	update := bson.M{"$set": bson.M{"deleted_at": now, "updated_at": now}}
	_, err = mgm.Coll(&domain.CollectionMethods{}).UpdateByID(ctx, Id, update)
	return
}

func (l CollectionMethodRepository) ListCollectionMethod(ctx context.Context, collectionMethodType string, userId primitive.ObjectID) (collectionMethods []domain.CollectionMethods, err error) {
	// Query the collection for the total count of documents
	err = l.collectionMethodCollection.SimpleFindWithCtx(ctx, &collectionMethods, bson.M{"deleted_at": nil, "user_id": userId, "type": collectionMethodType})
	if err != nil {
		return []domain.CollectionMethods{}, err
	}

	return collectionMethods, err

}
