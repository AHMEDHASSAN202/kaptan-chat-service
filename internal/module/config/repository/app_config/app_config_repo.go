package app_config

import (
	"context"
	"samm/internal/module/config/domain"
	"samm/internal/module/config/dto/app_config"
	"samm/pkg/logger"
	"samm/pkg/utils/dto"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type appConfigRepo struct {
	configCollection *mgm.Collection
	logger           logger.ILogger
}

func NewAppConfigRepository(dbs *mongo.Database, log logger.ILogger) domain.AppConfigRepository {
	configCollection := mgm.Coll(&domain.AppConfig{})
	return &appConfigRepo{
		configCollection: configCollection,
		logger:           log,
	}
}

func (i *appConfigRepo) Create(ctx context.Context, doc *domain.AppConfig) error {
	err := i.configCollection.Create(doc)
	if err != nil {
		return err
	}
	return nil
}

func (i *appConfigRepo) Update(ctx context.Context, id primitive.ObjectID, doc *domain.AppConfig) error {
	update := bson.M{"$set": doc}
	_, err := i.configCollection.UpdateByID(ctx, &id, update)
	if err != nil {
		return err
	}
	return nil
}

func (l appConfigRepo) FindById(ctx context.Context, Id primitive.ObjectID) (*domain.AppConfig, error) {
	var domainData domain.AppConfig
	filter := bson.M{"deleted_at": nil, "_id": Id}
	err := l.configCollection.FirstWithCtx(ctx, filter, &domainData)
	return &domainData, err
}

func (l appConfigRepo) FindByType(ctx context.Context, configType string) (*domain.AppConfig, error) {
	var domainData domain.AppConfig
	filter := bson.M{"deleted_at": nil, "type": configType}
	err := l.configCollection.FirstWithCtx(ctx, filter, &domainData)
	return &domainData, err
}

func (i *appConfigRepo) List(ctx context.Context, dto app_config.ListAppConfigDto) (configs []domain.AppConfig, err error) {
	configs = make([]domain.AppConfig, 0)
	filter := bson.M{"deleted_at": nil}

	if dto.Type != "" {
		filter["type"] = dto.Type
	}

	err = i.configCollection.SimpleFindWithCtx(ctx, &configs, filter)

	return
}

func (i *appConfigRepo) SoftDelete(ctx context.Context, id primitive.ObjectID, adminDetails dto.AdminDetails) error {
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}, "$push": bson.M{"admin_details": adminDetails}}
	_, err := i.configCollection.UpdateByID(ctx, &id, update)
	if err != nil {
		return err
	}
	return nil
}
