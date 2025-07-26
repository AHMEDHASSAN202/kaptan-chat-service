package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"kaptan/internal/module/transfer/domain"
	"kaptan/internal/module/transfer/dto"
	domain2 "kaptan/internal/module/user/domain"
	"kaptan/pkg/database/mysql/custom_types"
	"kaptan/pkg/logger"
	"time"
)

type Repository struct {
	logger           logger.ILogger
	db               *gorm.DB
	driverRepository domain2.DriverRepository
}

func NewTransferRepository(log logger.ILogger, db *gorm.DB, driverRepository domain2.DriverRepository) domain.TransferRepository {
	return &Repository{
		logger:           log,
		db:               db,
		driverRepository: driverRepository,
	}
}

func (r *Repository) Find(ctx *context.Context, id uint) (*domain.Transfer, error) {
	transfer := domain.Transfer{ID: id}
	result := r.db.First(&transfer)
	if result.Error != nil {
		return nil, result.Error
	}
	return &transfer, nil
}

func (r *Repository) AssignSellerToTransfer(ctx *context.Context, driverID uint, transferID uint) (*domain.Transfer, error) {
	// Step 1: Find the transfer
	transfer, err := r.Find(ctx, transferID)
	if err != nil {
		fmt.Println("Error finding transfer", "transferID", transferID, "error", err)
		return nil, err
	}

	var carId *uint // assuming CarID is a uint in the domain
	driverResult, err := r.driverRepository.Find(ctx, driverID)
	if err != nil {
		r.logger.Error("Failed to find driver", "driverID", driverID, "error", err)
	} else {
		carId = driverResult.VehicleId
	}

	type Vehicle struct {
		ID                 uint       `json:"id"`
		Sort               int        `json:"sort"`
		Color              string     `json:"color"`
		RefID              *uint      `json:"ref_id"`
		Status             string     `json:"status"`
		MaxBags            int        `json:"max_bags"`
		IsActive           bool       `json:"is_active"`
		MaxSeats           int        `json:"max_seats"`
		CompanyID          *uint      `json:"company_id"`
		CreatedAt          time.Time  `json:"created_at"`
		DeletedAt          *time.Time `json:"deleted_at"`
		UpdatedAt          time.Time  `json:"updated_at"`
		CarBrandID         uint       `json:"car_brand_id"`
		ModelNumber        string     `json:"model_number"`
		PlateNumber        string     `json:"plate_number"`
		LicencePlate       string     `json:"licence_plate"`
		ManufactureYear    int        `json:"manufacture_year"`
		HasChildrenSeat    int        `json:"has_children_seat"`
		LicenceExpiredDate string     `json:"licence_expired_date"`
	}

	var car *custom_types.JSONMap
	if carId != nil {
		var vehicle Vehicle
		if err := r.db.Table("vehicles").Where("id = ?", *carId).First(&vehicle).Error; err != nil {
			r.logger.Error("Failed to get vehicle ", "carID ", *carId, " error ", err)
		} else {
			vehicleMap := make(custom_types.JSONMap)
			b, _ := json.Marshal(vehicle)
			_ = json.Unmarshal(b, &vehicleMap)
			car = &vehicleMap
		}
	}

	// Step 2: Assign driverID as sellerID
	transfer.SellerID = transfer.DriverID // assuming SellerID is a *uint in the domain
	transfer.DriverID = &driverID         // assuming SellerID is a *uint in the domain
	transfer.CarID = carId                // assuming CarID is a *uint in the domain
	transfer.CarObject = car              // assuming transfer.Car is of type map[string]interface{}

	// Step 3: Save the updated transfer
	if err := r.db.Save(&transfer).Error; err != nil {
		return nil, err
	}

	// Step 4: Return updated transfer
	return transfer, nil
}

func (r *Repository) MarkTransferAsStart(ctx *context.Context, transferDto *dto.StartTransfer) (*domain.Transfer, error) {
	// Step 1: Find the transfer
	transfer, err := r.Find(ctx, transferDto.TransferId)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	// Step 2: Change status to "Start"
	transfer.Status = domain.TransferStatusStart // assuming Status is a string
	transfer.StartAt = &now

	// Step 3: Save the updated transfer
	if err := r.db.Save(&transfer).Error; err != nil {
		return nil, err
	}

	// Step 4: Return updated transfer
	return transfer, nil
}

func (r *Repository) MarkTransferAsEnd(ctx *context.Context, transferDto *dto.EndTransfer) (*domain.Transfer, error) {
	// Step 1: Find the transfer
	transfer, err := r.Find(ctx, transferDto.TransferId)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	// Step 2: Change status to "End"
	transfer.Status = domain.TransferStatusEnd // assuming Status is a string
	transfer.EndAt = &now

	// Step 3: Save the updated transfer
	if err := r.db.Save(&transfer).Error; err != nil {
		return nil, err
	}

	// Step 4: Return updated transfer
	return transfer, nil
}
