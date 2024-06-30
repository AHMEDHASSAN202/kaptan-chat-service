package jwt

import "context"

type JwtService interface {
	GenerateToken(ctx context.Context, id string) (string, error)
	ValidateToken(ctx context.Context, signedToken string) (interface{}, error)
}
