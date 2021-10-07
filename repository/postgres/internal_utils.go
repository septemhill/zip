package postgres

import (
	"sample/repository/errors"

	"github.com/jmoiron/sqlx"
)

func doCommit(tx *sqlx.Tx) error {
	return closeTx(tx, func(tx *sqlx.Tx) error { return tx.Commit() })
}

func doRollback(tx *sqlx.Tx) error {
	return closeTx(tx, func(tx *sqlx.Tx) error { return tx.Rollback() })
}

func closeTx(tx *sqlx.Tx, fn func(tx *sqlx.Tx) error) error {
	if tx != nil {
		defer func() {
			tx = nil
		}()
		return fn(tx)
	}
	return errors.ErrPostgresEmptyTx
}
