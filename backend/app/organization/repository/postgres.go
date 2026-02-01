package repository

import (
	"context"
	"errors"
	"time"

	"backend/app/organization/domain"
	basedomain "backend/domain"
	apperrors "backend/domain/errors"
	"backend/infra/dafi"
	"backend/infra/database"
	"backend/infra/sqlcraft"

	"github.com/jackc/pgx/v5"
	"github.com/samber/oops"
)

const tableName = "auth.organizations"

var columns = []string{
	"id",
	"name",
	"slug",
	"owner_id",
	"created_at",
	"updated_at",
}

var sqlColumnByDomainField = map[string]string{
	"id":        "id",
	"name":      "name",
	"slug":      "slug",
	"ownerId":   "owner_id",
	"createdAt": "created_at",
	"updatedAt": "updated_at",
}

type postgres struct {
	db     database.PoolInterface
	logger basedomain.Logger
}

func NewPostgres(db database.PoolInterface, logger basedomain.Logger) domain.Repository {
	return postgres{
		db:     db,
		logger: logger.With("component", "organization.repository"),
	}
}

func (r postgres) FindOne(ctx context.Context, criteria dafi.Criteria) (domain.Organization, error) {
	query := sqlcraft.Select(columns...).
		From(tableName).
		Where(criteria.Filters...).
		SQLColumnByDomainField(sqlColumnByDomainField).
		Limit(1)

	result, err := query.ToSQL()
	if err != nil {
		return domain.Organization{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)

	row := r.db.QueryRow(ctx, result.SQL, result.Args...)

	var organization domain.Organization
	err = row.Scan(
		&organization.ID,
		&organization.Name,
		&organization.Slug,
		&organization.OwnerID,
		&organization.CreatedAt,
		&organization.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Organization{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Code(apperrors.CodeNotFound).Wrap(err)
		}
		return domain.Organization{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	return organization, nil
}

func (r postgres) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[domain.Organization], error) {
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

	var organizations basedomain.List[domain.Organization]
	for rows.Next() {
		var organization domain.Organization
		err = rows.Scan(
			&organization.ID,
			&organization.Name,
			&organization.Slug,
			&organization.OwnerID,
			&organization.CreatedAt,
			&organization.UpdatedAt,
		)
		if err != nil {
			return nil, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
		}
		organizations = append(organizations, organization)
	}

	return organizations, nil
}

func (r postgres) Create(ctx context.Context, input domain.CreateOrganization) error {
	now := time.Now()

	query := sqlcraft.InsertInto(tableName).
		WithColumns(columns...).
		WithValues(input.ID, input.Name, input.Slug, input.OwnerID, now, now)

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

func (r postgres) CreateBulk(ctx context.Context, inputs basedomain.List[domain.CreateOrganization]) error {
	if inputs.IsEmpty() {
		return nil
	}

	now := time.Now()
	query := sqlcraft.InsertInto(tableName).WithColumns(columns...)

	for _, input := range inputs {
		query = query.WithValues(input.ID, input.Name, input.Slug, input.OwnerID, now, now)
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

func (r postgres) Update(ctx context.Context, input domain.UpdateOrganization, filters ...dafi.Filter) error {
	query := sqlcraft.Update(tableName).
		WithColumns("name", "slug", "updated_at").
		WithValues(input.Name, input.Slug, time.Now()).
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
