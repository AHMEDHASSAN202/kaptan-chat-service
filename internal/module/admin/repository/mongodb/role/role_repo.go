package role

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/admin/domain"
	"samm/internal/module/admin/dto/role"
	"samm/pkg/logger"
)

type roleRepo struct {
	roleCollection  *mgm.Collection
	adminRepository domain.AdminRepository
	logger          logger.ILogger
}

func NewRoleRepository(dbs *mongo.Database, log logger.ILogger, adminRepository domain.AdminRepository) domain.RoleRepository {
	roleCollection := mgm.Coll(&domain.Role{})
	createIndexes(roleCollection.Collection)
	return &roleRepo{
		roleCollection:  roleCollection,
		adminRepository: adminRepository,
		logger:          log,
	}
}

func (r *roleRepo) Create(ctx context.Context, domainData *domain.Role) (*domain.Role, error) {
	err := mgm.Coll(domainData).Create(domainData)
	if err != nil {
		r.logger.Error("roleRepo -> Create -> ", err)
	}
	return domainData, err
}

func (r *roleRepo) Update(ctx context.Context, domainData *domain.Role) (*domain.Role, error) {
	err := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		err := r.roleCollection.UpdateWithCtx(sc, domainData)
		if err != nil {
			r.logger.Error("roleRepo -> Update -> ", err)
			return err
		}
		err = r.adminRepository.SyncRole(sc, domainData)
		if err != nil {
			r.logger.Error("roleRepo -> SyncRole -> ", err)
			return err
		}
		return session.CommitTransaction(sc)
	})

	if err != nil {
		r.logger.Error("roleRepo -> Update -> ", err)
	}
	return domainData, err
}

func (r *roleRepo) Delete(ctx context.Context, domainData *domain.Role) error {
	err := mgm.Coll(domainData).Delete(domainData)
	if err != nil {
		r.logger.Error("roleRepo -> Delete -> ", err)
	}
	return err
}

func (r *roleRepo) Find(ctx context.Context, roleId primitive.ObjectID) (*domain.Role, error) {
	domainData := domain.Role{}
	result := mgm.Coll(&domain.Role{}).FindOne(ctx, bson.M{"_id": roleId})
	if err := result.Err(); err != nil {
		r.logger.Error("roleRepo -> Find -> ", err)
		return &domainData, err
	}
	if err := result.Decode(&domainData); err != nil {
		r.logger.Error("roleRepo -> Find -> ", err)
		return &domainData, err
	}
	return &domainData, nil
}

func (r *roleRepo) List(ctx context.Context, dto *role.ListRoleDTO) ([]domain.Role, *PaginationData, error) {
	models := make([]domain.Role, 0)

	query := New(r.roleCollection.Collection).Context(ctx).Limit(dto.Limit).Page(dto.Page).Sort("created_at", -1)

	var data *PaginatedData
	var err error

	matchStage := bson.M{}
	if dto.Query != "" {
		pattern := ".*" + dto.Query + ".*"
		matchStage = bson.M{
			"$or": []bson.M{
				{"name.ar": bson.M{"$regex": pattern, "$options": "i"}},
				{"name.en": bson.M{"$regex": pattern, "$options": "i"}},
			},
		}
	}

	if len(matchStage) > 0 {
		data, err = query.Aggregate(bson.M{"$match": matchStage})
	} else {
		data, err = query.Aggregate()
	}

	if err != nil {
		r.logger.Error("roleRepo -> List -> ", err)
		return models, nil, err
	}

	if data == nil || data.Data == nil {
		return models, nil, err
	}

	for _, raw := range data.Data {
		model := domain.Role{}
		errUnmarshal := bson.Unmarshal(raw, &model)
		if errUnmarshal != nil {
			r.logger.Error("roleRepo -> List -> ", errUnmarshal)
			break
		}
		models = append(models, model)
	}

	return models, &data.Pagination, err
}
