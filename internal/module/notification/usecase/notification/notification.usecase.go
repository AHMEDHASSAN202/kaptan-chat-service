package notification

import (
	"context"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/notification/consts"
	"samm/internal/module/notification/domain"
	"samm/internal/module/notification/dto/notification"
	"samm/internal/module/notification/external"
	consts2 "samm/internal/module/notification/gateways/onesignal/consts"
	"samm/internal/module/notification/gateways/onesignal/requests"
	"samm/internal/module/notification/responses"
	"samm/pkg/logger"
	"samm/pkg/utils"
	utilsDto "samm/pkg/utils/dto"
	"samm/pkg/validators"
	"time"
)

type NotificationUseCase struct {
	repo             domain.NotificationRepository
	logger           logger.ILogger
	extService       *external.ExtService
	oneSignalService domain.OneSignalService
}

const tag = " NotificationUseCase "

func NewNotificationUseCase(repo domain.NotificationRepository, logger logger.ILogger, extService *external.ExtService, oneSignalService domain.OneSignalService) domain.NotificationUseCase {
	return &NotificationUseCase{
		repo:             repo,
		logger:           logger,
		extService:       extService,
		oneSignalService: oneSignalService,
	}
}

func (l NotificationUseCase) CreateNotification(ctx context.Context, payload *notification.StoreNotificationDto) (err validators.ErrorResponse) {
	notificationDomain := domain.Notification{}
	copier.Copy(&notificationDomain, payload)
	notificationDomain.CountryId = payload.CountryId
	notificationDomain.CreatedAt = time.Now()
	notificationDomain.UpdatedAt = time.Now()
	notificationDomain.UserIds = make([]primitive.ObjectID, 0)
	if len(payload.UserIds) > 0 {
		notificationDomain.UserIds = utils.ConvertStringIdsToObjectIds(payload.UserIds)
	}
	if notificationDomain.Type == consts.REDIRECT_LOCATION {
		notificationDomain.RedirectData.LocationId = utils.ConvertStringIdToObjectId(payload.LocationId)
	}
	causerDetails := utilsDto.AdminDetails{Id: utils.ConvertStringIdToObjectId(payload.CauserId), Name: payload.CauserName, Type: payload.CauserType, Operation: "Create", UpdatedAt: time.Now()}
	notificationDomain.AdminDetails = causerDetails
	dbErr := l.repo.CreateNotification(&notificationDomain)
	if dbErr != nil {
		return validators.GetErrorResponseFromErr(dbErr)
	}

	// Send notification
	var notificationData notification.NotificationDto
	copier.Copy(&notificationData, payload)
	notificationData.Ids = payload.UserIds
	notificationData.ModelType = consts2.UserModelType

	l.SendPushNotification(ctx, &notificationData)
	return
}

func (l NotificationUseCase) FindNotification(ctx context.Context, Id string) (notification domain.Notification, err validators.ErrorResponse) {
	domainNotification, dbErr := l.repo.FindNotification(ctx, utils.ConvertStringIdToObjectId(Id))
	if dbErr != nil {
		return *domainNotification, validators.GetErrorResponseFromErr(dbErr)
	}
	return *domainNotification, validators.ErrorResponse{}
}

func (l NotificationUseCase) DeleteNotification(ctx context.Context, Id string) (err validators.ErrorResponse) {

	delErr := l.repo.DeleteNotification(ctx, utils.ConvertStringIdToObjectId(Id))
	if delErr != nil {
		return validators.GetErrorResponseFromErr(delErr)
	}
	return validators.ErrorResponse{}
}

func (l NotificationUseCase) List(ctx *context.Context, dto *notification.ListNotificationDto) (*responses.ListResponse, validators.ErrorResponse) {
	users, paginationMeta, resErr := l.repo.List(ctx, dto)
	if resErr != nil {
		return nil, validators.GetErrorResponseFromErr(resErr)
	}
	return responses.SetListResponse(users, paginationMeta), validators.ErrorResponse{}
}

func (l NotificationUseCase) ListMobile(ctx *context.Context, dto *notification.ListNotificationMobileDto) (*responses.ListResponse, validators.ErrorResponse) {
	users, paginationMeta, resErr := l.repo.ListMobile(ctx, dto)
	if resErr != nil {
		return nil, validators.GetErrorResponseFromErr(resErr)
	}
	return responses.SetListResponse(users, paginationMeta), validators.ErrorResponse{}
}

func (l NotificationUseCase) SendPushNotification(ctx context.Context, dto *notification.NotificationDto) (err validators.ErrorResponse) {
	// Get Player Ids Based On Type
	playerIDs := make([]string, 0)
	if dto.Type == consts.TYPE_PRIVATE && dto.ModelType == consts2.UserModelType {
		playerIDs, errRe := l.extService.UserService.GetUsersPlayerIds(ctx, dto.Ids)
		if errRe.IsError || len(playerIDs) == 0 {
			l.logger.Error(tag, errRe, len(playerIDs))
			return validators.ErrorResponse{}
		}
	}
	// For Kitchen
	if dto.Type == consts.TYPE_PRIVATE && dto.ModelType == consts2.KitchenModelType {
		//playerIDs, errRe := l.extService.UserService.GetUsersPlayerIds(ctx, dto.Ids)
		//if errRe.IsError || len(playerIDs) == 0 {
		//	l.logger.Error(tag, errRe, len(playerIDs))
		//	return validators.ErrorResponse{}
		//}
	}

	// Prepare dto to send with onesignal
	pushNotificationRequest := requests.PushNotificationPayload{
		Headings: requests.Translated{
			Ar: dto.Title.Ar,
			En: dto.Title.En,
		},
		Contents: requests.Translated{
			Ar: dto.Text.Ar,
			En: dto.Text.En,
		},
		IncludePlayerIds: playerIDs,
		Name:             dto.ModelType,
	}
	notificationResponse, errRe := l.oneSignalService.SendNotification(ctx, &pushNotificationRequest, dto.ModelType)
	if errRe.IsError {
		l.logger.Error(tag, errRe)
		return errRe
	}
	l.logger.Info(tag, notificationResponse)
	return
}
