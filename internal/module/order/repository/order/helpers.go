package mongodb

var (
	InProgressStatuses = []string{"initiated", "pending"}
	CompletedStatuses  = []string{"timeout", "cancelled", "rejected", "accepted", "ready_for_pickup", "pickedup", "no_show"}
)
