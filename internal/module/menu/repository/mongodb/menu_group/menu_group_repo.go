package menu_group

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/menu_group"
	menu_group2 "samm/internal/module/menu/repository/structs/menu_group"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/utils/dto"
	"samm/pkg/validators/localization"
	"strings"
)

type menuGroupRepo struct {
	menuGroupCollection *mgm.Collection
	menuGroupItemRepo   domain.MenuGroupItemRepository
	logger              logger.ILogger
}

func NewMenuGroupRepository(dbs *mongo.Database, menuGroupItemRepo domain.MenuGroupItemRepository, log logger.ILogger) domain.MenuGroupRepository {
	menuGroupCollection := mgm.Coll(&domain.MenuGroup{})
	return &menuGroupRepo{
		menuGroupCollection: menuGroupCollection,
		menuGroupItemRepo:   menuGroupItemRepo,
		logger:              log,
	}
}

func (r *menuGroupRepo) Create(ctx context.Context, domainData *domain.MenuGroup, menuGroupItems *[]domain.MenuGroupItem) (*domain.MenuGroup, error) {
	err := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		err := mgm.Coll(domainData).CreateWithCtx(sc, domainData)
		if err != nil {
			return err
		}
		if menuGroupItems != nil && len(*menuGroupItems) >= 1 {
			err = r.menuGroupItemRepo.CreateUpdateBulk(sc, menuGroupItems)
			if err != nil {
				return err
			}
		}
		return session.CommitTransaction(sc)
	})
	return domainData, err
}

func (r *menuGroupRepo) Update(ctx context.Context, domainData *domain.MenuGroup, menuGroupItems *[]domain.MenuGroupItem) (*domain.MenuGroup, error) {
	menuIds := []primitive.ObjectID{}
	if menuGroupItems != nil {
		for _, item := range *menuGroupItems {
			menuIds = append(menuIds, item.ID)
		}
	}
	err := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		err := mgm.Coll(domainData).UpdateWithCtx(sc, domainData)
		if err != nil {
			return err
		}
		if menuGroupItems != nil && len(*menuGroupItems) >= 1 {
			err = r.menuGroupItemRepo.CreateUpdateBulk(sc, menuGroupItems)
			if err != nil {
				return err
			}
		}
		err = r.menuGroupItemRepo.DeleteBulkByGroupMenuId(sc, domainData.ID, menuIds)
		if err != nil {
			return err
		}
		return session.CommitTransaction(sc)
	})
	return domainData, err
}

func (r *menuGroupRepo) Delete(ctx context.Context, domainData *domain.MenuGroup) error {
	err := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		err := mgm.Coll(&domain.MenuGroup{}).DeleteWithCtx(sc, domainData)
		if err != nil {
			return err
		}
		err = r.menuGroupItemRepo.DeleteBulkByGroupMenuId(sc, domainData.ID, nil)
		if err != nil {
			return err
		}
		return session.CommitTransaction(sc)
	})
	return err
}

func (r *menuGroupRepo) Find(ctx context.Context, menuGroupId primitive.ObjectID) (*domain.MenuGroup, error) {
	domainData := domain.MenuGroup{}
	err := mgm.Coll(&domain.MenuGroup{}).FindByID(menuGroupId, &domainData)
	return &domainData, err
}

func (r *menuGroupRepo) FindWithItems(ctx context.Context, menuGroupId primitive.ObjectID) (*menu_group2.FindMenuGroupWithItems, error) {
	domainData := menu_group2.FindMenuGroupWithItems{}
	itemsCollectionName := mgm.Coll(&domain.MenuGroupItem{}).Name()
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "_id", Value: menuGroupId}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$categories"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: itemsCollectionName},
			{Key: "localField", Value: "categories._id"},
			{Key: "foreignField", Value: "category._id"},
			{Key: "as", Value: "categories.menu_items"},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id"},
			{Key: "categories", Value: bson.D{{Key: "$push", Value: "$categories"}}},
			{Key: "name", Value: bson.D{{Key: "$first", Value: "$name"}}},
			{Key: "branch_ids", Value: bson.D{{Key: "$first", Value: "$branch_ids"}}},
			{Key: "availabilities", Value: bson.D{{Key: "$first", Value: "$availabilities"}}},
			{Key: "status", Value: bson.D{{Key: "$first", Value: "$status"}}},
			{Key: "account_id", Value: bson.D{{Key: "$first", Value: "$account_id"}}},
			{Key: "created_at", Value: bson.D{{Key: "$first", Value: "$created_at"}}},
			{Key: "updated_at", Value: bson.D{{Key: "$first", Value: "$updated_at"}}},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "categories", Value: bson.D{
				{Key: "$sortArray", Value: bson.D{
					{Key: "input", Value: "$categories"},
					{Key: "sortBy", Value: bson.D{{Key: "sort", Value: 1}}},
				}},
			}},
		}}},
	}
	cursor, err := mgm.Coll(&domain.MenuGroup{}).Aggregate(ctx, pipeline, options.Aggregate())
	if err != nil {
		r.logger.Error("menuGroupRepo -> FindWithItems -> ", err)
		return &domainData, errors.New(localization.GetTranslation(&ctx, localization.E1000, nil, ""))
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if err := cursor.Decode(&domainData); err != nil {
			r.logger.Error("menuGroupRepo -> FindWithItems -> ", err)
			return &domainData, err
		}
	} else {
		r.logger.Error("menuGroupRepo -> FindWithItems -> ", cursor.Err())
		return &domainData, errors.New(localization.GetTranslation(&ctx, localization.E1000, nil, ""))
	}

	return &domainData, err
}

