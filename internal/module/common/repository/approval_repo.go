package repository

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	builder "samm/internal/module/common/builder/approval"
	"samm/internal/module/common/domain"
	dto2 "samm/internal/module/common/dto"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"strings"
)

type approvalRepo struct {
	dbs                *mongo.Database
	approvalCollection *mgm.Collection
	logger             logger.ILogger
}

func NewApprovalRepository(dbs *mongo.Database, log logger.ILogger) domain.ApprovalRepository {
	approvalCollection := mgm.Coll(&domain.Approval{})
	return &approvalRepo{
		dbs:                dbs,
		approvalCollection: approvalCollection,
		logger:             log,
	}
}

func (r *approvalRepo) CreateOrUpdate(ctx context.Context, dto []dto2.CreateApprovalDto) error {
	var models []mongo.WriteModel
	for _, approvalDto := range dto {
		domainData := builder.CreateApprovalBuilder(&approvalDto)
		filter := bson.M{"entity_id": domainData.EntityId, "entity_type": domainData.EntityType, "status": utils.APPROVAL_STATUS.WAIT_FOR_APPROVAL}
		update := bson.M{"$set": domainData}
		model := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		models = append(models, model)
	}
	_, err := mgm.Coll(&domain.Approval{}).BulkWrite(ctx, models)
	if err != nil {
		r.logger.Error("approvalRepo -> CreateOrUpdate -> ", err)
	}
	return err
}

func (r *approvalRepo) FindByEntity(ctx context.Context, entityId primitive.ObjectID, entityType string) (*domain.Approval, error) {
	domainData := domain.Approval{}
	result := mgm.Coll(&domain.Approval{}).FindOne(ctx, bson.M{"entity_id": entityId, "entity_type": entityType, "status": utils.APPROVAL_STATUS.WAIT_FOR_APPROVAL})
	if err := result.Err(); err != nil {
		r.logger.Error("approvalRepo -> FindByEntity -> ", err)
		return &domainData, err
	}
	if err := result.Decode(&domainData); err != nil {
		r.logger.Error("approvalRepo -> FindByEntity -> ", err)
		return &domainData, err
	}
	return &domainData, nil
}

func (r *approvalRepo) FindById(ctx context.Context, approvalId primitive.ObjectID) (*domain.Approval, error) {
	domainData := domain.Approval{}
	result := mgm.Coll(&domain.Approval{}).FindOne(ctx, bson.M{"_id": approvalId})
	if err := result.Err(); err != nil {
		r.logger.Error("approvalRepo -> FindById -> ", err)
		return &domainData, err
	}
	if err := result.Decode(&domainData); err != nil {
		r.logger.Error("approvalRepo -> FindById -> ", err)
		return &domainData, err
	}
	return &domainData, nil
}

func (r *approvalRepo) DeleteByEntity(ctx context.Context, entityId primitive.ObjectID, entityType string) error {
	_, err := mgm.Coll(&domain.Approval{}).DeleteMany(ctx, bson.M{"entity_id": entityId, "entity_type": entityType})
	if err != nil {
		r.logger.Error("roleRepo -> Delete -> ", err)
	}
	return err
}

func (r *approvalRepo) List(ctx context.Context, dto *dto2.ListApprovalDto) ([]domain.Approval, *PaginationData, error) {
	models := make([]domain.Approval, 0)
	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.D{{"country_id", dto.CountryId}},
		bson.D{{"status", utils.APPROVAL_STATUS.WAIT_FOR_APPROVAL}},
	}}}
	if dto.Type != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"entity_type": strings.ToLower(dto.Type)})
	}
	data, err := New(r.approvalCollection.Collection).Context(ctx).Limit(dto.Limit).Page(dto.Page).Sort("created_at", -1).Aggregate(matching)
	if data == nil || data.Data == nil {
		return models, nil, err
	}
	for _, raw := range data.Data {
		model := domain.Approval{}
		errUnmarshal := bson.Unmarshal(raw, &model)
		if errUnmarshal != nil {
			r.logger.Error("approvalRepo -> List -> ", errUnmarshal)
			break
		}
		models = append(models, model)
	}
	return models, &data.Pagination, err
}

func (r *approvalRepo) ChangeStatus(ctx context.Context, domainData *domain.Approval) error {
	err := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		err := r.approvalCollection.UpdateWithCtx(sc, domainData)
		if err != nil {
			r.logger.Error("approvalRepo -> Update -> ", err)
			return err
		}
		updateData := bson.M{"approval_status": domainData.Status}
		if domainData.Status == utils.APPROVAL_STATUS.APPROVED {
			updateData["has_original"] = true
			for key, value := range domainData.Fields.New {
				updateData[key] = value
			}
		}
		_, err = r.dbs.Collection(domainData.EntityType).UpdateOne(sc, bson.M{"_id": domainData.EntityId}, bson.M{"$set": updateData})
		if err != nil {
			r.logger.Error("approvalRepo -> SyncRole -> ", err)
			return err
		}
		return session.CommitTransaction(sc)
	})
	if err != nil {
		r.logger.Error("approvalRepo -> Update -> ", err)
	}
	return err
}
