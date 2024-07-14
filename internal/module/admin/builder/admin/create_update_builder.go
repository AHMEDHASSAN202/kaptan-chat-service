package admin

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/admin/domain"
	dto "samm/internal/module/admin/dto/admin"
	"samm/pkg/utils"
	dto2 "samm/pkg/utils/dto"
	"strings"
)

func CreateUpdateAdminBuilder(admin *domain.Admin, input *dto.CreateAdminDTO, role domain.Role) (*domain.Admin, error) {
	if admin == nil {
		admin = &domain.Admin{}
		admin.ID = primitive.NewObjectID()
		admin.Tokens = make([]string, 0)
		admin.AdminDetails = []dto2.AdminDetails{}
		if input.Account != nil {
			admin.Account = &domain.Account{Id: utils.ConvertStringIdToObjectId(input.Account.Id), Name: domain.Name{Ar: input.Account.Name.Ar, En: input.Account.Name.En}}
		}
		if input.Kitchen != nil {
			admin.Kitchen = &domain.Kitchen{Id: utils.ConvertStringIdToObjectId(input.Kitchen.Id), Name: domain.Name{Ar: input.Kitchen.Name.Ar, En: input.Kitchen.Name.En}, AllowedStatus: input.Kitchen.AllowedStatus}
		}
	}
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
	admin.Status = strings.ToLower(input.Status)
	admin.Type = strings.ToLower(input.Type)
	admin.CountryIds = utils.ArrayToUpper(input.CountryIds)
	admin.Role = role
	admin.AdminDetails = append(admin.AdminDetails, input.AdminDetails)
	return admin, nil
}
