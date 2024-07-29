package menu_group_item

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/menu_group"
	menu_group2 "samm/internal/module/menu/repository/structs/menu_group"
	"samm/internal/module/menu/repository/structs/menu_group_item"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/utils/dto"
	"samm/pkg/validators/localization"
)

type menuGroupItemRepo struct {
	menuGroupItemCollection *mgm.Collection
	logger                  logger.ILogger
}

func NewMenuGroupItemRepository(dbs *mongo.Database, log logger.ILogger) domain.MenuGroupItemRepository {
	menuGroupItemCollection := mgm.Coll(&domain.MenuGroupItem{})
	return &menuGroupItemRepo{
		menuGroupItemCollection: menuGroupItemCollection,
		logger:                  log,
	}
}

func (i *menuGroupItemRepo) CreateUpdateBulk(ctx context.Context, models *[]domain.MenuGroupItem) error {
	var bulkOperations []mongo.WriteModel
	for _, update := range *models {
		filter := bson.M{"_id": update.ID}
		updateDoc := bson.M{"$set": update}
		updateModel := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(updateDoc).SetUpsert(true)
		bulkOperations = append(bulkOperations, updateModel)
	}
	_, err := i.menuGroupItemCollection.BulkWrite(ctx, bulkOperations)
	if err != nil {
		return err
	}
	return nil
}

func (i *menuGroupItemRepo) SyncMenuItemsChanges(ctx context.Context, itemDoc menu_group_item.MenuGroupItemSyncItemModel) error {
	filter := bson.M{"item_id": itemDoc.ItemId}
	_, err := i.menuGroupItemCollection.UpdateMany(ctx, filter, bson.M{"$set": itemDoc})
	if err != nil {
		return err
	}
	return nil
}

