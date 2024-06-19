package brand

import (
	"context"
	"errors"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/brand"
	"samm/pkg/database/mongodb"
	"samm/pkg/logger"
	"samm/pkg/utils"
)

type brandRepo struct {
	brandCollection *mgm.Collection
	locationRepo    domain.LocationRepository
	logger          logger.ILogger
}

func NewBrandRepository(dbs *mongo.Database, locationRepo domain.LocationRepository, log logger.ILogger) domain.BrandRepository {
	brandCollection := mgm.Coll(&domain.Brand{})
	mongodb.CreateIndex(brandCollection.Collection, false, bson.E{"name.ar", mongodb.IndexType.Text}, bson.E{"name.en", mongodb.IndexType.Text})
	return &brandRepo{
		brandCollection: brandCollection,
		locationRepo:    locationRepo,
		logger:          log,
	}
}

func (i *brandRepo) Create(doc *domain.Brand) (err error) {
	err = mgm.Coll(doc).Create(doc)
	if err != nil {
		return
	}
	return
}

func (i *brandRepo) Update(doc *domain.Brand) error {
	err := i.brandCollection.Update(doc)
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

func (i *brandRepo) UpdateBrandAndLocations(doc *domain.Brand) error {
	err := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		err := i.brandCollection.UpdateWithCtx(sc, doc)
		if err != nil {
			return err
		}
		brandDetails := domain.BrandDetails{
			Id:       doc.ID,
			Name:     doc.Name,
			Logo:     doc.Logo,
			IsActive: doc.IsActive,
		}
		err = i.locationRepo.UpdateBulkByBrand(sc, brandDetails)
		if err != nil {
			return err
		}
		return session.CommitTransaction(sc)
	})
	return err
}

func (i *brandRepo) SoftDelete(doc *domain.Brand) error {
	err := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		err := i.brandCollection.UpdateWithCtx(sc, doc)
		if err != nil {
			return err
		}
		err = i.locationRepo.SoftDeleteBulkByBrandId(sc, doc.ID)
		if err != nil {
			return err
		}
		return session.CommitTransaction(sc)
	})
	return err
}

func (i *brandRepo) List(ctx *context.Context, dto *brand.ListBrandDto) (brandsRes *[]domain.Brand, paginationMeta *PaginationData, err error) {
	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.D{{"deleted_at", nil}},
	}}}

	if dto.Query != "" {
		pattern := ".*" + dto.Query + ".*"
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"$or": []bson.M{{"name.ar": bson.M{"$regex": pattern, "$options": "i"}}, {"name.en": bson.M{"$regex": pattern, "$options": "i"}}}})
	}
	if len(dto.Ids) > 0 && dto.Ids[0] != "" {
		ids := utils.ConvertStringIdsToObjectIds(dto.Ids)
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"_id": bson.M{"$in": ids}})
	}

	data, err := New(i.brandCollection.Collection).Context(*ctx).Limit(dto.Limit).Page(dto.Page).Sort("created_at", -1).Aggregate(matching)

	if data == nil || data.Data == nil {
		return nil, nil, err
	}

	brands := make([]domain.Brand, 0)
	for _, raw := range data.Data {
		model := domain.Brand{}
		err = bson.Unmarshal(raw, &model)
		if err != nil {
			i.logger.Error("brands Repo -> List -> ", err)
			break
		}
		brands = append(brands, model)
	}
	paginationMeta = &data.Pagination
	brandsRes = &brands

	return
}

func (i brandRepo) FindWithCuisines(ctx *context.Context, Id primitive.ObjectID) (*domain.Brand, error) {
	var domainData domain.Brand
	var filters []bson.M
	filters = append(filters, bson.M{"$match": bson.M{"_id": Id, "deleted_at": nil}})
	filters = append(filters, bson.M{"$lookup": bson.M{"from": "cuisines", "localField": "cuisineids", "foreignField": "_id", "as": "cuisines"}})

	data, err := i.brandCollection.Aggregate(*ctx, filters)

	if err != nil {
		return nil, err
	}
	var tempData []bson.Raw

	if err = data.All(*ctx, &tempData); err != nil {
		return nil, err
	}

	if len(tempData) > 0 {
		if err := bson.Unmarshal(tempData[0], &domainData); err == nil {
			return &domainData, err
		}
	}
	return nil, errors.New("Not Found")
}
