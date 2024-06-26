package user

import (
	"github.com/jinzhu/copier"
	"samm/internal/module/user/domain"
	"samm/internal/module/user/dto/user"
	"samm/pkg/utils"
	"time"
)

func domainBuilderAtUpdateProfile(dto *user.UpdateUserProfileDto, domainData *domain.User) *domain.User {
	copier.Copy(&domainData, dto)
	domainData.ID = utils.ConvertStringIdToObjectId(dto.ID)
	domainData.UpdatedAt = time.Now()
	return domainData
}

//func domainBuilderChangeStatus(dto *brand.ChangeBrandStatusDto, domainData *domain.Brand) *domain.Brand {
//	brandDoc := domain.Brand{}
//	copier.Copy(&brandDoc, domainData)
//	brandDoc.IsActive = dto.IsActive
//	return &brandDoc
//}
