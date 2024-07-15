package domain

import (
	"context"
	mongopagination "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/admin/consts"
	"samm/internal/module/admin/dto/admin"
	"samm/internal/module/admin/dto/auth"
	admin2 "samm/internal/module/admin/responses/admin"
	"samm/pkg/utils"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"time"
)

type MetaData struct {
}

type Name struct {
	Ar string `json:"ar" bson:"ar"`
	En string `json:"en" bson:"en"`
}

type Account struct {
	Id   primitive.ObjectID `json:"id" bson:"_id"`
	Name Name               `json:"name" bson:"name"`
}
type Kitchen struct {
	Id            primitive.ObjectID `json:"id" bson:"_id"`
	Name          Name               `json:"name" bson:"name"`
	AllowedStatus []string           `json:"allowed_status" bson:"allowed_status"`
}

type Admin struct {
	mgm.DefaultModel  `bson:",inline"`
	Name              string             `json:"name" bson:"name"`
	Email             string             `json:"email" bson:"email"`
	Password          string             `json:"password" bson:"password,omitempty"`
	EncryptedPassword string             `json:"encrypted_password" bson:"encrypted_password,omitempty"`
	Type              string             `json:"type" bson:"type"`
	Role              Role               `json:"role" bson:"role"`
	CountryIds        []string           `json:"country_ids" bson:"country_ids"`
	Status            string             `json:"status" bson:"status"`
	Tokens            []string           `json:"tokens" bson:"tokens,omitempty"`
	MetaData          MetaData           `json:"meta_data" bson:"meta_data"`
	Account           *Account           `json:"account" bson:"account,omitempty"`
	Kitchen           *Kitchen           `json:"kitchen" bson:"kitchen,omitempty"`
	AdminDetails      []dto.AdminDetails `json:"admin_details" bson:"admin_details,omitempty"`
	DeletedAt         *time.Time         `json:"deleted_at" bson:"deleted_at"`
}

type AdminUseCase interface {
	Create(ctx context.Context, dto *admin.CreateAdminDTO) (string, validators.ErrorResponse)
	Update(ctx context.Context, dto *admin.CreateAdminDTO) (string, validators.ErrorResponse)
	Delete(ctx context.Context, adminId primitive.ObjectID, accountId string) validators.ErrorResponse
	DeleteBy(ctx context.Context, id primitive.ObjectID, key string) validators.ErrorResponse
	List(ctx context.Context, dto *admin.ListAdminDTO) (interface{}, validators.ErrorResponse)
	Find(ctx context.Context, adminId primitive.ObjectID, accountId string) (interface{}, validators.ErrorResponse)
	ChangeStatus(ctx context.Context, input *admin.ChangeAdminStatusDto) validators.ErrorResponse
	CheckEmailExists(ctx context.Context, email string, adminId primitive.ObjectID, adminType string) (bool, validators.ErrorResponse)
	CheckRoleExists(ctx context.Context, roleId primitive.ObjectID) (bool, validators.ErrorResponse)
	AdminLogin(ctx context.Context, dto *auth.AdminAuthDTO) (interface{}, string, validators.ErrorResponse)
	PortalLogin(ctx context.Context, dto *auth.PortalAuthDTO) (interface{}, string, validators.ErrorResponse)
	KitchenLogin(ctx context.Context, input *auth.KitchenAuthDTO) (interface{}, string, validators.ErrorResponse)
	Profile(ctx context.Context, profileDTO auth.ProfileDTO) (*admin2.AdminProfileResponse, validators.ErrorResponse)
	KitchenProfile(ctx context.Context, profileDTO auth.KitchenProfileDTO) (*admin2.AdminProfileResponse, validators.ErrorResponse)
	UpdateAdminProfile(ctx context.Context, dto *auth.UpdateAdminProfileDTO) (*admin2.AdminProfileResponse, validators.ErrorResponse)
	UpdatePortalProfile(ctx context.Context, dto *auth.UpdatePortalProfileDTO) (*admin2.AdminProfileResponse, validators.ErrorResponse)
	UpdateKitchenProfile(ctx context.Context, input *auth.UpdateKitchenProfileDTO) (*admin2.AdminProfileResponse, validators.ErrorResponse)
	SyncAccount(ctx context.Context, input admin.Account) validators.ErrorResponse
	LoginAsPortal(ctx context.Context, portalDto *admin.LoginAsPortalDto) (interface{}, string, validators.ErrorResponse)
}

type AdminRepository interface {
	Create(ctx context.Context, domainData *Admin) (*Admin, error)
	Update(ctx context.Context, domainData *Admin) (*Admin, error)
	SyncRole(ctx context.Context, domainData *Role) error
	Delete(ctx context.Context, domainData *Admin, adminDetails dto.AdminDetails) error
	Find(ctx context.Context, adminId primitive.ObjectID) (*Admin, error)
	FindByToken(ctx context.Context, token string, adminType []string) (*Admin, error)
	List(ctx context.Context, dto *admin.ListAdminDTO) ([]Admin, *mongopagination.PaginationData, error)
	ChangeStatus(ctx context.Context, model *Admin, input *admin.ChangeAdminStatusDto, adminDetails dto.AdminDetails) error
	CheckEmailExists(ctx context.Context, email string, adminId primitive.ObjectID, adminType string) (bool, error)
	CheckRoleExists(ctx context.Context, roleId primitive.ObjectID) (bool, error)
	FindByEmail(ctx context.Context, email string, adminType string) (*Admin, error)
	SyncAccount(ctx context.Context, input admin.Account) error
}

func (model *Admin) Creating(ctx context.Context) error {
	if err := model.DefaultModel.Creating(); err != nil {
		return err
	}
	model.Status = utils.If(model.Status != "", model.Status, consts.ADMIN_DEFUALT_STATUS).(string)
	return nil
}

func (model *Admin) SetSoftDelete(ctx context.Context) error {
	model.DeletedAt = utils.GetAsPointer(time.Now().UTC())
	return nil
}

func (model *Admin) IsActive() bool {
	return model.Status == "active"
}

func (model *Admin) Authorized(accountId string) bool {
	return model.Account.Id == utils.ConvertStringIdToObjectId(accountId)
}
