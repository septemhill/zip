package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type selectorContext interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type selector interface {
	Select(dest interface{}, query string, args ...interface{}) error
}

type getterContext interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type getter interface {
	Get(dest interface{}, query string, args ...interface{}) error
}

type sqlExecutor interface {
	sqlx.Ext
	sqlx.ExtContext
	selector
	selectorContext
	getter
	getterContext
}
