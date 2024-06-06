package mongodb

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/menu/domain"
)

type itemRepo struct {
	itemCollection *mgm.Collection
}

func NewItemRepository(dbs *mongo.Database) domain.ItemRepository {
	itemCollection := mgm.Coll(&domain.Item{})
	return &itemRepo{
		itemCollection: itemCollection,
	}
}

func (i *itemRepo) GetByIds(ctx context.Context, ids []primitive.ObjectID) ([]domain.Item, error) {
	items := []domain.Item{}
	err := mgm.Coll(&domain.Item{}).SimpleFind(&items, bson.M{"_id": bson.M{"$in": ids}})
	return items, err
}
