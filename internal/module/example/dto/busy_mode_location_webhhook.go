package dto

type BusyModeLocationWebhook struct {
	AccountID     string `json:"accountId"`
	LocationID    string `json:"locationId"`
	ChannelLinkID string `json:"channelLinkId"`
	Status        string `json:"status"`
}
