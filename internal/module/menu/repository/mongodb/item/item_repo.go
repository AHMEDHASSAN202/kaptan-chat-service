package item

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/item"
	"samm/pkg/database/mongodb"
	"samm/pkg/utils"
	"time"
)

type itemRepo struct {
	itemCollection *mgm.Collection
}

func NewItemRepository(dbs *mongo.Database) domain.ItemRepository {
	itemCollection := mgm.Coll(&domain.Item{})
	//text search menu cuisine
	mongodb.CreateIndex(itemCollection.Collection, false, bson.E{"name.ar", mongodb.IndexType.Text}, bson.E{"name.en", mongodb.IndexType.Text}, bson.E{"tags", mongodb.IndexType.Text},
		bson.E{"desc.ar", mongodb.IndexType.Text}, bson.E{"desc.en", mongodb.IndexType.Text})
	//make sure there are no duplicated menu cuisine
	mongodb.CreateIndex(itemCollection.Collection, true, bson.E{"name.ar", mongodb.IndexType.Asc}, bson.E{"name.en", mongodb.IndexType.Asc}, bson.E{"account_id", mongodb.IndexType.Asc}, bson.E{"deleted_at", mongodb.IndexType.Asc})
	return &itemRepo{
		itemCollection: itemCollection,
	}
}

func (i *itemRepo) GetByIds(ctx context.Context, ids []primitive.ObjectID) ([]domain.Item, error) {
	items := []domain.Item{}
	err := i.itemCollection.SimpleFind(&items, bson.M{"_id": bson.M{"$in": ids}})
	return items, err
}

func (i *itemRepo) Create(ctx context.Context, doc []domain.Item) error {
	_, err := i.itemCollection.InsertMany(ctx, utils.ConvertArrStructToInterfaceArr(doc))
	if err != nil {
		return err
	}
	return nil
}

func (i *itemRepo) Update(ctx context.Context, id *primitive.ObjectID, doc *domain.Item) error {
	update := bson.M{"$set": doc}
	_, err := i.itemCollection.UpdateByID(ctx, &id, update)
	if err != nil {
		return err
	}
	return nil
}

func (i *itemRepo) SoftDelete(ctx context.Context, id *primitive.ObjectID) error {
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	_, err := i.itemCollection.UpdateByID(ctx, &id, update)
	if err != nil {
		return err
	}
	return nil
}
func (i *itemRepo) ChangeStatus(ctx context.Context, id *primitive.ObjectID, dto *item.ChangeItemStatusDto) error {
	update := bson.M{"$set": bson.M{"status": dto.Status, "admin_details": dto.AdminDetails}}
	_, err := i.itemCollection.UpdateByID(ctx, &id, update)
	if err != nil {
		return err
	}
	return nil
}

func (i *itemRepo) List(ctx context.Context, query *item.ListItemsDto) ([]domain.Item, error) {
	filter := bson.M{"$text": bson.M{
		"$search": query.Query}}
	items := []domain.Item{}
	err := i.itemCollection.SimpleFindWithCtx(ctx, &items, filter)
	return items, err
}
