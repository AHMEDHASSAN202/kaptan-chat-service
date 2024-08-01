package domain

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/menu/dto/menu_group"
	menu_group2 "samm/internal/module/menu/repository/structs/menu_group"
	"samm/internal/module/menu/repository/structs/menu_group_item"
	"samm/pkg/utils/dto"
)

type MenuGroupItemCategory struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name   LocalizationText   `json:"name" bson:"name"`
	Icon   string             `json:"icon" bson:"icon"`
	Sort   int                `json:"sort" bson:"sort"`
	Status string             `json:"status" bson:"status"`
}

type ItemMenuGroup struct {
	ID             primitive.ObjectID      `json:"id" bson:"_id,omitempty"`
	LocationIds    []primitive.ObjectID    `json:"location_ids" bson:"location_ids"`
	Availabilities []MenuGroupAvailability `json:"availabilities" bson:"availabilities"`
	Status         string                  `json:"status" bson:"status"`
}

type MenuGroupItem struct {
	mgm.DefaultModel `bson:",inline"`
	ItemId           primitive.ObjectID    `json:"item_id" bson:"item_id"`
	AccountId        primitive.ObjectID    `json:"account_id" bson:"account_id"`
	Name             LocalizationText      `json:"name" bson:"name"`
	Desc             LocalizationText      `json:"desc" bson:"desc"`
	Calories         int                   `json:"calories" bson:"calories"`
	Price            float64               `json:"price" bson:"price"`
	ModifierGroupIds []primitive.ObjectID  `json:"modifier_group_ids" bson:"modifier_group_ids"`
	MenuGroup        ItemMenuGroup         `json:"menu_group" bson:"menu_group"`
	Category         MenuGroupItemCategory `json:"category" bson:"category"`
	Availabilities   []ItemAvailability    `json:"availabilities" bson:"availabilities"`
	Tags             []string              `json:"tags" bson:"tags"`
	Image            string                `json:"image" bson:"image"`
	AdminDetails     []dto.AdminDetails    `json:"admin_details" bson:"admin_details,omitempty"`
	Status           string                `json:"status" bson:"status"`
	Sort             int                   `json:"sort" bson:"sort"`
	dto.ApprovalData `bson:",inline"`
}

type MenuGroupItemRepository interface {
	CreateUpdateBulk(ctx context.Context, models *[]MenuGroupItem) error
	SyncMenuItemsChanges(ctx context.Context, itemId menu_group_item.MenuGroupItemSyncItemModel) error
	FindMenuGroupItem(ctx context.Context, id primitive.ObjectID) (MenuGroupItem, error)
	DeleteByItemId(ctx context.Context, itemId primitive.ObjectID) error
	DeleteBulkByGroupMenuId(ctx context.Context, groupMenuId primitive.ObjectID, exceptionIds []primitive.ObjectID) error
	ChangeStatusByItemId(ctx context.Context, itemId primitive.ObjectID, model MenuGroupItem) error
	ChangeMenuStatus(ctx context.Context, id primitive.ObjectID, dto *menu_group.ChangeMenuGroupStatusDto, adminDetails dto.AdminDetails) error
	ChangeCategoryStatus(ctx context.Context, id primitive.ObjectID, dto *menu_group.ChangeMenuGroupStatusDto, adminDetails dto.AdminDetails) error
	ChangeItemStatus(ctx context.Context, id primitive.ObjectID, dto *menu_group.ChangeMenuGroupStatusDto, adminDetails dto.AdminDetails) error
	DeleteByCategory(ctx context.Context, dto *menu_group.DeleteEntityFromMenuGroupDto) error
	Delete(ctx context.Context, dto *menu_group.DeleteEntityFromMenuGroupDto) error
	MobileGetMenuGroupItems(ctx context.Context, dto *menu_group.GetMenuGroupItemsDTO) (*[]menu_group2.MobileGetMenuGroupItems, error)
	MobileGetMenuGroupItem(ctx context.Context, dto *menu_group.GetMenuGroupItemDTO) (*menu_group2.MobileGetItem, error)
	MobileFilterMenuGroupItemForOrder(ctx context.Context, dto *menu_group.FilterMenuGroupItemsForOrder) ([]menu_group2.MobileGetItem, error)
}
