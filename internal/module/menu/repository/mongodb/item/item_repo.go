package item

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/jinzhu/copier"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/item"
	"samm/internal/module/menu/repository/structs/menu_group_item"
	responseItem "samm/internal/module/menu/responses/item"
	"samm/pkg/database/mongodb"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"time"
)

type itemRepo struct {
	itemCollection    *mgm.Collection
	logger            logger.ILogger
	menuGroupItemRepo domain.MenuGroupItemRepository
}

func NewItemRepository(dbs *mongo.Database, logger logger.ILogger, menuGroupItemRepo domain.MenuGroupItemRepository) domain.ItemRepository {
	itemCollection := mgm.Coll(&domain.Item{})
	//text search menu cuisine
	mongodb.CreateIndex(itemCollection.Collection, false, bson.E{"name.ar", mongodb.IndexType.Text}, bson.E{"name.en", mongodb.IndexType.Text}, bson.E{"tags", mongodb.IndexType.Text},
		bson.E{"desc.ar", mongodb.IndexType.Text}, bson.E{"desc.en", mongodb.IndexType.Text})
	//make sure there are no duplicated menu cuisine
	mongodb.CreateIndex(itemCollection.Collection, true, bson.E{"name.ar", mongodb.IndexType.Asc}, bson.E{"name.en", mongodb.IndexType.Asc}, bson.E{"account_id", mongodb.IndexType.Asc}, bson.E{"deleted_at", mongodb.IndexType.Asc})
	return &itemRepo{
		itemCollection:    itemCollection,
		logger:            logger,
		menuGroupItemRepo: menuGroupItemRepo,
	}
}

func (i *itemRepo) GetByIds(ctx context.Context, ids []primitive.ObjectID) ([]domain.Item, error) {
	items := []domain.Item{}
	filter := bson.M{"_id": bson.M{"$in": ids}, "deleted_at": nil}
	err := i.itemCollection.SimpleFind(&items, filter)
	if err != nil {
		return items, err
	}
	if len(items) <= 0 {
		return items, mongo.ErrNoDocuments
	}
	return items, nil
}
func (i *itemRepo) Find(ctx context.Context, id primitive.ObjectID) (responseItem.ItemResponse, error) {
	items := responseItem.ItemResponse{}
	filter := bson.M{"_id": id, "deleted_at": nil}

	_, err := i.itemCollection.SimpleAggregateFirst(&items, bson.M{"$match": filter}, bson.M{"$lookup": bson.M{
		"from":         "modifier_groups",
		"localField":   "modifier_groups_ids",
		"foreignField": "_id",
		"as":           "modifier_groups_ids",
	}}, bson.M{"$lookup": bson.M{
		"from":         "skus",
		"localField":   "sku",
		"foreignField": "_id",
		"as":           "sku_details",
		"pipeline":     bson.A{bson.M{"$project": bson.M{"name": 1}}},
	}}, bson.M{"$addFields": bson.M{"sku_details": bson.M{"$first": "$sku_details"}}})
	return items, err
}

func (i *itemRepo) Create(ctx context.Context, doc []domain.Item) error {
	err := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {

		_, err := i.itemCollection.InsertMany(sc, utils.ConvertArrStructToInterfaceArr(doc))
		if err != nil {
			return err
		}
		return session.CommitTransaction(sc)
	})
	if err != nil {
		i.logger.Error("Create -> transaction error in Create item -> ", err)
		return err
	}
	return nil
}

func (i *itemRepo) Update(ctx context.Context, id *primitive.ObjectID, doc *domain.Item) error {
	err := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		doc.ID = *id
		err := i.itemCollection.UpdateWithCtx(sc, doc)
		if err != nil {
			i.logger.Error("ItemRepository -> Update -> ", err)
			return err
		}
		menuGroupItem := menu_group_item.MenuGroupItemSyncItemModel{}
		copier.Copy(&menuGroupItem, &doc)
		menuGroupItem.UpdatedAt = time.Now()
		menuGroupItem.ItemId = *id
		err = i.menuGroupItemRepo.SyncMenuItemsChanges(sc, menuGroupItem)
		if err != nil {
			i.logger.Error("ItemRepository -> SyncMenuItemsChanges -> ", err)
			return err
		}
		return session.CommitTransaction(sc)
	})

	if err != nil {
		i.logger.Error("ItemRepository -> transaction error in update item -> ", err)
		return err
	}
	return nil
}

