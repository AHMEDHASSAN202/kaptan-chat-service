package mongodb

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/menu/domain"
	"samm/pkg/utils"
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
