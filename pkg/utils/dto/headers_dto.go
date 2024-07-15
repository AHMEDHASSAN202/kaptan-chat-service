package dto

import (
	"encoding/json"
	"fmt"
)

type MobileHeaders struct {
	CountryId  string `header:"Country-Id"`
	CauserId   string `header:"causer-id"`
	CauserType string `header:"causer-type"`
	Lat        string `header:"Lat"`
	Lng        string `header:"Lng"`
}

type AdminHeaders struct {
	CountryId         string   `header:"Country-Id" validate:"required"`
	CauserId          string   `header:"causer-id"`
	CauserType        string   `header:"causer-type"`
	CauserName        string   `header:"causer-name"`
	CauserPermissions []string `header:"causer-permissions"`
}

type PortalHeaders struct {
	AccountId         string   `header:"Account-Id" validate:"required,mongodb"`
	CauserId          string   `header:"causer-id"`
	CauserType        string   `header:"causer-type"`
	CauserName        string   `header:"causer-name"`
	CauserAccountId   string   `header:"causer-account-id"`
	CauserPermissions []string `header:"causer-permissions"`
	CauserDetails     string   `header:"causer-details"`
}

type LocalizeText struct {
	Ar string `json:"ar"`
	En string `json:"en"`
}

type CauserDetails struct {
	Id   string       `json:"id"`
	Name LocalizeText `json:"name"`
}

func (p *PortalHeaders) GetCauserDetails() *CauserDetails {
	causerDetails := CauserDetails{}
	err := json.Unmarshal([]byte(p.CauserDetails), &causerDetails)
	if err != nil {
		fmt.Println("Error unmarshaling JSON: %v", err)
	}
	return &causerDetails
}
func (p *PortalHeaders) GetCauserDetailsAsMap() map[string]interface{} {
	causerDetails := map[string]interface{}{}
	err := json.Unmarshal([]byte(p.CauserDetails), &causerDetails)
	if err != nil {
		fmt.Println("Error unmarshaling JSON: %v", err)
	}
	return causerDetails
}
