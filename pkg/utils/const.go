package utils

var Countries = []string{"SA", "EG"}

const ADMIN_TYPE = "admin"
const PORTAL_TYPE = "portal"
const KITCHEN_TYPE = "kitchen"

type APPROVAL_STATUS_STUCT struct {
	WAIT_FOR_APPROVAL string
	APPROVED          string
	REJECTED          string
}

var (
	APPROVAL_STATUS = APPROVAL_STATUS_STUCT{WAIT_FOR_APPROVAL: "wait_for_approval", APPROVED: "approved", REJECTED: "rejected"}
)
