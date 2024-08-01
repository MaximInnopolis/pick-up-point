//go:build integration

package tests

import "route/tests/postgresql"

var (
	db *postgresql.TDB
)

func init() {
	db = postgresql.NewFromEnv()
}
