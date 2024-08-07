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
	"samm/internal/module/order/dto/order/kitchen"
	"samm/internal/module/order/repository/structs"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"time"
)

type OrderRepository struct {
	orderCollection *mgm.Collection
	logger          logger.ILogger
}

const (
	MongoOrderRepositoryTag = "OrderMongoRepository"
	MaxLimitPerRound        = 500
)

func NewOrderMongoRepository(dbs *mongo.Database, logger logger.ILogger) domain.OrderRepository {
	orderDbCollection := mgm.Coll(&domain.Order{})
	return &OrderRepository{
		orderCollection: orderDbCollection,
		logger:          logger,
	}
}

func (l OrderRepository) StoreOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	_, err := mgm.Coll(&domain.Order{}).InsertOne(ctx, order)
	if err != nil {
		return nil, err
	}
	return order, nil

}

func (l OrderRepository) UserHasOrders(ctx context.Context, userId primitive.ObjectID, orderStatus []string, gt int64) (bool, error) {
	ordersCount, err := mgm.Coll(&domain.Order{}).CountDocuments(ctx, bson.M{"user._id": userId, "status": bson.M{"$in": orderStatus}})
	if err != nil {
		return false, err
	}
	return ordersCount >= gt, nil
}

func (l OrderRepository) FindOrder(ctx *context.Context, Id primitive.ObjectID) (*domain.Order, error) {
	var domainData domain.Order
	filter := bson.M{"_id": Id}
	err := l.orderCollection.FirstWithCtx(*ctx, filter, &domainData)
	return &domainData, err
}

func (l OrderRepository) UpdateOrder(ctx context.Context, order *domain.Order) (err error) {
	//upsert := true
	//opts := options.UpdateOptions{Upsert: &upsert}
	err = l.orderCollection.UpdateWithCtx(ctx, order)
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
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"location.account._id": utils.ConvertStringIdToObjectId(dto.AccountId)})
	}
	if dto.BrandId != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"location.brand_details._id": utils.ConvertStringIdToObjectId(dto.BrandId)})
	}
	if dto.CountryId != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"location.country._id": dto.CountryId})
	}
	if dto.Query != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"$or": []bson.M{{"serial_num": dto.Query}, {"user.phone_number": dto.Query}}})
	}
	if dto.IsFavourite {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"is_favourite": true})
	}

	if dto.From != "" {
		fromDate, dateErr := time.Parse(time.RFC3339, dto.From)
		if dateErr != nil {
			i.logger.Error("Order Repo -> parsing date -> ", dateErr)
			err = dateErr
			return
		}
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"created_at": bson.M{"$gte": fromDate}})
	}

	if dto.To != "" {
		toDate, dateErr := time.Parse(time.RFC3339, dto.To)
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

func (i *OrderRepository) ListInprogressOrdersForMobile(ctx *context.Context, dto *order.ListOrderDtoForMobile) (ordersRes *[]structs.MobileListOrders, paginationMeta *PaginationData, err error) {
	threeMonthsAgo := time.Now().UTC().AddDate(0, -3, 0)
	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.M{"created_at": bson.M{"$gte": threeMonthsAgo}},
		bson.M{"status": bson.M{"$in": InProgressStatuses}},
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
func (i *OrderRepository) ListRunningOrdersForKitchen(ctx *context.Context, dto *kitchen.ListRunningOrderDto) (ordersRes *[]structs.MobileListOrders, paginationMeta *PaginationData, err error) {

	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.M{"status": bson.M{"$in": dto.Status}},
		bson.M{"meta_data.target_kitchen_ids": utils.ConvertStringIdToObjectId(dto.CauserKitchenId)},
	}}}

	//add time limit condition
	if dto.NumberOfHoursLimit > 0 {
		eightHoursAgo := time.Now().UTC().Add(time.Duration(-1*dto.NumberOfHoursLimit) * time.Hour)
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"created_at": bson.M{"$gte": eightHoursAgo}})
	}

	if dto.Pagination.Pagination {
		return executeListWithPagination(ctx, i, dto, matching)
	} else {
		orders := make([]structs.MobileListOrders, 0)
		err = i.orderCollection.SimpleAggregate(&orders, matching)
		ordersRes = &orders
		return ordersRes, nil, err
	}

}

