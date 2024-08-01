package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBops interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
}

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, fx func(context.Context) error) error
	GetQueryEngine(ctx context.Context) DBops
}

type Database struct {
	pool *pgxpool.Pool
}

func NewDatabase(pool *pgxpool.Pool) *Database {
	return &Database{pool: pool}
}

func NewPool(dbUrl string) (*pgxpool.Pool, error) {
	return pgxpool.Connect(context.Background(), dbUrl)
}

type contextKey struct{}

// key is context key for storing transaction context
var key = &contextKey{}

// RunRepeatableRead runs  function within a repeatable read transaction
// If the function returns error, the transaction is rolled back
// Otherwise, the transaction is committed
func (d Database) RunRepeatableRead(ctx context.Context, fx func(context.Context) error) error {
	conn, err := d.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	txCtx := context.WithValue(ctx, key, tx)
	err = fx(txCtx)
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			return fmt.Errorf("transaction failed: %w, rollback also failed: %v", err, rollbackErr)
		}
		return fmt.Errorf("transaction failed and was rolled back: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		fmt.Printf("Commit failed with error: %v\n", err)
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			return fmt.Errorf("commit failed: %w, rollback also failed: %v", err, rollbackErr)
		}
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}

// GetQueryEngine returns DBops for current context
// If context contains transaction, it returns transaction as DBops
// Otherwise, it returns connection pool as DBops
func (d Database) GetQueryEngine(ctx context.Context) DBops {
	if tx, ok := ctx.Value(key).(pgx.Tx); ok && tx != nil {
		return tx
	}
	return d.pool
}
