package admin

import (
	"context"
	"net/http"
	admin2 "samm/internal/module/admin/builder/admin"
	"samm/internal/module/admin/consts"
	"samm/internal/module/admin/domain"
	admin3 "samm/internal/module/admin/dto/admin"
	dto "samm/internal/module/admin/dto/auth"
	"samm/internal/module/admin/responses/admin"
	"samm/pkg/utils"
	utilsDto "samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"time"
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
func (oRec *AdminUseCase) KitchenLogin(ctx context.Context, input *dto.KitchenAuthDTO) (admin.AuthKitchenResponse, validators.ErrorResponse) {
	//find admin
	adminDoc, err := oRec.LoginHelper(ctx, input.Email, input.Password, consts.KITCHEN_TYPE)
	if err.IsError {
		return admin.AuthKitchenResponse{}, err
	}

	//generate token
	token, errToken := oRec.KitchenJwtService.GenerateToken(ctx, utils.ConvertObjectIdToStringId(adminDoc.ID))
	if err.IsError {
		return admin.AuthKitchenResponse{}, validators.GetErrorResponseFromErr(errToken)
	}

	//update token
	adminDoc.Tokens = append(adminDoc.Tokens, token)
	_, errUpdate := oRec.repo.Update(ctx, adminDoc)
	if errUpdate != nil {
		return admin.AuthKitchenResponse{}, validators.GetErrorResponse(&ctx, localization.E1000, nil, utils.GetAsPointer(http.StatusBadRequest))
	}

	//get firebase token
	firebaseToken, err2 := oRec.authClient.CustomTokenWithClaims(ctx, adminDoc.Kitchen.Id.Hex(), nil)
	if err2 != nil {
		return admin.AuthKitchenResponse{}, validators.GetErrorResponseFromErr(err2)
	}

	return admin.AuthKitchenResponse{
		Profile:       admin2.AdminProfileBuilder(adminDoc),
		Token:         token,
		FirebaseToken: firebaseToken,
	}, validators.ErrorResponse{}
}

func (oRec *AdminUseCase) Profile(ctx context.Context, profileDTO dto.ProfileDTO) (*admin.AdminProfileResponse, validators.ErrorResponse) {
	admin, errFindAdmin := oRec.repo.Find(ctx, utils.ConvertStringIdToObjectId(profileDTO.AdminId))
	if errFindAdmin != nil {
		oRec.logger.Error("AdminUseCase -> Auth -> AdminLogin -> ", errFindAdmin)
		return nil, validators.GetErrorResponse(&ctx, localization.ErrLoginEmail, nil, utils.GetAsPointer(http.StatusBadRequest)) //change message
	}
	if admin == nil {
		oRec.logger.Error("AdminUseCase -> Auth -> AdminLogin -> Admin Is NULL")
		return nil, validators.GetErrorResponse(&ctx, localization.ErrLoginEmail, nil, utils.GetAsPointer(http.StatusBadRequest)) //change message
	}
	if profileDTO.AccountId != "" {
		accountId := utils.SafeMapGet(profileDTO.CauserDetails, "id", "").(string)
		name, ok := utils.SafeMapGet(profileDTO.CauserDetails, "name", nil).(map[string]interface{})
		if ok {
			admin.Account = &domain.Account{
				Id: utils.ConvertStringIdToObjectId(accountId),
				Name: domain.Name{
					Ar: utils.SafeMapGet(name, "ar", "").(string),
					En: utils.SafeMapGet(name, "en", "").(string),
				},
			}
		}
	}
	return admin2.AdminProfileBuilder(admin), validators.ErrorResponse{}
}
func (oRec *AdminUseCase) KitchenProfile(ctx context.Context, profileDTO dto.KitchenProfileDTO) (*admin.AdminProfileResponse, string, validators.ErrorResponse) {
	admin, errFindAdmin := oRec.repo.Find(ctx, utils.ConvertStringIdToObjectId(profileDTO.AdminId))
	if errFindAdmin != nil {
		oRec.logger.Error("AdminUseCase -> Auth -> AdminLogin -> ", errFindAdmin)
		return nil, "", validators.GetErrorResponse(&ctx, localization.ErrLoginEmail, nil, utils.GetAsPointer(http.StatusBadRequest)) //change message
	}
	if admin == nil {
		oRec.logger.Error("AdminUseCase -> Auth -> AdminLogin -> Admin Is NULL")
		return nil, "", validators.GetErrorResponse(&ctx, localization.ErrLoginEmail, nil, utils.GetAsPointer(http.StatusBadRequest)) //change message
	}

	//get firebase token
	firebaseToken, err2 := oRec.authClient.CustomTokenWithClaims(ctx, admin.Kitchen.Id.Hex(), nil)
	if err2 != nil {
		return nil, "", validators.GetErrorResponseFromErr(err2)
	}

	//if profileDTO.AccountId != "" {
	//	accountId := utils.SafeMapGet(profileDTO.CauserDetails, "id", "").(string)
	//	name, ok := utils.SafeMapGet(profileDTO.CauserDetails, "name", nil).(map[string]interface{})
	//	if ok {
	//		admin.Account = &domain.Account{
	//			Id: utils.ConvertStringIdToObjectId(accountId),
	//			Name: domain.Name{
	//				Ar: utils.SafeMapGet(name, "ar", "").(string),
	//				En: utils.SafeMapGet(name, "en", "").(string),
	//			},
	//		}
	//	}
	//}
	return admin2.AdminProfileBuilder(admin), firebaseToken, validators.ErrorResponse{}
}

func (oRec *AdminUseCase) UpdateAdminProfile(ctx context.Context, input *dto.UpdateAdminProfileDTO) (*admin.AdminProfileResponse, validators.ErrorResponse) {
	admin, errFind := oRec.repo.Find(ctx, utils.ConvertStringIdToObjectId(input.CauserId))
	if errFind != nil {
		oRec.logger.Error("AdminUseCase -> UpdateAdminProfile -> ", errFind)
		return nil, validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
	}

	//input.AdminDetails = utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: input.Name, Operation: "Update My Profile", UpdatedAt: time.Now()}
	input.AdminDetails = utilsDto.AdminDetails{Id: utils.ConvertStringIdToObjectId(input.CauserId), Name: input.CauserName, Type: input.CauserType, Operation: "Update My Profile", UpdatedAt: time.Now()}
	adminDomain, err := admin2.UpdateAdminProfileBuilder(admin, input)
	if err != nil {
		oRec.logger.Error("AdminUseCase -> UpdateAdminProfile -> ", err)
	}

	admin, errCreate := oRec.repo.Update(ctx, adminDomain)
	if errCreate != nil {
		oRec.logger.Error("AdminUseCase -> UpdateAdminProfile -> ", errCreate)
		return nil, validators.GetErrorResponse(&ctx, localization.E1000, nil, nil)
	}

	return admin2.AdminProfileBuilder(admin), validators.ErrorResponse{}
}

func (oRec *AdminUseCase) UpdatePortalProfile(ctx context.Context, input *dto.UpdatePortalProfileDTO) (*admin.AdminProfileResponse, validators.ErrorResponse) {
	admin, errFind := oRec.repo.Find(ctx, utils.ConvertStringIdToObjectId(input.CauserId))
	if errFind != nil {
		oRec.logger.Error("AdminUseCase -> UpdatePortalProfile -> ", errFind)
		return nil, validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
	}

	//input.AdminDetails = utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: input.Name, Operation: "Update My Profile", UpdatedAt: time.Now()}
	input.AdminDetails = utilsDto.AdminDetails{Id: utils.ConvertStringIdToObjectId(input.CauserId), Name: input.CauserName, Type: input.CauserType, Operation: "Update My Profile", UpdatedAt: time.Now()}
	adminDomain, err := admin2.UpdatePortalProfileBuilder(admin, input)
	if err != nil {
		oRec.logger.Error("AdminUseCase -> UpdatePortalProfile -> ", err)
	}

	admin, errCreate := oRec.repo.Update(ctx, adminDomain)
	if errCreate != nil {
		oRec.logger.Error("AdminUseCase -> UpdatePortalProfile -> ", errCreate)
		return nil, validators.GetErrorResponse(&ctx, localization.E1000, nil, nil)
	}

	return admin2.AdminProfileBuilder(admin), validators.ErrorResponse{}
}
func (oRec *AdminUseCase) UpdateKitchenProfile(ctx context.Context, input *dto.UpdateKitchenProfileDTO) (*admin.AdminProfileResponse, validators.ErrorResponse) {
	admin, errFind := oRec.repo.Find(ctx, utils.ConvertStringIdToObjectId(input.CauserId))
	if errFind != nil {
		oRec.logger.Error("AdminUseCase -> UpdatePortalProfile -> ", errFind)
		return nil, validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
	}

	//input.AdminDetails = utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: input.Name, Operation: "Update My Profile", UpdatedAt: time.Now()}
	input.AdminDetails = utilsDto.AdminDetails{Id: utils.ConvertStringIdToObjectId(input.CauserId), Name: input.CauserName, Type: input.CauserType, Operation: "Update My Profile", UpdatedAt: time.Now()}
	adminDomain, err := admin2.UpdateKitchenProfileBuilder(admin, input)
	if err != nil {
		oRec.logger.Error("AdminUseCase -> UpdatePortalProfile -> ", err)
	}

	admin, errCreate := oRec.repo.Update(ctx, adminDomain)
	if errCreate != nil {
		oRec.logger.Error("AdminUseCase -> UpdatePortalProfile -> ", errCreate)
		return nil, validators.GetErrorResponse(&ctx, localization.E1000, nil, nil)
	}

	return admin2.AdminProfileBuilder(admin), validators.ErrorResponse{}
}

func (oRec *AdminUseCase) LoginAsPortal(ctx context.Context, input *admin3.LoginAsPortalDto) (interface{}, string, validators.ErrorResponse) {
	admin, errFindAdmin := oRec.repo.Find(ctx, utils.ConvertStringIdToObjectId(input.CauserId))
	if errFindAdmin != nil {
		oRec.logger.Error("AdminUseCase -> LoginAsPortal -> ", errFindAdmin)
		return admin, "", validators.GetErrorResponse(&ctx, localization.ErrLoginEmail, nil, utils.GetAsPointer(http.StatusBadRequest)) //change message
	}
	if admin == nil {
		oRec.logger.Error("AdminUseCase -> LoginAsPortal -> Admin Is NULL")
		return admin, "", validators.GetErrorResponse(&ctx, localization.ErrLoginEmail, nil, utils.GetAsPointer(http.StatusBadRequest)) //change message
	}

	//generate token
	data := map[string]interface{}{
		"id":   input.Id,
		"name": input.Name,
	}
	token, errToken := oRec.PortalJwtService.GenerateTokenByAdmin(ctx, utils.ConvertObjectIdToStringId(admin.ID), data)
	if errToken != nil {
		return admin, "", validators.GetErrorResponseFromErr(errToken)
	}

	//update token
	admin.Tokens = append(admin.Tokens, token)
	_, errUpdate := oRec.repo.Update(ctx, admin)
	if errUpdate != nil {
		return admin, "", validators.GetErrorResponse(&ctx, localization.E1000, nil, utils.GetAsPointer(http.StatusBadRequest))
	}

	admin.Account = &domain.Account{
		Id:   utils.ConvertStringIdToObjectId(input.Id),
		Name: domain.Name{Ar: input.Name.Ar, En: input.Name.En},
	}
	return admin2.AdminProfileBuilder(admin), token, validators.ErrorResponse{}
}
