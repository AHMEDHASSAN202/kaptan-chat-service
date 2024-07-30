package approval_helpers

import (
	"samm/internal/module/menu/domain"
	"samm/pkg/utils"
	"strings"
)

func NeedToApproveItem(doc *domain.Item, oldDoc *domain.Item) (bool, map[string]interface{}, map[string]interface{}) {
	if doc.ApprovalStatus != utils.APPROVAL_STATUS.WAIT_FOR_APPROVAL {
		return false, nil, nil
	}
	n := map[string]interface{}{}
	o := map[string]interface{}{}
	if strings.TrimSpace(doc.Name.Ar) != strings.TrimSpace(oldDoc.Name.Ar) || strings.TrimSpace(doc.Name.En) != strings.TrimSpace(oldDoc.Name.En) {
		n["name"] = utils.StructToMap(doc.Name, "bson")
		o["name"] = utils.StructToMap(oldDoc.Name, "bson")
	}
	if strings.TrimSpace(doc.Desc.Ar) != strings.TrimSpace(oldDoc.Desc.Ar) || strings.TrimSpace(doc.Desc.En) != strings.TrimSpace(oldDoc.Desc.En) {
		n["desc"] = utils.StructToMap(doc.Desc, "bson")
		o["desc"] = utils.StructToMap(oldDoc.Desc, "bson")
	}
	if doc.Price != oldDoc.Price {
		n["price"] = doc.Price
		o["price"] = oldDoc.Price
	}
	if strings.TrimSpace(doc.Image) != strings.TrimSpace(oldDoc.Image) {
		n["image"] = doc.Image
		o["image"] = oldDoc.Image
	}
	return len(n) >= 1, n, o
}
