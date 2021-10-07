package postgres

import (
	"context"
	"log"
	"sample/models"
	"sample/repository/inf"
	"strings"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct {
	// postgres repository internal information.
	db      *sqlx.DB
	builder sq.StatementBuilderType
	dsnOpts []string

	// sub-repositories.
	mutant inf.MutantRepository
}

type PostgresRepositoryOption func(*postgresRepository)

func WithUser(username string) PostgresRepositoryOption {
	return func(repo *postgresRepository) {
		repo.dsnOpts = append(repo.dsnOpts, "user="+username)
	}
}

func WithPassword(password string) PostgresRepositoryOption {
	return func(repo *postgresRepository) {
		repo.dsnOpts = append(repo.dsnOpts, "password="+password)
	}
}

func WithDBName(dbname string) PostgresRepositoryOption {
	return func(repo *postgresRepository) {
		repo.dsnOpts = append(repo.dsnOpts, "dbname="+dbname)
	}
}

func WithDisabledSSL() PostgresRepositoryOption {
	return func(repo *postgresRepository) {
		repo.dsnOpts = append(repo.dsnOpts, "sslmode=disable")
	}
}

func NewPostgresRepository(opts ...PostgresRepositoryOption) *postgresRepository {
	repo := &postgresRepository{}

	for _, opt := range opts {
		opt(repo)
	}

	db, err := sqlx.Connect("pgx", strings.Join(repo.dsnOpts, " "))
	if err != nil {
		log.Fatal("failed to connect to postgres: ", err)
	}

	return &postgresRepository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		mutant:  NewMutantRepotory(context.Background(), db),
	}
}

func (repo *postgresRepository) CreateMutant(ctx context.Context, m *models.Mutant) error {
	return repo.mutant.CreateMutant(ctx, m)
}

func (repo *postgresRepository) DeleteMutant(ctx context.Context, m *models.Mutant) error {
	return repo.mutant.DeleteMutant(ctx, m)
}

func (repo *postgresRepository) UpdateMutant(ctx context.Context, m *models.Mutant) error {
	return repo.mutant.UpdateMutant(ctx, m)
}

func (repo *postgresRepository) GetMutant(ctx context.Context, m *models.Mutant) ([]models.Mutant, error) {
	return repo.mutant.GetMutant(ctx, m)
}

func (repo *postgresRepository) NewMutantTxHolderRepository(ctx context.Context) inf.MutantTxHolderRepository {
	return NewMutantRepotory(ctx, repo.db, WithTxHolder(true))
}

func (repo *postgresRepository) NewOrderTxHolderRepository(ctx context.Context) inf.OrderTxHolderRepository {
	return NewOrderRepository(ctx, repo.db, WithTxHolder(true))
}

var _ inf.Repository = (*postgresRepository)(nil)