func (i *itemRepo) SoftDelete(ctx context.Context, doc *domain.Item) error {
	err := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		err := i.itemCollection.UpdateWithCtx(ctx, doc)
		if err != nil {
			i.logger.Error("ItemRepository -> SoftDelete -> ", err)
			return err
		}
		menuGroupItem := menu_group_item.MenuGroupItemSyncItemModel{}
		copier.Copy(&menuGroupItem, &doc)
		menuGroupItem.UpdatedAt = time.Now()
		menuGroupItem.ItemId = doc.ID
		err = i.menuGroupItemRepo.SyncMenuItemsChanges(sc, menuGroupItem)
		if err != nil {
			i.logger.Error("ItemRepository -> SyncMenuItemsChanges -> ", err)
			return err
		}
		return session.CommitTransaction(sc)
	})

	if err != nil {
		i.logger.Error("ItemRepository -> transaction error in update item -> ", err)
		return err
	}
	return nil
}
func (i *itemRepo) ChangeStatus(ctx context.Context, doc *domain.Item) error {
	err := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		err := i.itemCollection.UpdateWithCtx(ctx, doc)
		if err != nil {
			i.logger.Error("ItemRepository -> ChangeStatus -> ", err)
			return err
		}
		menuGroupItem := menu_group_item.MenuGroupItemSyncItemModel{}
		copier.Copy(&menuGroupItem, &doc)
		menuGroupItem.UpdatedAt = time.Now()
		menuGroupItem.ItemId = doc.ID
		err = i.menuGroupItemRepo.SyncMenuItemsChanges(sc, menuGroupItem)
		if err != nil {
			i.logger.Error("ItemRepository -> SyncMenuItemsChanges -> ", err)
			return err
		}
		return session.CommitTransaction(sc)
	})

	if err != nil {
		i.logger.Error("ItemRepository -> transaction error in update item -> ", err)
		return err
	}
	return nil
}

func (i *itemRepo) List(ctx context.Context, query *item.ListItemsDto) (items []domain.Item, paginationMeta *PaginationData, err error) {
	filter := bson.M{"account_id": utils.ConvertStringIdToObjectId(query.AccountId), "deleted_at": nil}
	if query.Query != "" {
		filter["$text"] = bson.M{
			"$search": query.Query}
	}
	if query.Type != "" {
		filter["type"] = query.Type
	}
	data, err := New(i.itemCollection.Collection).Context(ctx).Limit(query.Limit).Page(query.Page).Sort("_id", -1).Aggregate(bson.M{"$match": filter})

	items = make([]domain.Item, 0)
	if data == nil || data.Data == nil {
		return items, &PaginationData{}, err
	}

	for _, raw := range data.Data {
		model := domain.Item{}
		err = bson.Unmarshal(raw, &model)
		if err != nil {
			i.logger.Error("brands Repo -> List -> ", err)
			break
		}
		items = append(items, model)
	}
	paginationMeta = &data.Pagination

	return items, paginationMeta, err
}

func (i *itemRepo) CheckExists(ctx context.Context, accountId, name string, _exceptProductIds ...string) (bool, error) {
	exceptProductIds := make([]string, 0)
	for _, id := range _exceptProductIds {
		exceptProductIds = append(exceptProductIds, id)
	}
	filter := bson.M{"$or": bson.A{bson.M{"name.ar": name}, bson.M{"name.en": name}}, "account_id": utils.ConvertStringIdToObjectId(accountId), "deleted_at": nil,
		"_id": bson.M{"$nin": utils.ConvertStringIdsToObjectIds(exceptProductIds)}}
	c, err := i.itemCollection.CountDocuments(ctx, filter)

	return c > 0, err
}
