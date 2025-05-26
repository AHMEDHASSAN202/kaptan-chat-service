package app

import (
	"kaptan/internal/module/transfer/domain"
	"time"
)

type TransferResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`

	// Client information
	ClientID   uint   `json:"client_id"`
	ClientType string `json:"client_type"`

	// Location information
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`

	// Status information
	Status     domain.TransferStatus     `json:"status"`
	HostStatus domain.TransferHostStatus `json:"host_status"`
	Type       domain.TransferDirection  `json:"type"`
	TransType  domain.TransferType       `json:"trans_type"`

	// Transfer details
	Number      string     `json:"number"`
	NumberSeats string     `json:"number_seats"`
	Date        time.Time  `json:"date"`
	StartAt     *time.Time `json:"start_at,omitempty"`
	EndAt       *time.Time `json:"end_at,omitempty"`

	// Contact information
	Phone       string  `json:"phone"`
	Nationality string  `json:"nationality"`
	Email       *string `json:"email,omitempty"`

	// Basic transfer info
	IsAirport    bool    `json:"is_airport"`
	HasChildSeat bool    `json:"has_child_seat"`
	BagsCount    int     `json:"bags_count"`
	Price        float64 `json:"price"`
	Distance     float64 `json:"distance"`

	// Notes
	Notes      *string `json:"notes,omitempty"`
	GuestNotes *string `json:"guest_notes,omitempty"`
}
