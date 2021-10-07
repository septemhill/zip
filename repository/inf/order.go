package inf

import "context"

// OrderRepository handles order object CRUD only.
type OrderRepository interface {
	CreateOrder(context.Context) error
	UpdateOrderStatus(context.Context) error
	CancelOrder(context.Context) error
}

// OrderTxHolderRepository handles order object CRUD and add Commitable and Rollbackable.
type OrderTxHolderRepository interface {
	OrderRepository
	TransactionHolderRepository
}
