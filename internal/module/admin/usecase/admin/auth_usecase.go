package admin

import (
	"context"
	"fmt"
	"net/http"
	admin2 "samm/internal/module/admin/builder/admin"
	"samm/internal/module/admin/consts"
	dto "samm/internal/module/admin/dto/auth"
	"samm/internal/module/admin/responses/admin"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

func (oRec *AdminUseCase) AdminLogin(ctx context.Context, input *dto.AdminAuthDTO) (interface{}, string, validators.ErrorResponse) {
	//find admin
	admin, err := oRec.LoginHelper(ctx, input.Email, input.Password, consts.ADMIN_TYPE)
	if err.IsError {
		return admin, "", err
	}

	//generate token
	token, errToken := oRec.AdminJwtService.GenerateToken(ctx, utils.ConvertObjectIdToStringId(admin.ID))
	if err.IsError {
		return admin, "", validators.GetErrorResponseFromErr(errToken)
	}

	//update token
	admin.Tokens = append(admin.Tokens, token)
	_, errUpdate := oRec.repo.Update(ctx, admin)
	if errUpdate != nil {
		return admin, "", validators.GetErrorResponse(&ctx, localization.E1000, nil, utils.GetAsPointer(http.StatusBadRequest))
	}

	return admin2.AdminProfileBuilder(admin), token, validators.ErrorResponse{}
}

func (oRec *AdminUseCase) PortalLogin(ctx context.Context, input *dto.PortalAuthDTO) (interface{}, string, validators.ErrorResponse) {
	//find admin
	admin, err := oRec.LoginHelper(ctx, input.Email, input.Password, consts.PORTAL_TYPE)
	if err.IsError {
		return admin, "", err
	}

	//generate token
	token, errToken := oRec.PortalJwtService.GenerateToken(ctx, utils.ConvertObjectIdToStringId(admin.ID))
	if err.IsError {
		return admin, "", validators.GetErrorResponseFromErr(errToken)
	}

	//update token
	admin.Tokens = append(admin.Tokens, token)
	_, errUpdate := oRec.repo.Update(ctx, admin)
	if errUpdate != nil {
		return admin, "", validators.GetErrorResponse(&ctx, localization.E1000, nil, utils.GetAsPointer(http.StatusBadRequest))
	}

	return admin2.AdminProfileBuilder(admin), token, validators.ErrorResponse{}
}

func (oRec *AdminUseCase) Profile(ctx context.Context, adminId string) (*admin.AdminProfileResponse, validators.ErrorResponse) {
	fmt.Println("CauserId => ", adminId)
	admin, errFindAdmin := oRec.repo.Find(ctx, utils.ConvertStringIdToObjectId(adminId))
	if errFindAdmin != nil {
		oRec.logger.Error("AdminUseCase -> Auth -> AdminLogin -> ", errFindAdmin)
		return nil, validators.GetErrorResponse(&ctx, localization.ErrLoginEmail, nil, utils.GetAsPointer(http.StatusBadRequest)) //change message
	}
	if admin == nil {
		oRec.logger.Error("AdminUseCase -> Auth -> AdminLogin -> Admin Is NULL")
		return nil, validators.GetErrorResponse(&ctx, localization.ErrLoginEmail, nil, utils.GetAsPointer(http.StatusBadRequest)) //change message
	}
	return admin2.AdminProfileBuilder(admin), validators.ErrorResponse{}
}
