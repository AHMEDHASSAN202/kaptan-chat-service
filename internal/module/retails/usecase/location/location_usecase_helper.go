package location

import (
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/retails/consts"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/location"
	"samm/pkg/utils"
	"time"
)

func LocationBuilder(payload *location.StoreLocationDto) *domain.Location {
	var locationDomain domain.Location

	copier.Copy(&locationDomain, payload)

	locationDomain.ID = primitive.NewObjectID()
	locationDomain.City.Id = utils.ConvertStringIdToObjectId(payload.City.Id)
	locationDomain.BrandDetails.Id = utils.ConvertStringIdToObjectId(payload.BrandDetails.Id)
	locationDomain.AccountId = utils.ConvertStringIdToObjectId(payload.AccountId)
	locationDomain.Coordinate = domain.Coordinate{
		Type:        "Point",
		Coordinates: []float64{payload.Lng, payload.Lat},
	}

	locationDomain.Status = consts.LocationStatusInActive
	// Set Branch Signature
	//locationDomain.BranchSignature = ""

	locationDomain.CreatedAt = time.Now().UTC()
	locationDomain.UpdatedAt = time.Now().UTC()
	return &locationDomain
}
func UpdateLocationBuilder(payload *location.StoreLocationDto, locationDomain *domain.Location) *domain.Location {

	copier.Copy(&locationDomain, payload)
	locationDomain.City.Id = utils.ConvertStringIdToObjectId(payload.City.Id)
	locationDomain.BrandDetails.Id = utils.ConvertStringIdToObjectId(payload.BrandDetails.Id)
	locationDomain.AccountId = utils.ConvertStringIdToObjectId(payload.AccountId)
	locationDomain.Coordinate = domain.Coordinate{
		Type:        "Point",
		Coordinates: []float64{payload.Lng, payload.Lat},
	}

	// Set Branch Signature
	//locationDomain.BranchSignature = ""

	locationDomain.UpdatedAt = time.Now().UTC()
	return locationDomain
}