func executeListWithPagination(ctx *context.Context, i *OrderRepository, dto *kitchen.ListRunningOrderDto, matching bson.M) (ordersRes *[]structs.MobileListOrders, paginationMeta *PaginationData, err error) {

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

func (i *OrderRepository) ListCompletedOrdersForMobile(ctx *context.Context, dto *order.ListOrderDtoForMobile) (ordersRes *[]structs.MobileListOrders, paginationMeta *PaginationData, err error) {
	threeMonthsAgo := time.Now().UTC().AddDate(0, -3, 0)
	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.M{"created_at": bson.M{"$gte": threeMonthsAgo}},
		bson.M{"status": bson.M{"$in": CompletedStatuses}},
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

func (i *OrderRepository) ListLastOrdersForMobile(ctx *context.Context, dto *order.ListOrderDtoForMobile) (ordersRes *[]structs.MobileListOrders, paginationMeta *PaginationData, err error) {
	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.M{"status": bson.M{"$in": CompletedStatuses}},
		bson.M{"user._id": utils.ConvertStringIdToObjectId(dto.UserId)},
	}}}

	data, err := New(i.orderCollection.Collection).Context(*ctx).Limit(3).Page(dto.Page).Sort("is_favourite", -1).Sort("created_at", -1).Aggregate(matching)

	if data == nil || data.Data == nil {
		return nil, nil, err
	}

	orders := make([]structs.MobileListOrders, 0)
	for _, raw := range data.Data {
		model := structs.MobileListOrders{}
		err = bson.Unmarshal(raw, &model)
		if err != nil {
			i.logger.Error("Order Repo -> Mobile List Last Orders -> ", err)
			break
		}
		orders = append(orders, model)
	}
	paginationMeta = &data.Pagination
	ordersRes = &orders

	return
}

func (l OrderRepository) FindOrderByUser(ctx *context.Context, id string, userId string) (order *domain.Order, err error) {
	var orderDomain domain.Order
	filter := bson.M{"user._id": utils.ConvertStringIdToObjectId(userId), "_id": utils.ConvertStringIdToObjectId(id)}
	err = l.orderCollection.FirstWithCtx(*ctx, filter, &orderDomain)
	return &orderDomain, err
}

func (l OrderRepository) UpdateOrderStatus(ctx *context.Context, orderDomain *domain.Order, previousStatus []string, statusLog *domain.StatusLog, updateSet interface{}) (order *domain.Order, err error) {
	filter := bson.M{"_id": orderDomain.ID}

	if len(previousStatus) > 0 {
		filter = bson.M{"_id": orderDomain.ID, "status": bson.M{"$in": previousStatus}}
	}
	if statusLog != nil {
		_, err = l.orderCollection.UpdateOne(*ctx, filter, bson.M{
			"$set":  updateSet,
			"$push": bson.M{"status_logs": statusLog},
		})
	} else {
		_, err = l.orderCollection.UpdateOne(*ctx, filter, bson.M{
			"$set": updateSet,
		})
	}

	if err != nil {
		return nil, err
	}
	return l.FindOrderByUser(ctx, utils.ConvertObjectIdToStringId(orderDomain.ID), utils.ConvertObjectIdToStringId(orderDomain.User.ID))

}

func (l OrderRepository) UpdateUserAllOrdersFavorite(ctx context.Context, userId string) (err error) {
	filter := bson.M{"user._id": utils.ConvertStringIdToObjectId(userId), "is_favourite": true}
	_, err = l.orderCollection.UpdateMany(ctx, filter, bson.M{
		"$set": bson.M{"is_favourite": false},
	})
	return
}

func (i *OrderRepository) GetAllOrdersForCronJobs(ctx *context.Context, filters bson.M) (ordersRes *[]domain.Order, paginationMeta *PaginationData, err error) {
	data, err := New(i.orderCollection.Collection).Context(*ctx).Limit(MaxLimitPerRound).Page(0).Sort("created_at", -1).Aggregate(filters)

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