func (i *menuGroupItemRepo) FindMenuGroupItem(ctx context.Context, id primitive.ObjectID) (domain.MenuGroupItem, error) {
	result := domain.MenuGroupItem{}
	err := i.menuGroupItemCollection.FindByIDWithCtx(ctx, id, &result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (i *menuGroupItemRepo) DeleteByItemId(ctx context.Context, itemId primitive.ObjectID) error {
	filter := bson.M{"item_id": itemId}
	_, err := i.menuGroupItemCollection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (i *menuGroupItemRepo) ChangeStatusByItemId(ctx context.Context, itemId primitive.ObjectID, model domain.MenuGroupItem) error {
	_, err := i.menuGroupItemCollection.UpdateByID(ctx, itemId, bson.M{"$set": model})
	if err != nil {
		return err
	}
	return nil
}

func (i *menuGroupItemRepo) DeleteBulkByGroupMenuId(ctx context.Context, groupMenuId primitive.ObjectID, exceptionIds []primitive.ObjectID) error {
	_, err := i.menuGroupItemCollection.DeleteMany(ctx, bson.M{"menu_group._id": groupMenuId, "_id": bson.M{"$nin": utils.If(exceptionIds != nil, exceptionIds, make([]primitive.ObjectID, 0)).([]primitive.ObjectID)}})
	if err != nil {
		return err
	}
	return nil
}

func (i *menuGroupItemRepo) ChangeMenuStatus(ctx context.Context, id primitive.ObjectID, dto *menu_group.ChangeMenuGroupStatusDto, adminDetails dto.AdminDetails) error {
	update := bson.M{"$set": bson.M{"menu_group.status": dto.Status}, "$push": bson.M{"admin_details": adminDetails}}
	_, err := i.menuGroupItemCollection.UpdateMany(ctx, bson.M{"menu_group._id": id}, update)
	if err != nil {
		return err
	}
	return nil
}

func (i *menuGroupItemRepo) ChangeCategoryStatus(ctx context.Context, id primitive.ObjectID, dto *menu_group.ChangeMenuGroupStatusDto, adminDetails dto.AdminDetails) error {
	update := bson.M{"$set": bson.M{"category.status": dto.Status}, "$push": bson.M{"admin_details": adminDetails}}
	_, err := i.menuGroupItemCollection.UpdateMany(ctx, bson.M{"menu_group._id": id, "category._id": utils.ConvertStringIdToObjectId(dto.Id)}, update)
	if err != nil {
		return err
	}
	return nil
}

func (i *menuGroupItemRepo) ChangeItemStatus(ctx context.Context, id primitive.ObjectID, dto *menu_group.ChangeMenuGroupStatusDto, adminDetails dto.AdminDetails) error {
	update := bson.M{"$set": bson.M{"status": dto.Status}, "$push": bson.M{"admin_details": adminDetails}}
	_, err := i.menuGroupItemCollection.UpdateOne(ctx, bson.M{"menu_group._id": id, "_id": utils.ConvertStringIdToObjectId(dto.Id)}, update)
	if err != nil {
		return err
	}
	return nil
}

func (i *menuGroupItemRepo) DeleteByCategory(ctx context.Context, dto *menu_group.DeleteEntityFromMenuGroupDto) error {
	_, err := i.menuGroupItemCollection.DeleteMany(ctx, bson.M{"menu_group._id": utils.ConvertStringIdToObjectId(dto.Id), "category._id": utils.ConvertStringIdToObjectId(dto.EntityId)})
	if err != nil {
		return err
	}
	return nil
}

func (i *menuGroupItemRepo) Delete(ctx context.Context, dto *menu_group.DeleteEntityFromMenuGroupDto) error {
	_, err := i.menuGroupItemCollection.DeleteOne(ctx, bson.M{"menu_group._id": utils.ConvertStringIdToObjectId(dto.Id), "_id": utils.ConvertStringIdToObjectId(dto.EntityId)})
	if err != nil {
		return err
	}
	return nil
}

func (i *menuGroupItemRepo) MobileGetMenuGroupItems(ctx context.Context, dto *menu_group.GetMenuGroupItemsDTO) (*[]menu_group2.MobileGetMenuGroupItems, error) {
	var pipeline []interface{}

	matching := bson.M{
		"$match": bson.M{"$and": []interface{}{
			bson.D{{"menu_group.status", "active"}},
			bson.D{{"category.status", "active"}},
			bson.D{{"status", "active"}},
			AddAvailabilityQuery(dto.CountryId, "$menu_group.availabilities"),
			AddAvailabilityQuery(dto.CountryId, "$availabilities"),
		}},
	}

	if dto.LocationId != primitive.NilObjectID.Hex() {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.D{{"menu_group.location_ids", utils.ConvertStringIdToObjectId(dto.LocationId)}})
	}

	if dto.Query != "" {
		pattern := ".*" + dto.Query + ".*"
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.D{{"$or", []bson.M{{"name.ar": bson.M{"$regex": pattern, "$options": "i"}}, {"name.en": bson.M{"$regex": pattern, "$options": "i"}}}}})
	}

	sort := bson.M{"$sort": bson.M{"category.sort": 1, "category._id": 1, "sort": 1, "_id": 1}}

	group := bson.M{
		"$group": bson.D{
			{Key: "_id", Value: "$category._id"},
			{Key: "name", Value: bson.D{{Key: "$first", Value: "$category.name"}}},
			{Key: "icon", Value: bson.D{{Key: "$first", Value: "$category.icon"}}},
			{Key: "sort", Value: bson.D{{Key: "$first", Value: "$category.sort"}}},
			{Key: "items", Value: bson.D{{Key: "$push", Value: bson.D{
				{Key: "_id", Value: "$_id"},
				{Key: "item_id", Value: "$item_id"},
				{Key: "name", Value: "$name"},
				{Key: "image", Value: "$image"},
				{Key: "calories", Value: "$calories"},
				{Key: "price", Value: "$price"},
				{Key: "sort", Value: "$sort"},
			}}}},
		},
	}

	project := bson.M{
		"$project": bson.D{
			{Key: "_id", Value: 1},
			{Key: "name", Value: 1},
			{Key: "icon", Value: 1},
			{Key: "sort", Value: 1},
			{Key: "items", Value: bson.D{{Key: "$map", Value: bson.D{
				{Key: "input", Value: bson.D{{Key: "$sortArray", Value: bson.D{
					{Key: "input", Value: "$items"},
					{Key: "sortBy", Value: bson.D{{Key: "sort", Value: 1}}},
				}}}},
				{Key: "as", Value: "item"},
				{Key: "in", Value: "$$item"},
			}}}},
		},
	}

	pipeline = append(pipeline, matching, group, project, sort)

	items := make([]menu_group2.MobileGetMenuGroupItems, 0)

	cursor, err := i.menuGroupItemCollection.Aggregate(ctx, pipeline)
	if err != nil {
		i.logger.Error("menuGroupItemRepo -> MobileGetMenuGroupItems -> ", err)
		return &items, err
	}

	if err = cursor.All(ctx, &items); err != nil {
		i.logger.Error("menuGroupItemRepo -> MobileGetMenuGroupItems -> ", err)
		return &items, err
	}

	return &items, nil
}

func (i *menuGroupItemRepo) MobileGetMenuGroupItem(ctx context.Context, dto *menu_group.GetMenuGroupItemDTO) (*menu_group2.MobileGetItem, error) {
	var pipeline []interface{}

	matching := bson.M{
		"$match": bson.M{"$and": []interface{}{
			bson.D{{"_id", utils.ConvertStringIdToObjectId(dto.ID)}},
			bson.D{{"menu_group.location_ids", utils.ConvertStringIdToObjectId(dto.LocationId)}},
			bson.D{{"menu_group.status", "active"}},
			bson.D{{"category.status", "active"}},
			bson.D{{"status", "active"}},
			AddAvailabilityQuery(dto.CountryId, "$menu_group.availabilities"),
			AddAvailabilityQuery(dto.CountryId, "$availabilities"),
		}},
	}

	modifierGroupLookup := bson.M{
		"$lookup": bson.M{
			"from":         "modifier_groups",
			"localField":   "modifier_group_ids",
			"foreignField": "_id",
			"pipeline": []bson.M{
				{"$match": bson.M{"status": "active", "deleted_at": nil}},
			},
			"as": "modifier_groups",
		},
	}

	addonsLookup := bson.M{
		"$lookup": bson.M{
			"from":         "items",
			"localField":   "modifier_groups.product_ids",
			"foreignField": "_id",
			"pipeline": []bson.M{
				{"$match": bson.M{"status": "active", "deleted_at": nil}},
			},
			"as": "addons",
		},
	}

	limit := bson.M{"$limit": 1}

	pipeline = append(pipeline, matching, modifierGroupLookup, addonsLookup, limit)

	item := menu_group2.MobileGetItem{}

	cursor, errDB := i.menuGroupItemCollection.Aggregate(ctx, pipeline)
	if errDB != nil {
		i.logger.Error("menuGroupItemRepo -> MobileGetMenuGroupItem -> ", errDB)
		return &item, errDB
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&item); err != nil {
			i.logger.Error("menuGroupItemRepo -> MobileGetMenuGroupItem -> Decode -> ", err)
			return &item, err
		}
	} else {
		i.logger.Error("menuGroupRepo -> MobileGetMenuGroupItem -> ", cursor.Err())
		return &item, errors.New(localization.GetTranslation(&ctx, localization.E1000, nil, ""))
	}

	return &item, nil
}

