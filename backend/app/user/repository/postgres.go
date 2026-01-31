package repository

import (
	"context"
	"errors"
	"time"

	basedomain "backend/domain"
	apperrors "backend/domain/errors"
	"backend/infra/dafi"
	"backend/infra/database"
	"backend/infra/sqlcraft"
	"backend/app/user/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/samber/oops"
)

const tableName = "users"

var columns = []string{
	"id",
	"name",
	"email",
	"created_at",
	"updated_at",
}

var sqlColumnByDomainField = map[string]string{
	"id":        "id",
	"name":      "name",
	"email":     "email",
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
		logger: logger.With("component", "user.repository"),
	}
}

func (r postgres) FindOne(ctx context.Context, criteria dafi.Criteria) (domain.User, error) {
	query := sqlcraft.Select(columns...).
		From(tableName).
		Where(criteria.Filters...).
		SQLColumnByDomainField(sqlColumnByDomainField).
		Limit(1)

	result, err := query.ToSQL()
	if err != nil {
		return domain.User{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)

	row := r.db.QueryRow(ctx, result.SQL, result.Args...)

	var user domain.User
	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Code(apperrors.CodeNotFound).Wrap(err)
		}
		return domain.User{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	return user, nil
}

func (r postgres) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[domain.User], error) {
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

	var users basedomain.List[domain.User]
	for rows.Next() {
		var user domain.User
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (r postgres) Create(ctx context.Context, input domain.CreateUser) error {
	now := time.Now()
	id := uuid.New().String()

	query := sqlcraft.InsertInto(tableName).
		WithColumns(columns...).
		WithValues(id, input.Name, input.Email, now, now)

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

func (r postgres) CreateBulk(ctx context.Context, inputs basedomain.List[domain.CreateUser]) error {
	if inputs.IsEmpty() {
		return nil
	}

	now := time.Now()
	query := sqlcraft.InsertInto(tableName).WithColumns(columns...)

	for _, input := range inputs {
		id := uuid.New().String()
		query = query.WithValues(id, input.Name, input.Email, now, now)
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

func (r postgres) Update(ctx context.Context, input domain.UpdateUser, filters ...dafi.Filter) error {
	cols := []string{}
	vals := []any{}

	if input.Name.Valid {
		cols = append(cols, "name")
		vals = append(vals, input.Name.String)
	}
	if input.Email.Valid {
		cols = append(cols, "email")
		vals = append(vals, input.Email.String)
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
