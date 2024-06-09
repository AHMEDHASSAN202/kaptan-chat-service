package cuisine

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/cuisine"
	"samm/pkg/utils"
	"time"
)

type cuisineRepo struct {
	cuisineCollection *mgm.Collection
}

func NewCuisineRepository(dbs *mongo.Database) domain.CuisineRepository {
	cuisineCollection := mgm.Coll(&domain.Cuisine{})
	//text search menu cuisine
	//mongodb.CreateIndex(cuisineCollection.Collection, false, bson.E{"name.ar", mongodb.IndexType.Text}, bson.E{"name.en", mongodb.IndexType.Text}, bson.E{"tags", mongodb.IndexType.Text},
	//	bson.E{"desc.ar", mongodb.IndexType.Text}, bson.E{"desc.en", mongodb.IndexType.Text})
	//make sure there are no duplicated menu cuisine
	//mongodb.CreateIndex(cuisineCollection.Collection, true, bson.E{"name.ar", mongodb.IndexType.Asc}, bson.E{"name.en", mongodb.IndexType.Asc}, bson.E{"account_id", mongodb.IndexType.Asc}, bson.E{"deleted_at", mongodb.IndexType.Asc})
	return &cuisineRepo{
		cuisineCollection: cuisineCollection,
	}
}

func (i *cuisineRepo) Create(ctx *context.Context, doc *[]domain.Cuisine) error {
	_, err := i.cuisineCollection.InsertMany(*ctx, utils.ConvertArrStructToInterfaceArr(*doc))
	if err != nil {
		return err
	}
	return nil
}

func (i *cuisineRepo) Update(ctx *context.Context, id primitive.ObjectID, doc *domain.Cuisine) error {
	update := bson.M{"$set": doc}
	_, err := i.cuisineCollection.UpdateByID(*ctx, id, update)
	if err != nil {
		return err
	}
	return nil
}

func (l cuisineRepo) Find(ctx *context.Context, Id primitive.ObjectID) (*domain.Cuisine, error) {
	var domainData domain.Cuisine
	filter := bson.M{"deleted_at": nil, "_id": Id}
	err := l.cuisineCollection.FirstWithCtx(*ctx, filter, &domainData)
	return &domainData, err
}

func (i *cuisineRepo) GetByIds(ctx *context.Context, ids *[]primitive.ObjectID) (*[]domain.Cuisine, error) {
	var cuisines []domain.Cuisine
	err := i.cuisineCollection.SimpleFind(&cuisines, bson.M{"_id": bson.M{"$in": *ids}})
	return &cuisines, err
}

func (i *cuisineRepo) SoftDelete(ctx *context.Context, id primitive.ObjectID) error {
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	_, err := i.cuisineCollection.UpdateByID(*ctx, id, update)
	if err != nil {
		return err
	}
	return nil
}
func (i *cuisineRepo) ChangeStatus(ctx *context.Context, dto *cuisine.ChangeCuisineStatusDto) error {
	update := bson.M{"$set": bson.M{"is_hidden": dto.IsHidden}}
	_, err := i.cuisineCollection.UpdateByID(*ctx, dto.Id, update)
	if err != nil {
		return err
	}
	return nil
}

func (i *cuisineRepo) List(ctx *context.Context, query *cuisine.ListCuisinesDto) (*[]domain.Cuisine, *utils.PaginationResult, error) {
	filter := bson.M{}
	var cuisines []domain.Cuisine

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

	// Query the collection for the total count of documents
	totalItems, err := i.cuisineCollection.CountDocuments(*ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalItems) / float64(query.Limit)))

	err = i.cuisineCollection.SimpleFind(&cuisines, filter, options)

	return &cuisines, &utils.PaginationResult{Page: query.Page, TotalPages: int64(totalPages), TotalItems: totalItems}, err
}
