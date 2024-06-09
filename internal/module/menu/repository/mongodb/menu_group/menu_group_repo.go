package menu_group

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/menu_group"
	"samm/pkg/logger"
	"samm/pkg/utils"
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

func (r *menuGroupRepo) ListPortal(ctx context.Context, dto menu_group.ListMenuGroupDTO) ([]domain.MenuGroup, *PaginationData, error) {
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
