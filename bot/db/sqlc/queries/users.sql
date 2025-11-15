-- name: CreateUser :one
INSERT INTO users (
    id,
    username,
    name,
    role,
    state,
    is_blocked,
    location_lat,
    location_lon
) VALUES (
    sqlc.arg(id),
    sqlc.arg(username),
    sqlc.arg(name),
    sqlc.arg(role),
    sqlc.arg(state),
    COALESCE(sqlc.arg(is_blocked), false),
    sqlc.arg(location_lat),
    sqlc.arg(location_lon)
)
RETURNING *;

-- name: UpsertUser :one
INSERT INTO users (
    id,
    username,
    name,
    role,
    state,
    is_blocked,
    location_lat,
    location_lon
) VALUES (
    sqlc.arg(id),
    sqlc.arg(username),
    sqlc.arg(name),
    sqlc.arg(role),
    sqlc.arg(state),
    COALESCE(sqlc.arg(is_blocked), false),
    sqlc.arg(location_lat),
    sqlc.arg(location_lon)
)
ON CONFLICT (id) DO UPDATE
SET
    username = EXCLUDED.username,
    name = EXCLUDED.name,
    role = EXCLUDED.role,
    state = EXCLUDED.state,
    is_blocked = EXCLUDED.is_blocked,
    location_lat = EXCLUDED.location_lat,
    location_lon = EXCLUDED.location_lon,
    updated_at = NOW()
RETURNING *;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = sqlc.arg(id);

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = sqlc.arg(username);

-- name: ListUsersByRole :many
SELECT *
FROM users
WHERE role = sqlc.arg(role)
ORDER BY created_at DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListUsersByState :many
SELECT *
FROM users
WHERE state = sqlc.arg(state)
ORDER BY updated_at DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: SearchUsers :many
SELECT *
FROM users
WHERE (username ILIKE CONCAT('%', sqlc.arg(query), '%') OR name ILIKE CONCAT('%', sqlc.arg(query), '%'))
ORDER BY updated_at DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListBlockedUsers :many
SELECT *
FROM users
WHERE is_blocked = TRUE
ORDER BY updated_at DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: UpdateUserProfile :one
UPDATE users
SET
    username = sqlc.arg(username),
    name = sqlc.arg(name),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateUserRole :one
UPDATE users
SET
    role = sqlc.arg(role),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateUserState :one
UPDATE users
SET
    state = sqlc.arg(state),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateUserLocation :one
UPDATE users
SET
    location_lat = sqlc.arg(location_lat),
    location_lon = sqlc.arg(location_lon),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: BlockUser :exec
UPDATE users
SET
    is_blocked = TRUE,
    updated_at = NOW()
WHERE id = sqlc.arg(id);

-- name: UnblockUser :exec
UPDATE users
SET
    is_blocked = FALSE,
    updated_at = NOW()
WHERE id = sqlc.arg(id);

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = sqlc.arg(id);

-- name: ListUsersNearLocation :many
SELECT *
FROM users
WHERE is_blocked = FALSE
  AND location_lat IS NOT NULL
  AND location_lon IS NOT NULL
  AND ABS(location_lat - sqlc.arg(target_lat)) <= sqlc.arg(lat_delta)
  AND ABS(location_lon - sqlc.arg(target_lon)) <= sqlc.arg(lon_delta)
ORDER BY POWER(location_lat - sqlc.arg(target_lat), 2) + POWER(location_lon - sqlc.arg(target_lon), 2)
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListUsersByIDs :many
SELECT *
FROM users
WHERE id = ANY(sqlc.arg(ids)::bigint[])
ORDER BY updated_at DESC;
