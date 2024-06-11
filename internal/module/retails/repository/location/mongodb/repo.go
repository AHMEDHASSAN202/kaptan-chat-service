package mongodb

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/location"
	"samm/pkg/utils"
	"time"
)

type locationRepository struct {
	locationCollection *mgm.Collection
}

const mongoLocationRepositoryTag = "LocationMongoRepository"

func NewLocationMongoRepository(dbs *mongo.Database) domain.LocationRepository {
	locationDbCollection := mgm.Coll(&domain.Location{})

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

func (l locationRepository) DeleteLocationByAccountId(ctx context.Context, accountId primitive.ObjectID) (err error) {
	now := time.Now().UTC()

	filter := bson.M{"deleted_at": nil, "account_id": accountId}
	update := bson.M{"$set": bson.M{"deleted_at": now}}
	_, err = l.locationCollection.UpdateMany(ctx, filter, update)
	return
}

func (l locationRepository) ListLocation(ctx context.Context, payload *location.ListLocationDto) (locations []domain.Location, paginationResult utils.PaginationResult, err error) {

	offset := (payload.Page - 1) * payload.Limit
	findOptions := options.Find().SetLimit(payload.Limit).SetSkip(offset)

	filter := bson.M{}
	match := []bson.M{}
	match = append(match, bson.M{"deleted_at": nil, "brand_details.is_active": true})

	if payload.Query != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"name.ar": bson.M{"$regex": payload.Query, "$options": "i"}},
				{"name.en": bson.M{"$regex": payload.Query, "$options": "i"}},
				{"tags": bson.M{"$regex": payload.Query, "$options": "i"}},
				{"phone": bson.M{"$regex": payload.Query, "$options": "i"}},
			},
		}
	}
	if payload.AccountId != "" {
		match = append(match, bson.M{"account_id": utils.ConvertStringIdToObjectId(payload.AccountId)})
	}
	if payload.BrandId != "" {
		match = append(match, bson.M{"brand_details._id": utils.ConvertStringIdToObjectId(payload.BrandId)})
	}
	if len(match) > 0 {
		filter["$and"] = match
	}

	// Query the collection for the total count of documents
	collection := mgm.Coll(&domain.Location{})
	totalItems, err := collection.CountDocuments(ctx, filter)

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalItems) / float64(payload.Limit)))

	var data []domain.Location
	err = l.locationCollection.SimpleFind(&data, filter, findOptions)

	return data, utils.PaginationResult{Page: payload.Page, TotalPages: int64(totalPages), TotalItems: totalItems}, err

}

func (i *locationRepository) UpdateBulkByBrand(ctx context.Context, brand domain.BrandDetails) error {
	_, err := i.locationCollection.UpdateMany(ctx, bson.M{"brand_details._id": brand.Id}, bson.M{"$set": bson.M{"brand_details": brand}})
	if err != nil {
		return err
	}
	return nil
}

func (i *locationRepository) SoftDeleteBulkByBrandId(ctx context.Context, brandId primitive.ObjectID) error {
	_, err := i.locationCollection.UpdateMany(ctx, bson.M{"brand_details._id": brandId}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	if err != nil {
		return err
	}
	return nil
}
