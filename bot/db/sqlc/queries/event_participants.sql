-- name: AddEventParticipant :one
INSERT INTO event_participants (
    event_id,
    volunteer_id,
    application_id,
    joined_chat_at
) VALUES (
    sqlc.arg(event_id),
    sqlc.arg(volunteer_id),
    sqlc.arg(application_id),
    COALESCE(sqlc.arg(joined_chat_at), NOW())
)
RETURNING *;

-- name: RemoveEventParticipant :exec
DELETE FROM event_participants
WHERE event_id = sqlc.arg(event_id)
  AND volunteer_id = sqlc.arg(volunteer_id);

-- name: DeleteParticipantsByEvent :exec
DELETE FROM event_participants
WHERE event_id = sqlc.arg(event_id);

-- name: GetEventParticipant :one
SELECT *
FROM event_participants
WHERE event_id = sqlc.arg(event_id)
  AND volunteer_id = sqlc.arg(volunteer_id);

-- name: ListEventParticipants :many
SELECT *
FROM event_participants
WHERE event_id = sqlc.arg(event_id)
ORDER BY joined_chat_at DESC, id DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListEventParticipantsWithUsers :many
SELECT ep.*, u.username, u.name, u.state
FROM event_participants ep
JOIN users u ON u.id = ep.volunteer_id
WHERE ep.event_id = sqlc.arg(event_id)
ORDER BY ep.joined_chat_at DESC, ep.id DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListParticipantEvents :many
SELECT ep.*
FROM event_participants ep
WHERE ep.volunteer_id = sqlc.arg(volunteer_id)
ORDER BY ep.joined_chat_at DESC, ep.id DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: CountParticipantsForEvent :one
SELECT COUNT(*) AS count
FROM event_participants
WHERE event_id = sqlc.arg(event_id);
