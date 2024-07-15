package mongodb

import (
	"context"
	"fmt"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/dto/order"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"time"
)

type OrderRepository struct {
	orderCollection *mgm.Collection
	logger          logger.ILogger
}

const mongoOrderRepositoryTag = "OrderMongoRepository"

func NewOrderMongoRepository(dbs *mongo.Database, logger logger.ILogger) domain.OrderRepository {
	orderDbCollection := mgm.Coll(&domain.Order{})
	return &OrderRepository{
		orderCollection: orderDbCollection,
		logger:          logger,
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

	if dto.Status != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"status": dto.Status})
	}

	if dto.LocationId != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"location._id": utils.ConvertStringIdToObjectId(dto.LocationId)})
	}
	if dto.AccountId != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"location.account_id": utils.ConvertStringIdToObjectId(dto.AccountId)})
	}
	if dto.BrandId != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"location.brand_details._id": utils.ConvertStringIdToObjectId(dto.BrandId)})
	}
	if dto.CountryId != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"location.country._id": dto.CountryId})
	}
	if dto.SerialNum != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"serial_num": dto.SerialNum})
	}
	if dto.UserId != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"user._id": utils.ConvertStringIdToObjectId(dto.UserId)})
	}

	if dto.From != "" {
		fromDate, dateErr := time.Parse(time.DateTime, dto.From)
		if dateErr != nil {
			i.logger.Error("Order Repo -> parsing date -> ", dateErr)
			err = dateErr
			return
		}
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"created_at": bson.M{"$gte": fromDate}})
	}

	if dto.To != "" {
		toDate, dateErr := time.Parse(time.DateTime, dto.To)
		if dateErr != nil {
			i.logger.Error("Order Repo -> parsing date -> ", dateErr)
			err = dateErr
			return
		}
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"created_at": bson.M{"$lte": toDate}})
	}

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

func (l OrderRepository) FindOrder(ctx *context.Context, id string, userId string) (order *domain.Order, err error) {
	var orderDomain domain.Order
	filter := bson.M{"user._id": utils.ConvertStringIdToObjectId(userId), "_id": utils.ConvertStringIdToObjectId(id)}
	err = l.orderCollection.FirstWithCtx(*ctx, filter, &orderDomain)
	return &orderDomain, err
}

func (l OrderRepository) UpdateOrderStatus(ctx *context.Context, orderDomain *domain.Order, previousStatus []string, statusLog domain.StatusLog, updateSet interface{}) (order *domain.Order, err error) {
	fmt.Println("previousStatus => ", previousStatus)

	filter := bson.M{"_id": orderDomain.ID, "status": bson.M{"$in": previousStatus}}
	_, err = l.orderCollection.UpdateOne(*ctx, filter, bson.M{
		"$set":  updateSet,
		"$push": bson.M{"status_logs": statusLog},
	})
	fmt.Println("error Repo => ", err)
	if err != nil {
		return nil, err
	}
	return l.FindOrder(ctx, utils.ConvertObjectIdToStringId(orderDomain.ID), utils.ConvertObjectIdToStringId(orderDomain.User.ID))

}
