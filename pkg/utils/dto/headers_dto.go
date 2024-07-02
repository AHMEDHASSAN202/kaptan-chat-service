package dto

type MobileHeaders struct {
	CountryId string `header:"Country-Id"`
	CauserId  string `header:"causer-id"`
	Lat       string `header:"Lat"`
	Lng       string `header:"Lng"`
}

type AdminHeaders struct {
	CountryId         string   `header:"Country-Id" validate:"required"`
	CauserId          string   `header:"causer-id"`
	CauserType        string   `header:"causer-type"`
	CauserDetails     string   `header:"causer-details"`
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
}
