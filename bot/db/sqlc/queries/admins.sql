-- name: CreateAdmin :one
INSERT INTO admins (id)
VALUES (sqlc.arg(id))
RETURNING *;

-- name: DeleteAdmin :exec
DELETE FROM admins
WHERE id = sqlc.arg(id);

-- name: GetAdmin :one
SELECT *
FROM admins
WHERE id = sqlc.arg(id);

-- name: ListAdmins :many
SELECT *
FROM admins
ORDER BY created_at DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListAdminsWithUsers :many
SELECT a.*, u.username, u.name, u.role
FROM admins a
JOIN users u ON u.id = a.id
ORDER BY a.created_at DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;
