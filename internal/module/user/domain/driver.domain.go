package domain

import (
	"context"
	"fmt"
	"kaptan/pkg/utils"
	"time"
)

type Driver struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"column:name;not null" json:"name"`
	Phone     string    `gorm:"column:phone;not null;index:idx_drivers_phone" json:"phone"`
	Address   string    `gorm:"column:address" json:"address"`
	CreatedAt time.Time `gorm:"column:created_at;index:idx_drivers_created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Rating    float32   `gorm:"column:rating;index:idx_drivers_rating" json:"rating"`
	SoldTrips int       `gorm:"column:sold_trips;index:idx_drivers_sold_trips" json:"sold_trips"`

	// Media relationship
	Media []Media `gorm:"foreignKey:ModelID;references:ID" json:"media,omitempty"`
}

// Media represents the media table (similar to Spatie Media Library)
type Media struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	ModelType       string    `gorm:"column:model_type;not null" json:"model_type"`
	ModelID         uint      `gorm:"column:model_id;not null" json:"model_id"`
	UUID            string    `gorm:"column:uuid;uniqueIndex" json:"uuid"`
	CollectionName  string    `gorm:"column:collection_name;not null" json:"collection_name"`
	Name            string    `gorm:"column:name;not null" json:"name"`
	FileName        string    `gorm:"column:file_name;not null" json:"file_name"`
	MimeType        string    `gorm:"column:mime_type" json:"mime_type"`
	Disk            string    `gorm:"column:disk;not null" json:"disk"`
	ConversionsDisk string    `gorm:"column:conversions_disk" json:"conversions_disk"`
	Size            uint64    `gorm:"column:size;not null" json:"size"`
	OrderColumn     int       `gorm:"column:order_column" json:"order_column"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TableName specifies the table name for Media
func (Media) TableName() string {
	return "media"
}

// GetFullUrl returns the full URL of the media file
func (m *Media) GetFullUrl() string {
	return utils.Assets("/storage/" + fmt.Sprintf("%d", m.ID) + "/" + m.FileName)
}

type DriverResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	Rating    float32   `json:"rating"`
	SoldTrips int       `json:"sold_trips"`
	Avatar    *string   `json:"avatar"`
}

// ToResponse converts Driver to DriverResponse with media
func (d *Driver) ToResponse() *DriverResponse {
	response := &DriverResponse{
		ID:        d.ID,
		Name:      d.Name,
		Phone:     d.Phone,
		Address:   d.Address,
		CreatedAt: d.CreatedAt,
		Rating:    d.Rating,
		SoldTrips: d.SoldTrips,
	}

	// Process media files
	var images []string
	for _, media := range d.Media {
		url := media.GetFullUrl()
		images = append(images, url)

		// Set first image as avatar if collection is 'avatar' or first image
		if response.Avatar == nil && (media.CollectionName == "avatar" || media.CollectionName == "images") {
			response.Avatar = &url
		}
	}

	return response
}

type DriverRepository interface {
	Find(ctx *context.Context, id uint) (domainData *Driver, err error)
	FindWithMedia(ctx *context.Context, id uint) (domainData *Driver, err error)
	FindByAccessTokenId(ctx *context.Context, id uint) (*Driver, error)
	IncrementSoldTripsByValue(ctx *context.Context, id uint, value int) error
}
