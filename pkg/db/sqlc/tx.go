package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Tx interface {
	ExecTx(ctx context.Context, fn func(*Queries) error) error
}

type tx struct {
	DB *sql.DB
}

func NewTx(DB *sql.DB) Tx {
	return &tx{DB: DB}
}

func (t *tx) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := t.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	if err = fn(q); err != nil {
		if rbErr := tx.Rollback(); err != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
