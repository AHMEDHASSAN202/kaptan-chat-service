package custom_validators

import (
	"context"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/admin/domain"
	"samm/internal/module/admin/dto/admin"
	"samm/internal/module/admin/dto/auth"
	"samm/pkg/utils"
	"strings"
)

type AdminCustomValidator struct {
	adminUseCase domain.AdminUseCase
	roleUseCase  domain.RoleUseCase
}

func InitNewCustomValidatorsAdmin(adminUseCase domain.AdminUseCase, roleUseCase domain.RoleUseCase) AdminCustomValidator {
	return AdminCustomValidator{
		adminUseCase: adminUseCase,
		roleUseCase:  roleUseCase,
	}
}

func (i *AdminCustomValidator) ValidateEmailIsUnique() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		val := fl.Field().Interface().(string)
		//read the value of the top struct to get accountId
		var adminId primitive.ObjectID

		switch fl.Top().Interface().(type) {
		case *admin.CreateAdminDTO:
			adminDto, ok := fl.Top().Interface().(*admin.CreateAdminDTO)
			if !ok {
				return false
			}
			adminId = adminDto.ID
		case *auth.UpdateAdminProfileDTO:
			adminDto, ok := fl.Top().Interface().(*auth.UpdateAdminProfileDTO)
			if !ok {
				return false
			}
			adminId = adminDto.ID
		case *auth.UpdatePortalProfileDTO:
			adminDto, ok := fl.Top().Interface().(*auth.UpdatePortalProfileDTO)
			if !ok {
				return false
			}
			adminId = adminDto.ID
			//default:
			//return false
		}
		isExists, err := i.adminUseCase.CheckEmailExists(context.Background(), strings.ToLower(val), adminId)
		if err.IsError {
			return false
		}
		return !isExists
	}
}

func (i *AdminCustomValidator) PasswordRequiredIfIdIsZero() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		var adminId primitive.ObjectID
		switch fl.Top().Interface().(type) {
		case *admin.CreateAdminDTO:
			adminDto, ok := fl.Top().Interface().(*admin.CreateAdminDTO)
			if !ok {
				return false
			}
			adminId = adminDto.ID
		case *auth.UpdateAdminProfileDTO:
			adminDto, ok := fl.Top().Interface().(*auth.UpdateAdminProfileDTO)
			if !ok {
				return false
			}
			adminId = adminDto.ID
		case *auth.UpdatePortalProfileDTO:
			adminDto, ok := fl.Top().Interface().(*auth.UpdatePortalProfileDTO)
			if !ok {
				return false
			}
			adminId = adminDto.ID
		default:
			return false
		}
		if adminId.IsZero() && fl.Field().String() == "" {
			return false
		}
		return true
	}
}

func (i *AdminCustomValidator) ValidateRoleExists() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		val := fl.Field().Interface().(string)
		role, err := i.roleUseCase.Find(context.Background(), utils.ConvertStringIdToObjectId(val))
		if err.IsError {
			return false
		}
		if role == nil {
			return false
		}
		return true
	}
}
