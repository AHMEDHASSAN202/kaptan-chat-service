package helper

import (
	"context"
	"encoding/json"
	"fmt"
	. "github.com/ahmetb/go-linq/v3"
	"math/rand"
	"os"
	"path/filepath"
	"samm/internal/module/order/consts"
	"samm/internal/module/order/domain"
	"samm/pkg/logger"
	"samm/pkg/validators"
	"time"
)

func GenerateSerialNumber() string {
	currentTime := time.Now().Format("020106")
	randomNumber := rand.Intn(900000) + 10000
	serialNumber := fmt.Sprintf("%s-%06d", currentTime, randomNumber)
	return serialNumber
}
func GetNextAndPreviousStatusByType(actor string, currentStatus string, nextStatus string) (nextStatuses []string, previousStatus []string) {

	var orderStatus map[string]domain.OrderStatusJson
	dir, err := os.Getwd()
	if err != nil {
		return nextStatuses, previousStatus
	}

	path := filepath.Join(dir, "internal", "module", "order", "assets", "order_status.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nextStatuses, previousStatus
	}

	if errRe := json.Unmarshal(data, &orderStatus); errRe != nil {
		return nextStatuses, previousStatus
	}
	statusRule := orderStatus[currentStatus]
	nextStatusRule := orderStatus[nextStatus]

	switch actor {
	case consts.ActorAdmin:
		return statusRule.AllowAdminToChange, nextStatusRule.PreviousStatus
	case consts.ActorKitchen:
		return statusRule.AllowKitchenToChange, nextStatusRule.PreviousStatus
	case consts.ActorUser:
		return statusRule.AllowUserToChange, nextStatusRule.PreviousStatus
	}
	return nextStatuses, previousStatus

}

func KitchenRejectionReasons(ctx context.Context, status string, id string) ([]domain.KitchenRejectionReason, validators.ErrorResponse) {
	kitchenRejectionReason := make([]domain.KitchenRejectionReason, 0)
	dir, err := os.Getwd()
	if err != nil {
		return kitchenRejectionReason, validators.GetErrorResponseFromErr(err)
	}

	path := filepath.Join(dir, "internal", "module", "order", "assets", "kitchen_cancel_reasons.json")
	data, err := os.ReadFile(path)
	if err != nil {
		logger.Logger.Error(err)
		return kitchenRejectionReason, validators.GetErrorResponseFromErr(err)
	}

	if errRe := json.Unmarshal(data, &kitchenRejectionReason); errRe != nil {
		logger.Logger.Error(err)
		return kitchenRejectionReason, validators.GetErrorResponseFromErr(errRe)
	}

	// Handle Status
	if status != "" {
		From(kitchenRejectionReason).Where(func(c interface{}) bool {
			return c.(domain.KitchenRejectionReason).Status == status || c.(domain.KitchenRejectionReason).Status == "all"
		}).ToSlice(&kitchenRejectionReason)
	}
	if id != "" {
		From(kitchenRejectionReason).Where(func(c interface{}) bool {
			return c.(domain.KitchenRejectionReason).Id == id
		}).ToSlice(&kitchenRejectionReason)
	}

	return kitchenRejectionReason, validators.ErrorResponse{}
}
func UserRejectionReasons(ctx context.Context, status string, id string) ([]domain.UserRejectionReason, validators.ErrorResponse) {
	userRejectionReason := make([]domain.UserRejectionReason, 0)
	dir, err := os.Getwd()
	if err != nil {
		return userRejectionReason, validators.GetErrorResponseFromErr(err)
	}

	path := filepath.Join(dir, "internal", "module", "order", "assets", "user_cancel_reasons.json")
	data, err := os.ReadFile(path)
	if err != nil {
		logger.Logger.Error("Read Json File -> Error -> ", err)
		return userRejectionReason, validators.GetErrorResponseFromErr(err)
	}

	if errRe := json.Unmarshal(data, &userRejectionReason); errRe != nil {
		logger.Logger.Error("ListPermissions -> Error -> ", errRe)
		return userRejectionReason, validators.GetErrorResponseFromErr(errRe)
	}

	// Handle Status
	if status != "" {
		From(userRejectionReason).Where(func(c interface{}) bool {
			return c.(domain.UserRejectionReason).Status == status || c.(domain.UserRejectionReason).Status == "all"
		}).ToSlice(&userRejectionReason)
	}
	if id != "" {
		From(userRejectionReason).Where(func(c interface{}) bool {
			return c.(domain.UserRejectionReason).Id == id
		}).ToSlice(&userRejectionReason)
	}

	return userRejectionReason, validators.ErrorResponse{}
}
