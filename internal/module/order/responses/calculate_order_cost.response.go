package responses

type LocalizationText struct {
	Ar string `json:"ar"`
	En string `json:"en"`
}

type LocationDoc struct {
	Id   string           `json:"id"`
	Name LocalizationText `json:"name"`
	Desc LocalizationText `json:"desc"`
}

type PriceSummary struct {
	Qty                      int     `json:"qty"`
	UnitPrice                float64 `json:"unit_price"`
	TotalPriceBeforeDiscount float64 `json:"total_price_before_discount"`
	TotalPriceAfterDiscount  float64 `json:"total_price_after_discount"`
}

type TotalPriceSummary struct {
	Fees                     float64 `json:"unit_price"`
	TotalPriceBeforeDiscount float64 `json:"total_price_before_discount"`
	TotalPriceAfterDiscount  float64 `json:"total_price_after_discount"`
}

type MenuDoc struct {
	Id            string           `json:"id"`
	Name          LocalizationText `json:"name"`
	Desc          LocalizationText `json:"desc"`
	Image         string           `json:"image"`
	PriceSummary  PriceSummary     `json:"price_summary"`
	ModifierItems []MenuDoc        `json:"modifier_items"`
}

type CalculateOrderCostResp struct {
	Location          LocationDoc       `json:"location"`
	MenuItems         []MenuDoc         `json:"menu_items"`
	TotalPriceSummary TotalPriceSummary `json:"total_price_summary"`
}
