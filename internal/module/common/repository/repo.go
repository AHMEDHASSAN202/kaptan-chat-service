package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/common/domain"
)

const mongoLocationRepositoryTag = "CommonMongoRepository"

type CommonRepository struct {
}

func NewCommonMongoRepository(dbs *mongo.Database) domain.CommonRepository {

	return &CommonRepository{}
}
