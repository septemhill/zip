package inf

import "context"

// TransactionHolderRepositoryCreator creates repositories which could hold the
// transaction by caller. Caller could commit/rollback after receiving some
// specified event, e.g. http request/response.
type TransactionHolderRepositoryCreator interface {
	NewOrderTxHolderRepository(context.Context) OrderTxHolderRepository
	NewMutantTxHolderRepository(context.Context) MutantTxHolderRepository
}

// TransactionHolderRepository
type TransactionHolderRepository interface {
	CommitableRepository
	RollbackableRepository
}

// CommitableRepository makes transaction holder repostiory could commit.
type CommitableRepository interface {
	Commit() error
}

// RollbackableRepository makes transaction holder repostiory could rollback.
type RollbackableRepository interface {
	Rollback() error
}
