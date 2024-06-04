package mysql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/rrgmc/litsql-db-tests/utils"
	"github.com/rrgmc/litsql/dialect/mysql"
	"github.com/rrgmc/litsql/dialect/mysql/sm"
	"gotest.tools/v3/assert"
)

func TestSakila1(t *testing.T) {
	ctx := context.Background()

	runDBTest(t, ctx, func(db *sql.DB) {
		query := mysql.Select(
			sm.Columns("film.film_id", "film.title", "store.store_id", "inventory.inventory_id"),
			sm.From("inventory"),
			sm.InnerJoin("store").Using("store_id"),
			sm.InnerJoin("film").Using("film_id"),
			sm.WhereC("film.title = ? and store.store_id = ?", "Academy Dinosaur", 1),
		)

		var ct int

		utils.DBExecute(t, db, query, nil,
			func(row map[string]any) {
				spew.Dump(row)
				ct++
			})

		assert.Assert(t, ct > 0)
	})
}
