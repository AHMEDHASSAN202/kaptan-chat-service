package mongodb

import (
	"context"
	"example.com/fxdemo/internal/module/menu/domain"
	"github.com/kamva/mgm/v3"
	"github.com/n-goo/ngo-menu-service/pkg/database/mongodb"
	"github.com/n-goo/ngo-menu-service/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	//db                 *mongo.Database
	locationCollection *mgm.Collection
	tokenCollection    *mgm.Collection
}

const mongoMenuRepositoryTag = "MenuMongoRepository"

func NewMenuMongoRepository(dbs *mongo.Database) domain.MenuRepository {
	locationDbCollection := mgm.Coll(&domain.Location{})
	tokenCollection := mgm.Coll(&domain.Token{})

	mongodb.CreateIndex(locationDbCollection.Collection, false, "channel_link_id")

	return &mongoRepository{
		locationCollection: locationDbCollection,
		tokenCollection:    tokenCollection,
	}
}

func (r *mongoRepository) UpdateTokens(ctx context.Context, domainData *domain.Token) (err error) {
	filter := bson.D{{"scope", domainData.Scope}}
	update := bson.D{{"$set", domainData}}
	upsert := true
	opts := options.UpdateOptions{Upsert: &upsert}
	_, err = r.tokenCollection.UpdateOne(ctx, filter, update, &opts)
	if err != nil {
		return err
	}

	return nil
}

func (r *mongoRepository) FindToken(ctx context.Context) (result domain.Token, err error) {
	filter := bson.M{} //bson.M{"scope": scope}

	if err = r.tokenCollection.First(filter, &result); err != nil {
		logger.Error(ctx, err)
		return
	}

	return
}
