package postgres

import (
	"context"
	"errors"
	"testing"

	basedomain "backend/port"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type scanRow struct {
	val int64
	err error
}

func (r scanRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	p, ok := dest[0].(*int64)
	if !ok {
		return errors.New("expected *int64 dest")
	}
	*p = r.val
	return nil
}

type countQueryDB struct {
	row pgx.Row
}

func (c countQueryDB) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	panic("unexpected Exec")
}

func (c countQueryDB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	panic("unexpected Query")
}

func (c countQueryDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return c.row
}

func TestPostgres_CountByAccountID(t *testing.T) {
	t.Parallel()

	id := uuid.MustParse("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")

	tests := []struct {
		name    string
		row     scanRow
		want    int64
		wantErr bool
	}{
		{
			name: "returns count",
			row:  scanRow{val: 42},
			want: 42,
		},
		{
			name:    "scan error",
			row:     scanRow{err: errors.New("boom")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := postgres{
				db:     countQueryDB{row: tt.row},
				logger: noopTestLogger{},
			}

			got, err := r.CountByAccountID(context.Background(), id)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

type noopTestLogger struct{}

func (noopTestLogger) Debug(string, ...any)                     {}
func (noopTestLogger) Info(string, ...any)                      {}
func (noopTestLogger) Warn(string, ...any)                      {}
func (noopTestLogger) Error(string, ...any)                    {}
func (noopTestLogger) With(...any) basedomain.Logger           { return noopTestLogger{} }
func (noopTestLogger) WithContext(context.Context) basedomain.Logger {
	return noopTestLogger{}
}
