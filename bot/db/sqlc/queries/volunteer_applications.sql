-- name: CreateVolunteerApplication :one
INSERT INTO volunteer_applications (
    event_id,
    volunteer_id,
    status,
    rejection_reason,
    reviewed_by,
    reviewed_at
) VALUES (
    sqlc.arg(event_id),
    sqlc.arg(volunteer_id),
    COALESCE(sqlc.arg(status), 'pending'),
    sqlc.arg(rejection_reason),
    sqlc.arg(reviewed_by),
    sqlc.arg(reviewed_at)
)
RETURNING *;

-- name: UpsertVolunteerApplication :one
INSERT INTO volunteer_applications (
    event_id,
    volunteer_id,
    status,
    rejection_reason,
    reviewed_by,
    reviewed_at
) VALUES (
    sqlc.arg(event_id),
    sqlc.arg(volunteer_id),
    COALESCE(sqlc.arg(status), 'pending'),
    sqlc.arg(rejection_reason),
    sqlc.arg(reviewed_by),
    sqlc.arg(reviewed_at)
)
ON CONFLICT (event_id, volunteer_id) DO UPDATE
SET
    status = EXCLUDED.status,
    rejection_reason = EXCLUDED.rejection_reason,
    reviewed_by = EXCLUDED.reviewed_by,
    reviewed_at = EXCLUDED.reviewed_at
RETURNING *;

-- name: DeleteVolunteerApplication :exec
DELETE FROM volunteer_applications
WHERE id = sqlc.arg(id);

-- name: GetVolunteerApplicationByID :one
SELECT *
FROM volunteer_applications
WHERE id = sqlc.arg(id);

-- name: GetVolunteerApplication :one
SELECT *
FROM volunteer_applications
WHERE event_id = sqlc.arg(event_id)
  AND volunteer_id = sqlc.arg(volunteer_id);

-- name: ListApplicationsByVolunteer :many
SELECT *
FROM volunteer_applications
WHERE volunteer_id = sqlc.arg(volunteer_id)
ORDER BY applied_at DESC, id DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListApplicationsByEvent :many
SELECT *
FROM volunteer_applications
WHERE event_id = sqlc.arg(event_id)
ORDER BY applied_at DESC, id DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListPendingApplicationsByEvent :many
SELECT *
FROM volunteer_applications
WHERE event_id = sqlc.arg(event_id)
  AND status = 'pending'
ORDER BY applied_at ASC, id ASC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListApplicationsByStatus :many
SELECT *
FROM volunteer_applications
WHERE status = sqlc.arg(status)
ORDER BY applied_at DESC, id DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListApplicationsForOrganizer :many
SELECT va.*
FROM volunteer_applications va
JOIN events e ON e.id = va.event_id
WHERE e.organizer_id = sqlc.arg(organizer_id)
ORDER BY va.applied_at DESC, va.id DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: UpdateVolunteerApplicationStatus :one
UPDATE volunteer_applications
SET
    status = sqlc.arg(status),
    rejection_reason = sqlc.arg(rejection_reason),
    reviewed_by = sqlc.arg(reviewed_by),
    reviewed_at = COALESCE(sqlc.arg(reviewed_at), NOW())
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: ResetVolunteerApplicationReview :one
UPDATE volunteer_applications
SET
    status = 'pending',
    rejection_reason = NULL,
    reviewed_by = NULL,
    reviewed_at = NULL
WHERE id = sqlc.arg(id)
RETURNING *;
