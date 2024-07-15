package consts

type OrderStates struct {
	Initiated      string
	Pending        string
	TimedOut       string
	Accepted       string
	ReadyForPickup string
	PickedUp       string
	NoShow         string
	Cancelled      string
	Rejected       string
}

var (
	OrderStatus  = OrderStates{Initiated: "initiated", Pending: "pending", TimedOut: "timedOut", Accepted: "accepted", Cancelled: "cancelled", Rejected: "rejected", ReadyForPickup: "ready_for_pickup", PickedUp: "pickedup", NoShow: "no_show"}
	ActorAdmin   = "admin"
	ActorUser    = "user"
	ActorKitchen = "kitchen"
)
