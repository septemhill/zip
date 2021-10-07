package cassandra

import (
	"context"
	"log"
	"sample/models"
	"sample/repository/inf"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/table"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
)

type CassandraRepositoryOption func(*gocql.ClusterConfig)

type cassandraRepository struct {
	stmts statement
	sess  gocqlx.Session
}

func WithKeyspace(keyspace string) CassandraRepositoryOption {
	return func(config *gocql.ClusterConfig) {
		config.Keyspace = keyspace
	}
}

func NewCassandraRepository(hosts []string, opts ...CassandraRepositoryOption) *cassandraRepository {
	cluster := gocql.NewCluster(hosts...)
	for _, opt := range opts {
		opt(cluster)
	}

	sess, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		log.Fatal(err)
	}

	return &cassandraRepository{
		sess:  sess,
		stmts: buildMutantStatements(),
	}
}

func buildMutantStatements() statement {
	meta := table.Metadata{
		Name:    "mutant_data",
		Columns: []string{"first_name", "last_name", "address", "location"},
		PartKey: []string{"first_name", "last_name"},
	}
	tbl := table.New(meta)
	deleteStmt, deleteNames := tbl.Delete()
	insertStmt, insertNames := tbl.Insert()
	selectStmt, selectNames := qb.Select(meta.Name).Columns(meta.Columns...).Where(qb.Eq("first_name"), qb.Eq("last_name")).ToCql()
	updateStmt, updateNames := qb.Update(meta.Name).Set("address", "location").Where(qb.Eq("first_name"), qb.Eq("last_name")).ToCql()

	stmt := statement{
		del: query{stmt: deleteStmt, names: deleteNames},
		ins: query{stmt: insertStmt, names: insertNames},
		sel: query{stmt: selectStmt, names: selectNames},
		up:  query{stmt: updateStmt, names: updateNames},
	}

	return stmt
}

func (repo *cassandraRepository) CreateMutant(ctx context.Context, m *models.Mutant) error {
	return repo.sess.ContextQuery(ctx, repo.stmts.ins.stmt, repo.stmts.ins.names).BindStruct(m).ExecRelease()
}

func (repo *cassandraRepository) DeleteMutant(ctx context.Context, m *models.Mutant) error {
	return repo.sess.ContextQuery(ctx, repo.stmts.del.stmt, repo.stmts.del.names).BindStruct(m).ExecRelease()
}

func (repo *cassandraRepository) UpdateMutant(ctx context.Context, m *models.Mutant) error {
	return repo.sess.ContextQuery(ctx, repo.stmts.up.stmt, repo.stmts.up.names).BindStruct(m).ExecRelease()
}

func (repo *cassandraRepository) GetMutant(ctx context.Context, m *models.Mutant) ([]models.Mutant, error) {
	muts := make([]models.Mutant, 0)
	if err := repo.sess.ContextQuery(ctx, repo.stmts.sel.stmt, repo.stmts.sel.names).BindStruct(m).SelectRelease(&muts); err != nil {
		return nil, err
	}
	return muts, nil
}

func (repo *cassandraRepository) NewOrderTxHolderRepository(ctx context.Context) inf.OrderTxHolderRepository {
	return nil
}

func (repo *cassandraRepository) NewMutantTxHolderRepository(ctx context.Context) inf.MutantTxHolderRepository {
	return nil
}

var _ inf.Repository = (*cassandraRepository)(nil)
