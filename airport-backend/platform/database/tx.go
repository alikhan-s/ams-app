package database

import (
	"context"
	"database/sql"
	"fmt"
)

type txKey struct{}

// Executor defines the interface for database operations (satisfied by *sql.DB and *sql.Tx).
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// TxManager handles database transactions.
type TxManager interface {
	Run(ctx context.Context, fn func(ctx context.Context) error) error
}

type txManager struct {
	db *sql.DB
}

// NewTxManager creates a new TxManager.
func NewTxManager(db *sql.DB) TxManager {
	return &txManager{db: db}
}

// Run executes the given function within a database transaction.
func (tm *txManager) Run(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Inject transaction into context
	ctxWithTx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback failed: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetTx retrieves the transaction from the context if it exists.
// Returns Executor interface which can be *sql.Tx or nil.
// NOTE: For now, we return *sql.Tx specifically as requested by some consumers,
// but casting to Executor is also fine.
func GetTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}
