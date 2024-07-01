package admin

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/admin/domain"
	"samm/internal/module/admin/dto/admin"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/utils/dto"
	"strings"
)

type adminRepo struct {
	adminCollection *mgm.Collection
	logger          logger.ILogger
}

func NewAdminRepository(dbs *mongo.Database, log logger.ILogger) domain.AdminRepository {
	adminCollection := mgm.Coll(&domain.Admin{})
	createIndexes(adminCollection.Collection)
	return &adminRepo{
		adminCollection: adminCollection,
		logger:          log,
	}
}

func (r *adminRepo) Create(ctx context.Context, domainData *domain.Admin) (*domain.Admin, error) {
	err := mgm.Coll(domainData).Create(domainData)
	if err != nil {
		r.logger.Error("adminRepo -> Create -> ", err)
	}
	return domainData, err
}

func (r *adminRepo) Update(ctx context.Context, domainData *domain.Admin) (*domain.Admin, error) {
	err := mgm.Coll(domainData).Update(domainData)
	if err != nil {
		r.logger.Error("adminRepo -> Update -> ", err)
	}
	return domainData, err
}

func (r *adminRepo) SyncRole(ctx context.Context, domainData *domain.Role) error {
	_, err := mgm.Coll(&domain.Admin{}).UpdateMany(ctx, bson.M{"role._id": domainData.ID}, bson.M{"$set": bson.M{"role": domainData}})
	if err != nil {
		r.logger.Error("adminRepo -> SyncRole -> ", err)
	}
	return err
}

func (r *adminRepo) Delete(ctx context.Context, domainData *domain.Admin, adminDetails dto.AdminDetails) error {
	domainData.SetSoftDelete(ctx)
	domainData.AdminDetails = append(domainData.AdminDetails, adminDetails)
	err := mgm.Coll(domainData).Update(domainData)
	if err != nil {
		r.logger.Error("adminRepo -> Delete -> ", err)
	}
	return err
}

func (r *adminRepo) Find(ctx context.Context, adminId primitive.ObjectID) (*domain.Admin, error) {
	domainData := domain.Admin{}
	result := mgm.Coll(&domain.Admin{}).FindOne(ctx, bson.M{"_id": adminId, "deleted_at": nil})
	if err := result.Err(); err != nil {
		r.logger.Error("adminRepo -> Find -> ", err)
		return &domainData, err
	}
	if err := result.Decode(&domainData); err != nil {
		r.logger.Error("adminRepo -> Find -> ", err)
		return &domainData, err
	}
	return &domainData, nil
}

func (r *adminRepo) List(ctx context.Context, dto *admin.ListAdminDTO) ([]domain.Admin, *PaginationData, error) {
	models := make([]domain.Admin, 0)

	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.D{{"deleted_at", nil}},
	}}}

	if dto.CountryId != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"country_ids": dto.CountryId})
	}

	if dto.Query != "" {
		pattern := ".*" + dto.Query + ".*"
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"$or": []bson.M{{"name": bson.M{"$regex": pattern, "$options": "i"}}, {"email": bson.M{"$regex": pattern, "$options": "i"}}}})
	}

	if dto.Status != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"status": strings.ToLower(dto.Status)})
	}

	if dto.Type != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"type": strings.ToLower(dto.Type)})
	}

	if dto.Role != "" {
		matching["$match"].(bson.M)["$and"] = append(matching["$match"].(bson.M)["$and"].([]interface{}), bson.M{"role": strings.ToLower(dto.Role)})
	}

	data, err := New(r.adminCollection.Collection).Context(ctx).Limit(dto.Limit).Page(dto.Page).Sort("created_at", -1).Aggregate(matching)

	if data == nil || data.Data == nil {
		return models, nil, err
	}

	for _, raw := range data.Data {
		model := domain.Admin{}
		errUnmarshal := bson.Unmarshal(raw, &model)
		if errUnmarshal != nil {
			r.logger.Error("adminRepo -> List -> ", errUnmarshal)
			break
		}
		models = append(models, model)
	}

	return models, &data.Pagination, err
}

func (r *adminRepo) ChangeStatus(ctx context.Context, model *domain.Admin, input *admin.ChangeAdminStatusDto, adminDetails dto.AdminDetails) error {
	model.Status = input.Status
	model.AdminDetails = utils.If(model.AdminDetails == nil, []dto.AdminDetails{adminDetails}, append(model.AdminDetails, adminDetails)).([]dto.AdminDetails)
	err := mgm.Coll(model).Update(model)
	if err != nil {
		r.logger.Error("adminRepo -> ChangeStatus -> ", err)
	}
	return err
}

func (r *adminRepo) CheckEmailExists(ctx context.Context, email string, adminId primitive.ObjectID) (bool, error) {
	filter := bson.M{"email": email, "deleted_at": nil}
	if !adminId.IsZero() {
		filter["_id"] = bson.M{"$ne": adminId}
	}
	c, err := r.adminCollection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("adminRepo -> CheckEmailExists -> ", err)
	}
	return c > 0, err
}

func (r *adminRepo) CheckRoleExists(ctx context.Context, roleId primitive.ObjectID) (bool, error) {
	filter := bson.M{"role._id": roleId, "deleted_at": nil}
	c, err := r.adminCollection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("adminRepo -> CheckRoleExists -> ", err)
	}
	return c > 0, err
}

func (r *adminRepo) FindByEmail(ctx context.Context, email, adminType string) (*domain.Admin, error) {
	domainData := domain.Admin{}
	result := mgm.Coll(&domain.Admin{}).FindOne(ctx, bson.M{"email": email, "type": adminType, "deleted_at": nil})
	if err := result.Err(); err != nil {
		r.logger.Error("adminRepo -> FindByEmail -> ", err)
		return &domainData, err
	}
	if err := result.Decode(&domainData); err != nil {
		r.logger.Error("adminRepo -> FindByEmail -> ", err)
		return &domainData, err
	}
	return &domainData, nil
}

func (r *adminRepo) FindByToken(ctx context.Context, token string, adminType []string) (*domain.Admin, error) {
	domainData := domain.Admin{}
	result := mgm.Coll(&domain.Admin{}).FindOne(ctx, bson.M{"tokens": token, "type": bson.M{"$in": adminType}, "deleted_at": nil})
	if err := result.Err(); err != nil {
		r.logger.Error("adminRepo -> FindByToken -> ", err)
		return &domainData, err
	}
	if err := result.Decode(&domainData); err != nil {
		r.logger.Error("adminRepo -> FindFindByTokenByEmail -> ", err)
		return &domainData, err
	}
	return &domainData, nil
}
