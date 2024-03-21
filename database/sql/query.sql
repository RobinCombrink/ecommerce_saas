-- name: GetProduct :one
SELECT * FROM Products
WHERE id = ?;

-- name: GetProducts :many
SELECT * FROM Products;

-- name: CreateProduct :one
INSERT INTO Products (name, description, price)
VALUES (?, ?, ?)
RETURNING id, name, description, price;