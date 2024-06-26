package dto

type MobileHeaders struct {
	CountryId string `header:"Country-Id"`
	Lat       string `header:"Lat"`
	Lng       string `header:"Lng"`
}

type AdminHeaders struct {
	CountryId string `header:"Country-Id" validate:"required"`
}

type PortalHeaders struct {
	AccountId string `header:"Account-Id" validate:"required,mongodb"`
}
