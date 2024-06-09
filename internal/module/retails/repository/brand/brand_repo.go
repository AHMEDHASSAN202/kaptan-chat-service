package brand

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/brand"
	"samm/pkg/database/mongodb"
	"time"
)

type brandRepo struct {
	brandCollection *mgm.Collection
}

func NewBrandRepository(dbs *mongo.Database) domain.BrandRepository {
	brandCollection := mgm.Coll(&domain.Brand{})
	mongodb.CreateIndex(brandCollection.Collection, false, bson.E{"name.ar", mongodb.IndexType.Text}, bson.E{"name.en", mongodb.IndexType.Text})
	return &brandRepo{
		brandCollection: brandCollection,
	}
}

func (i *brandRepo) Create(ctx *context.Context, doc *domain.Brand) error {
	_, err := i.brandCollection.InsertOne(*ctx, doc)
	if err != nil {
		return err
	}
	return nil
}

func (i *brandRepo) Update(ctx *context.Context, doc *domain.Brand) error {
	update := bson.M{"$set": doc}
	_, err := i.brandCollection.UpdateByID(*ctx, doc.ID, update)
	if err != nil {
		return err
	}
	return nil
}

func (i *brandRepo) GetByIds(ctx *context.Context, ids *[]primitive.ObjectID) (*[]domain.Brand, error) {
	var cuisines []domain.Brand
	err := i.brandCollection.SimpleFind(&cuisines, bson.M{"_id": bson.M{"$in": *ids}})
	return &cuisines, err
}

func (i *brandRepo) SoftDelete(ctx *context.Context, id primitive.ObjectID) error {
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	_, err := i.brandCollection.UpdateByID(*ctx, id, update)
	if err != nil {
		return err
	}
	return nil
}

//func (i *brandRepo) ChangeStatus(ctx *context.Context, dto *cuisine.ChangeCuisineStatusDto) error {
//	update := bson.M{"$set": bson.M{"is_hidden": dto.IsHidden}}
//	_, err := i.brandCollection.UpdateByID(*ctx, dto.Id, update)
//	if err != nil {
//		return err
//	}
//	return nil
//}

func (i *brandRepo) List(ctx *context.Context, query *brand.ListBrandDto) (*[]domain.Brand, error) {
	options := options.Find()
	filter := bson.M{}

	offset := (query.Page - 1) * query.Limit
	options.SetLimit(query.Limit)
	options.SetSkip(offset)

	if query.Query != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"name.ar": bson.M{"$regex": query.Query, "$options": "i"}},
				{"name.en": bson.M{"$regex": query.Query, "$options": "i"}},
			},
		}
	}
	var cuisines []domain.Brand
	err := i.brandCollection.SimpleFind(&cuisines, filter, options)
	return &cuisines, err
}
