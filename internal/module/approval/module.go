package approval

import (
	"go.uber.org/fx"
	"samm/internal/module/approval/delivery"
	approvalRepo "samm/internal/module/approval/repository/approval"
	"samm/internal/module/approval/usecase/approval"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		approvalRepo.NewApprovalRepository,
		approval.NewApprovalUseCase,
	),
	fx.Invoke(delivery.InitApprovalController),
)
