package mongodb

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/retails/domain"
)

type locationRepository struct {
	locationCollection *mgm.Collection
}

const mongoLocationRepositoryTag = "LocationMongoRepository"

func NewLocationMongoRepository(dbs *mongo.Database) domain.LocationRepository {
	locationDbCollection := mgm.Coll(&domain.Location{})

	return &locationRepository{
		locationCollection: locationDbCollection,
	}
}
