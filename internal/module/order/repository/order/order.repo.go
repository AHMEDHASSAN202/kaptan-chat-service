package mongodb

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/order/domain"
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

func (l OrderRepository) StoreOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	_, err := mgm.Coll(&domain.Order{}).InsertOne(ctx, order)
	if err != nil {
		return nil, err
	}
	return order, nil

}

func (l OrderRepository) UserHasOrders(ctx context.Context, userId primitive.ObjectID, orderStatus []string) (bool, error) {
	ordersCount, err := mgm.Coll(&domain.Order{}).CountDocuments(ctx, bson.M{"user._id": userId, "status": bson.M{"$in": orderStatus}})
	if err != nil {
		return false, err
	}
	return ordersCount >= 1, nil
}
