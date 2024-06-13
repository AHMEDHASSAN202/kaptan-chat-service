package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
	location "samm/internal/module/common/dto"
	"samm/pkg/validators"
)

type Name struct {
	Ar string `json:"ar"`
	En string `json:"en"`
}

type City struct {
	Id   primitive.ObjectID `json:"_id" bson:"id"`
	Name Name               `json:"name" bson:"name"`
}

type CommonUseCase interface {
	ListCities(ctx context.Context, payload *location.ListCitiesDto) (data interface{}, err validators.ErrorResponse)
	ListCountries(ctx context.Context) (data interface{}, err validators.ErrorResponse)
	UploadFile(ctx context.Context, file *multipart.FileHeader, filePath string) (string, validators.ErrorResponse)
	ReadFile(ctx context.Context, filePath string) (string, validators.ErrorResponse)
}

type CommonRepository interface {
}
