package domain

import (
	"context"
	mongopagination "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/admin/dto/role"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type LocalizeText struct {
	Ar string `json:"ar" bson:"ar"`
	En string `json:"en" bson:"en"`
}

type Role struct {
	mgm.DefaultModel `bson:",inline"`
	Name             LocalizeText       `json:"name" bson:"name"`
	Permissions      []string           `json:"permissions" bson:"permissions"`
	AdminDetails     []dto.AdminDetails `json:"admin_details" bson:"admin_details,omitempty"`
}

type RoleUseCase interface {
	Create(ctx context.Context, dto *role.CreateRoleDTO) (string, validators.ErrorResponse)
	Update(ctx context.Context, dto *role.CreateRoleDTO) (string, validators.ErrorResponse)
	Delete(ctx context.Context, roleId primitive.ObjectID) validators.ErrorResponse
	List(ctx context.Context, dto *role.ListRoleDTO) (interface{}, validators.ErrorResponse)
	Find(ctx context.Context, roleId primitive.ObjectID) (interface{}, validators.ErrorResponse)
}

type RoleRepository interface {
	Create(ctx context.Context, domainData *Role) (*Role, error)
	Update(ctx context.Context, domainData *Role) (*Role, error)
	Delete(ctx context.Context, domainData *Role) error
	Find(ctx context.Context, roleId primitive.ObjectID) (*Role, error)
	List(ctx context.Context, dto *role.ListRoleDTO) ([]Role, *mongopagination.PaginationData, error)
}
