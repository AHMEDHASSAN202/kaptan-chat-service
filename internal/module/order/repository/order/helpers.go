package mongodb

var (
	InProgressStatuses = []string{"initiated", "pending", "accepted"}
	CompletedStatuses  = []string{"timeout", "cancelled", "rejected", "ready_for_pickup", "pickedup", "no_show"}
)