func (r *menuGroupRepo) List(ctx context.Context, dto menu_group.ListMenuGroupDTO) ([]domain.MenuGroup, *PaginationData, error) {
	models := make([]domain.MenuGroup, 0)

	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.D{{"account_id", utils.ConvertStringIdToObjectId(dto.AccountId)}},
	}}}

	if dto.Query != "" {
		pattern := ".*" + dto.Query + ".*"
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"$or": []bson.M{{"name.ar": bson.M{"$regex": pattern, "$options": "i"}}, {"name.en": bson.M{"$regex": pattern, "$options": "i"}}}})
	}

	if dto.BranchId != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"branch_ids": utils.ConvertStringIdToObjectId(dto.BranchId)})
	}

	if dto.Status != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"status": strings.ToLower(dto.Status)})
	}

	data, err := New(r.menuGroupCollection.Collection).Context(ctx).Limit(dto.Limit).Page(dto.Page).Sort("created_at", -1).Aggregate(matching)

	if data == nil || data.Data == nil {
		return models, nil, err
	}

	for _, raw := range data.Data {
		model := domain.MenuGroup{}
		errUnmarshal := bson.Unmarshal(raw, &model)
		if errUnmarshal != nil {
			r.logger.Error("menuGroupRepo -> List -> ", errUnmarshal)
			break
		}
		models = append(models, model)
	}

	return models, &data.Pagination, err
}

func (r *menuGroupRepo) ChangeMenuStatus(ctx context.Context, model *domain.MenuGroup, input *menu_group.ChangeMenuGroupStatusDto, adminDetails dto.AdminDetails) error {
	model.Status = input.Status
	model.AdminDetails = utils.If(model.AdminDetails == nil, []dto.AdminDetails{adminDetails}, append(model.AdminDetails, adminDetails)).([]dto.AdminDetails)
	err := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		err := mgm.Coll(model).UpdateWithCtx(sc, model)
		if err != nil {
			r.logger.Error("menuGroupRepo -> ChangeMenuStatus -> ", err)
			return err
		}
		err = r.menuGroupItemRepo.ChangeMenuStatus(sc, model.ID, input, adminDetails)
		if err != nil {
			r.logger.Error("menuGroupRepo -> ChangeMenuStatus -> ", err)
			return err
		}
		return session.CommitTransaction(sc)
	})
	return err
}

func (r *menuGroupRepo) ChangeCategoryStatus(ctx context.Context, model *domain.MenuGroup, input *menu_group.ChangeMenuGroupStatusDto, adminDetails dto.AdminDetails) error {
	model.AdminDetails = utils.If(model.AdminDetails == nil, []dto.AdminDetails{adminDetails}, append(model.AdminDetails, adminDetails)).([]dto.AdminDetails)
	exists := false
	if model.Categories != nil {
		for i, category := range model.Categories {
			if utils.ConvertObjectIdToStringId(category.ID) == input.Id {
				model.Categories[i].Status = input.Status
				exists = true
			}
		}
	}
	if !exists {
		return errors.New(localization.GetTranslation(&ctx, localization.E1000, nil, ""))
	}
	err := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		err := mgm.Coll(model).UpdateWithCtx(sc, model)
		if err != nil {
			return err
		}
		err = r.menuGroupItemRepo.ChangeCategoryStatus(sc, model.ID, input, adminDetails)
		if err != nil {
			return err
		}
		return session.CommitTransaction(sc)
	})
	return err
}

func (r *menuGroupRepo) DeleteCategory(ctx context.Context, model *domain.MenuGroup, input *menu_group.DeleteEntityFromMenuGroupDto, adminDetails dto.AdminDetails) error {
	model.AdminDetails = utils.If(model.AdminDetails == nil, []dto.AdminDetails{adminDetails}, append(model.AdminDetails, adminDetails)).([]dto.AdminDetails)
	var index *int
	if model.Categories == nil {
		return errors.New(localization.GetTranslation(&ctx, localization.E1002, nil, ""))
	}
	for i, category := range model.Categories {
		if utils.ConvertObjectIdToStringId(category.ID) == input.EntityId {
			index = &i
			break
		}
	}
	if index == nil {
		return errors.New(localization.GetTranslation(&ctx, localization.E1002, nil, ""))
	}
	model.Categories = utils.RemoveItemByIndex[domain.Category](model.Categories, *index)
	err := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		err := mgm.Coll(model).UpdateWithCtx(sc, model)
		if err != nil {
			return err
		}
		err = r.menuGroupItemRepo.DeleteByCategory(sc, input)
		if err != nil {
			return err
		}
		return session.CommitTransaction(sc)
	})
	return err
}

func (r *menuGroupRepo) DeleteItem(ctx context.Context, model *domain.MenuGroup, input *menu_group.DeleteEntityFromMenuGroupDto, adminDetails dto.AdminDetails) error {
	model.AdminDetails = utils.If(model.AdminDetails == nil, []dto.AdminDetails{adminDetails}, append(model.AdminDetails, adminDetails)).([]dto.AdminDetails)
	exists := false
	entityId := utils.ConvertStringIdToObjectId(input.EntityId)
	if model.Categories != nil {
		for i, category := range model.Categories {
			if category.MenuItemIds != nil && utils.Contains(category.MenuItemIds, entityId) {
				model.Categories[i].MenuItemIds = utils.RemoveItemByValue[primitive.ObjectID](category.MenuItemIds, entityId)
				exists = true
				break
			}
		}
	}
	if !exists {
		return errors.New(localization.GetTranslation(&ctx, localization.E1002, nil, ""))
	}
	err := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		err := mgm.Coll(model).UpdateWithCtx(sc, model)
		if err != nil {
			return err
		}
		err = r.menuGroupItemRepo.Delete(sc, input)
		if err != nil {
			return err
		}
		return session.CommitTransaction(sc)
	})
	return err
}
