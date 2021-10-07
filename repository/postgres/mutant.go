package postgres

import (
	"context"
	"database/sql"
	"log"
	"runtime"
	"sample/models"
	"sample/repository/inf"
	"sync"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type mutantTxHolderRepository struct {
	db       *sqlx.DB
	txHolder *sqlx.Tx
	mutex    sync.Mutex
	builder  sq.StatementBuilderType
	opts     options
}

func NewMutantRepotory(ctx context.Context, db *sqlx.DB, opts ...Option) *mutantTxHolderRepository {
	repo := &mutantTxHolderRepository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}

	for _, opt := range opts {
		opt.apply(&repo.opts)
	}

	if err := repo.applyOptions(ctx, repo.opts); err != nil {
		log.Fatal(err)
	}

	// Make sure the repository would release the transaction when garbage collection.
	runtime.SetFinalizer(repo, func(repo *mutantTxHolderRepository) { repo.Rollback() })
	return repo
}

func (repo *mutantTxHolderRepository) applyOptions(ctx context.Context, opts options) error {
	if opts.txHolder {
		tx, err := repo.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		if err != nil {
			return err
		}
		repo.txHolder = tx
	}
	return nil
}

func (repo *mutantTxHolderRepository) CreateMutant(ctx context.Context, m *models.Mutant) error {
	if repo.txHolder != nil {
		repo.mutex.Lock()
		defer repo.mutex.Unlock()
	}

	ext := repo.getSQLExecutor()

	stmt, vals, err := repo.builder.Insert("mutant_data").
		Columns("first_name", "last_name", "address", "location").
		Values(m.FirstName, m.LastName, m.Address, m.Location).ToSql()
	if err != nil {
		return err
	}

	if _, err := ext.ExecContext(ctx, stmt, vals...); err != nil {
		return err
	}

	return nil
}

func (repo *mutantTxHolderRepository) DeleteMutant(ctx context.Context, m *models.Mutant) error {
	if repo.txHolder != nil {
		repo.mutex.Lock()
		defer repo.mutex.Unlock()
	}

	ext := repo.getSQLExecutor()

	stmt, vals, err := repo.builder.Delete("mutant_data").
		Where(sq.Eq{"first_name": m.FirstName, "last_name": m.LastName}).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := ext.ExecContext(ctx, stmt, vals...); err != nil {
		return err
	}
	return nil
}

func (repo *mutantTxHolderRepository) UpdateMutant(ctx context.Context, m *models.Mutant) error {
	if repo.txHolder != nil {
		repo.mutex.Lock()
		defer repo.mutex.Unlock()
	}

	ext := repo.getSQLExecutor()

	stmt, vals, err := repo.builder.Update("mutant_data").
		SetMap(map[string]interface{}{"address": m.Address, "location": m.Location}).
		Where(sq.Eq{"first_name": m.FirstName, "last_name": m.LastName}).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := ext.ExecContext(ctx, stmt, vals...); err != nil {
		return err
	}
	return nil
}

func (repo *mutantTxHolderRepository) GetMutant(ctx context.Context, m *models.Mutant) ([]models.Mutant, error) {
	if repo.txHolder != nil {
		repo.mutex.Lock()
		defer repo.mutex.Unlock()
	}

	ext := repo.getSQLExecutor()

	stmt, vals, err := repo.builder.Select("mutant").
		Columns("*").
		ToSql()
	if err != nil {
		return nil, err
	}

	muts := make([]models.Mutant, 0)
	if err := ext.SelectContext(ctx, &muts, stmt, vals...); err != nil {
		return nil, err
	}

	return muts, nil
}

func (repo *mutantTxHolderRepository) Commit() error {
	return doCommit(repo.txHolder)
}

func (repo *mutantTxHolderRepository) Rollback() error {
	return doRollback(repo.txHolder)
}

func (repo *mutantTxHolderRepository) getSQLExecutor() sqlExecutor {
	if repo.txHolder != nil {
		return repo.txHolder
	}
	return repo.db
}

var _ inf.MutantTxHolderRepository = (*mutantTxHolderRepository)(nil)
