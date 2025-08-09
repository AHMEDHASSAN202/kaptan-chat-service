package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"kaptan/internal/module/chat/builder"
	"kaptan/internal/module/chat/consts"
	"kaptan/internal/module/chat/domain"
	"kaptan/internal/module/chat/dto"
	"kaptan/internal/module/chat/responses/app"
	builder2 "kaptan/internal/module/transfer/builder"
	domain3 "kaptan/internal/module/transfer/domain"
	domain2 "kaptan/internal/module/user/domain"
	"kaptan/internal/module/user/usecase/notification"
	"kaptan/pkg/fcm_notification"
	"kaptan/pkg/gate"
	"kaptan/pkg/logger"
	"kaptan/pkg/utils"
	"kaptan/pkg/validators"
	"kaptan/pkg/websocket"
	"sync"
)

type ChatUseCase struct {
	repo             domain.ChatRepository
	logger           logger.ILogger
	gate             *gate.Gate
	websocketManager *websocket.ChannelManager
	driverRepo       domain2.DriverRepository
	transferRepo     domain3.TransferRepository
	fcmClient        *fcm_notification.FCMClient
	notification     *notification.UseCase
}

func NewChatUseCase(repo domain.ChatRepository, driverRepo domain2.DriverRepository, transferRepo domain3.TransferRepository, gate *gate.Gate, websocketManager *websocket.ChannelManager, logger logger.ILogger, fcmClient *fcm_notification.FCMClient, notification *notification.UseCase) domain.ChatUseCase {
	return &ChatUseCase{
		repo:             repo,
		logger:           logger,
		gate:             gate,
		driverRepo:       driverRepo,
		websocketManager: websocketManager,
		transferRepo:     transferRepo,
		fcmClient:        fcmClient,
		notification:     notification,
	}
}

func (u ChatUseCase) GetChats(ctx context.Context, dto *dto.GetChats) (app.ListChatResponse, validators.ErrorResponse) {
	privateChats := u.repo.PrivateChats(ctx, dto)
	chatsResponse := builder.ChatsResponseBuilder(privateChats)
	return chatsResponse, validators.ErrorResponse{}
}

func (u ChatUseCase) UnreadMessages(ctx context.Context, dto *dto.UnreadMessages) (int, validators.ErrorResponse) {
	return u.repo.UnreadMessages(ctx, dto)
}

func (u ChatUseCase) GetChatMessages(ctx context.Context, dto *dto.GetChatMessage) (*app.MessagesResponse, validators.ErrorResponse) {
	messages, pagination := u.repo.GetChatMessages(ctx, dto)
	messagesResponse := builder.MessagesResponseBuilder(messages, pagination)
	return messagesResponse, validators.ErrorResponse{}
}

func (u ChatUseCase) AddPrivateChat(ctx context.Context, dto *dto.AddPrivateChat) (*app.ChatResponse, validators.ErrorResponse) {
	chat, message, err := u.repo.AddPrivateChat(ctx, dto)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}

	driver, err := u.driverRepo.Find(&ctx, uint(chat.UserId))
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}

	chatResponse := builder.ChatResponseBuilder(chat, driver)

	go func() {
		contentJson, _ := json.Marshal(chatResponse)
		myClient := u.websocketManager.GetClient(utils.GetClientUserId(dto.CauserType, dto.CauserId))
		anotherClient := u.websocketManager.GetClient(utils.GetClientUserId(message.SenderType, fmt.Sprintf("%v", message.SenderId)))
		u.websocketManager.JoinChannel(myClient, chatResponse.Channel)
		u.websocketManager.JoinChannel(anotherClient, chatResponse.Channel)
		u.logger.Info("Joining channel | ", "channel | ", chatResponse.Channel, " | userId", utils.GetClientUserId(dto.CauserType, dto.CauserId), " | userId2", utils.GetClientUserId(message.SenderType, fmt.Sprintf("%v", message.SenderId)))
		u.logger.Info("Broadcasting chat response | ", "myClient | ", myClient, " | anotherClient", anotherClient)
		u.websocketManager.Broadcast <- websocket.Message{
			ChannelID: chatResponse.Channel,
			Content:   string(contentJson),
			Action:    consts.START_CHAT_ACTION,
		}
		u.addUnreadMessage(myClient, chatResponse.Channel)
	}()

	go func() {
		messageTitle := "New Private Chat"
		messageBody := fmt.Sprintf("New private chat with %s", message.User["name"].(string))
		contentJson, _ := json.Marshal(chatResponse)
		chatJson, _ := json.Marshal(websocket.Message{
			ChannelID: chatResponse.Channel,
			Content:   string(contentJson),
			Action:    consts.START_CHAT_ACTION,
		})
		anotherSideChat, errAnotherSideChat := u.repo.GetAnotherSideChat(ctx, chatResponse.Channel, cast.ToInt(dto.CauserId))
		if errAnotherSideChat != nil {
			u.logger.Info("AnotherSideChat Error => ", err)
			return
		}
		_ = u.notification.SendNotificationToUsers(ctx, []uint{cast.ToUint(anotherSideChat.UserId)}, messageTitle, messageBody, map[string]string{
			"chat": string(chatJson),
		})
	}()

	return chatResponse, validators.ErrorResponse{}
}

