package menu_group_item

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/menu_group"
	"samm/internal/module/menu/repository/structs/menu_group_item"
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

func (i *menuGroupItemRepo) SyncMenuItemsChanges(ctx context.Context, itemDoc menu_group_item.MenuGroupItemSyncItemModel) error {
	filter := bson.M{"item_id": itemDoc.ItemId}
	_, err := i.menuGroupItemCollection.UpdateMany(ctx, filter, bson.M{"$set": itemDoc})
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
