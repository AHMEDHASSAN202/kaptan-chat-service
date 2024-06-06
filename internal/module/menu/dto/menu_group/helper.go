package menu_group

import (
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils"
)

func ValidateBranchIDs(fl validator.FieldLevel) bool {
	branchIDs := fl.Field().Interface().([]string)
	if len(branchIDs) == 0 {
		return true
	}
	for _, id := range branchIDs {
		if !utils.IsValidateObjectId(id) {
			return false
		}
	}
	return true
}
