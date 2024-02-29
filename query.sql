-- name: GetProduct :one
SELECT * FROM Products
WHERE id = ?;

-- name: GetProducts :many
SELECT * FROM Products;
