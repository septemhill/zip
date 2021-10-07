package inf

import (
	"context"
	"sample/models"
)

// MutantRepository handles mutant object CRUD only.
type MutantRepository interface {
	CreateMutant(context.Context, *models.Mutant) error
	DeleteMutant(context.Context, *models.Mutant) error
	UpdateMutant(context.Context, *models.Mutant) error
	GetMutant(context.Context, *models.Mutant) ([]models.Mutant, error)
}

// MutantTxHolderRepository handles mutant object CRUD and add Commitable and Rollbackable.
type MutantTxHolderRepository interface {
	MutantRepository
	TransactionHolderRepository
}
