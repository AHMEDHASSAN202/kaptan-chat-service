package jwt

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"kaptan/pkg/localization"
	"kaptan/pkg/logger"
	"kaptan/pkg/validators"
	"time"
)

type PortalJwtService struct {
	secretKey    string
	ExpiredHours time.Duration
	logger       logger.ILogger
}

type PortalJwtClaim struct {
	jwt.RegisteredClaims
	CauserId      string                 `json:"causer_id"`
	CauserType    string                 `json:"causer_type"`
	CauserDetails map[string]interface{} `json:"causer_data"`
}

func (jwtService *PortalJwtService) GenerateToken(ctx context.Context, id string, isTempToken ...bool) (token string, err error) {
	expiredAt := time.Now().Add(time.Duration(jwtService.ExpiredHours.Hours()) * time.Hour)
	claims := &PortalJwtClaim{
		CauserId:   id,
		CauserType: "portal",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Ktha",
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        primitive.NewObjectID().Hex(),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = jwtToken.SignedString([]byte(jwtService.secretKey))
	if err != nil {
		jwtService.logger.Error(ctx, err)
		err = validators.GetError(&ctx, localization.JwtSigningError, nil)
		return
	}
	return
}

func (jwtService *PortalJwtService) GenerateTokenByAdmin(ctx context.Context, id string, data map[string]interface{}) (token string, err error) {
	expiredAt := time.Now().Add(time.Duration(jwtService.ExpiredHours.Hours()) * time.Hour)
	claims := &PortalJwtClaim{
		CauserId:      id,
		CauserType:    "portal",
		CauserDetails: data,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Ktha",
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        primitive.NewObjectID().Hex(),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = jwtToken.SignedString([]byte(jwtService.secretKey))
	if err != nil {
		jwtService.logger.Error(ctx, err)
		err = validators.GetError(&ctx, localization.JwtSigningError, nil)
		return
	}
	return
}

func (jwtService *PortalJwtService) ValidateToken(ctx context.Context, signedToken string, isTempToken ...bool) (interface{}, error) {
	token, err := jwt.ParseWithClaims(
		signedToken, &PortalJwtClaim{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtService.secretKey), nil
		},
	)
	if err != nil {
		jwtService.logger.Error(ctx, err)
		return nil, validators.GetError(&ctx, localization.JwtTokenInvalidError, nil)
	}

	claims, ok := token.Claims.(*PortalJwtClaim)
	if ok && token.Valid {
		return claims, nil
	}

	if ve, ok := err.(*jwt.ValidationError); ok {
		jwtService.logger.Error(ctx, err)
		switch {
		case ve.Errors&jwt.ValidationErrorMalformed != 0:
			return nil, validators.GetError(&ctx, localization.JwtTokenInvalidError, nil)
		case ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0:
			return nil, validators.GetError(&ctx, localization.JwtTokenExpiredError, nil)
		default:
			return nil, validators.GetError(&ctx, localization.JwtTokenParsingError, nil)
		}
	}

	jwtService.logger.Error(ctx, err)
	return claims, validators.GetError(&ctx, localization.JwtTokenParsingError, nil)
}
