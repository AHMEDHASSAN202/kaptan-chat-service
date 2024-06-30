package dto

type MobileHeaders struct {
	CountryId string `header:"Country-Id"`
	Lat       string `header:"Lat"`
	Lng       string `header:"Lng"`
}

type AdminHeaders struct {
	CountryId     string `header:"Country-Id" validate:"required"`
	CauserId      string `header:"causer-id"`
	CauserType    string `header:"causer-type"`
	CauserDetails string `header:"causer-details"`
}

type PortalHeaders struct {
	AccountId     string `header:"Account-Id" validate:"required,mongodb"`
	CauserId      string `header:"causer-id"`
	CauserType    string `header:"causer-type"`
	CauserDetails string `header:"causer-details"`
}

type CauserAdminDetails struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Type        string   `json:"type"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	CountryIds  []string `json:"country_ids"`
}
