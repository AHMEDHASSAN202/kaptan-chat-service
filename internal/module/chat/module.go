package chat

import (
	"go.uber.org/fx"
	"kaptan/internal/module/chat/delivery"
	"kaptan/internal/module/chat/migrations"
	"kaptan/internal/module/chat/policies"
	chat_repo "kaptan/internal/module/chat/repository/chat"
	"kaptan/internal/module/chat/subscribers"
	chat_usecase "kaptan/internal/module/chat/usecase/chat"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		chat_repo.NewChatRepository,
		chat_usecase.NewChatUseCase,
	),
	fx.Invoke(
		policies.NewIPolicy, delivery.InitChatController,
		subscribers.Message,
	),
	migrations.Module,
)
