package inf

// Repository
type Repository interface {
	MutantRepository
	TransactionHolderRepositoryCreator
}
