package domain

import (
	"context"
	"gorm.io/gorm"
	"kaptan/internal/module/transfer/dto"
	"kaptan/internal/module/transfer/types"
	"kaptan/pkg/database/mysql/custom_types"
	"kaptan/pkg/validators"
	"time"
)

type Transfer struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Polymorphic relationship for client
	ClientID   uint   `gorm:"not null" json:"client_id"`
	ClientType string `gorm:"not null" json:"client_type"`

	// Polymorphic relationship for from (nullable)
	FromID   *uint   `json:"from_id,omitempty"`
	FromType *string `json:"from_type,omitempty"`

	// Polymorphic relationship for to (nullable)
	ToID   *uint   `json:"to_id,omitempty"`
	ToType *string `json:"to_type,omitempty"`

	// Address and location fields
	FromAddress string `gorm:"not null" json:"from_address"`
	FromLat     string `gorm:"not null" json:"from_lat"`
	FromLng     string `gorm:"not null" json:"from_lng"`
	ToAddress   string `gorm:"not null" json:"to_address"`
	ToLat       string `gorm:"not null" json:"to_lat"`
	ToLng       string `gorm:"not null" json:"to_lng"`

	// Brand relationship
	BrandID     uint                  `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"brand_id"`
	BrandObject *custom_types.JSONMap `gorm:"type:json;comment:Stores brand details as JSON" json:"brand_object,omitempty"`

	// Status fields
	Status       types.TransferStatus `gorm:"type:enum('pending','accepted','rejected');default:'pending'" json:"status"`
	HostStatus   types.HostStatus     `gorm:"type:enum('pending','accepted');default:'pending'" json:"host_status"`
	HideForHost  bool                 `gorm:"default:false" json:"hide_for_host"`
	GetMoneyFrom MoneySource          `gorm:"type:enum('client','guest');default:'client'" json:"get_money_from"`
	CashReceived float64              `gorm:"type:decimal(10,2);default:0" json:"cash_received"`
	TransType    types.TransferType   `gorm:"type:enum('normal','fast_carrier');default:'normal'" json:"trans_type"`

	// Car relationship
	CarID     *uint                 `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"car_id,omitempty"`
	CarObject *custom_types.JSONMap `gorm:"type:json;comment:Stores car details as JSON" json:"car_object,omitempty"`

	// Airport and type fields
	IsAirport bool               `gorm:"default:false" json:"is_airport"`
	Type      types.TransferType `gorm:"type:enum('arrival','departure');default:'arrival'" json:"type"`

	// Timestamps for transfer
	StartAt *time.Time `json:"start_at,omitempty"`
	EndAt   *time.Time `json:"end_at,omitempty"`

	// Foreign key relationships
	HostID    *uint `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"host_id,omitempty"`
	CompanyID *uint `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"company_id,omitempty"`
	DriverID  *uint `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"driver_id,omitempty"`
	AdminID   *uint `gorm:"comment:created by;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"admin_id,omitempty"`

	// Transfer details
	Number      string    `gorm:"not null;comment:room number or air number" json:"number"`
	NumberSeats string    `gorm:"not null;comment:number_seats" json:"number_seats"`
	Date        time.Time `gorm:"not null;comment:datetime" json:"date"`

	// Contact and passenger information
	Phone        string  `gorm:"not null" json:"phone"`
	Nationality  string  `gorm:"not null" json:"nationality"`
	Email        *string `json:"email,omitempty"`
	HasChildSeat bool    `gorm:"default:false" json:"has_child_seat"`
	BagsCount    int     `gorm:"default:0" json:"bags_count"`

	// Financial fields
	Price          float64        `gorm:"type:decimal(12,2);default:0" json:"price"`
	Distance       float64        `gorm:"type:decimal(12,2);default:0" json:"distance"`
	DriverCashPaid float64        `gorm:"type:decimal(10,2);default:0" json:"driver_cash_paid"`
	PaymentSource  *PaymentSource `gorm:"type:enum('client','hotel_owner')" json:"payment_source,omitempty"`

	// Notes
	Notes      *string `gorm:"type:text" json:"notes,omitempty"`
	GuestNotes *string `gorm:"type:text" json:"guest_notes,omitempty"`

	SellerID *uint `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"seller_id,omitempty"`
}

const (
	TransferStatusPending  types.TransferStatus = "pending"
	TransferStatusAccepted types.TransferStatus = "accepted"
	TransferStatusRejected types.TransferStatus = "rejected"
	TransferStatusStart    types.TransferStatus = "start"
	TransferStatusEnd      types.TransferStatus = "end"
)

type TransferHostStatus string

const (
	TransferHostStatusPending  TransferHostStatus = "pending"
	TransferHostStatusAccepted TransferHostStatus = "accepted"
)

type MoneySource string

const (
	MoneySourceClient MoneySource = "client"
	MoneySourceGuest  MoneySource = "guest"
)

type TransferType string

const (
	TransferTypeNormal      TransferType = "normal"
	TransferTypeFastCarrier TransferType = "fast_carrier"
)

type TransferDirection string

const (
	TransferDirectionArrival   TransferDirection = "arrival"
	TransferDirectionDeparture TransferDirection = "departure"
)

type PaymentSource string

const (
	PaymentSourceClient     PaymentSource = "client"
	PaymentSourceHotelOwner PaymentSource = "hotel_owner"
)

type TransferRepository interface {
	Find(ctx *context.Context, id uint) (domainData *Transfer, err error)
	AssignSellerToTransfer(ctx *context.Context, driverID uint, transferID uint) (*Transfer, error)
	MarkTransferAsStart(ctx *context.Context, transferDto *dto.StartTransfer) (*Transfer, error)
	MarkTransferAsEnd(ctx *context.Context, transferDto *dto.EndTransfer) (*Transfer, error)
}

type UseCase interface {
	StartTransfer(ctx *context.Context, dto *dto.StartTransfer) validators.ErrorResponse
	EndTransfer(ctx *context.Context, dto *dto.EndTransfer) validators.ErrorResponse
}
