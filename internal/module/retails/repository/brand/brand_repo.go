package brand

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/brand"
	"samm/pkg/database/mongodb"
	"samm/pkg/utils"
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

func (i *brandRepo) Update(ctx *context.Context, id primitive.ObjectID, doc *domain.Brand) error {
	update := bson.M{"$set": doc}
	_, err := i.brandCollection.UpdateByID(*ctx, id, update)
	if err != nil {
		return err
	}
	return nil
}

func (l brandRepo) FindBrand(ctx *context.Context, Id primitive.ObjectID) (*domain.Brand, error) {
	var domainData domain.Brand
	filter := bson.M{"deleted_at": nil, "_id": Id}
	err := l.brandCollection.FirstWithCtx(*ctx, filter, &domainData)
	return &domainData, err
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

func (i *brandRepo) List(ctx *context.Context, query *brand.ListBrandDto) (*[]domain.Brand, *utils.PaginationResult, error) {
	filter := bson.M{}
	var match []bson.M
	match = append(match, bson.M{"deleted_at": nil})

	var brands []domain.Brand

	offset := (query.Page - 1) * query.Limit
	options := options.Find().SetLimit(query.Limit).SetSkip(offset)

	if query.Query != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"name.ar": bson.M{"$regex": query.Query, "$options": "i"}},
				{"name.en": bson.M{"$regex": query.Query, "$options": "i"}},
			},
		}
	}

	if len(match) > 0 {
		filter["$and"] = match
	}

	// Query the collection for the total count of documents
	totalItems, err := i.brandCollection.CountDocuments(*ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	totalPages := int(math.Ceil(float64(totalItems) / float64(query.Limit)))

	err = i.brandCollection.SimpleFind(&brands, filter, options)
	return &brands, &utils.PaginationResult{Page: query.Page, TotalPages: int64(totalPages), TotalItems: totalItems}, err
}
