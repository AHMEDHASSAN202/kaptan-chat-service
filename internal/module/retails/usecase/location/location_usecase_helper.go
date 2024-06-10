package location

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/uber/h3-go/v4"
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

	// Convert latitude and longitude to H3 index
	latLng := h3.NewLatLng(payload.Lat, payload.Lng)
	locationDomain.Index = h3.LatLngToCell(latLng, consts.H3Resolution).String()

	// Set Branch Signature

	locationDomain.BranchSignature = GenerateLocationSignature(payload)

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
	locationDomain.BranchSignature = GenerateLocationSignature(payload)

	// Convert latitude and longitude to H3 index
	latLng := h3.NewLatLng(payload.Lat, payload.Lng)
	locationDomain.Index = h3.LatLngToCell(latLng, consts.H3Resolution).String()

	locationDomain.UpdatedAt = time.Now().UTC()
	return locationDomain
}
func GenerateLocationSignature(payload *location.StoreLocationDto) string {

	jsonBytes, err := json.Marshal(payload.Name)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	nameString := string(jsonBytes)

	jsonBytes, err = json.Marshal(payload.Street)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	streetString := string(jsonBytes)

	return utils.Encrypt("", nameString+streetString+payload.Logo+payload.CoverImage)

}
func domainBuilderToggleSnooze(dto *location.LocationToggleSnoozeDto, domainData *domain.Location) *domain.Location {
	var snoozedTill time.Time
	if dto.IsSnooze && dto.SnoozeMinutesInterval > 0 {
		snoozedTill = time.Now().UTC().Add(time.Duration(dto.SnoozeMinutesInterval) * time.Minute)
	}
	domainData.UpdatedAt = time.Now()
	domainData.SnoozeTo = &snoozedTill
	return domainData
}
