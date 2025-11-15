-- name: CreateVolunteer :one
INSERT INTO volunteers (
    id,
    cv,
    search_radius,
    category_ids
) VALUES (
    sqlc.arg(id),
    sqlc.arg(cv),
    COALESCE(sqlc.arg(search_radius), 10),
    sqlc.arg(category_ids)
)
RETURNING *;

-- name: UpsertVolunteer :one
INSERT INTO volunteers (
    id,
    cv,
    search_radius,
    category_ids
) VALUES (
    sqlc.arg(id),
    sqlc.arg(cv),
    COALESCE(sqlc.arg(search_radius), 10),
    sqlc.arg(category_ids)
)
ON CONFLICT (id) DO UPDATE
SET
    cv = EXCLUDED.cv,
    search_radius = EXCLUDED.search_radius,
    category_ids = EXCLUDED.category_ids
RETURNING *;

-- name: DeleteVolunteer :exec
DELETE FROM volunteers
WHERE id = sqlc.arg(id);

-- name: GetVolunteer :one
SELECT *
FROM volunteers
WHERE id = sqlc.arg(id);

-- name: GetVolunteerWithUser :one
SELECT v.*, u.username, u.name, u.role, u.state, u.location_lat, u.location_lon
FROM volunteers v
JOIN users u ON u.id = v.id
WHERE v.id = sqlc.arg(id);

-- name: UpdateVolunteerProfile :one
UPDATE volunteers
SET
    cv = sqlc.arg(cv),
    search_radius = sqlc.arg(search_radius)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateVolunteerCategories :one
UPDATE volunteers
SET
    category_ids = sqlc.arg(category_ids)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: ListVolunteers :many
SELECT *
FROM volunteers
ORDER BY id
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListVolunteersByCategory :many
SELECT *
FROM volunteers
WHERE category_ids && sqlc.arg(category_ids)
ORDER BY id
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListVolunteersWithUsers :many
SELECT v.*, u.username, u.name, u.role, u.state, u.location_lat, u.location_lon
FROM volunteers v
JOIN users u ON u.id = v.id
ORDER BY u.updated_at DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListVolunteersByCategoryWithUsers :many
SELECT v.*, u.username, u.name, u.role, u.state, u.location_lat, u.location_lon
FROM volunteers v
JOIN users u ON u.id = v.id
WHERE category_ids && sqlc.arg(category_ids)
ORDER BY u.updated_at DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListVolunteersNearLocation :many
SELECT v.*, u.username, u.name, u.role, u.state, u.location_lat, u.location_lon
FROM volunteers v
JOIN users u ON u.id = v.id
WHERE u.is_blocked = FALSE
  AND u.location_lat IS NOT NULL
  AND u.location_lon IS NOT NULL
  AND ABS(u.location_lat - sqlc.arg(target_lat)) <= sqlc.arg(lat_delta)
  AND ABS(u.location_lon - sqlc.arg(target_lon)) <= sqlc.arg(lon_delta)
  AND (v.search_radius IS NULL OR v.search_radius >= sqlc.arg(required_radius))
ORDER BY POWER(u.location_lat - sqlc.arg(target_lat), 2) + POWER(u.location_lon - sqlc.arg(target_lon), 2)
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListVolunteersByIDs :many
SELECT *
FROM volunteers
WHERE id = ANY(sqlc.arg(ids)::bigint[])
ORDER BY id;
