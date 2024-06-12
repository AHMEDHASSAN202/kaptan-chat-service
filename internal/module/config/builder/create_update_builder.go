package builder

import (
	"samm/internal/module/config/domain"
	"samm/internal/module/config/dto/app_config"
	utilsDto "samm/pkg/utils/dto"
	"time"

	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertDtoToCorrespondingDomain(dto app_config.CreateUpdateAppConfigDto, oldDoc *domain.AppConfig) domain.AppConfig {
	var appConfig domain.AppConfig
	copier.Copy(&appConfig, &dto)
	appConfig.UpdatedAt = time.Now()
	if oldDoc == nil {
		appConfig.CreatedAt = time.Now()
		appConfig.AdminDetails = make([]utilsDto.AdminDetails, 0)
		appConfig.AdminDetails = append(appConfig.AdminDetails, utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: "Malhat", Operation: "Create", UpdatedAt: time.Now()})
	} else {
		appConfig.CreatedAt = oldDoc.CreatedAt
		appConfig.AdminDetails = oldDoc.AdminDetails
		appConfig.AdminDetails = append(appConfig.AdminDetails, utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: "Malhat", Operation: "Update", UpdatedAt: time.Now()})
	}
	return appConfig
}
