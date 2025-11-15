-- name: CreateOrganizerVerificationRequest :one
INSERT INTO organizer_verification_requests (
    organizer_id,
    status,
    organizer_comment,
    admin_comment,
    reviewed_by,
    reviewed_at
) VALUES (
    sqlc.arg(organizer_id),
    COALESCE(sqlc.arg(status), 'pending'),
    sqlc.arg(organizer_comment),
    sqlc.arg(admin_comment),
    sqlc.arg(reviewed_by),
    sqlc.arg(reviewed_at)
)
RETURNING *;

-- name: ListOrganizerVerificationRequests :many
SELECT *
FROM organizer_verification_requests
WHERE organizer_id = sqlc.arg(organizer_id)
ORDER BY submitted_at DESC, id DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: GetLatestPendingOrganizerVerificationRequest :one
SELECT *
FROM organizer_verification_requests
WHERE organizer_id = sqlc.arg(organizer_id)
    AND status = 'pending'
ORDER BY submitted_at DESC, id DESC
LIMIT 1;

-- name: UpdateOrganizerVerificationRequestComment :one
UPDATE organizer_verification_requests
SET organizer_comment = sqlc.arg(organizer_comment),
        submitted_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;
