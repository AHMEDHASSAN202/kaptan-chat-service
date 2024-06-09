package dto

type SnoozeMenuWebhook struct {
	AccountID     string       `json:"accountId"`
	LocationID    string       `json:"locationId"`
	ChannelLinkID string       `json:"channelLinkId"`
	Operations    []Operations `json:"operations"`
}
type Items struct {
	Plu         string `json:"plu"`
	SnoozeStart string `json:"snoozeStart"`
	SnoozeEnd   string `json:"snoozeEnd"`
}
type Data struct {
	Items []Items `json:"cuisine"`
}
type Operations struct {
	Action string `json:"action"`
	Data   Data   `json:"data"`
}
