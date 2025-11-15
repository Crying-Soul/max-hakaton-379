-- name: CreateOrganizer :one
INSERT INTO organizers (
    id,
    organization_name,
    about,
    verified_at,
    verified_by
) VALUES (
    sqlc.arg(id),
    sqlc.arg(organization_name),
    sqlc.arg(about),
    sqlc.arg(verified_at),
    sqlc.arg(verified_by)
)
RETURNING *;

-- name: UpsertOrganizer :one
INSERT INTO organizers (
    id,
    organization_name,
    about,
    verified_at,
    verified_by
) VALUES (
    sqlc.arg(id),
    sqlc.arg(organization_name),
    sqlc.arg(about),
    sqlc.arg(verified_at),
    sqlc.arg(verified_by)
)
ON CONFLICT (id) DO UPDATE
SET
    organization_name = EXCLUDED.organization_name,
    about = EXCLUDED.about,
    verified_at = EXCLUDED.verified_at,
    verified_by = EXCLUDED.verified_by,
    updated_at = NOW()
RETURNING *;

-- name: DeleteOrganizer :exec
DELETE FROM organizers
WHERE id = sqlc.arg(id);

-- name: GetOrganizer :one
SELECT *
FROM organizers
WHERE id = sqlc.arg(id);

-- name: GetOrganizerWithUser :one
SELECT o.*, u.username, u.name, u.role, u.state
FROM organizers o
JOIN users u ON u.id = o.id
WHERE o.id = sqlc.arg(id);

-- name: UpdateOrganizerProfile :one
UPDATE organizers
SET
    organization_name = sqlc.arg(organization_name),
    about = sqlc.arg(about),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: SetOrganizerVerification :one
UPDATE organizers
SET
    verified_at = sqlc.arg(verified_at),
    verified_by = sqlc.arg(verified_by),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: ListOrganizers :many
SELECT *
FROM organizers
ORDER BY created_at DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListVerifiedOrganizers :many
SELECT *
FROM organizers
WHERE verified_at IS NOT NULL
ORDER BY verified_at DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListUnverifiedOrganizers :many
SELECT *
FROM organizers
WHERE verified_at IS NULL
ORDER BY created_at ASC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;