func (u ChatUseCase) SaleTransferChat(ctx context.Context, dto *dto.SaleTransferChat) (*app.ChatResponse, validators.ErrorResponse) {
	chat, err := u.repo.SaleTransferChat(ctx, dto)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}

	chatResponse := builder.ChatResponseBuilder(chat, nil)

	go func() {
		contentJson, _ := json.Marshal(chatResponse)
		u.websocketManager.Broadcast <- websocket.Message{
			ChannelID: chatResponse.Channel,
			Content:   string(contentJson),
			Action:    consts.SALE_CHAT_ACTION,
		}
		u.addUnreadMessage(u.websocketManager.GetClient(utils.GetClientUserId(dto.CauserType, dto.CauserId)), chatResponse.Channel)
	}()

	go func() {
		contentJson, _ := json.Marshal(chatResponse)
		u.websocketManager.Broadcast <- websocket.Message{
			ChannelID: consts.GENERAL_CHAT,
			Content:   string(contentJson),
			Action:    consts.CLOSED_CHAT_ACTION,
		}
		u.addUnreadMessage(nil, consts.GENERAL_CHAT)
	}()

	go func() {
		soldMap := map[string]interface{}{
			"channel":    chatResponse.Channel,
			"message_id": chat.OpenedBy,
		}
		contentJson, _ := json.Marshal(soldMap)
		u.websocketManager.Broadcast <- websocket.Message{
			ChannelID: consts.GENERAL_CHAT,
			Content:   string(contentJson),
			Action:    consts.SOLD_MESSAGE_ACTION,
		}
	}()

	go func() {
		messageTitle := "Sold Transfer"
		messageBody := fmt.Sprintf("New sold transfer sold with %s", chatResponse.User["name"].(string))
		soldMap := map[string]interface{}{
			"channel":    chatResponse.Channel,
			"message_id": chat.OpenedBy,
		}
		chatContentJson, _ := json.Marshal(chatResponse)
		contentJson, _ := json.Marshal(soldMap)
		chatJson, _ := json.Marshal(websocket.Message{
			ChannelID: chatResponse.Channel,
			Content:   string(chatContentJson),
			Action:    consts.SOLD_MESSAGE_ACTION,
		})
		anotherSideChat, errAnotherSideChat := u.repo.GetAnotherSideChat(ctx, chatResponse.Channel, cast.ToInt(dto.CauserId))
		if errAnotherSideChat != nil {
			u.logger.Info("AnotherSideChat Error => ", err)
			return
		}
		_ = u.notification.SendNotificationToUsers(ctx, []uint{cast.ToUint(anotherSideChat.UserId)}, messageTitle, messageBody, map[string]string{
			"message": string(contentJson),
			"chat":    string(chatJson),
		})
	}()

	return chatResponse, validators.ErrorResponse{}
}

func (u ChatUseCase) RejectOffer(ctx context.Context, dto *dto.RejectOffer) (*app.MessageResponse, validators.ErrorResponse) {
	message, err := u.repo.RejectOffer(ctx, dto)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}
	messageResponse := builder.MessageResponseBuilder(message)

	go func() {
		contentJson, _ := json.Marshal(messageResponse)
		u.websocketManager.Broadcast <- websocket.Message{
			ChannelID: messageResponse.Channel,
			Content:   string(contentJson),
			Action:    consts.REJECT_OFFER_ACTION,
		}
		myClient := u.websocketManager.GetClient(utils.GetClientUserId(dto.CauserType, dto.CauserId))
		u.addUnreadMessage(myClient, messageResponse.Channel)
	}()

	go func() {
		messageTitle := "Rejected Offer"
		messageBody := fmt.Sprintf("Offer rejected by %s", message.User["name"].(string))
		contentJson, _ := json.Marshal(messageResponse)
		messageJson, _ := json.Marshal(websocket.Message{
			ChannelID: messageResponse.Channel,
			Content:   string(contentJson),
			Action:    consts.REJECT_OFFER_ACTION,
		})
		anotherSideChat, errAnotherSideChat := u.repo.GetAnotherSideChat(ctx, message.Channel, cast.ToInt(dto.CauserId))
		if errAnotherSideChat != nil {
			u.logger.Info("AnotherSideChat Error => ", err)
			return
		}
		_ = u.notification.SendNotificationToUsers(ctx, []uint{cast.ToUint(anotherSideChat.UserId)}, messageTitle, messageBody, map[string]string{
			"message": string(messageJson),
		})
	}()

	return messageResponse, validators.ErrorResponse{}
}

