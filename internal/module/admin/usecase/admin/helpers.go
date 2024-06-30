package admin

import (
	"context"
	"net/http"
	"samm/internal/module/admin/domain"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

func (oRec *AdminUseCase) LoginHelper(ctx context.Context, email, password, adminType string) (*domain.Admin, validators.ErrorResponse) {
	admin, errFindAdmin := oRec.repo.FindByEmail(ctx, email, adminType)
	if errFindAdmin != nil {
		oRec.logger.Error("AdminUseCase -> AdminLogin -> ", errFindAdmin)
		return admin, validators.GetErrorResponse(&ctx, localization.ErrLoginEmail, nil, utils.GetAsPointer(http.StatusBadRequest)) //change message
	}
	if admin == nil {
		oRec.logger.Error("AdminUseCase -> AdminLogin -> Admin Is NULL")
		return admin, validators.GetErrorResponse(&ctx, localization.ErrLoginEmail, nil, utils.GetAsPointer(http.StatusBadRequest)) //change message
	}
	if !utils.CheckPasswordHash(password, admin.Password) {
		oRec.logger.Error("AdminUseCase -> AdminLogin -> Password incorrect")
		return admin, validators.GetErrorResponse(&ctx, localization.ErrLoginPassword, nil, utils.GetAsPointer(http.StatusBadRequest)) //change message
	}
	if !admin.IsActive() {
		oRec.logger.Error("AdminUseCase -> AdminLogin -> Admin Is Not Active")
		return admin, validators.GetErrorResponse(&ctx, localization.ErrLoginInActive, nil, utils.GetAsPointer(http.StatusBadRequest)) //change message
	}
	return admin, validators.ErrorResponse{}
}
