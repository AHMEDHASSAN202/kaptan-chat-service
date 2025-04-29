package dto

type MobileHeaders struct {
	CauserId   string `header:"causer-id"`
	CauserType string `header:"causer-type"`
	Lat        string `header:"Lat"`
	Lng        string `header:"Lng"`
}
