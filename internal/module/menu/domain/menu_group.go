package domain

import (
	"context"
	mongopagination "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/menu/consts"
	"samm/internal/module/menu/dto/menu_group"
	menu_group2 "samm/internal/module/menu/repository/structs/menu_group"
	"samm/pkg/utils"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type Category struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name        LocalizationText     `json:"name" bson:"name"`
	Icon        string               `json:"icon" bson:"icon"`
	Sort        int                  `json:"sort" bson:"sort"`
	Status      string               `json:"status" bson:"status"`
	MenuItemIds []primitive.ObjectID `json:"menu_item_ids" bson:"menu_item_ids"`
}

type MenuGroupAvailability struct {
	Day  string `json:"day" bson:"day"`
	From string `json:"from" bson:"from"`
	To   string `json:"to" bson:"to"`
}

type MenuGroup struct {
	mgm.DefaultModel `bson:",inline"`
	AccountId        primitive.ObjectID      `json:"account_id" bson:"account_id"`
	Name             LocalizationText        `json:"name" bson:"name"`
	BranchIds        []primitive.ObjectID    `json:"branch_ids" bson:"branch_ids"`
	Categories       []Category              `json:"categories" bson:"categories"`
	Availabilities   []MenuGroupAvailability `json:"availabilities" bson:"availabilities"`
	Status           string                  `json:"status" bson:"status,omitempty"`
	AdminDetails     []dto.AdminDetails      `json:"admin_details" bson:"admin_details,omitempty"`
}

type MenuGroupUseCase interface {
	Create(ctx context.Context, dto *menu_group.CreateMenuGroupDTO) (string, validators.ErrorResponse)
	Update(ctx context.Context, dto *menu_group.CreateMenuGroupDTO) (string, validators.ErrorResponse)
	Delete(ctx context.Context, menuGroupId primitive.ObjectID) validators.ErrorResponse
	List(ctx context.Context, dto menu_group.ListMenuGroupDTO) (interface{}, validators.ErrorResponse)
	Find(ctx context.Context, id primitive.ObjectID) (interface{}, validators.ErrorResponse)
	ChangeStatus(ctx context.Context, id primitive.ObjectID, input *menu_group.ChangeMenuGroupStatusDto) validators.ErrorResponse
	DeleteEntity(ctx context.Context, input *menu_group.DeleteEntityFromMenuGroupDto) validators.ErrorResponse
}

type MenuGroupRepository interface {
	Create(ctx context.Context, menuGroup *MenuGroup, menuGroupItems *[]MenuGroupItem) (*MenuGroup, error)
	Update(ctx context.Context, menuGroup *MenuGroup, menuGroupItems *[]MenuGroupItem) (*MenuGroup, error)
	Delete(ctx context.Context, domainData *MenuGroup) error
	Find(ctx context.Context, menuGroupId primitive.ObjectID) (*MenuGroup, error)
	List(ctx context.Context, dto menu_group.ListMenuGroupDTO) ([]MenuGroup, *mongopagination.PaginationData, error)
	FindWithItems(ctx context.Context, menuGroupId primitive.ObjectID) (*menu_group2.FindMenuGroupWithItems, error)
	ChangeMenuStatus(ctx context.Context, domainData *MenuGroup, dto *menu_group.ChangeMenuGroupStatusDto, adminDetails dto.AdminDetails) error
	ChangeCategoryStatus(ctx context.Context, domainData *MenuGroup, dto *menu_group.ChangeMenuGroupStatusDto, adminDetails dto.AdminDetails) error
	DeleteCategory(ctx context.Context, domainData *MenuGroup, input *menu_group.DeleteEntityFromMenuGroupDto, adminDetails dto.AdminDetails) error
	DeleteItem(ctx context.Context, domainData *MenuGroup, input *menu_group.DeleteEntityFromMenuGroupDto, adminDetails dto.AdminDetails) error
}

func (model *MenuGroup) Creating(ctx context.Context) error {
	if err := model.DefaultModel.Creating(); err != nil {
		return err
	}
	model.Status = utils.If(model.Status != "", model.Status, consts.MENU_GROUP_DEFUALT_STATUS).(string)
	if model.Categories != nil {
		for i, category := range model.Categories {
			model.Categories[i].Status = utils.If(category.Status != "", category.Status, consts.MENU_GROUP_CATEGORY_DEFUALT_STATUS).(string)
		}
	}
	return nil
}

func (model *MenuGroup) Authorized(accountId primitive.ObjectID) bool {
	return model.AccountId == accountId
}
