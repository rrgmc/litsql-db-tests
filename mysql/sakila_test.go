package mysql

import (
	"context"
	"database/sql"
	"testing"

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
			sm.WhereC("film.title = ? AND store.store_id = ?", "Academy Dinosaur", 1),
		)

		var ct int

		utils.DBExecute(t, db, query, nil,
			func(row map[string]any) {
				ct++
			})

		assert.Assert(t, ct > 0)
	})
}

func TestSakila2(t *testing.T) {
	ctx := context.Background()

	runDBTest(t, ctx, func(db *sql.DB) {
		query := mysql.Select(
			sm.Columns("inventory.inventory_id"),
			sm.From("inventory"),
			sm.InnerJoin("store").Using("store_id"),
			sm.InnerJoin("film").Using("film_id"),
			sm.InnerJoin("rental").Using("inventory_id"),
			sm.WhereC("film.title = ? AND store.store_id = ?", "Academy Dinosaur", 1),
			sm.WhereC("NOT EXISTS ?", mysql.Select(
				sm.Columns("*"),
				sm.From("rental"),
				sm.WhereC("rental.inventory_id = inventory.inventory_id AND rental.return_date IS NULL"),
			)),
		)

		var ct int

		utils.DBExecute(t, db, query, nil,
			func(row map[string]any) {
				ct++
			})

		assert.Assert(t, ct > 0)
	})
}

func TestSakila3(t *testing.T) {
	ctx := context.Background()

	runDBTest(t, ctx, func(db *sql.DB) {
		query := mysql.Select(
			sm.Columns("rental_date"),
			sm.ColumnsC("rental_date + interval ? day AS due_date", mysql.Select(
				sm.Columns("rental_duration"),
				sm.From("film"),
				sm.WhereC("film_id = ?", 1),
			)),
			sm.From("rental"),
			sm.WhereC("rental_id = ?", mysql.Select(
				sm.Columns("rental_id"),
				sm.From("rental"),
				sm.OrderBy("rental_id DESC"),
				sm.Limit(1),
			)),
		)

		var ct int

		utils.DBExecute(t, db, query, nil,
			func(row map[string]any) {
				ct++
			})

		assert.Assert(t, ct > 0)
	})
}

func TestSakila4(t *testing.T) {
	ctx := context.Background()

	runDBTest(t, ctx, func(db *sql.DB) {
		query := mysql.Select(
			sm.Columns("category.name", "avg(length)"),
			sm.From("film"),
			sm.InnerJoin("film_category").Using("film_id"),
			sm.InnerJoin("category").Using("category_id"),
			sm.GroupBy("category.name"),
			sm.HavingC("avg(length) > ?", mysql.Select(
				sm.Columns("avg(length)"),
				sm.From("film"),
			)),
			sm.OrderBy("avg(length) desc"),
		)

		var ct int

		utils.DBExecute(t, db, query, nil,
			func(row map[string]any) {
				ct++
			})

		assert.Assert(t, ct > 0)
	})
}
