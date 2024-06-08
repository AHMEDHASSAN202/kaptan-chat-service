package modifier_group

import (
	"context"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/modifier_group"
	"samm/pkg/database/mongodb"
	"samm/pkg/utils/dto"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type modifierGroupRepo struct {
	modifierGroupCollection *mgm.Collection
}

func NewModifierGroupRepository(dbs *mongo.Database) domain.ModifierGroupRepository {
	modifierGroupCollection := mgm.Coll(&domain.ModifierGroup{})
	//text search modifier group
	mongodb.CreateIndex(modifierGroupCollection.Collection, false, bson.E{"name.ar", mongodb.IndexType.Text}, bson.E{"name.en", mongodb.IndexType.Text})
	//make sure there are no duplicated modifier group
	mongodb.CreateIndex(modifierGroupCollection.Collection, true, bson.E{"name.ar", mongodb.IndexType.Asc}, bson.E{"name.en", mongodb.IndexType.Asc}, bson.E{"account_id", mongodb.IndexType.Asc}, bson.E{"deleted_at", mongodb.IndexType.Asc})
	return &modifierGroupRepo{
		modifierGroupCollection: modifierGroupCollection,
	}
}

func (i *modifierGroupRepo) Create(ctx context.Context, doc domain.ModifierGroup) error {
	_, err := i.modifierGroupCollection.InsertOne(ctx, doc)
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

func (i *modifierGroupRepo) GetByIds(ctx context.Context, ids []primitive.ObjectID) ([]domain.ModifierGroup, error) {
	modifierGroups := []domain.ModifierGroup{}
	err := i.modifierGroupCollection.SimpleFind(&modifierGroups, bson.M{"_id": bson.M{"$in": ids}})
	return modifierGroups, err
}

func (i *modifierGroupRepo) List(ctx context.Context, query *modifier_group.ListModifierGroupsDto) ([]domain.ModifierGroup, error) {
	filter := bson.M{"deleted_at": nil}
	if query.Query != "" {
		filter["$text"] = bson.M{"$search": query.Query}
	}
	modifierGroups := []domain.ModifierGroup{}
	err := i.modifierGroupCollection.SimpleFindWithCtx(ctx, &modifierGroups, filter)
	return modifierGroups, err
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
