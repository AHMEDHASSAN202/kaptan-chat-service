package mongodb

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"samm/internal/module/retails/consts"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/location"
	"samm/pkg/database/mongodb"
	"samm/pkg/utils"
	utilsDto "samm/pkg/utils/dto"
	"strconv"
	"time"
)

type locationRepository struct {
	locationCollection *mgm.Collection
}

const mongoLocationRepositoryTag = "LocationMongoRepository"

func NewLocationMongoRepository(dbs *mongo.Database) domain.LocationRepository {
	locationDbCollection := mgm.Coll(&domain.Location{})

	mongodb.CreateIndex(locationDbCollection.Collection, false,
		bson.E{"status", mongodb.IndexType.Asc},
		bson.E{"brand_details.is_active", mongodb.IndexType.Asc},
		bson.E{"brand_details._id", mongodb.IndexType.Asc},
		bson.E{"brand_details.cuisines._id", mongodb.IndexType.Asc},
		bson.E{"country._id", mongodb.IndexType.Asc},
		bson.E{"coordinate", mongodb.IndexType.Spatial},
		bson.E{"deleted_at", mongodb.IndexType.Asc})

	mongodb.CreateIndex(locationDbCollection.Collection, false,
		bson.E{"status", mongodb.IndexType.Asc},
		bson.E{"brand_details.is_active", mongodb.IndexType.Asc},
		bson.E{"country._id", mongodb.IndexType.Asc},
		bson.E{"coordinate", mongodb.IndexType.Spatial},
		bson.E{"deleted_at", mongodb.IndexType.Asc})

	return &locationRepository{
		locationCollection: locationDbCollection,
	}
}

func (l locationRepository) StoreLocation(ctx context.Context, location *domain.Location) (err error) {
	_, err = mgm.Coll(&domain.Location{}).InsertOne(ctx, location)
	if err != nil {
		return err
	}
	return nil

}
func (l locationRepository) BulkStoreLocation(ctx context.Context, data []domain.Location) (err error) {
	_, err = mgm.Coll(&domain.Location{}).InsertMany(ctx, utils.ConvertArrStructToInterfaceArr(data))
	if err != nil {
		return err
	}
	return nil

}

func (l locationRepository) UpdateLocation(ctx context.Context, location *domain.Location) (err error) {
	update := bson.M{"$set": location}
	_, err = mgm.Coll(&domain.Location{}).UpdateByID(ctx, location.ID, update)
	return
}
func (l locationRepository) FindLocation(ctx context.Context, Id primitive.ObjectID) (location *domain.Location, err error) {
	domainData := domain.Location{}
	filter := bson.M{"deleted_at": nil, "_id": Id}
	err = l.locationCollection.FirstWithCtx(ctx, filter, &domainData)

	return &domainData, err
}

func (l locationRepository) DeleteLocation(ctx context.Context, Id primitive.ObjectID) (err error) {
	locationData, err := l.FindLocation(ctx, Id)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	locationData.DeletedAt = &now
	locationData.UpdatedAt = now
	return l.UpdateLocation(ctx, locationData)
}

func (l locationRepository) DeleteLocationByAccountId(ctx context.Context, accountId primitive.ObjectID, causer *utilsDto.AdminDetails) (err error) {
	now := time.Now().UTC()

	filter := bson.M{"deleted_at": nil, "account_id": accountId}
	update := bson.M{"$set": bson.M{"deleted_at": now}, "$push": bson.M{"admin_details": causer}}
	_, err = l.locationCollection.UpdateMany(ctx, filter, update)
	return
}

func (l locationRepository) ListLocation(ctx context.Context, payload *location.ListLocationDto) (locations []domain.Location, paginationResult PaginationData, err error) {
	models := make([]domain.Location, 0)

	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.M{"deleted_at": nil},
	}}}

	if payload.Query != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{
			"$or": []bson.M{
				{"name.ar": bson.M{"$regex": payload.Query, "$options": "i"}},
				{"name.en": bson.M{"$regex": payload.Query, "$options": "i"}},
				{"tags": bson.M{"$regex": payload.Query, "$options": "i"}},
				{"phone": bson.M{"$regex": payload.Query, "$options": "i"}},
			},
		})

	}
	if payload.AccountId != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{
			"account_id": utils.ConvertStringIdToObjectId(payload.AccountId),
		})
	}
	if payload.BrandId != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{
			"brand_details._id": utils.ConvertStringIdToObjectId(payload.BrandId),
		})
	}
	if len(payload.Ids) > 0 && payload.Ids[0] != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{
			"_id": bson.M{"$in": utils.ConvertStringIdsToObjectIds(payload.Ids)},
		})
	}
	data, err := New(l.locationCollection.Collection).Context(ctx).Limit(payload.Limit).Page(payload.Page).Sort("created_at", -1).Aggregate(matching)
	if data == nil || data.Data == nil {
		return models, paginationResult, err
	}

	for _, raw := range data.Data {
		model := domain.Location{}
		errUnmarshal := bson.Unmarshal(raw, &model)
		if errUnmarshal != nil {
			break
		}
		models = append(models, model)
	}
	return models, data.Pagination, err

}

