package notification

import (
	"samm/pkg/utils/dto"
)

type ListNotificationMobileDto struct {
	dto.Pagination
	dto.MobileHeaders
}
