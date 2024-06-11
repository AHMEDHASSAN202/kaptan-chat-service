package cuisine

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/cuisine"
	"samm/pkg/logger"
	"time"
)

type cuisineRepo struct {
	cuisineCollection *mgm.Collection
	logger            logger.ILogger
}

func NewCuisineRepository(dbs *mongo.Database, log logger.ILogger) domain.CuisineRepository {
	cuisineCollection := mgm.Coll(&domain.Cuisine{})
	//text search menu cuisine
	//mongodb.CreateIndex(cuisineCollection.Collection, false, bson.E{"name.ar", mongodb.IndexType.Text}, bson.E{"name.en", mongodb.IndexType.Text}, bson.E{"tags", mongodb.IndexType.Text},
	//	bson.E{"desc.ar", mongodb.IndexType.Text}, bson.E{"desc.en", mongodb.IndexType.Text})
	//make sure there are no duplicated menu cuisine
	//mongodb.CreateIndex(cuisineCollection.Collection, true, bson.E{"name.ar", mongodb.IndexType.Asc}, bson.E{"name.en", mongodb.IndexType.Asc}, bson.E{"account_id", mongodb.IndexType.Asc}, bson.E{"deleted_at", mongodb.IndexType.Asc})
	return &cuisineRepo{
		cuisineCollection: cuisineCollection,
		logger:            log,
	}
}

func (i *cuisineRepo) Create(doc *domain.Cuisine) error {
	err := i.cuisineCollection.Create(doc)
	if err != nil {
		return err
	}
	return nil
}

func (i *cuisineRepo) Update(doc *domain.Cuisine) error {
	err := i.cuisineCollection.Update(doc)
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
	err := i.cuisineCollection.SimpleFind(&cuisines, bson.M{"_id": bson.M{"$in": *ids}, "deleted_at": nil})
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

func (i *cuisineRepo) List(ctx *context.Context, dto *cuisine.ListCuisinesDto) (cuisinesRes *[]domain.Cuisine, paginationMeta *PaginationData, err error) {
	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.D{{"deleted_at", nil}},
	}}}

	if dto.Query != "" {
		pattern := ".*" + dto.Query + ".*"
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"$or": []bson.M{{"name.ar": bson.M{"$regex": pattern, "$options": "i"}}, {"name.en": bson.M{"$regex": pattern, "$options": "i"}}}})
	}

	data, err := New(i.cuisineCollection.Collection).Context(*ctx).Limit(dto.Limit).Page(dto.Page).Sort("created_at", -1).Aggregate(matching)

	if data == nil || data.Data == nil {
		return nil, nil, err
	}

	cuisines := make([]domain.Cuisine, 0)
	for _, raw := range data.Data {
		model := domain.Cuisine{}
		err = bson.Unmarshal(raw, &model)
		if err != nil {
			i.logger.Error("cuisine Repo -> List -> ", err)
			break
		}
		cuisines = append(cuisines, model)
	}
	paginationMeta = &data.Pagination
	cuisinesRes = &cuisines

	return
}
