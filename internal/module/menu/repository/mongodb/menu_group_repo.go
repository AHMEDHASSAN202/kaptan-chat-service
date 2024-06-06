package mongodb

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/menu/domain"
)

type menuGroupRepo struct {
	menuGroupCollection *mgm.Collection
}

func NewMenuGroupRepository(dbs *mongo.Database) domain.MenuGroupRepository {
	menuGroupCollection := mgm.Coll(&domain.MenuGroup{})
	return &menuGroupRepo{
		menuGroupCollection: menuGroupCollection,
	}
}

func (r *menuGroupRepo) Create(ctx context.Context, domainData *domain.MenuGroup, menuGroupItems *[]domain.MenuGroupItem) (*domain.MenuGroup, error) {
	err := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		err := mgm.Coll(domainData).CreateWithCtx(sc, domainData)
		if err != nil {
			return err
		}
		//add to menu group item
		return session.CommitTransaction(sc)
	})
	return domainData, err
}
