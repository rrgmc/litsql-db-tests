package mysql

import (
	"context"
	"database/sql"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"gotest.tools/v3/assert"
)

func runDBTest(t *testing.T, ctx context.Context, f func(db *sql.DB)) {
	dbName := "sakila"
	dbUser := "root"
	dbPassword := "password"

	mysqlContainer, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("ghcr.io/rrgmc/litsql-dbtest-sakila-mysql:latest"),
		mysql.WithDatabase(dbName),
		mysql.WithUsername(dbUser),
		mysql.WithPassword(dbPassword),
	)
	assert.NilError(t, err)

	// Clean up the container
	defer func() {
		if err := mysqlContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	connStr, err := mysqlContainer.ConnectionString(ctx, "tls=skip-verify")
	assert.NilError(t, err)

	db, err := sql.Open("mysql", connStr)
	assert.NilError(t, err)

	f(db)
}
