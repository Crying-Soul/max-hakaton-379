-- name: CreateCategory :one
INSERT INTO categories (
    name,
    description,
    is_active
) VALUES (
    sqlc.arg(name),
    sqlc.arg(description),
    COALESCE(sqlc.arg(is_active), TRUE)
)
RETURNING *;

-- name: GetCategory :one
SELECT *
FROM categories
WHERE id = sqlc.arg(id);

-- name: GetCategoryByName :one
SELECT *
FROM categories
WHERE name = sqlc.arg(name);

-- name: UpdateCategory :one
UPDATE categories
SET
    name = sqlc.arg(name),
    description = sqlc.arg(description)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: SetCategoryActive :one
UPDATE categories
SET
    is_active = sqlc.arg(is_active)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: ListCategories :many
SELECT *
FROM categories
ORDER BY name
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListActiveCategories :many
SELECT *
FROM categories
WHERE is_active = TRUE
ORDER BY name
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: SearchCategories :many
SELECT *
FROM categories
WHERE name ILIKE CONCAT('%', sqlc.arg(query), '%')
   OR description ILIKE CONCAT('%', sqlc.arg(query), '%')
ORDER BY name
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: CountActiveCategories :one
SELECT COUNT(*)
FROM categories
WHERE is_active = TRUE;
