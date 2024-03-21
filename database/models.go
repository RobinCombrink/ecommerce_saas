// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import ()

type Customer struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Order struct {
	ID         int64   `json:"id"`
	Paid       int64   `json:"paid"`
	Customerid int64   `json:"customerid"`
	Total      float64 `json:"total"`
}

type OrderItem struct {
	ID        int64 `json:"id"`
	Orderid   int64 `json:"orderid"`
	Productid int64 `json:"productid"`
}

type Product struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Price       float64 `json:"price"`
}
