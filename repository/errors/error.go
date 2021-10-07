package errors

import "fmt"

type repositoryError struct {
	Code    int
	Message string
}

func (err repositoryError) Error() string {
	return fmt.Sprintf("ERROR-[%05d]: %s", err.Code, err.Message)
}

var (
	// Cassandra
	// Postgres
	ErrPostgresEmptyTx = repositoryError{Code: 3, Message: "transaction holder is nil"}
)
