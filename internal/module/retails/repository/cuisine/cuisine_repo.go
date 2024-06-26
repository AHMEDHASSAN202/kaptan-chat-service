package cuisine

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/jinzhu/copier"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/cuisine"
	"samm/pkg/database/mongodb"
	"samm/pkg/logger"
	"time"
)

type cuisineRepo struct {
	cuisineCollection *mgm.Collection
	locationRepo      domain.LocationRepository
	logger            logger.ILogger
}

func NewCuisineRepository(dbs *mongo.Database, locationRepo domain.LocationRepository, log logger.ILogger) domain.CuisineRepository {
	cuisineCollection := mgm.Coll(&domain.Cuisine{})
	mongodb.CreateIndex(cuisineCollection.Collection, false, bson.E{"name.ar", mongodb.IndexType.Text}, bson.E{"name.en", mongodb.IndexType.Text})
	return &cuisineRepo{
		cuisineCollection: cuisineCollection,
		locationRepo:      locationRepo,
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

func (i *cuisineRepo) UpdateCuisineAndLocations(doc *domain.Cuisine) error {
	err := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		err := i.cuisineCollection.UpdateWithCtx(sc, doc)
		if err != nil {
			return err
		}
		var cuisineDetails domain.CuisineDetails
		copier.Copy(&cuisineDetails, doc)
		err = i.locationRepo.UpdateBulkByBrandCuisine(sc, cuisineDetails)
		if err != nil {
			return err
		}
		return session.CommitTransaction(sc)
	})
	return err
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
