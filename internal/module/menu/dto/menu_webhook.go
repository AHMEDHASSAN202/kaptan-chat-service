package dto

type MenuWebhook struct {
	Menu            string                    `json:"menu"`
	MenuId          string                    `json:"menuId"`
	Availabilities  []Availability            `json:"availabilities"`
	Categories      []Category                `json:"categories"`
	Products        map[string]Product        `json:"products"`
	Bundles         map[string]Bundle         `json:"bundles"`
	ModifierGroups  map[string]ModifierGroup  `json:"modifierGroups"`
	Modifiers       map[string]Modifier       `json:"modifiers"`
	SnoozedProducts map[string]SnoozedProduct `json:"snoozedProducts"`
	ChannelLinkID   string                    `json:"channelLinkId"`
}

type MenusWebhook []MenuWebhook
type Availability struct {
	DayOfWeek int    `json:"dayOfWeek"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type Category struct {
	ID                      string            `json:"_id"`
	Account                 string            `json:"account"`
	Name                    string            `json:"name"`
	NameTranslations        map[string]string `json:"nameTranslations"`
	Description             string            `json:"description"`
	DescriptionTranslations map[string]string `json:"descriptionTranslations"`
	Availabilities          []Availability    `json:"availabilities"` //todo: check this
	ImageURL                string            `json:"imageUrl"`
	Products                []string          `json:"products"`
	Menu                    string            `json:"menu"`
	Level                   int               `json:"level,omitempty"`
}
type Product struct {
	ID                      string            `json:"_id"`
	Account                 string            `json:"account"`
	Location                string            `json:"location"`
	ProductType             int               `json:"productType"`
	Plu                     string            `json:"plu"`
	Price                   float64           `json:"price"`
	Name                    string            `json:"name"`
	NameTranslations        map[string]string `json:"nameTranslations"`
	Calories                float64           `json:"calories"`
	DeliveryTax             float64           `json:"deliveryTax"`
	SubProducts             []string          `json:"subProducts"`
	ImageURL                string            `json:"imageUrl"`
	UniqueKey               string            `json:"uniqueKey"`
	Description             string            `json:"description"`
	DescriptionTranslations map[string]string `json:"descriptionTranslations"`
	Max                     int               `json:"max"`
	Min                     int               `json:"min"`
	Channel                 int               `json:"channel"`
}
type Bundle struct {
	ID                      string            `json:"_id"`
	Name                    string            `json:"name"`
	NameTranslations        map[string]string `json:"nameTranslations"`
	DescriptionTranslations map[string]string `json:"descriptionTranslations"`
	Description             string            `json:"description"`
	Account                 string            `json:"account"`
	CapacityUsages          []any             `json:"capacityUsages"`
	DeliveryTax             int               `json:"deliveryTax"`
	EatInTax                int               `json:"eatInTax"`
	Location                string            `json:"location"`
	Max                     int               `json:"max"`
	Min                     int               `json:"min"`
	Multiply                int               `json:"multiply"`
	Plu                     string            `json:"plu"`
	PosCategoryIds          []any             `json:"posCategoryIds"`
	PosProductCategoryID    string            `json:"posProductCategoryId"`
	PosProductID            string            `json:"posProductId"`
	ProductTags             []any             `json:"productTags"`
	ProductType             int               `json:"productType"`
	SubProducts             []string          `json:"subProducts"`
	TakeawayTax             int               `json:"takeawayTax"`
	ParentID                string            `json:"parentId"`
	Snoozed                 bool              `json:"snoozed"`
	SubProductSortOrder     []any             `json:"subProductSortOrder"`
}
type ModifierGroup struct {
	ID               string            `json:"_id"`
	Account          string            `json:"account"`
	Location         string            `json:"location"`
	ProductType      int               `json:"productType"`
	Plu              string            `json:"plu"`
	Price            float64           `json:"price"`
	Name             string            `json:"name"`
	NameTranslations map[string]string `json:"nameTranslations"`
	DeliveryTax      float64           `json:"deliveryTax"`
	SubProducts      []string          `json:"subProducts"`
	ImageURL         string            `json:"imageUrl"`
	Description      string            `json:"description"`
	Max              int               `json:"max"`
	Min              int               `json:"min"`
	MultiMax         int               `json:"multiMax"`
	Channel          int               `json:"channel"`
}

type Modifier struct {
	ID                      string            `json:"_id"`
	Account                 string            `json:"account"`
	Location                string            `json:"location"`
	ProductType             int               `json:"productType"`
	Plu                     string            `json:"plu"`
	Price                   float64           `json:"price"`
	Name                    string            `json:"name"`
	NameTranslations        map[string]string `json:"nameTranslations"`
	DescriptionTranslations map[string]string `json:"descriptionTranslations"`
	Calories                float64           `json:"calories"`
	DeliveryTax             float64           `json:"deliveryTax"`
	SubProducts             []string          `json:"subProducts"`
	ImageURL                string            `json:"imageUrl"`
	Description             string            `json:"description"`
	Max                     int               `json:"max"`
	Min                     int               `json:"min"`
	Channel                 int               `json:"channel"`
}

type SnoozedProduct struct {
	Location    string `json:"location"`
	Plu         string `json:"plu"`
	Name        string `json:"name"`
	SnoozeStart string `json:"snoozeStart"`
	SnoozeEnd   string `json:"snoozeEnd"`
}

type MailData struct {
	VendorId   string
	Operation  string
	VendorName string
	BranchName string
	Status     string
	Reason     string
	Time       string
}

type BrandData struct {
	VendorId   string
	VendorName string
	BranchName string
}
