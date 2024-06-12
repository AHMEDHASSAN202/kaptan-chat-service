package app_config

import (
	"context"
	"samm/internal/module/config/domain"
	"samm/internal/module/config/dto/app_config"
	"samm/pkg/logger"
	"samm/pkg/utils"
	utilsDto "samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppConfigUseCase struct {
	repo   domain.AppConfigRepository
	logger logger.ILogger
}

func NewAppConfigUseCase(repo domain.AppConfigRepository, logger logger.ILogger) domain.AppConfigUseCase {
	return &AppConfigUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (oRec *AppConfigUseCase) Create(ctx context.Context, dto app_config.CreateUpdateAppConfigDto) validators.ErrorResponse {
	var doc domain.AppConfig
	copier.Copy(&doc, &dto)
	doc.AdminDetails = make([]utilsDto.AdminDetails, 0)
	doc.AdminDetails = append(doc.AdminDetails, utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: "Malhat", Operation: "Create", UpdatedAt: time.Now()})
	err := oRec.repo.Create(ctx, &doc)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *AppConfigUseCase) Update(ctx context.Context, dto app_config.CreateUpdateAppConfigDto) validators.ErrorResponse {
	var doc domain.AppConfig
	copier.Copy(&doc, &dto)
	id := utils.ConvertStringIdToObjectId(dto.Id)
	err := oRec.repo.Update(ctx, id, &doc)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *AppConfigUseCase) List(ctx context.Context, dto app_config.ListAppConfigDto) (configs []domain.AppConfig, err validators.ErrorResponse) {
	configs, resErr := oRec.repo.List(ctx, dto)
	if resErr != nil {
		err = validators.GetErrorResponseFromErr(resErr)
		return
	}
	return
}

func (oRec *AppConfigUseCase) FindById(ctx context.Context, id string) (*domain.AppConfig, validators.ErrorResponse) {
	config, err := oRec.repo.FindById(ctx, utils.ConvertStringIdToObjectId(id))
	if err != nil {
		return config, validators.GetErrorResponseFromErr(err)
	}
	if config == nil {
		return nil, validators.GetErrorResponseFromErr(errors.New(localization.E1002))
	}

	return config, validators.ErrorResponse{}
}

func (oRec *AppConfigUseCase) FindByType(ctx context.Context, configType string) (*domain.AppConfig, validators.ErrorResponse) {
	config, err := oRec.repo.FindByType(ctx, configType)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}
	if config == nil {
		return nil, validators.GetErrorResponseFromErr(errors.New(localization.E1002))
	}

	return config, validators.ErrorResponse{}
}

func (oRec *AppConfigUseCase) SoftDelete(ctx context.Context, id string) validators.ErrorResponse {
	idDoc := utils.ConvertStringIdToObjectId(id)
	adminDetails := utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: "Malhat", Operation: "Delete", UpdatedAt: time.Now()}
	err := oRec.repo.SoftDelete(ctx, idDoc, adminDetails)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *AppConfigUseCase) CheckExists(ctx context.Context, configType string, exceptIds ...string) (bool, validators.ErrorResponse) {
	isExists, err := oRec.repo.CheckExists(ctx, configType, exceptIds...)
	if err != nil {
		return isExists, validators.GetErrorResponseFromErr(err)
	}
	return isExists, validators.ErrorResponse{}
}
