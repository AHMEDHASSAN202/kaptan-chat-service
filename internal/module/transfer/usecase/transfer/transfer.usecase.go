package transfer

import (
	"context"
	"fmt"
	domain4 "kaptan/internal/module/chat/domain"
	domain3 "kaptan/internal/module/transfer/domain"
	"kaptan/internal/module/transfer/dto"
	domain2 "kaptan/internal/module/user/domain"
	"kaptan/pkg/fcm_notification"
	"kaptan/pkg/gate"
	"kaptan/pkg/logger"
	"kaptan/pkg/validators"
	"kaptan/pkg/websocket"
)

type UseCase struct {
	logger           logger.ILogger
	gate             *gate.Gate
	websocketManager *websocket.ChannelManager
	driverRepo       domain2.DriverRepository
	transferRepo     domain3.TransferRepository
	chatUseCase      domain4.ChatUseCase
	fcmClient        *fcm_notification.FCMClient
}

func NewTransferUseCase(chatUseCase domain4.ChatUseCase, driverRepo domain2.DriverRepository, transferRepo domain3.TransferRepository, gate *gate.Gate, websocketManager *websocket.ChannelManager, logger logger.ILogger, fcmClient *fcm_notification.FCMClient) domain3.UseCase {
	return &UseCase{
		chatUseCase:      chatUseCase,
		logger:           logger,
		gate:             gate,
		driverRepo:       driverRepo,
		websocketManager: websocketManager,
		transferRepo:     transferRepo,
		fcmClient:        fcmClient,
	}
}

func (u UseCase) StartTransfer(ctx *context.Context, dto *dto.StartTransfer) validators.ErrorResponse {
	chat, err := u.chatUseCase.GetAcceptedChatByTransferId(*ctx, dto.TransferId, dto.CauserId)
	if err.IsError {
		return err
	}

	//transfer, errTransfer := u.transferRepo.MarkTransferAsStart(ctx, dto)
	//if errTransfer != nil {
	//	return validators.GetErrorResponseFromErr(errTransfer)
	//}
	//
	//transferResponse := builder.TransferResponseBuilder(transfer)

	go func() {
		//push message here
		fmt.Println(chat)
	}()

	//contentJson, _ := json.Marshal(transferResponse)
	//
	//u.websocketManager.Broadcast <- websocket.Message{
	//	ChannelID: chatResponse.Channel,
	//	Content:   string(contentJson),
	//	Action:    consts.CHANCE_CHAT_STATUS_ACTION,
	//}

	return validators.ErrorResponse{}
}

func (u UseCase) EndTransfer(ctx *context.Context, dto *dto.EndTransfer) validators.ErrorResponse {
	chat, err := u.chatUseCase.GetAcceptedChatByTransferId(*ctx, dto.TransferId, dto.CauserId)
	if err.IsError {
		return err
	}

	//transfer, errTransfer := u.transferRepo.MarkTransferAsEnd(ctx, dto)
	//if errTransfer != nil {
	//	return validators.GetErrorResponseFromErr(errTransfer)
	//}
	//
	//transferResponse := builder.TransferResponseBuilder(transfer)

	go func() {
		//push message here
		fmt.Println(chat)
	}()

	//contentJson, _ := json.Marshal(transferResponse)
	//
	//u.websocketManager.Broadcast <- websocket.Message{
	//	ChannelID: chatResponse.Channel,
	//	Content:   string(contentJson),
	//	Action:    consts.CHANCE_CHAT_STATUS_ACTION,
	//}

	return validators.ErrorResponse{}
}
