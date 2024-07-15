package admin

import (
	"samm/internal/module/admin/domain"
	"samm/internal/module/admin/dto/auth"
	"samm/pkg/utils"
	"strings"
)

func UpdateAdminProfileBuilder(admin *domain.Admin, input *auth.UpdateAdminProfileDTO) (*domain.Admin, error) {
	if input.Password != "" {
		//hash password
		password, err := utils.HashPassword(input.Password)
		if err != nil {
			return admin, err
		}
		admin.Password = password
	}
	admin.Name = input.Name
	admin.Email = strings.ToLower(input.Email)
	admin.AdminDetails = append(admin.AdminDetails, input.AdminDetails)
	return admin, nil
}

func UpdatePortalProfileBuilder(admin *domain.Admin, input *auth.UpdatePortalProfileDTO) (*domain.Admin, error) {
	if input.Password != "" {
		password, err := utils.HashPassword(input.Password)
		if err != nil {
			return admin, err
		}
		admin.Password = password
	}
	admin.Name = input.Name
	admin.Email = strings.ToLower(input.Email)
	admin.AdminDetails = append(admin.AdminDetails, input.AdminDetails)
	return admin, nil
}

func UpdateKitchenProfileBuilder(admin *domain.Admin, input *auth.UpdateKitchenProfileDTO) (*domain.Admin, error) {
	if input.Password != "" {
		password, err := utils.HashPassword(input.Password)
		if err != nil {
			return admin, err
		}
		admin.Password = password
	}
	admin.Name = input.Name
	admin.Email = strings.ToLower(input.Email)
	admin.AdminDetails = append(admin.AdminDetails, input.AdminDetails)
	return admin, nil
}
