package builder

import (
	"kaptan/internal/module/transfer/domain"
	"kaptan/internal/module/transfer/responses/app"
)

func TransferResponseBuilder(chat *domain.Transfer) *app.TransferResponse {
	return &app.TransferResponse{
		ID:           chat.ID,
		CreatedAt:    chat.CreatedAt,
		UpdatedAt:    chat.UpdatedAt,
		ClientID:     chat.ClientID,
		ClientType:   chat.ClientType,
		FromAddress:  chat.FromAddress,
		ToAddress:    chat.ToAddress,
		Status:       chat.Status,
		HostStatus:   chat.HostStatus,
		Type:         chat.Type,
		TransType:    chat.TransType,
		Number:       chat.Number,
		NumberSeats:  chat.NumberSeats,
		Date:         chat.Date,
		StartAt:      chat.StartAt,
		EndAt:        chat.EndAt,
		Phone:        chat.Phone,
		Nationality:  chat.Nationality,
		Email:        chat.Email,
		IsAirport:    chat.IsAirport,
		HasChildSeat: chat.HasChildSeat,
		BagsCount:    chat.BagsCount,
		Price:        chat.Price,
		Distance:     chat.Distance,
		BrandID:      chat.BrandID,
		CarID:        chat.CarID,
		Notes:        chat.Notes,
		GuestNotes:   chat.GuestNotes,
	}
}
