package main

import (
	"context"
	"fmt"
	"log"
	"sample/models"
	"sample/repository/inf"
	"sample/repository/postgres"
)

func MutantOperations(ctx context.Context, repo inf.Repository) {
	r := &models.Mutant{
		FirstName: "Hazel",
		LastName:  "Chen",
		Address:   "DingPu",
		Location:  "New Taipei",
	}

	if err := repo.CreateMutant(ctx, r); err != nil {
		log.Fatal("failed to create mutant: ", err)
	}

	r = &models.Mutant{
		FirstName: "Hazel",
		LastName:  "Chen",
	}

	if err := repo.DeleteMutant(ctx, r); err != nil {
		log.Fatal(err)
	}

	r = &models.Mutant{
		FirstName: "Septem",
		LastName:  "Li",
		Address:   "Address after update",
		Location:  "Location after update",
	}

	if err := repo.UpdateMutant(ctx, r); err != nil {
		log.Fatal(err)
	}

	r = &models.Mutant{
		FirstName: "Septem 71",
		LastName:  "Shen",
	}

	muts, err := repo.GetMutant(ctx, r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(muts)
}

func main() {
	ctx := context.Background()

	// Cassandra repository sample
	// repo := cassandra.NewCassandraRepository(
	// 	[]string{"127.0.0.1"},
	// 	cassandra.WithKeyspace("people"),
	// )

	// Postgres repository sample
	repo := postgres.NewPostgresRepository(
		postgres.WithUser("postgres"), postgres.WithPassword("gintamaed3op2"),
		postgres.WithDBName("postgres"), postgres.WithDisabledSSL())

	MutantOperations(ctx, repo)
}
