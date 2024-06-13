package dto

type MobileHeaders struct {
	CountryId string `header:"Country-Id"`
	Lat       string `header:"Lat"`
	Lng       string `header:"Lng"`
}
