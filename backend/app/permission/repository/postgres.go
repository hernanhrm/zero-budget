package repository

import (
	"context"
	"errors"

	"backend/app/permission/domain"
	basedomain "backend/domain"
	apperrors "backend/domain/errors"
	"backend/infra/dafi"
	"backend/infra/database"
	"backend/infra/sqlcraft"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/samber/oops"
)

const tableName = "auth.permissions"

var columns = []string{
	"id",
	"slug",
	"description",
}

var sqlColumnByDomainField = map[string]string{
	"id":          "id",
	"slug":        "slug",
	"description": "description",
	"createdAt":   "created_at",
}

type postgres struct {
	db     database.PoolInterface
	logger basedomain.Logger
}

func NewPostgres(db database.PoolInterface, logger basedomain.Logger) domain.Repository {
	return postgres{
		db:     db,
		logger: logger.With("component", "permission.repository"),
	}
}

func (r postgres) FindOne(ctx context.Context, criteria dafi.Criteria) (domain.Permission, error) {
	query := sqlcraft.Select(columns...).
		From(tableName).
		Where(criteria.Filters...).
		SQLColumnByDomainField(sqlColumnByDomainField).
		Limit(1)

	result, err := query.ToSQL()
	if err != nil {
		return domain.Permission{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)

	row := r.db.QueryRow(ctx, result.SQL, result.Args...)

	var item domain.Permission
	err = row.Scan(&item.ID, &item.Slug, &item.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Permission{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Code(apperrors.CodeNotFound).Wrap(err)
		}
		return domain.Permission{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}
	return item, nil
}

func (r postgres) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[domain.Permission], error) {
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

	var items basedomain.List[domain.Permission]
	for rows.Next() {
		var item domain.Permission
		err = rows.Scan(&item.ID, &item.Slug, &item.Description)
		if err != nil {
			return nil, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (r postgres) Create(ctx context.Context, input domain.CreatePermission) error {
	id := uuid.New().String()

	query := sqlcraft.InsertInto(tableName).
		WithColumns(columns...).
		WithValues(id, input.Slug, input.Description)

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

func (r postgres) CreateBulk(ctx context.Context, inputs basedomain.List[domain.CreatePermission]) error {
	if inputs.IsEmpty() {
		return nil
	}
	query := sqlcraft.InsertInto(tableName).WithColumns(columns...)

	for _, input := range inputs {
		id := uuid.New().String()
		query = query.WithValues(id, input.Slug, input.Description)
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

func (r postgres) Update(ctx context.Context, input domain.UpdatePermission, filters ...dafi.Filter) error {
	cols := []string{}
	vals := []any{}
	if input.Description.Valid {
		cols = append(cols, "description")
		vals = append(vals, input.Description.String)
	}

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
