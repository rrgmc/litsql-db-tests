package psql

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rrgmc/litsql/dialect/psql"
	"github.com/rrgmc/litsql/dialect/psql/sm"
	"github.com/rrgmc/litsql/sq"
	"gotest.tools/v3/assert"
)

func TestSelect(t *testing.T) {
	ctx := context.Background()

	runDBTest(t, ctx, func(db *sql.DB) {
		query := psql.Select(
			sm.Columns("film_id", "title", "length"),
			sm.From("film"),
			sm.WhereC("length > ?", sq.NamedArg("length")),
			sm.Limit(10),
		)

		squery, params, err := query.Build()
		assert.NilError(t, err)

		args, err := sq.ParseArgs(params, map[string]any{
			"length": 100,
		})
		assert.NilError(t, err)

		rows, err := db.QueryContext(ctx, squery, args...)
		assert.NilError(t, err)
		defer rows.Close()

		var ct int

		for rows.Next() {
			var id, length int
			var title string
			err := rows.Scan(&id, &title, &length)
			assert.NilError(t, err)
			ct++
		}

		assert.NilError(t, rows.Err())
		assert.Equal(t, 10, ct)
	})
}
