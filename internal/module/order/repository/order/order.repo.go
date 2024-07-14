package mongodb

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/dto/order"
	"samm/pkg/logger"
	"time"
)

type OrderRepository struct {
	orderCollection *mgm.Collection
	logger          logger.ILogger
}

const mongoOrderRepositoryTag = "OrderMongoRepository"

func NewOrderMongoRepository(dbs *mongo.Database) domain.OrderRepository {
	orderDbCollection := mgm.Coll(&domain.Order{})

	return &OrderRepository{
		orderCollection: orderDbCollection,
	}
}

func (l OrderRepository) CreateOrder(ctx context.Context, order *domain.Order) (err error) {
	_, err = mgm.Coll(&domain.Order{}).InsertOne(ctx, order)
	if err != nil {
		return err
	}
	return nil

}

func (l OrderRepository) UpdateOrder(ctx context.Context, order *domain.Order) (err error) {
	update := bson.M{"$set": order}
	//upsert := false
	//opts := options.UpdateOptions{Upsert: &upsert}
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

func (i *OrderRepository) ListOrderForDashboard(ctx *context.Context, dto *order.ListOrderDto) (ordersRes *[]domain.OrderV2, paginationMeta *PaginationData, err error) {
	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.D{{"deleted_at", nil}},
	}}}

	if dto.Query != "" {
		pattern := ".*" + dto.Query + ".*"
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"$or": []bson.M{{"name": bson.M{"$regex": pattern, "$options": "i"}}, {"phone_number": bson.M{"$regex": pattern, "$options": "i"}}}})
	}

	//if dto.Status != "" {
	//	if dto.Status == "active" {
	//		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"is_active": true})
	//	} else if dto.Status == "inactive" {
	//		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"is_active": false})
	//	}
	//}
	//if dto.Dob != "" {
	//	matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"dob": dto.Dob})
	//}

	data, err := New(i.orderCollection.Collection).Context(*ctx).Limit(dto.Limit).Page(dto.Page).Sort("created_at", -1).Aggregate(matching)

	if data == nil || data.Data == nil {
		return nil, nil, err
	}

	orders := make([]domain.OrderV2, 0)
	for _, raw := range data.Data {
		model := domain.OrderV2{}
		err = bson.Unmarshal(raw, &model)
		if err != nil {
			i.logger.Error("Order Repo -> List -> ", err)
			break
		}
		orders = append(orders, model)
	}
	paginationMeta = &data.Pagination
	ordersRes = &orders

	return
}
