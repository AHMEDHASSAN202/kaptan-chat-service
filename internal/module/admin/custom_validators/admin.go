package custom_validators

import (
	"context"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/admin/domain"
	"samm/internal/module/admin/dto/admin"
)

type AdminCustomValidator struct {
	adminUseCase domain.AdminUseCase
}

func InitNewCustomValidatorsAdmin(adminUseCase domain.AdminUseCase) AdminCustomValidator {
	return AdminCustomValidator{
		adminUseCase: adminUseCase,
	}
}

func (i *AdminCustomValidator) ValidateEmailIsUnique() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		val := fl.Field().Interface().(string)
		//read the value of the top struct to get accountId
		var adminId primitive.ObjectID

		adminDto, ok := fl.Top().Interface().(*admin.CreateAdminDTO)
		if !ok {
			return false
		}
		adminId = adminDto.ID
		isExists, err := i.adminUseCase.CheckEmailExists(context.Background(), val, adminId)
		if err.IsError {
			return false
		}
		return !isExists
	}
}
