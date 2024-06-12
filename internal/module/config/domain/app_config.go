package domain

import (
	"context"
	"samm/internal/module/config/dto/app_config"
	responses "samm/internal/module/config/responses/app_config"
	utilsDto "samm/pkg/utils/dto"
	"samm/pkg/validators"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppConfig struct {
	mgm.DefaultModel    `bson:",inline"`
	ID                  primitive.ObjectID      `json:"id" bson:"_id,omitempty"`
	Type                string                  `json:"type" bson:"type"`
	MinIOSVersion       int64                   `json:"min_ios_version" bson:"min_ios_version"`
	AppStoreLink        string                  `json:"app_store_link" bson:"app_store_link"`
	MinAndroidVersion   int64                   `json:"min_android_version" bson:"min_android_version"`
	PlayStoreLink       string                  `json:"play_store_link" bson:"play_store_link"`
	MinHuaweiVersion    int64                   `json:"min_huawei_version" bson:"min_huawei_version"`
	AppGalleryLink      string                  `json:"app_gallery_link" bson:"app_gallery_link"`
	LocalizationVersion int64                   `json:"localization_version" bson:"localization_version"`
	StartupImage        string                  `json:"stratup_image" bson:"stratup_image"`
	AdminDetails        []utilsDto.AdminDetails `json:"admin_details" bson:"admin_details"`
	DeletedAt           *time.Time              `json:"deleted_at" bson:"deleted_at"`
}

type AppConfigUseCase interface {
	// Admin
	Create(ctx context.Context, dto app_config.CreateUpdateAppConfigDto) validators.ErrorResponse
	Update(ctx context.Context, dto app_config.CreateUpdateAppConfigDto) validators.ErrorResponse
	FindById(ctx context.Context, id string) (*AppConfig, validators.ErrorResponse)
	FindByType(ctx context.Context, configType string) (*AppConfig, validators.ErrorResponse)
	List(ctx context.Context, dto app_config.ListAppConfigDto) ([]AppConfig, validators.ErrorResponse)
	SoftDelete(ctx context.Context, id string) validators.ErrorResponse
	CheckExists(ctx context.Context, appType string, exceptIds ...string) (bool, validators.ErrorResponse)

	// Mobile
	FindMobileConfig(ctx context.Context, dto app_config.FindMobileConfigDto) (responses.FindMobileConfigResponse, validators.ErrorResponse)
}

type AppConfigRepository interface {
	// Admin
	Create(ctx context.Context, doc *AppConfig) error
	Update(ctx context.Context, id primitive.ObjectID, doc *AppConfig) error
	FindById(ctx context.Context, id primitive.ObjectID) (*AppConfig, error)
	FindByType(ctx context.Context, configType string) (*AppConfig, error)
	List(ctx context.Context, dto app_config.ListAppConfigDto) ([]AppConfig, error)
	SoftDelete(ctx context.Context, id primitive.ObjectID, adminDetails utilsDto.AdminDetails) error
	CheckExists(ctx context.Context, configType string, exceptIds ...string) (bool, error)
}
