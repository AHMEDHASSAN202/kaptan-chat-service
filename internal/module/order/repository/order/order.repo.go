package mongodb

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/dto/order"
	"samm/pkg/logger"
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

func (l OrderRepository) StoreOrder(ctx *context.Context, order *domain.Order) (err error) {
	_, err = l.orderCollection.InsertOne(*ctx, order)
	if err != nil {
		return err
	}
	return nil
}

func (i *OrderRepository) ListOrderForDashboard(ctx *context.Context, dto *order.ListOrderDto) (ordersRes *[]domain.Order, paginationMeta *PaginationData, err error) {
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

	orders := make([]domain.Order, 0)
	for _, raw := range data.Data {
		model := domain.Order{}
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