func (i *locationRepository) UpdateBulkByBrand(ctx context.Context, brand domain.BrandDetails) error {
	_, err := i.locationCollection.UpdateMany(ctx, bson.M{"brand_details._id": brand.Id}, bson.M{"$set": bson.M{"brand_details": brand}})
	if err != nil {
		return err
	}
	return nil
}
func (i *locationRepository) UpdateBulkByBrandCuisine(ctx context.Context, cuisine domain.CuisineDetails) error {

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.M{"elem._id": cuisine.ID}},
	})
	_, err := i.locationCollection.UpdateMany(ctx, bson.M{"brand_details.cuisines._id": cuisine.ID}, bson.M{"$set": bson.M{"brand_details.cuisines.$[elem]": cuisine}}, arrayFilters)
	if err != nil {
		return err
	}
	return nil
}

func (i *locationRepository) SoftDeleteBulkByBrandId(ctx context.Context, brandId primitive.ObjectID, causer *utilsDto.AdminDetails) error {
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	if causer.Id != primitive.NilObjectID {
		update["$push"] = bson.M{"admin_details": causer}
	}
	_, err := i.locationCollection.UpdateMany(ctx, bson.M{"brand_details._id": brandId}, update)
	if err != nil {
		return err
	}
	return nil
}

func (l locationRepository) ListMobileLocation(ctx context.Context, payload *location.ListLocationMobileDto) (locations []domain.LocationMobile, paginationResult *PaginationData, err error) {
	models := make([]domain.LocationMobile, 0)

	latPoint, _ := strconv.ParseFloat(payload.MobileHeaders.Lat, 64)
	lngPoint, _ := strconv.ParseFloat(payload.MobileHeaders.Lng, 64)

	var pipeline []interface{}
	matching := bson.M{"$and": []bson.D{
		{{"deleted_at", nil}},
		{{"country._id", payload.MobileHeaders.CountryId}},
		{{"brand_details.is_active", true}},
		{{"status", consts.LocationStatusActive}},
	}}
	if payload.Query != "" {
		matching["$and"] = append(matching["$and"].([]bson.D), bson.D{
			{
				"$or", []bson.M{
					{"name.ar": bson.M{"$regex": payload.Query, "$options": "i"}},
					{"name.en": bson.M{"$regex": payload.Query, "$options": "i"}},
					{"tags": bson.M{"$regex": payload.Query, "$options": "i"}},
					{"phone": bson.M{"$regex": payload.Query, "$options": "i"}},
				},
			}})
	}

	if payload.BrandId != "" {
		matching["$and"] = append(matching["$and"].([]bson.D), bson.D{
			{"brand_details._id", utils.ConvertStringIdToObjectId(payload.BrandId)}})
	}

	if payload.CuisineIds != nil && len(payload.CuisineIds) > 0 && payload.CuisineIds[0] != "" {
		matching["$and"] = append(matching["$and"].([]bson.D), bson.D{{"brand_details.cuisines._id", bson.M{"$in": utils.ConvertStringIdsToObjectIds(payload.CuisineIds)}}})
	}

	nearbyStage := bson.M{"$geoNear": bson.M{
		"near": bson.M{
			"coordinates": []float64{
				lngPoint,
				latPoint,
			},
			"type": "Point",
		},
		"key":                "coordinate",
		"maxDistance":        payload.Distance * 1000,
		"spherical":          true,
		"distanceMultiplier": 0.001,
		"distanceField":      "distance",
		"query":              matching,
	}}
	pipeline = append(pipeline, nearbyStage)

	pipeline = append(pipeline,
		AddIsOpenFieldPipeline(payload.CountryId),
	)
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"is_open": -1}})

	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"name":             1,
			"city":             1,
			"street":           1,
			"cover_image":      1,
			"account_id":       1,
			"logo":             1,
			"phone":            1,
			"coordinate":       1,
			"brand_details":    1,
			"preparation_time": 1,
			"country":          1,
			"is_open":          1,
			"distance":         1,
			"created_at":       1,
			"updated_at":       1,
		},
	})

	data, err := New(l.locationCollection.Collection).Context(ctx).Limit(payload.Limit).Page(payload.Page).Sort("created_at", -1).Aggregate(pipeline...)
	if data == nil || data.Data == nil {
		return models, nil, err
	}

	for _, raw := range data.Data {
		model := domain.LocationMobile{}
		errUnmarshal := bson.Unmarshal(raw, &model)
		if errUnmarshal != nil {
			break
		}
		model.WorkingHour = make([]domain.WorkingHour, 0)
		model.SetOpenStatus()
		models = append(models, model)
	}
	return models, &data.Pagination, err

}

func (l locationRepository) FindMobileLocation(ctx context.Context, Id primitive.ObjectID, payload *location.FindLocationMobileDto) (location *domain.LocationMobile, err error) {
	domainData := &domain.LocationMobile{}
	var pipeline []interface{}
	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.M{"_id": Id},
		bson.M{"deleted_at": nil},
		bson.M{"brand_details.is_active": true},
		bson.M{"status": consts.LocationStatusActive},
	}}}
	pipeline = append(pipeline, matching)

	pipeline = append(pipeline,
		AddIsOpenFieldPipeline(payload.CountryId),
	)
	exists, err := l.locationCollection.SimpleAggregateFirstWithCtx(ctx, domainData, pipeline...)
	if !exists {
		return domainData, errors.New("Not Found")
	}
	domainData.SetOpenStatus()

	latPoint, _ := strconv.ParseFloat(payload.MobileHeaders.Lat, 64)
	lngPoint, _ := strconv.ParseFloat(payload.MobileHeaders.Lng, 64)
	domainData.SetDistance(latPoint, lngPoint)
	return domainData, err
}
