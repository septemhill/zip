package models

import "time"

type Order struct {
	OrderID   int
	Products  []Product
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Product struct {
	ProductID int
	Count     int
	Price     int
}
