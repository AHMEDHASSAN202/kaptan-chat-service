package domain

import (
	"context"
	"samm/internal/module/config/dto/app_config"
	utilsDto "samm/pkg/utils/dto"
	"samm/pkg/validators"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppConfig struct {
	mgm.DefaultModel `bson:",inline"`
	ID               primitive.ObjectID      `json:"id" bson:"_id,omitempty"`
	ForceUpdate      bool                    `json:"force_update" bson:"force_update"`
	Type             string                  `json:"type" bson:"type"`
	StartupImage     string                  `json:"stratup_image" bson:"stratup_image"`
	AdminDetails     []utilsDto.AdminDetails `json:"admin_details" bson:"admin_details"`
	DeletedAt        *time.Time              `json:"deleted_at" bson:"deleted_at"`
}

type AppConfigUseCase interface {
	Create(ctx context.Context, dto app_config.CreateUpdateAppConfigDto) validators.ErrorResponse
	Update(ctx context.Context, dto app_config.CreateUpdateAppConfigDto) validators.ErrorResponse
	FindById(ctx context.Context, id string) (*AppConfig, validators.ErrorResponse)
	FindByType(ctx context.Context, configType string) (*AppConfig, validators.ErrorResponse)
	List(ctx context.Context, dto app_config.ListAppConfigDto) ([]AppConfig, validators.ErrorResponse)
	SoftDelete(ctx context.Context, id string) validators.ErrorResponse
	CheckExists(ctx context.Context, appType string, exceptIds ...string) (bool, validators.ErrorResponse)
}

type AppConfigRepository interface {
	Create(ctx context.Context, doc *AppConfig) error
	Update(ctx context.Context, id primitive.ObjectID, doc *AppConfig) error
	FindById(ctx context.Context, id primitive.ObjectID) (*AppConfig, error)
	FindByType(ctx context.Context, configType string) (*AppConfig, error)
	List(ctx context.Context, dto app_config.ListAppConfigDto) ([]AppConfig, error)
	SoftDelete(ctx context.Context, id primitive.ObjectID, adminDetails utilsDto.AdminDetails) error
	CheckExists(ctx context.Context, configType string, exceptIds ...string) (bool, error)
}
