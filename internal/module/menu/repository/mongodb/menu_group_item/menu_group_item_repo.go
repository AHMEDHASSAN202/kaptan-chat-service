package menu_group_item

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/menu_group"
	menu_group2 "samm/internal/module/menu/repository/structs/menu_group"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/utils/dto"
)

type menuGroupItemRepo struct {
	menuGroupItemCollection *mgm.Collection
}

func NewMenuGroupItemRepository(dbs *mongo.Database) domain.MenuGroupItemRepository {
	menuGroupItemCollection := mgm.Coll(&domain.MenuGroupItem{})
	return &menuGroupItemRepo{
		menuGroupItemCollection: menuGroupItemCollection,
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

func (i *menuGroupItemRepo) MobileGetMenuGroupItems(ctx context.Context, dto *menu_group.GetMenuGroupItemDTO) (*[]menu_group2.MobileGetMenuGroupItems, error) {
	var pipeline []interface{}

	matching := bson.M{
		"$match": bson.M{"$and": []interface{}{
			bson.D{{"menu_group.branch_ids", utils.ConvertStringIdToObjectId(dto.BranchId)}},
			bson.D{{"menu_group.status", "active"}},
			bson.D{{"category.status", "active"}},
			bson.D{{"status", "active"}},
			AddAvailabilityQuery(dto.CountryId, "$menu_group.availabilities"),
			AddAvailabilityQuery(dto.CountryId, "$availabilities"),
		}},
	}

	if dto.Query != "" {
		pattern := ".*" + dto.Query + ".*"
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.D{{"$or", []bson.M{{"name.ar": bson.M{"$regex": pattern, "$options": "i"}}, {"name.en": bson.M{"$regex": pattern, "$options": "i"}}}}})
	}

	sort := bson.M{"$sort": bson.M{"sort": 1}}

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
		logger.Logger.Error("menuGroupItemRepo -> MobileGetMenuGroupItems -> ", err)
		return &items, err
	}

	if err = cursor.All(ctx, &items); err != nil {
		logger.Logger.Error("menuGroupItemRepo -> MobileGetMenuGroupItems -> ", err)
		return &items, err
	}

	return &items, nil
}