func (u ChatUseCase) GetChat(ctx context.Context, dto *dto.GetChat) (*app.ChatResponse, interface{}, validators.ErrorResponse) {
	chat, err := u.repo.GetChat(ctx, dto)
	if err != nil {
		return nil, nil, validators.GetErrorResponseFromErr(err)
	}

	var (
		driver   *domain2.Driver
		transfer *domain3.Transfer
	)

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	// Get driver
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	var dErr error
	//	driver, dErr = u.driverRepo.Find(&ctx, uint(chat.UserId))
	//	if dErr != nil {
	//		errChan <- dErr
	//	}
	//}()

	// Get transfer if needed
	if chat.TransferId != nil {
		wg.Add(1)
		go func(transferID uint) {
			defer wg.Done()
			var tErr error
			transfer, tErr = u.transferRepo.Find(&ctx, transferID)
			if tErr != nil {
				//errChan <- tErr
				u.logger.Info("Transfer Error => ", tErr)
			}
		}(uint(*chat.TransferId))
	}

	// Wait for both goroutines to finish
	wg.Wait()
	close(errChan)

	// Check for errors
	for e := range errChan {
		return nil, nil, validators.GetErrorResponseFromErr(e)
	}

	return builder.ChatResponseBuilder(chat, driver), builder2.TransferResponseBuilder(transfer), validators.ErrorResponse{}
}

func (u ChatUseCase) GetAcceptedChatByTransferId(ctx context.Context, transferId uint, userId string) (*app.ChatResponse, validators.ErrorResponse) {
	chat, err := u.repo.GetAcceptedChatByTransferId(ctx, transferId, userId)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}

	return builder.ChatResponseBuilder(chat, nil), validators.ErrorResponse{}
}

func (u ChatUseCase) SendMessage(ctx context.Context, dto *dto.SendMessage) (*app.MessageResponse, validators.ErrorResponse) {
	message, err := u.repo.StoreMessage(ctx, dto)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}

	chat, _ := u.repo.GetChatOnly(ctx, dto.Channel, cast.ToInt(dto.CauserId))
	message.Chat = chat

	messageResponse := builder.MessageResponseBuilder(message)

	go func() {
		contentJson, _ := json.Marshal(messageResponse)
		u.websocketManager.Broadcast <- websocket.Message{
			ChannelID: messageResponse.Channel,
			Content:   string(contentJson),
			Action:    consts.ADD_MESSAGE_ACTION,
		}
		myClient := u.websocketManager.GetClient(utils.GetClientUserId(dto.CauserType, dto.CauserId))
		u.addUnreadMessage(myClient, messageResponse.Channel)
	}()

	if dto.BrandId != nil {
		go func() {
			contentJson, _ := json.Marshal(messageResponse)
			u.websocketManager.Broadcast <- websocket.Message{
				ChannelID: consts.GENERAL_CHAT,
				Content:   string(contentJson),
				Action:    consts.ADD_MESSAGE_ACTION,
			}
			u.addUnreadMessage(nil, consts.GENERAL_CHAT)
		}()
	}

	go func() {
		messageTitle := message.User["name"].(string) + " | New Message"
		contentJson, _ := json.Marshal(messageResponse)
		messageJson, _ := json.Marshal(websocket.Message{
			ChannelID: consts.GENERAL_CHAT,
			Content:   string(contentJson),
			Action:    consts.ADD_MESSAGE_ACTION,
		})
		if message.IsPrivate {
			anotherSideChat, errAnotherSideChat := u.repo.GetAnotherSideChat(ctx, messageResponse.Channel, cast.ToInt(dto.CauserId))
			if errAnotherSideChat != nil {
				u.logger.Info("AnotherSideChat Error => ", err)
				return
			}
			_ = u.notification.SendNotificationToUsers(ctx, []uint{cast.ToUint(anotherSideChat.UserId)}, messageTitle, dto.Message, map[string]string{
				"message": string(messageJson),
			})
		} else {
			condition := fmt.Sprintf("'drivers' in topics && !'user_%s' in topics", dto.CauserId)
			_ = u.notification.SendNotificationToTopic(ctx, "drivers", messageTitle, dto.Message, map[string]string{
				"message": string(contentJson),
			}, condition)
		}
	}()

	return messageResponse, validators.ErrorResponse{}
}

