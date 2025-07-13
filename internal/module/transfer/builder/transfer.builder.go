package builder

import (
	"fmt"
	"kaptan/internal/module/transfer/domain"
	"kaptan/internal/module/transfer/helper"
	"kaptan/internal/module/transfer/responses/app"
	"kaptan/pkg/utils"
)

func TransferResponseBuilder(transfer *domain.Transfer) *app.TransferResponse {
	if transfer == nil {
		return nil
	}

	clientId, clientType := helper.GenerateClientId(fmt.Sprintf("%d", transfer.ClientID), transfer.ClientType)

	transferResponse := &app.TransferResponse{
		ID:           transfer.ID,
		CreatedAt:    transfer.CreatedAt,
		ClientID:     clientId,
		ClientType:   clientType,
		FromAddress:  transfer.FromAddress,
		ToAddress:    transfer.ToAddress,
		Status:       transfer.Status,
		HostStatus:   transfer.HostStatus,
		Type:         transfer.Type,
		TransType:    transfer.TransType,
		Number:       transfer.Number,
		NumberSeats:  transfer.NumberSeats,
		Date:         transfer.Date,
		StartAt:      transfer.StartAt,
		EndAt:        transfer.EndAt,
		Phone:        transfer.Phone,
		Nationality:  transfer.Nationality,
		Email:        transfer.Email,
		IsAirport:    transfer.IsAirport,
		HasChildSeat: transfer.HasChildSeat,
		BagsCount:    transfer.BagsCount,
		Price:        transfer.Price,
		Distance:     transfer.Distance,
		Notes:        transfer.Notes,
		GuestNotes:   transfer.GuestNotes,
		CarObject:    transfer.CarObject,
	}

	_ = utils.TryCatch(func() {
		car := *transfer.CarObject
		if car != nil {
			car["has_children_seat"] = utils.IsTrue(car["has_children_seat"])
			transferResponse.CarObject = &car
		}
	})()

	return transferResponse
}
