package responses

import (
	"github.com/kamva/mgm/v3"
	"time"
)

type VerifyOtpResp struct {
	Token              string `json:"token"`
	IsProfileCompleted bool   `json:"is_profile_completed"`
}

type MobileListOrders struct {
	mgm.DefaultModel `bson:",inline"`
	SerialNum        string     `json:"serial_num"`
	TotalPrice       float64    `json:"total_price"`
	IsFavourite      bool       `json:"is_favourite"`
	PickedupAt       *time.Time `json:"pickedup_at"`
	User             struct {
		Id               string `json:"id"`
		Name             string `json:"name"`
		Phone            string `json:"phone"`
		Country          string `json:"country"`
		CollectionMethod struct {
		} `json:"collection_method"`
	} `json:"user"` // need to update
	Items []struct {
		ItemDetails struct {
			Id     string `json:"_id"`
			Name   string `json:"name"`
			Qty    string `json:"qty"`
			Price  string `json:"price"`
			Addons struct {
				Id    string `json:"_id"`
				Name  string `json:"name"`
				Qty   string `json:"qty"`
				Price string `json:"price"`
			} `json:"addons"`
			Category struct {
			} `json:"category"`
		} `json:"item_details"`
	} `json:"items"` // need to update
	LocationDetails struct {
	} `json:"location_details"` // need to update
}
