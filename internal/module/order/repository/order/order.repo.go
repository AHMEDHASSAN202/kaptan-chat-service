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
	"samm/internal/module/order/repository/structs"
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

func (l OrderRepository) FindOrder(ctx *context.Context, Id primitive.ObjectID) (*domain.Order, error) {
	var domainData domain.Order
	filter := bson.M{"_id": Id}
	err := l.orderCollection.FirstWithCtx(*ctx, filter, &domainData)
	return &domainData, err
}

func (l OrderRepository) FindOrderForMobile(ctx *context.Context, Id primitive.ObjectID) (*structs.MobileFindOrder, error) {
	var orderData structs.MobileFindOrder
	filter := bson.M{"_id": Id}
	err := l.orderCollection.FirstWithCtx(*ctx, filter, &orderData)
	return &orderData, err
}

func (l OrderRepository) UpdateOrder(order *domain.Order) (err error) {
	//upsert := true
	//opts := options.UpdateOptions{Upsert: &upsert}
	err = l.orderCollection.Update(order)
	return
}

func (i *OrderRepository) ListOrderForDashboard(ctx *context.Context, dto *order.ListOrderDtoForDashboard) (ordersRes *[]domain.Order, paginationMeta *PaginationData, err error) {
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
	if dto.IsFavourite {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"is_favourite": true})
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

func (i *OrderRepository) ListOrderForMobile(ctx *context.Context, dto *order.ListOrderDtoForMobile) (ordersRes *[]structs.MobileListOrders, paginationMeta *PaginationData, err error) {
	threeMonthsAgo := time.Now().UTC().AddDate(0, -3, 0)
	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.M{"created_at": bson.M{"$gte": threeMonthsAgo}},
		bson.M{"user._id": utils.ConvertStringIdToObjectId(dto.UserId)},
	}}}

	data, err := New(i.orderCollection.Collection).Context(*ctx).Limit(dto.Limit).Page(dto.Page).Sort("created_at", -1).Aggregate(matching)

	if data == nil || data.Data == nil {
		return nil, nil, err
	}

	orders := make([]structs.MobileListOrders, 0)
	for _, raw := range data.Data {
		model := structs.MobileListOrders{}
		err = bson.Unmarshal(raw, &model)
		if err != nil {
			i.logger.Error("Order Repo -> Mobile List -> ", err)
			break
		}
		orders = append(orders, model)
	}
	paginationMeta = &data.Pagination
	ordersRes = &orders

	return
}