func (i *menuGroupItemRepo) MobileFilterMenuGroupItemForOrder(ctx context.Context, dto *menu_group.FilterMenuGroupItemsForOrder) ([]menu_group2.MobileGetItem, error) {
	productIds := getProductIds(dto)

	matching := bson.M{
		"$match": bson.M{"$and": []interface{}{
			bson.D{{"_id", bson.M{"$in": productIds}}},
			bson.D{{"menu_group.location_ids", utils.ConvertStringIdToObjectId(dto.LocationId)}},
			bson.D{{"menu_group.status", "active"}},
			bson.D{{"category.status", "active"}},
			bson.D{{"status", "active"}},
			AddAvailabilityQuery(dto.CountryId, "$menu_group.availabilities"),
			AddAvailabilityQuery(dto.CountryId, "$availabilities"),
		}},
	}

	modifierGroupLookup := bson.M{
		"$lookup": bson.M{
			"from":         "modifier_groups",
			"localField":   "modifier_group_ids",
			"foreignField": "_id",
			"pipeline": []bson.M{
				{"$match": bson.M{"status": "active", "deleted_at": nil}},
			},
			"as": "modifier_groups",
		},
	}

	addonsLookup := bson.M{
		"$lookup": bson.M{
			"from":         "items",
			"localField":   "modifier_groups.product_ids",
			"foreignField": "_id",
			"pipeline": []bson.M{
				{"$match": bson.M{"status": "active", "deleted_at": nil}},
			},
			"as": "addons",
		},
	}

	products := make([]menu_group2.MobileGetItem, 0)

	err := i.menuGroupItemCollection.SimpleAggregateWithCtx(ctx, &products, matching, modifierGroupLookup, addonsLookup)
	if err != nil {
		i.logger.Error("menuGroupItemRepo -> MobileGetMenuGroupItem -> ", err)
		return nil, err
	}
	return eliminateUnusedModifiers(dto, &products), nil
}

func eliminateUnusedModifiers(order *menu_group.FilterMenuGroupItemsForOrder, products *[]menu_group2.MobileGetItem) []menu_group2.MobileGetItem {
	for _, dtoItem := range order.MenuItems {
		for _, product := range *products {
			if dtoItem.Id == product.ID.Hex() {
				//get all the modifier inside one product
				dtoProductModifierIds := make([]string, 0)
				for _, modifier := range dtoItem.ModifierIds {
					dtoProductModifierIds = append(dtoProductModifierIds, modifier.Id)
				}
				for productModifierGroupIndex, group := range product.ModifierGroups {
					//remove modifierIds in group
					groupProductIds, _ := utils.EqualizeSlices(utils.ConvertObjectIdsToStringIds(group.ProductIds), dtoProductModifierIds)
					//product.ModifierGroups[productModifierGroupIndex].ProductIds = utils.ConvertStringIdsToObjectIds(groupProductIds)
					//check if modifier group is empty
					if len(groupProductIds) <= 0 {
						product.ModifierGroups = append(product.ModifierGroups[:productModifierGroupIndex], product.ModifierGroups[productModifierGroupIndex+1:]...)
					}

				}
			}
		}
	}
	return *products
}
