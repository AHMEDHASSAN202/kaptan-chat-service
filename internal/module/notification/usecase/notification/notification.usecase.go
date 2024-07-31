package notification

import (
	"context"
	"github.com/asaskevich/EventBus"
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
	"strings"
	"time"
)

type NotificationUseCase struct {
	repo             domain.NotificationRepository
	logger           logger.ILogger
	extService       *external.ExtService
	oneSignalService domain.OneSignalService
	bus              EventBus.Bus
}

const tag = " NotificationUseCase "

func NewNotificationUseCase(repo domain.NotificationRepository, logger logger.ILogger, extService *external.ExtService, oneSignalService domain.OneSignalService, bus EventBus.Bus) domain.NotificationUseCase {
	l := &NotificationUseCase{
		repo:             repo,
		logger:           logger,
		extService:       extService,
		oneSignalService: oneSignalService,
		bus:              bus,
	}
	bus.SubscribeAsync(consts.SEND_NOTIFICATION, l.SendPushNotificationV2, false)
	return l
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
	if !causerDetails.Id.IsZero() {
		notificationDomain.AdminDetails = &causerDetails

	}
	dbErr := l.repo.CreateNotification(&notificationDomain)
	if dbErr != nil {
		return validators.GetErrorResponseFromErr(dbErr)
	}

	//Send notification
	var notificationData notification.NotificationDto
	copier.Copy(&notificationData, payload)
	notificationData.Ids = payload.UserIds
	notificationData.ModelType = consts2.UserModelType
	l.SendPushNotification(notificationData)

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

func (l NotificationUseCase) SendPushNotification(dto notification.NotificationDto) (err validators.ErrorResponse) {
	ctx := context.Background()
	// Get Player Ids Based On Type
	playerIDs := make([]string, 0)
	if dto.Type == consts.TYPE_PRIVATE && dto.ModelType == consts2.UserModelType {
		_playerIDs, errRe := l.extService.UserService.GetUsersPlayerIds(ctx, dto.Ids)
		if errRe.IsError || len(_playerIDs) == 0 {
			l.logger.Error(tag+" => Error Get Player Ids", errRe, len(_playerIDs))
			return validators.ErrorResponse{}
		}
		playerIDs = _playerIDs
	}
	// For Kitchen
	if dto.Type == consts.TYPE_PRIVATE && dto.ModelType == consts2.LocationModelType {
		_playerIDs, errRe := l.extService.KitchenService.GetKitchenPlayerIds(ctx, dto.Ids, dto.AccountIds)
		if errRe.IsError || len(_playerIDs) == 0 {
			l.logger.Error(tag+" => Error Get Player Ids", errRe, len(_playerIDs))
			return validators.ErrorResponse{}
		}
		playerIDs = _playerIDs
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
	if dto.Type == consts.TYPE_PUBLIC {
		pushNotificationRequest.Filters = []map[string]string{{"field": "country", "relation": "=", "value": strings.ToUpper(dto.CountryId)}}

	}
	notificationResponse, errRe := l.oneSignalService.SendNotification(ctx, &pushNotificationRequest, dto.ModelType)
	if errRe.IsError {
		l.logger.Error(tag+" => Error Response When Send", errRe)
		return errRe
	}
	l.logger.Info(tag+" => Error Response When Send", notificationResponse)
	return
}

func (l NotificationUseCase) SendPushNotificationV2(dto notification.GeneralNotification) (err validators.ErrorResponse) {
	//ctx := context.Background()
	notificationCode, err := GetNotificationByCode(dto.NotificationCode, dto.NotificationData)
	if err.IsError {
		l.logger.Error(err)
		return
	}

	for _, toModel := range dto.To {
		// Prepare Dto
		notificationMessage := notificationCode.User
		if toModel.Model == consts2.LocationModelType {
			notificationMessage = notificationCode.Kitchen
		}
		var notificationData notification.NotificationDto
		notificationData.Title.Ar = notificationMessage.Title.Ar
		notificationData.Title.En = notificationMessage.Title.En
		notificationData.Text.En = notificationMessage.Description.Ar
		notificationData.Text.Ar = notificationMessage.Description.En
		notificationData.Type = consts.TYPE_PRIVATE
		notificationData.Ids = []string{toModel.Id}
		notificationData.AccountIds = []string{toModel.AccountId}
		notificationData.ModelType = toModel.Model
		notificationData.CountryId = dto.Country

		if toModel.Model == consts2.UserModelType && toModel.LogNotification {
			var storeNotificationDto notification.StoreNotificationDto
			copier.Copy(&storeNotificationDto, notificationData)
			storeNotificationDto.UserIds = notificationData.Ids
			go l.CreateNotification(context.Background(), &storeNotificationDto)
		} else {
			l.SendPushNotification(notificationData)
		}
	}
	return
}