func (u ChatUseCase) UpdateMessage(ctx context.Context, dto *dto.UpdateMessage) (*app.MessageResponse, validators.ErrorResponse) {
	message, err := u.repo.UpdateMessage(ctx, dto)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}

	messageResponse := builder.MessageResponseBuilder(message)

	go func() {
		contentJson, _ := json.Marshal(messageResponse)
		u.websocketManager.Broadcast <- websocket.Message{
			ChannelID: messageResponse.Channel,
			Content:   string(contentJson),
			Action:    consts.UPDATE_MESSAGE_ACTION,
		}
		//u.addUnreadMessage(u.websocketManager.GetClient(utils.GetClientUserId(dto.CauserType, dto.CauserId)), messageResponse.Channel)
	}()

	if message.BrandId != nil {
		go func() {
			contentJson, _ := json.Marshal(messageResponse)
			u.websocketManager.Broadcast <- websocket.Message{
				ChannelID: consts.GENERAL_CHAT,
				Content:   string(contentJson),
				Action:    consts.UPDATE_MESSAGE_ACTION,
			}
		}()
	}

	go func() {
		messageTitle := message.User["name"].(string) + " | Update Message"
		contentJson, _ := json.Marshal(messageResponse)
		messageJson, _ := json.Marshal(websocket.Message{
			ChannelID: consts.GENERAL_CHAT,
			Content:   string(contentJson),
			Action:    consts.UPDATE_MESSAGE_ACTION,
		})
		if message.IsPrivate {
			anotherSideChat, errAnotherSideChat := u.repo.GetAnotherSideChat(ctx, messageResponse.Channel, cast.ToInt(dto.CauserId))
			if errAnotherSideChat != nil {
				u.logger.Info("AnotherSideChat Error => ", err)
				return
			}
			_ = u.notification.SendNotificationToUsers(ctx, []uint{cast.ToUint(anotherSideChat.UserId)}, messageTitle, dto.Message, map[string]string{
				"message": string(messageJson),
			})
		} else {
			condition := fmt.Sprintf("'drivers' in topics && !'user_%s' in topics", dto.CauserId)
			_ = u.notification.SendNotificationToTopic(ctx, "drivers", messageTitle, dto.Message, map[string]string{
				"message": string(contentJson),
			}, condition)
		}
	}()

	return messageResponse, validators.ErrorResponse{}
}

func (u ChatUseCase) DeleteMessage(ctx context.Context, dto *dto.DeleteMessage) (*app.MessageResponse, validators.ErrorResponse) {
	message, err := u.repo.DeleteMessage(ctx, dto)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}

	messageResponse := builder.MessageResponseBuilder(message)

	go func() {
		contentJson, _ := json.Marshal(messageResponse)
		u.websocketManager.Broadcast <- websocket.Message{
			ChannelID: messageResponse.Channel,
			Content:   string(contentJson),
			Action:    consts.DELETE_MESSAGE_ACTION,
		}
		//u.addUnreadMessage(u.websocketManager.GetClient(utils.GetClientUserId(dto.CauserType, dto.CauserId)), messageResponse.Channel)
	}()

	if message.BrandId != nil {
		go func() {
			contentJson, _ := json.Marshal(messageResponse)
			u.websocketManager.Broadcast <- websocket.Message{
				ChannelID: consts.GENERAL_CHAT,
				Content:   string(contentJson),
				Action:    consts.DELETE_MESSAGE_ACTION,
			}
		}()
	}

	go func() {
		messageTitle := message.User["name"].(string) + " | Delete Message"
		contentJson, _ := json.Marshal(messageResponse)
		messageJson, _ := json.Marshal(websocket.Message{
			ChannelID: consts.GENERAL_CHAT,
			Content:   string(contentJson),
			Action:    consts.DELETE_MESSAGE_ACTION,
		})
		if message.IsPrivate {
			anotherSideChat, errAnotherSideChat := u.repo.GetAnotherSideChat(ctx, messageResponse.Channel, cast.ToInt(dto.CauserId))
			if errAnotherSideChat != nil {
				u.logger.Info("AnotherSideChat Error => ", err)
				return
			}
			_ = u.notification.SendNotificationToUsers(ctx, []uint{cast.ToUint(anotherSideChat.UserId)}, messageTitle, "---------", map[string]string{
				"message": string(messageJson),
			})
		} else {
			condition := fmt.Sprintf("'drivers' in topics && !'user_%s' in topics", dto.CauserId)
			_ = u.notification.SendNotificationToTopic(ctx, "drivers", messageTitle, "---------", map[string]string{
				"message": string(contentJson),
			}, condition)
		}
	}()

	return messageResponse, validators.ErrorResponse{}
}

func (u ChatUseCase) addUnreadMessage(exceptClient *websocket.Client, channel string) {
	unreadJson, _ := json.Marshal(map[string]string{"channel": channel})
	u.websocketManager.Broadcast <- websocket.Message{
		ChannelID:    channel,
		Content:      string(unreadJson),
		Action:       consts.UNREAD_MESSAGE_ACTION,
		ExceptClient: exceptClient,
	}
}
