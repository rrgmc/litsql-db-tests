package mysql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/rrgmc/litsql-db-tests/utils"
	"github.com/rrgmc/litsql/dialect/mysql"
	"github.com/rrgmc/litsql/dialect/mysql/sm"
	"github.com/rrgmc/litsql/sq"
	"gotest.tools/v3/assert"
)

func TestSelect(t *testing.T) {
	ctx := context.Background()

	runDBTest(t, ctx, func(db *sql.DB) {
		query := mysql.Select(
			sm.Columns("film_id", "title", "length"),
			sm.From("film"),
			sm.WhereClause("length > ?", sq.NamedArg("length")),
			sm.Limit(10),
		)

		var ct int

		utils.DBExecute(t, db, query, map[string]any{
			"length": 100,
		}, func(row map[string]any) {
			ct++
		})

		assert.Equal(t, 10, ct)
	})
}
