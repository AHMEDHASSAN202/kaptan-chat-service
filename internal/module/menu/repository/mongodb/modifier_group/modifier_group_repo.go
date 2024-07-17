package modifier_group

import (
	"context"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/modifier_group"
	modifier_group_resp "samm/internal/module/menu/responses/modifier_group"
	"samm/pkg/database/mongodb"
	"samm/pkg/utils"
	"samm/pkg/utils/dto"
	"time"

	. "github.com/gobeam/mongo-go-pagination"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type modifierGroupRepo struct {
	modifierGroupCollection *mgm.Collection
	itemCollection          *mgm.Collection
}

func NewModifierGroupRepository(dbs *mongo.Database) domain.ModifierGroupRepository {
	modifierGroupCollection := mgm.Coll(&domain.ModifierGroup{})
	itemCollection := mgm.Coll(&domain.Item{})
	//text search modifier group
	mongodb.CreateIndex(modifierGroupCollection.Collection, false, bson.E{"name.ar", mongodb.IndexType.Text}, bson.E{"name.en", mongodb.IndexType.Text})
	mongodb.CreateIndex(modifierGroupCollection.Collection, false, bson.E{"status", mongodb.IndexType.Asc}, bson.E{"deleted_at", mongodb.IndexType.Asc})
	//make sure there are no duplicated modifier group
	mongodb.CreateIndex(modifierGroupCollection.Collection, true, bson.E{"name.ar", mongodb.IndexType.Asc}, bson.E{"name.en", mongodb.IndexType.Asc}, bson.E{"account_id", mongodb.IndexType.Asc}, bson.E{"deleted_at", mongodb.IndexType.Asc})
	return &modifierGroupRepo{
		modifierGroupCollection: modifierGroupCollection,
		itemCollection:          itemCollection,
	}
}

func (i *modifierGroupRepo) Create(ctx context.Context, docs []domain.ModifierGroup) error {
	_, err := i.modifierGroupCollection.InsertMany(ctx, utils.ConvertArrStructToInterfaceArr(docs))
	if err != nil {
		return err
	}
	return nil
}

func (i *modifierGroupRepo) Update(ctx context.Context, id *primitive.ObjectID, doc *domain.ModifierGroup) error {
	update := bson.M{"$set": doc}
	_, err := i.modifierGroupCollection.UpdateByID(ctx, &id, update)
	if err != nil {
		return err
	}
	return nil
}

func (i *modifierGroupRepo) GetByIds(ctx context.Context, ids []primitive.ObjectID) ([]modifier_group_resp.ModifierGroupResp, error) {
	modifierGroups := []modifier_group_resp.ModifierGroupResp{}
	err := i.modifierGroupCollection.SimpleAggregate(&modifierGroups, bson.M{"$match": bson.M{"_id": bson.M{"$in": ids}, "deleted_at": nil}}, bson.M{"$lookup": bson.M{"foreignField": "_id", "as": "products", "from": i.itemCollection.Name(), "localField": "product_ids", "pipeline": bson.A{bson.M{"$project": bson.M{"name": 1, "min": 1, "max": 1, "desc": 1, "image": 1}}}}})
	//err := i.modifierGroupCollection.SimpleFind(&modifierGroups, bson.M{"_id": bson.M{"$in": ids}})
	return modifierGroups, err
}

func (i *modifierGroupRepo) List(ctx context.Context, dto *modifier_group.ListModifierGroupsDto) ([]domain.ModifierGroup, *PaginationData, error) {
	models := make([]domain.ModifierGroup, 0)
	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.D{{"deleted_at", nil}},
		bson.D{{"account_id", utils.ConvertStringIdToObjectId(dto.AccountId)}},
	}}}

	if dto.Query != "" {
		pattern := ".*" + dto.Query + ".*"
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"$or": []bson.M{{"name.ar": bson.M{"$regex": pattern, "$options": "i"}}, {"name.en": bson.M{"$regex": pattern, "$options": "i"}}}})
	}

	data, err := New(i.modifierGroupCollection.Collection).Context(ctx).Limit(dto.Limit).Page(dto.Page).Sort("created_at", -1).Aggregate(matching)

	if data == nil || data.Data == nil {
		return models, nil, err
	}

	for _, raw := range data.Data {
		model := domain.ModifierGroup{}
		errUnmarshal := bson.Unmarshal(raw, &model)
		if errUnmarshal != nil {
			break
		}
		models = append(models, model)
	}

	return models, &data.Pagination, err

}

func (i *modifierGroupRepo) ChangeStatus(ctx context.Context, id *primitive.ObjectID, dto *modifier_group.ChangeModifierGroupStatusDto, adminDetails dto.AdminDetails) error {
	update := bson.M{"$set": bson.M{"status": dto.Status}, "$push": bson.M{"admin_details": adminDetails}}
	_, err := i.modifierGroupCollection.UpdateByID(ctx, &id, update)
	if err != nil {
		return err
	}
	return nil
}

func (i *modifierGroupRepo) SoftDelete(ctx context.Context, id *primitive.ObjectID, adminDetails dto.AdminDetails) error {
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}, "$push": bson.M{"admin_details": adminDetails}}
	_, err := i.modifierGroupCollection.UpdateByID(ctx, &id, update)
	if err != nil {
		return err
	}
	return nil
}
