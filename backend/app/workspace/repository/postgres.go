package repository

import (
	"context"
	"errors"
	"time"

	"backend/app/workspace/domain"
	basedomain "backend/domain"
	apperrors "backend/domain/errors"
	"backend/infra/dafi"
	"backend/infra/database"
	"backend/infra/sqlcraft"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/samber/oops"
)

const tableName = "auth.workspaces"

var columns = []string{
	"id",
	"organization_id",
	"name",
	"slug",
	"created_at",
	"updated_at",
}

var sqlColumnByDomainField = map[string]string{
	"id":             "id",
	"organizationId": "organization_id",
	"name":           "name",
	"slug":           "slug",
	"createdAt":      "created_at",
	"updatedAt":      "updated_at",
}

type postgres struct {
	db     database.PoolInterface
	logger basedomain.Logger
}

func NewPostgres(db database.PoolInterface, logger basedomain.Logger) domain.Repository {
	return postgres{
		db:     db,
		logger: logger.With("component", "workspace.repository"),
	}
}

func (r postgres) FindOne(ctx context.Context, criteria dafi.Criteria) (domain.Workspace, error) {
	query := sqlcraft.Select(columns...).
		From(tableName).
		Where(criteria.Filters...).
		SQLColumnByDomainField(sqlColumnByDomainField).
		Limit(1)

	result, err := query.ToSQL()
	if err != nil {
		return domain.Workspace{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)

	row := r.db.QueryRow(ctx, result.SQL, result.Args...)

	var item domain.Workspace
	err = row.Scan(&item.ID, &item.OrganizationID, &item.Name, &item.Slug, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Workspace{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Code(apperrors.CodeNotFound).Wrap(err)
		}
		return domain.Workspace{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}
	return item, nil
}

func (r postgres) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[domain.Workspace], error) {
	query := sqlcraft.Select(columns...).
		From(tableName).
		Where(criteria.Filters...).
		OrderBy(criteria.Sorts...).
		Limit(criteria.Pagination.PageSize).
		Page(criteria.Pagination.PageNumber).
		SQLColumnByDomainField(sqlColumnByDomainField)

	result, err := query.ToSQL()
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)

	rows, err := r.db.Query(ctx, result.SQL, result.Args...)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}
	defer rows.Close()

	var items basedomain.List[domain.Workspace]
	for rows.Next() {
		var item domain.Workspace
		err = rows.Scan(&item.ID, &item.OrganizationID, &item.Name, &item.Slug, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (r postgres) Create(ctx context.Context, input domain.CreateWorkspace) error {
	now := time.Now()
	id := uuid.New().String()

	query := sqlcraft.InsertInto(tableName).
		WithColumns(columns...).
		WithValues(id, input.OrganizationID, input.Name, input.Slug, now, now)

	result, err := query.ToSQL()
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)

	_, err = r.db.Exec(ctx, result.SQL, result.Args...)
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}
	return nil
}

func (r postgres) CreateBulk(ctx context.Context, inputs basedomain.List[domain.CreateWorkspace]) error {
	if inputs.IsEmpty() {
		return nil
	}
	now := time.Now()
	query := sqlcraft.InsertInto(tableName).WithColumns(columns...)

	for _, input := range inputs {
		id := uuid.New().String()
		query = query.WithValues(id, input.OrganizationID, input.Name, input.Slug, now, now)
	}

	result, err := query.ToSQL()
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	r.logger.WithContext(ctx).Debug("executing bulk insert", "sql", result.SQL, "count", len(inputs))
	_, err = r.db.Exec(ctx, result.SQL, result.Args...)
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}
	return nil
}

func (r postgres) Update(ctx context.Context, input domain.UpdateWorkspace, filters ...dafi.Filter) error {
	cols := []string{}
	vals := []any{}
	if input.Name.Valid {
		cols = append(cols, "name")
		vals = append(vals, input.Name.String)
	}
	if input.Slug.Valid {
		cols = append(cols, "slug")
		vals = append(vals, input.Slug.String)
	}
	cols = append(cols, "updated_at")
	vals = append(vals, time.Now())

	query := sqlcraft.Update(tableName).
		WithColumns(cols...).
		WithValues(vals...).
		Where(filters...).
		SQLColumnByDomainField(sqlColumnByDomainField)

	result, err := query.ToSQL()
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}
	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)
	_, err = r.db.Exec(ctx, result.SQL, result.Args...)
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}
	return nil
}

func (r postgres) Delete(ctx context.Context, filters ...dafi.Filter) error {
	query := sqlcraft.DeleteFrom(tableName).
		Where(filters...).
		SQLColumnByDomainField(sqlColumnByDomainField)

	result, err := query.ToSQL()
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}
	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)
	_, err = r.db.Exec(ctx, result.SQL, result.Args...)
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}
	return nil
}
