package mongodb

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/menu/domain"
)

type menuGroupRepo struct {
	menuGroupCollection *mgm.Collection
	menuGroupItemRepo   domain.MenuGroupItemRepository
}

func NewMenuGroupRepository(dbs *mongo.Database, menuGroupItemRepo domain.MenuGroupItemRepository) domain.MenuGroupRepository {
	menuGroupCollection := mgm.Coll(&domain.MenuGroup{})
	return &menuGroupRepo{
		menuGroupCollection: menuGroupCollection,
		menuGroupItemRepo:   menuGroupItemRepo,
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
