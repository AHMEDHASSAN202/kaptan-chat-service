package approval_helper

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	domain2 "samm/internal/module/approval/domain"
	"samm/internal/module/approval/dto"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/external"
	"samm/internal/module/menu/responses/item"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"strings"
	"time"
)

type ApprovalItemHelper struct {
	approvalRepo   domain2.ApprovalRepository
	itemCollection *mgm.Collection
	logger         logger.ILogger
	extService     external.ExtService
}

func NewApprovalItemHelper(approvalRepo domain2.ApprovalRepository, logger logger.ILogger, extService external.ExtService) *ApprovalItemHelper {
	return &ApprovalItemHelper{
		approvalRepo:   approvalRepo,
		itemCollection: mgm.Coll(&domain.Item{}),
		logger:         logger,
		extService:     extService,
	}
}

func (a *ApprovalItemHelper) CreateApprovalAsArray(sc mongo.SessionContext, docs []domain.Item) error {
	if docs[0].ApprovalStatus != utils.APPROVAL_STATUS.WAIT_FOR_APPROVAL {
		return nil
	}
	//get account
	account, errAccount := a.extService.RetailsIService.GetAccountById(context.Background(), utils.ConvertObjectIdToStringId(docs[0].AccountId))
	if errAccount.IsError {
		a.logger.Error("UpdateApproval -> ErrAccount -> ", errAccount.ErrorMessageObject)
		return errors.New(errAccount.ErrorMessageObject.Text)
	}
	return a.approvalRepo.CreateOrUpdate(sc, a.CreateItemsApprovalBuilder(docs, *account))
}

func (a *ApprovalItemHelper) UpdateApproval(sc mongo.SessionContext, doc *domain.Item, oldDoc *domain.Item) (bool, error) {
	// Check if approval is needed
	if needToApprove, n, o := a.NeedToApproveItem(doc, oldDoc); needToApprove {
		//get account
		account, errAccount := a.extService.RetailsIService.GetAccountById(context.Background(), utils.ConvertObjectIdToStringId(doc.AccountId))
		if errAccount.IsError {
			a.logger.Error("UpdateApproval -> ErrAccount -> ", errAccount.ErrorMessageObject)
			return false, errors.New(errAccount.ErrorMessageObject.Text)
		}
		// Create or update approval
		err := a.approvalRepo.CreateOrUpdate(sc, []dto.CreateApprovalDto{a.UpdateItemApprovalBuilder(doc, n, o, *account)})
		if err != nil {
			return false, err
		}
		// Update item with approval status and updated time
		_, err = a.itemCollection.Collection.UpdateOne(sc, bson.M{"_id": doc.ID}, bson.M{"$set": bson.M{"approval_status": utils.APPROVAL_STATUS.WAIT_FOR_APPROVAL, "updated_at": time.Now().UTC()}})
		if err != nil {
			a.logger.Error("ItemRepository -> UpdateOne -> ", err)
			return false, err
		}
		return false, nil
	}

	// If item is updated by admin, approve previous change
	if doc.ApprovalStatus == utils.APPROVAL_STATUS.APPROVED {
		err := a.approvalRepo.ApprovePreviousChange(sc, doc.ID, "items", doc.AdminDetails[len(doc.AdminDetails)-1])
		if err != nil {
			return true, err
		}
	}

	return true, nil
}

func (a *ApprovalItemHelper) DeleteApproval(sc mongo.SessionContext, doc *domain.Item) error {
	return a.approvalRepo.DeleteByEntity(sc, doc.ID, "items")
}

func (a *ApprovalItemHelper) NeedToApproveItem(doc *domain.Item, oldDoc *domain.Item) (bool, map[string]interface{}, map[string]interface{}) {
	if doc.ApprovalStatus != utils.APPROVAL_STATUS.WAIT_FOR_APPROVAL {
		return false, nil, nil
	}
	return a.AreAnyChanges(doc, oldDoc)
}

func (a *ApprovalItemHelper) AreAnyChanges(doc *domain.Item, oldDoc *domain.Item) (bool, map[string]interface{}, map[string]interface{}) {
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

func (a *ApprovalItemHelper) UpdateItemByApproval(ctx context.Context, doc *item.ItemResponse) validators.ErrorResponse {
	if doc.ApprovalStatus != utils.APPROVAL_STATUS.WAIT_FOR_APPROVAL {
		return validators.ErrorResponse{}
	}

	//find approval
	approvalDoc, err := a.approvalRepo.FindByEntity(ctx, doc.ID, "items")
	if err != nil {
		a.logger.Error("ApprovalItemHelper -> UpdateItemByApproval -> ", err)
		return validators.GetErrorResponse(&ctx, localization.E1000, nil, utils.GetAsPointer(http.StatusBadRequest))
	}

	//if entity not has approval
	if approvalDoc == nil {
		return validators.ErrorResponse{}
	}

	//apply changes to doc
	err = utils.CopyMapToStruct(&doc, approvalDoc.Fields.New)
	if err != nil {
		a.logger.Error("ApprovalItemHelper -> UpdateItemByApproval -> Copier ERROR -> ", err)
		return validators.GetErrorResponse(&ctx, localization.E1000, nil, utils.GetAsPointer(http.StatusBadRequest))
	}

	return validators.ErrorResponse{}
}
