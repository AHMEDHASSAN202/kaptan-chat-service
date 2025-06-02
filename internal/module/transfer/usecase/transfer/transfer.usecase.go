package transfer

import (
	"context"
	"fmt"
	domain4 "kaptan/internal/module/chat/domain"
	"kaptan/internal/module/transfer/builder"
	domain3 "kaptan/internal/module/transfer/domain"
	"kaptan/internal/module/transfer/dto"
	"kaptan/internal/module/transfer/responses/app"
	domain2 "kaptan/internal/module/user/domain"
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
}

func NewTransferUseCase(chatUseCase domain4.ChatUseCase, driverRepo domain2.DriverRepository, transferRepo domain3.TransferRepository, gate *gate.Gate, websocketManager *websocket.ChannelManager, logger logger.ILogger) domain3.UseCase {
	return &UseCase{
		chatUseCase:      chatUseCase,
		logger:           logger,
		gate:             gate,
		driverRepo:       driverRepo,
		websocketManager: websocketManager,
		transferRepo:     transferRepo,
	}
}

func (u UseCase) StartTransfer(ctx *context.Context, dto *dto.StartTransfer) (*app.TransferResponse, validators.ErrorResponse) {
	chat, err := u.chatUseCase.GetAcceptedChatByTransferId(*ctx, dto.TransferId, dto.CauserId)
	if err.IsError {
		return nil, err
	}

	transfer, errTransfer := u.transferRepo.MarkTransferAsStart(ctx, dto)
	if errTransfer != nil {
		return nil, validators.GetErrorResponseFromErr(errTransfer)
	}

	transferResponse := builder.TransferResponseBuilder(transfer)

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

	return transferResponse, validators.ErrorResponse{}
}

func (u UseCase) EndTransfer(ctx *context.Context, dto *dto.EndTransfer) (*app.TransferResponse, validators.ErrorResponse) {
	chat, err := u.chatUseCase.GetAcceptedChatByTransferId(*ctx, dto.TransferId, dto.CauserId)
	if err.IsError {
		return nil, err
	}

	transfer, errTransfer := u.transferRepo.MarkTransferAsEnd(ctx, dto)
	if errTransfer != nil {
		return nil, validators.GetErrorResponseFromErr(errTransfer)
	}

	transferResponse := builder.TransferResponseBuilder(transfer)

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

	return transferResponse, validators.ErrorResponse{}
}
