package app_config

type FindMobileConfigResponse struct {
	ForceUpdate         bool   `json:"force_update"`
	Type                string `json:"type"`
	MinVersion          int64  `json:"min_version"`
	AppLink             string `json:"app_link"`
	LocalizationVersion int64  `json:"localization_version"`
	StartupImage        string `json:"stratup_image"`
}
