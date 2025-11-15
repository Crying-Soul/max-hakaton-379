-- name: AddEventMedia :one
INSERT INTO event_media (
    event_id,
    token,
    uploaded_at,
    uploaded_by
) VALUES (
    sqlc.arg(event_id),
    sqlc.arg(token),
    COALESCE(sqlc.arg(uploaded_at), NOW()),
    sqlc.arg(uploaded_by)
)
RETURNING *;

-- name: DeleteEventMedia :exec
DELETE FROM event_media
WHERE id = sqlc.arg(id);

-- name: DeleteEventMediaByEvent :exec
DELETE FROM event_media
WHERE event_id = sqlc.arg(event_id);

-- name: GetEventMediaByID :one
SELECT *
FROM event_media
WHERE id = sqlc.arg(id);

-- name: GetEventMediaByToken :one
SELECT *
FROM event_media
WHERE token = sqlc.arg(token);

-- name: ListEventMedia :many
SELECT *
FROM event_media
WHERE event_id = sqlc.arg(event_id)
ORDER BY uploaded_at DESC, id DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListEventMediaByUploader :many
SELECT *
FROM event_media
WHERE uploaded_by = sqlc.arg(uploaded_by)
ORDER BY uploaded_at DESC, id DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;
