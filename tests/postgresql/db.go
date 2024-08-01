package postgresql

import (
	"context"
	"fmt"
	"os"
	"testing"

	"route/internal/app/repository/database"
)

type TDB struct {
	DB database.TransactionManager
}

func NewFromEnv() *TDB {

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Println("DATABASE_URL не задан")
		os.Exit(1)
	}

	pool, err := database.NewPool(dbURL)
	if err != nil {
		panic(err)
	}
	return &TDB{DB: database.NewDatabase(pool)}
}

func (d *TDB) SetUp(t *testing.T) {
	t.Helper()
	_, err := d.DB.GetQueryEngine(context.Background()).Exec(context.Background(), "TRUNCATE TABLE orders CASCADE")
	if err != nil {
		t.Fatalf("Не удалось очистить таблицу orders: %v", err)
	}
	_, err = d.DB.GetQueryEngine(context.Background()).Exec(context.Background(), "TRUNCATE TABLE packaging_types CASCADE")
	if err != nil {
		t.Fatalf("Не удалось очистить таблицу packaging_types: %v", err)
	}
}

func (d *TDB) TearDown(t *testing.T) {
	t.Helper()
}
