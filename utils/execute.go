package utils

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/rrgmc/litsql/sq"
	"gotest.tools/v3/assert"
)

func DBExecute(t *testing.T, db *sql.DB, query sq.BuildQuery, parseArgs map[string]any,
	f func(row map[string]any)) {
	squery, args, err := query.Build()
	assert.NilError(t, err)

	fmt.Println(squery)

	if parseArgs != nil {
		args, err = sq.ParseArgs(args, sq.MapArgValues{
			"length": 100,
		})
		assert.NilError(t, err)
	}

	rows, err := db.Query(squery, args...)
	assert.NilError(t, err)
	defer rows.Close()

	srows := sqlscan.NewRowScanner(rows)
	for rows.Next() {
		r := map[string]any{}
		err = srows.Scan(&r)
		assert.NilError(t, err)
		f(r)
	}
	assert.NilError(t, rows.Err())
}
