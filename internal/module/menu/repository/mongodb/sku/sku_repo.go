package sku

import (
	"context"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/sku"
	"samm/pkg/database/mongodb"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type skuRepo struct {
	skuCollection *mgm.Collection
}

func NewSkuRepository(dbs *mongo.Database) domain.SKURepository {
	skuCollection := mgm.Coll(&domain.SKU{})
	//make sure there are no duplicated sku
	mongodb.CreateIndex(skuCollection.Collection, true, bson.E{"name", mongodb.IndexType.Asc})

	return &skuRepo{
		skuCollection: skuCollection,
	}
}

func (i *skuRepo) Create(ctx context.Context, doc domain.SKU) error {
	_, err := i.skuCollection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	return nil
}

func (i *skuRepo) List(ctx context.Context, query *sku.ListSKUDto) ([]domain.SKU, error) {
	filter := bson.M{}
	if query.Query != "" {
		filter = bson.M{
			"name": bson.M{"$regex": query.Query, "$options": "i"},
		}
	}
	skus := []domain.SKU{}
	err := i.skuCollection.SimpleFindWithCtx(ctx, &skus, filter)
	return skus, err
}
