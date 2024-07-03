package jwt

import "context"

type JwtService interface {
	GenerateToken(ctx context.Context, id string, isTempToken ...bool) (string, error)
	ValidateToken(ctx context.Context, signedToken string, isTempToken ...bool) (interface{}, error)
	GenerateTokenByAdmin(ctx context.Context, id string, data map[string]interface{}) (token string, err error)
}
