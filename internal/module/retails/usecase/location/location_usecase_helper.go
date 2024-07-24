package location

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/uber/h3-go/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/retails/consts"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/location"
	"samm/internal/module/retails/responses"
	"samm/pkg/utils"
	utilsDto "samm/pkg/utils/dto"
	"strings"
	"time"
)

func LocationBuilder(ctx context.Context, payload *location.StoreLocationDto, l LocationUseCase) *domain.Location {
	var locationDomain domain.Location
	var brandDetails domain.BrandDetails

	copier.Copy(&locationDomain, payload)

	// Set Brand Details
	brandDomain, _ := l.brandUseCase.FindWithCuisines(&ctx, payload.BrandDetails.Id)
	copier.Copy(&brandDetails, brandDomain)
	locationDomain.BrandDetails = brandDetails
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
	locationDomain.WorkingHour = mapWorkingHours(locationDomain.WorkingHour)
	locationDomain.WorkingHourEid = mapWorkingHours(locationDomain.WorkingHourEid)
	locationDomain.WorkingHourRamadan = mapWorkingHours(locationDomain.WorkingHourRamadan)
	locationDomain.AllowedCollectionMethodIds = payload.AllowedCollectionMethodIds
	locationDomain.AdminDetails = append(locationDomain.AdminDetails, utilsDto.AdminDetails{Id: utils.ConvertStringIdToObjectId(payload.CauserId), Name: payload.CauserName, Type: payload.CauserType, Operation: "Create Location", UpdatedAt: time.Now()})
	locationDomain.CreatedAt = time.Now().UTC()
	locationDomain.UpdatedAt = time.Now().UTC()
	return &locationDomain
}
func LocationBulkBuilder(ctx context.Context, payload location.LocationDto, dto location.StoreBulkLocationDto, l LocationUseCase) *domain.Location {
	var locationDomain domain.Location
	var brandDetails domain.BrandDetails

	copier.Copy(&locationDomain, payload)

	locationDomain.ID = primitive.NewObjectID()
	locationDomain.City.Id = utils.ConvertStringIdToObjectId(payload.City.Id)
	locationDomain.BrandDetails.Id = utils.ConvertStringIdToObjectId(dto.BrandDetails.Id)
	locationDomain.AccountId = utils.ConvertStringIdToObjectId(dto.AccountId)
	locationDomain.Coordinate = domain.Coordinate{
		Type:        "Point",
		Coordinates: []float64{payload.Lng, payload.Lat},
	}
	locationDomain.Country.Id = dto.Country.Id
	locationDomain.Country.Currency = dto.Country.Currency
	locationDomain.Country.Name.Ar = dto.Country.Name.Ar
	locationDomain.Country.Name.En = dto.Country.Name.En
	locationDomain.Country.PhonePrefix = dto.Country.PhonePrefix
	locationDomain.Country.Timezone = dto.Country.Timezone

	// Set Brand Details
	brandDomain, _ := l.brandUseCase.FindWithCuisines(&ctx, dto.BrandDetails.Id)
	copier.Copy(&brandDetails, brandDomain)
	locationDomain.BrandDetails = brandDetails

	locationDomain.WorkingHour = mapWorkingHours(locationDomain.WorkingHour)
	locationDomain.WorkingHourEid = mapWorkingHours(locationDomain.WorkingHourEid)
	locationDomain.WorkingHourRamadan = mapWorkingHours(locationDomain.WorkingHourRamadan)
	locationDomain.Status = consts.LocationStatusInActive

	// Convert latitude and longitude to H3 index
	latLng := h3.NewLatLng(payload.Lat, payload.Lng)
	locationDomain.Index = h3.LatLngToCell(latLng, consts.H3Resolution).String()

	// Set Branch Signature

	locationDomain.BranchSignature = GenerateLocationBulkSignature(payload)
	locationDomain.AdminDetails = append(locationDomain.AdminDetails, utilsDto.AdminDetails{Id: utils.ConvertStringIdToObjectId(dto.CauserId), Name: dto.CauserName, Type: dto.CauserType, Operation: "Create Bulk Location", UpdatedAt: time.Now()})
	locationDomain.CreatedAt = time.Now().UTC()
	locationDomain.UpdatedAt = time.Now().UTC()
	return &locationDomain
}
func UpdateLocationBuilder(ctx context.Context, payload *location.StoreLocationDto, locationDomain *domain.Location, l LocationUseCase) *domain.Location {
	var brandDetails domain.BrandDetails

	copier.Copy(&locationDomain, payload)

	// Set Brand Details
	brandDomain, _ := l.brandUseCase.FindWithCuisines(&ctx, payload.BrandDetails.Id)
	copier.Copy(&brandDetails, brandDomain)
	locationDomain.BrandDetails = brandDetails

	locationDomain.City.Id = utils.ConvertStringIdToObjectId(payload.City.Id)
	locationDomain.BrandDetails.Id = utils.ConvertStringIdToObjectId(payload.BrandDetails.Id)
	locationDomain.AllowedCollectionMethodIds = payload.AllowedCollectionMethodIds
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
	locationDomain.WorkingHour = mapWorkingHours(locationDomain.WorkingHour)
	locationDomain.WorkingHourEid = mapWorkingHours(locationDomain.WorkingHourEid)
	locationDomain.WorkingHourRamadan = mapWorkingHours(locationDomain.WorkingHourRamadan)
	locationDomain.AdminDetails = append(locationDomain.AdminDetails, utilsDto.AdminDetails{Id: utils.ConvertStringIdToObjectId(payload.CauserId), Name: payload.CauserName, Type: payload.CauserType, Operation: "Update Location", UpdatedAt: time.Now()})

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
func GenerateLocationBulkSignature(payload location.LocationDto) string {

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
	domainData.AdminDetails = append(domainData.AdminDetails, utilsDto.AdminDetails{Id: utils.ConvertStringIdToObjectId(dto.CauserId), Name: dto.CauserName, Type: dto.CauserType, Operation: "Location Toggle Snooze", UpdatedAt: time.Now()})
	return domainData
}
func mapWorkingHours(workingHours []domain.WorkingHour) []domain.WorkingHour {
	items := make([]domain.WorkingHour, 0)
	for _, hour := range workingHours {
		hour.Day = strings.ToLower(hour.Day)
		items = append(items, hour)
	}
	return items
}

func populateCollectionMethods(ctx context.Context, l LocationUseCase, domainLocation *domain.Location) responses.LocationResp {
	location := responses.LocationResp{}
	copier.Copy(&location, &domainLocation)
	location.AllowedCollectionMethod = make([]map[string]any, 0)
	for _, CollectionMethodId := range domainLocation.AllowedCollectionMethodIds {
		CollectionMethod, _ := l.commonUseCase.FindCollectionMethodByType(ctx, CollectionMethodId)
		location.AllowedCollectionMethod = append(location.AllowedCollectionMethod, CollectionMethod)
	}
	return location
}
