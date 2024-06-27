package user

import (
	"crypto/rand"
	"encoding/base32"
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

func generateOTP() (otp string, err error) {
	randomBytes := make([]byte, 4)
	_, err = rand.Read(randomBytes)
	if err != nil {
		return
	}
	otp = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	return
}

func otpTrialsPerDaySetter() string {

	return ""
}

func otpTrialsPerDayGetter() int {

	return 0
}

//func domainBuilderChangeStatus(dto *brand.ChangeBrandStatusDto, domainData *domain.Brand) *domain.Brand {
//	brandDoc := domain.Brand{}
//	copier.Copy(&brandDoc, domainData)
//	brandDoc.IsActive = dto.IsActive
//	return &brandDoc
//}
