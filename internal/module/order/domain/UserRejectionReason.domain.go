package domain

type UserRejectionReason struct {
	Id   string `json:"id"`
	Name struct {
		Ar string `json:"ar"`
		En string `json:"en"`
	} `json:"name"`
	Status string `json:"status"`
}

type OrderStatusJson struct {
	AllowUserToChange    []string `json:"allow_user_to_change"`
	AllowAdminToChange   []string `json:"allow_admin_to_change"`
	AllowKitchenToChange []string `json:"allow_kitchen_to_change"`
	PreviousStatus       []string `json:"previous_status"`
}
