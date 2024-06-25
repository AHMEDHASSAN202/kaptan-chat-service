package user

import (
	"github.com/jinzhu/copier"
	"samm/internal/module/user/domain"
	user "samm/internal/module/user/dto/User"
	"samm/pkg/utils"
)

func domainBuilderAtUpdate(dto *user.UpdateUserProfileDto, domainData *domain.User) *domain.User {
	userDoc := domain.User{}
	copier.Copy(&userDoc, dto)
	userDoc.ID = utils.ConvertStringIdToObjectId(dto.ID)
	userDoc.CreatedAt = domainData.CreatedAt
	return &userDoc
}

//func domainBuilderChangeStatus(dto *brand.ChangeBrandStatusDto, domainData *domain.Brand) *domain.Brand {
//	brandDoc := domain.Brand{}
//	copier.Copy(&brandDoc, domainData)
//	brandDoc.IsActive = dto.IsActive
//	return &brandDoc
//}
