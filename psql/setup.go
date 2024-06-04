package psql

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gotest.tools/v3/assert"
)

func runDBTest(t *testing.T, ctx context.Context, f func(db *sql.DB)) {
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

	f(db)
}
