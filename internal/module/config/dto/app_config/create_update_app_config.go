package app_config

import (
	"context"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"

	"github.com/go-playground/validator/v10"
)

type CreateUpdateAppConfigDto struct {
	Id                  string             `json:"_"`
	Type                string             `json:"type" validate:"required,oneof=user merchant,App_type_is_unique_rules_validation"`
	MinIOSVersion       int64              `json:"min_ios_version"`
	AppStoreLink        string             `json:"app_store_link"`
	MinAndroidVersion   int64              `json:"min_android_version"`
	PlayStoreLink       string             `json:"play_store_link"`
	MinHuaweiVersion    int64              `json:"min_huawei_version"`
	AppGalleryLink      string             `json:"app_gallery_link"`
	LocalizationVersion int64              `json:"localization_version"`
	StartupImage        string             `json:"stratup_image"`
	AdminDetails        []dto.AdminDetails `json:"-"`
}

func (input *CreateUpdateAppConfigDto) Validate(ctx context.Context, validate *validator.Validate, validateAppTypeUnique func(fl validator.FieldLevel) bool) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input, validators.CustomErrorTags{
		ValidationTag:          localization.App_type_is_unique_rules_validation,
		RegisterValidationFunc: validateAppTypeUnique,
	})
}
