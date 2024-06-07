package psql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/rrgmc/litsql-db-tests/utils"
	"github.com/rrgmc/litsql/dialect/psql"
	"github.com/rrgmc/litsql/dialect/psql/sm"
	"gotest.tools/v3/assert"
)

func TestSakila1(t *testing.T) {
	ctx := context.Background()

	runDBTest(t, ctx, func(db *sql.DB) {
		query := psql.Select(
			sm.Columns("film.film_id", "film.title", "store.store_id", "inventory.inventory_id"),
			sm.From("inventory"),
			sm.InnerJoin("store").Using("store_id"),
			sm.InnerJoin("film").Using("film_id"),
			sm.WhereClause("film.title = ? AND store.store_id = ?", "ACADEMY DINOSAUR", 1),
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
		query := psql.Select(
			sm.Columns("inventory.inventory_id"),
			sm.From("inventory"),
			sm.InnerJoin("store").Using("store_id"),
			sm.InnerJoin("film").Using("film_id"),
			sm.InnerJoin("rental").Using("inventory_id"),
			sm.WhereClause("film.title = ? AND store.store_id = ?", "ACADEMY DINOSAUR", 1),
			sm.WhereClause("NOT EXISTS ?", psql.Select(
				sm.Columns("*"),
				sm.From("rental"),
				sm.WhereClause("rental.inventory_id = inventory.inventory_id AND rental.return_date IS NULL"),
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
		query := psql.Select(
			sm.Columns("rental_date"),
			sm.ColumnsClause("rental_date + ? * INTERVAL '1 DAY' AS due_date", psql.Select(
				sm.Columns("rental_duration"),
				sm.From("film"),
				sm.WhereClause("film_id = ?", 1),
			)),
			sm.From("rental"),
			sm.WhereClause("rental_id = ?", psql.Select(
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
		query := psql.Select(
			sm.Columns("category.name", "avg(length)"),
			sm.From("film"),
			sm.InnerJoin("film_category").Using("film_id"),
			sm.InnerJoin("category").Using("category_id"),
			sm.GroupBy("category.name"),
			sm.HavingClause("avg(length) > ?", psql.Select(
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

func TestSakila5(t *testing.T) {
	ctx := context.Background()

	runDBTest(t, ctx, func(db *sql.DB) {
		query := psql.Select(
			sm.With("cte_src").As(psql.Select(
				sm.Columns("*", "ROW_NUMBER() OVER (PARTITION BY email ORDER BY email) AS ROWNUM"),
				sm.From("customer c"),
			)),
			sm.Columns("*"),
			sm.From("cte_src"),
			sm.Limit(10),
		)

		var ct int

		utils.DBExecute(t, db, query, nil,
			func(row map[string]any) {
				ct++
			})

		assert.Assert(t, ct > 0)
	})
}

func TestSakila6(t *testing.T) {
	ctx := context.Background()

	runDBTest(t, ctx, func(db *sql.DB) {
		query := psql.Select(
			sm.Columns("f.film_id", "f.title", "fc.category_id"),
			sm.From("film f"),
			sm.LeftJoinExpr(psql.Select(
				sm.Columns("*"),
				sm.From("film_category"),
				sm.WhereClause("film_id > ?", 3),
			)).As("fc").On("f.film_id = fc.film_id"),
			sm.OrderBy("f.film_id"),
		)

		var ct int

		utils.DBExecute(t, db, query, nil,
			func(row map[string]any) {
				ct++
			})

		assert.Assert(t, ct > 0)
	})
}

func TestSakila7(t *testing.T) {
	ctx := context.Background()

	runDBTest(t, ctx, func(db *sql.DB) {
		query := psql.Select(
			sm.Columns("first_name", "last_name"),
			sm.From("actor"),
			sm.WhereClause("actor_id IN ?", psql.Select(
				sm.Columns("actor_id"),
				sm.From("film_actor"),
				sm.WhereClause("film_id IN ?", psql.Select(
					sm.Columns("film_id"),
					sm.From("film"),
					sm.WhereClause("title = ?", "ALONE TRIP"),
				)),
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

func TestSakila8(t *testing.T) {
	ctx := context.Background()

	runDBTest(t, ctx, func(db *sql.DB) {
		query := psql.Select(
			sm.With("cte").As(psql.Select(
				sm.Columns("film.film_id", "film.title", "COUNT(film_actor.actor_id) AS actor_count"),
				sm.From("film_actor"),
				sm.InnerJoin("film").On("film_actor.film_id = film.film_id"),
				sm.GroupBy("film.film_id", "film.title"),
			)),
			sm.Columns("cte.title", "cte.actor_count"),
			sm.From("cte"),
			sm.WhereClause("cte.actor_count > ?", psql.Select(
				sm.Columns("avg(actor_count)"),
				sm.From("cte"),
			)),
			sm.OrderBy("cte.title"),
			sm.Limit(10),
		)

		var ct int

		utils.DBExecute(t, db, query, nil,
			func(row map[string]any) {
				ct++
			})

		assert.Assert(t, ct > 0)
	})
}

func TestSakila9(t *testing.T) {
	ctx := context.Background()

	runDBTest(t, ctx, func(db *sql.DB) {
		query := psql.Select(
			sm.With("table1").As(psql.Select(
				sm.Columns("f.film_id", "f.title AS \"Film\""),
				sm.From("film f"),
			)),
			sm.With("table2").As(psql.Select(
				sm.Columns("COUNT(r.rental_id) rental_count", "i.film_id"),
				sm.From("inventory i"),
				sm.InnerJoin("rental r").On("i.inventory_id = r.inventory_id"),
				sm.GroupBy("i.film_id"),
			)),
			sm.Columns("table2.rental_count", "table1.\"Film\""),
			sm.From("table1"),
			sm.InnerJoin("table2").On("table1.film_id = table2.film_id"),
			sm.OrderBy("rental_count DESC"),
			sm.Limit(10),
		)

		var ct int

		utils.DBExecute(t, db, query, nil,
			func(row map[string]any) {
				ct++
			})

		assert.Assert(t, ct > 0)
	})
}
