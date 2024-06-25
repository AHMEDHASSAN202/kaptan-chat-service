package admin

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/admin/domain"
	dto "samm/internal/module/admin/dto/admin"
	"samm/pkg/utils"
	dto2 "samm/pkg/utils/dto"
	"strings"
)

func CreateUpdateAdminBuilder(admin *domain.Admin, input *dto.CreateAdminDTO) *domain.Admin {
	if admin == nil {
		admin = &domain.Admin{}
		admin.ID = primitive.NewObjectID()
		admin.Tokens = make([]string, 0)
		admin.AdminDetails = []dto2.AdminDetails{}
	}
	if input.Password != "" {
		//hash password
		admin.Password = input.Password
	}
	admin.Name = input.Name
	admin.Email = input.Email
	admin.Status = strings.ToLower(input.Status)
	admin.Role = strings.ToLower(input.Role)
	admin.Permissions = utils.ArrayToLower(input.Permissions)
	admin.Type = strings.ToLower(input.Type)
	admin.CountryIds = utils.ArrayToUpper(input.CountryIds)
	admin.MetaData = domain.MetaData{AccountId: input.AccountId}
	admin.AdminDetails = append(admin.AdminDetails, input.AdminDetails)
	return admin
}
