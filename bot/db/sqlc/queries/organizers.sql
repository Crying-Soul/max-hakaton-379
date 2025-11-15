-- name: CreateOrganizer :one
INSERT INTO organizers (
    id,
    organization_name,
    verification_status,
    rejection_reason,
    contacts,
    verified_at,
    verified_by
) VALUES (
    sqlc.arg(id),
    sqlc.arg(organization_name),
    COALESCE(sqlc.arg(verification_status), 'pending'),
    sqlc.arg(rejection_reason),
    sqlc.arg(contacts),
    sqlc.arg(verified_at),
    sqlc.arg(verified_by)
)
RETURNING *;

-- name: UpsertOrganizer :one
INSERT INTO organizers (
    id,
    organization_name,
    verification_status,
    rejection_reason,
    contacts,
    verified_at,
    verified_by
) VALUES (
    sqlc.arg(id),
    sqlc.arg(organization_name),
    COALESCE(sqlc.arg(verification_status), 'pending'),
    sqlc.arg(rejection_reason),
    sqlc.arg(contacts),
    sqlc.arg(verified_at),
    sqlc.arg(verified_by)
)
ON CONFLICT (id) DO UPDATE
SET
    organization_name = EXCLUDED.organization_name,
    verification_status = EXCLUDED.verification_status,
    rejection_reason = EXCLUDED.rejection_reason,
    contacts = EXCLUDED.contacts,
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
    contacts = sqlc.arg(contacts),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: SetOrganizerVerification :one
UPDATE organizers
SET
    verification_status = sqlc.arg(verification_status),
    rejection_reason = sqlc.arg(rejection_reason),
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

-- name: ListOrganizersByStatus :many
SELECT *
FROM organizers
WHERE verification_status = sqlc.arg(status)
ORDER BY updated_at DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListPendingOrganizersWithUsers :many
SELECT o.*, u.username, u.name, u.role, u.state
FROM organizers o
JOIN users u ON u.id = o.id
WHERE o.verification_status = 'pending'
ORDER BY o.created_at ASC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;
