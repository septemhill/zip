package postgres

import (
	"context"
	"database/sql"
	"sample/repository/inf"
	"sync"

	"github.com/jmoiron/sqlx"
)

type orderTxHolderRepository struct {
	db       *sqlx.DB
	txHolder *sqlx.Tx
	mutex    sync.Mutex
}

func NewOrderRepository(ctx context.Context, db *sqlx.DB, opts ...Option) *orderTxHolderRepository {
	return &orderTxHolderRepository{
		db: db,
	}
}

func (repo *orderTxHolderRepository) CreateOrder(ctx context.Context) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx, err := repo.getTransaction(ctx)
	if err != nil {
		return err
	}

	_ = tx

	return nil
}

func (repo *orderTxHolderRepository) UpdateOrderStatus(ctx context.Context) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx, err := repo.getTransaction(ctx)
	if err != nil {
		return err
	}

	_ = tx

	return nil
}

func (repo *orderTxHolderRepository) CancelOrder(ctx context.Context) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx, err := repo.getTransaction(ctx)
	if err != nil {
		return err
	}

	_ = tx

	return nil
}

func (repo *orderTxHolderRepository) Commit() error {
	return doCommit(repo.txHolder)
}

func (repo *orderTxHolderRepository) Rollback() error {
	return doRollback(repo.txHolder)
}

func (repo *orderTxHolderRepository) getTransaction(ctx context.Context) (*sqlx.Tx, error) {
	if repo.txHolder != nil {
		return repo.txHolder, nil
	}

	tx, err := repo.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}

	repo.txHolder = tx
	return repo.txHolder, nil
}

var _ inf.OrderTxHolderRepository = (*orderTxHolderRepository)(nil)
