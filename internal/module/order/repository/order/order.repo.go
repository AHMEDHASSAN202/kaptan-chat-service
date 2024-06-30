package mongodb

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/dto/order"
	"samm/pkg/utils"
	"time"
)

type OrderRepository struct {
	orderCollection *mgm.Collection
}

const mongoOrderRepositoryTag = "OrderMongoRepository"

func NewOrderMongoRepository(dbs *mongo.Database) domain.OrderRepository {
	orderDbCollection := mgm.Coll(&domain.Order{})

	return &OrderRepository{
		orderCollection: orderDbCollection,
	}
}

func (l OrderRepository) StoreOrder(ctx context.Context, order *domain.Order) (err error) {
	_, err = mgm.Coll(&domain.Order{}).InsertOne(ctx, order)
	if err != nil {
		return err
	}
	return nil

}

func (l OrderRepository) UpdateOrder(ctx context.Context, order *domain.Order) (err error) {
	update := bson.M{"$set": order}
	_, err = mgm.Coll(&domain.Order{}).UpdateByID(ctx, order.ID, update)
	return
}
func (l OrderRepository) FindOrder(ctx context.Context, Id primitive.ObjectID) (order *domain.Order, err error) {
	domainData := domain.Order{}
	filter := bson.M{"deleted_at": nil, "_id": Id}
	err = l.orderCollection.FirstWithCtx(ctx, filter, &domainData)

	return &domainData, err
}

func (l OrderRepository) DeleteOrder(ctx context.Context, Id primitive.ObjectID) (err error) {
	orderData, err := l.FindOrder(ctx, Id)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	orderData.DeletedAt = &now
	orderData.UpdatedAt = now
	return l.UpdateOrder(ctx, orderData)
}

func (l OrderRepository) ListOrder(ctx context.Context, payload *order.ListOrderDto) (orders []domain.Order, paginationResult utils.PaginationResult, err error) {

	offset := (payload.Page - 1) * payload.Limit
	findOptions := options.Find().SetLimit(payload.Limit).SetSkip(offset)

	filter := bson.M{}
	match := []bson.M{}
	match = append(match, bson.M{"deleted_at": nil})
	if payload.Query != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"name.ar": bson.M{"$regex": payload.Query, "$options": "i"}},
				{"name.en": bson.M{"$regex": payload.Query, "$options": "i"}},
				{"email": bson.M{"$regex": payload.Query, "$options": "i"}},
			},
		}
	}
	filter["$and"] = match

	// Query the collection for the total count of documents
	collection := mgm.Coll(&domain.Order{})
	totalItems, err := collection.CountDocuments(ctx, filter)

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalItems) / float64(payload.Limit)))

	var data []domain.Order
	err = l.orderCollection.SimpleFind(&data, filter, findOptions)

	return data, utils.PaginationResult{Page: payload.Page, TotalPages: int64(totalPages), TotalItems: totalItems}, err

}
