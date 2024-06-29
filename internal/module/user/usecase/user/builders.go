package user

import (
	"crypto/rand"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"math/big"
	"samm/internal/module/user/domain"
	"samm/internal/module/user/dto/user"
	"samm/pkg/utils"
	"strconv"
	"strings"
	"time"
)

const maxNumOfTrialsPerDay = 3
const jwtSecretKey = "jwtSecret"

type UserClaim struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func domainBuilderAtCreateProfile(userDomain *domain.User, payload *user.SendUserOtpDto, otp, ctr string) *domain.User {
	expiry := time.Now().Add(5 * time.Minute)
	userDomain.Otp = otp
	userDomain.OtpCounter = ctr
	userDomain.ExpiryOtpDate = &expiry
	userDomain.PhoneNumber = payload.PhoneNumber
	userDomain.CountryCode = payload.CountryCode
	userDomain.Country = payload.CountryId
	return userDomain
}

func domainBuilderAtUpdateProfile(dto *user.UpdateUserProfileDto, domainData *domain.User) *domain.User {
	copier.Copy(&domainData, dto)
	domainData.ID = utils.ConvertStringIdToObjectId(dto.ID)
	domainData.UpdatedAt = time.Now()
	return domainData
}

func domainBuilderAtSignUp(dto *user.UserSignUpDto, userToken string, domainData *domain.User) *domain.User {
	domainData.UpdatedAt = time.Now()
	domainData.Name = dto.Name
	domainData.Email = dto.Email
	domainData.IsActive = true
	domainData.Tokens = append(domainData.Tokens, userToken)
	return domainData
}

// generateOTP generates a 4-digit numeric OTP
func generateOTP() (string, error) {
	const otpLength = 4
	const otpDigits = "0123456789"

	otp := make([]byte, otpLength)
	for i := 0; i < otpLength; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(otpDigits))))
		if err != nil {
			return "", err
		}
		otp[i] = otpDigits[num.Int64()]
	}

	return string(otp), nil
}

// 2024-06-10@1
func otpTrialsPerDaySetter(otpCounter string) (string, error) {
	toDay := time.Now().Format(utils.DefaultDateFormat)
	day, count := otpTrialsPerDayGetter(otpCounter)
	if day == toDay {
		if count >= maxNumOfTrialsPerDay {
			return "", errors.New("you have exceeded the limit per day")
		}
		count++
		return fmt.Sprintf("%s@%d", day, count), nil
	}

	return fmt.Sprintf("%s@1", toDay), nil
}

func otpTrialsPerDayGetter(otpCounter string) (day string, counter int) {
	// Split the otpCounter string by '@'
	parts := strings.Split(otpCounter, "@")
	if len(parts) != 2 {
		fmt.Println("Empty")
		return
	}

	day = parts[0]
	counter, _ = strconv.Atoi(parts[1])

	return
}

// todo update user claim
func generateUserToken(userDomain *domain.User) (string, error) {
	userClaim := UserClaim{
		UserID: userDomain.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(1, 0, 0).Unix(),
		},
	}

	token, err := utils.CreateJWT(userClaim, jwtSecretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func validateUserToken(token string) (claims *UserClaim, err error) {
	err = utils.ValidateJWT(token, jwtSecretKey, claims)
	if err != nil {
		return
	}
	return
}
