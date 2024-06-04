package psql

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rrgmc/litsql/dialect/psql"
	"github.com/rrgmc/litsql/dialect/psql/sm"
	"github.com/rrgmc/litsql/sq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gotest.tools/v3/assert"
)

func TestSelect(t *testing.T) {
	ctx := context.Background()

	dbName := "sakila"
	dbUser := "postgres"
	dbPassword := "password"

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("ghcr.io/rrgmc/litsql-dbtest-sakila-postgres:latest"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(1).
				WithStartupTimeout(5*time.Second)),
	)
	assert.NilError(t, err)

	// Clean up the container
	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NilError(t, err)

	db, err := sql.Open("pgx", connStr)
	assert.NilError(t, err)

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
}
