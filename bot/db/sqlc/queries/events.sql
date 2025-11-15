-- name: CreateEvent :one
INSERT INTO events (
    title,
    description,
    chat,
    date,
    duration_hours,
    location,
    location_lat,
    location_lon,
    category_id,
    organizer_id,
    contacts,
    max_volunteers,
    current_volunteers,
    status,
    cancelled_reason,
    completed_at
) VALUES (
    sqlc.arg(title),
    sqlc.arg(description),
    sqlc.arg(chat),
    sqlc.arg(date),
    sqlc.arg(duration_hours),
    sqlc.arg(location),
    sqlc.arg(location_lat),
    sqlc.arg(location_lon),
    sqlc.arg(category_id),
    sqlc.arg(organizer_id),
    sqlc.arg(contacts),
    sqlc.arg(max_volunteers),
    COALESCE(sqlc.arg(current_volunteers), 0),
    COALESCE(sqlc.arg(status), 'open'),
    sqlc.arg(cancelled_reason),
    sqlc.arg(completed_at)
)
RETURNING *;

-- name: UpdateEvent :one
UPDATE events
SET
    title = sqlc.arg(title),
    description = sqlc.arg(description),
    chat = sqlc.arg(chat),
    date = sqlc.arg(date),
    duration_hours = sqlc.arg(duration_hours),
    location = sqlc.arg(location),
    location_lat = sqlc.arg(location_lat),
    location_lon = sqlc.arg(location_lon),
    category_id = sqlc.arg(category_id),
    contacts = sqlc.arg(contacts),
    max_volunteers = sqlc.arg(max_volunteers),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateEventStatus :one
UPDATE events
SET
    status = sqlc.arg(status),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: CancelEvent :one
UPDATE events
SET
    status = 'cancelled',
    cancelled_reason = sqlc.arg(reason),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: CompleteEvent :one
UPDATE events
SET
    status = 'completed',
    completed_at = NOW(),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: IncrementEventVolunteers :one
UPDATE events
SET
    current_volunteers = current_volunteers + sqlc.arg(delta),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING current_volunteers;

-- name: SetEventVolunteerCounts :one
UPDATE events
SET
    current_volunteers = sqlc.arg(current_volunteers),
    max_volunteers = sqlc.arg(max_volunteers),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING current_volunteers, max_volunteers;

-- name: GetEventByID :one
SELECT *
FROM events
WHERE id = sqlc.arg(id);

-- name: GetEventWithOrganizer :one
SELECT e.*, o.organization_name, o.verification_status
FROM events e
JOIN organizers o ON o.id = e.organizer_id
WHERE e.id = sqlc.arg(id);

-- name: DeleteEvent :exec
DELETE FROM events
WHERE id = sqlc.arg(id);

-- name: ListEvents :many
SELECT *
FROM events
ORDER BY date DESC, id DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: CountEvents :one
SELECT COUNT(*) FROM events;

-- name: CountAvailableEventsForVolunteer :one
SELECT COUNT(*)
FROM events e
WHERE e.status = 'open'
  AND NOT EXISTS (
    SELECT 1
    FROM volunteer_applications va
    WHERE va.event_id = e.id
      AND va.volunteer_id = sqlc.arg(volunteer_id)
  )
  AND NOT EXISTS (
    SELECT 1
    FROM event_participants ep
    WHERE ep.event_id = e.id
      AND ep.volunteer_id = sqlc.arg(volunteer_id)
  );

-- name: ListAvailableEventsForVolunteer :many
SELECT *
FROM events e
WHERE e.status = 'open'
  AND NOT EXISTS (
    SELECT 1
    FROM volunteer_applications va
    WHERE va.event_id = e.id
      AND va.volunteer_id = sqlc.arg(volunteer_id)
  )
  AND NOT EXISTS (
    SELECT 1
    FROM event_participants ep
    WHERE ep.event_id = e.id
      AND ep.volunteer_id = sqlc.arg(volunteer_id)
  )
ORDER BY e.date DESC, e.id DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListEventsByOrganizer :many
SELECT *
FROM events
WHERE organizer_id = sqlc.arg(organizer_id)
ORDER BY date DESC, id DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListEventsByStatus :many
SELECT *
FROM events
WHERE status = sqlc.arg(status)
ORDER BY date DESC, id DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListUpcomingEvents :many
SELECT *
FROM events
WHERE status = 'open'
  AND date >= sqlc.arg(start_date)
ORDER BY date ASC, id ASC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListEventsByCategory :many
SELECT *
FROM events
WHERE category_id = sqlc.arg(category_id)
ORDER BY date DESC, id DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListEventsNearLocation :many
SELECT *
FROM events
WHERE status = 'open'
  AND date >= sqlc.arg(start_date)
  AND ABS(location_lat - sqlc.arg(target_lat)) <= sqlc.arg(lat_delta)
  AND ABS(location_lon - sqlc.arg(target_lon)) <= sqlc.arg(lon_delta)
ORDER BY POWER(location_lat - sqlc.arg(target_lat), 2) + POWER(location_lon - sqlc.arg(target_lon), 2)
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListEventsForVolunteer :many
SELECT e.*
FROM events e
JOIN event_participants ep ON ep.event_id = e.id
WHERE ep.volunteer_id = sqlc.arg(volunteer_id)
ORDER BY e.date DESC, e.id DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListEventsWithPendingApplications :many
SELECT e.*
FROM events e
WHERE EXISTS (
    SELECT 1
    FROM volunteer_applications va
    WHERE va.event_id = e.id
      AND va.status = 'pending'
)
ORDER BY e.updated_at DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListEventsForMap :many
WITH candidate AS (
    SELECT
        e.id,
        e.title,
        e.description,
        e.date,
        e.duration_hours,
        e.location,
        e.location_lat,
        e.location_lon,
        e.category_id,
        e.organizer_id,
        e.contacts,
        e.chat,
        e.max_volunteers,
        COALESCE(e.current_volunteers, 0) AS current_volunteers,
        e.status,
        e.cancelled_reason,
        e.completed_at,
        e.created_at,
        e.updated_at,
        c.name AS category_name,
        GREATEST(e.max_volunteers - COALESCE(e.current_volunteers, 0), 0)::int4 AS slots_left,
        CAST(6371 * acos(
            LEAST(
                1,
                GREATEST(
                    -1,
                    cos(radians(sqlc.arg('lat'))) * cos(radians(e.location_lat::float8)) *
                    cos(radians(e.location_lon::float8) - radians(sqlc.arg('lon'))) +
                    sin(radians(sqlc.arg('lat'))) * sin(radians(e.location_lat::float8))
                )
            )
        ) AS double precision) AS distance_km
    FROM events e
    LEFT JOIN categories c ON c.id = e.category_id
    WHERE e.status = 'open'
      AND e.location_lat IS NOT NULL
      AND e.location_lon IS NOT NULL
      AND (
            e.max_volunteers = 0 OR
            e.current_volunteers IS NULL OR
            e.current_volunteers < e.max_volunteers
      )
      AND (
        sqlc.narg('category_ids')::int[] IS NULL OR
        e.category_id = ANY(sqlc.narg('category_ids')::int[])
      )
      AND CAST(6371 * acos(
            LEAST(
                1,
                GREATEST(
                    -1,
                    cos(radians(sqlc.arg('lat'))) * cos(radians(e.location_lat::float8)) *
                    cos(radians(e.location_lon::float8) - radians(sqlc.arg('lon'))) +
                    sin(radians(sqlc.arg('lat'))) * sin(radians(e.location_lat::float8))
                )
            )
    ) AS double precision) <= sqlc.arg('radius_km')::double precision
)
SELECT
    id,
    title,
    description,
    date,
    duration_hours,
    location,
    location_lat,
    location_lon,
    category_id,
    organizer_id,
    contacts,
    chat,
    max_volunteers,
    current_volunteers,
    status,
    cancelled_reason,
    completed_at,
    created_at,
    updated_at,
    category_name,
    slots_left::int4 AS slots_left,
    distance_km::double precision AS distance_km
FROM candidate
ORDER BY distance_km ASC, date ASC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: ListEventsForMapByVolunteer :many
WITH user_apps AS (
    SELECT
        va.event_id,
        va.status
    FROM volunteer_applications va
    WHERE va.volunteer_id = sqlc.arg(volunteer_id)
),
candidate AS (
    SELECT
        e.id,
        e.title,
        e.description,
        e.date,
        e.duration_hours,
        e.location,
        e.location_lat,
        e.location_lon,
        e.category_id,
        e.organizer_id,
        e.contacts,
        e.chat,
        e.max_volunteers,
        COALESCE(e.current_volunteers, 0) AS current_volunteers,
        e.status,
        e.cancelled_reason,
        e.completed_at,
        e.created_at,
        e.updated_at,
        c.name AS category_name,
        GREATEST(e.max_volunteers - COALESCE(e.current_volunteers, 0), 0)::int4 AS slots_left,
        CAST(6371 * acos(
            LEAST(
                1,
                GREATEST(
                    -1,
                    cos(radians(sqlc.arg('lat'))) * cos(radians(e.location_lat::float8)) *
                    cos(radians(e.location_lon::float8) - radians(sqlc.arg('lon'))) +
                    sin(radians(sqlc.arg('lat'))) * sin(radians(e.location_lat::float8))
                )
            )
        ) AS double precision) AS distance_km,
        ua.status AS application_status
    FROM events e
    JOIN user_apps ua ON ua.event_id = e.id
    LEFT JOIN categories c ON c.id = e.category_id
    WHERE e.location_lat IS NOT NULL
      AND e.location_lon IS NOT NULL
      AND (
        sqlc.narg('category_ids')::int[] IS NULL OR
        e.category_id = ANY(sqlc.narg('category_ids')::int[])
      )
      AND CAST(6371 * acos(
            LEAST(
                1,
                GREATEST(
                    -1,
                    cos(radians(sqlc.arg('lat'))) * cos(radians(e.location_lat::float8)) *
                    cos(radians(e.location_lon::float8) - radians(sqlc.arg('lon'))) +
                    sin(radians(sqlc.arg('lat'))) * sin(radians(e.location_lat::float8))
                )
            )
    ) AS double precision) <= sqlc.arg('radius_km')::double precision
)
SELECT
    id,
    title,
    description,
    date,
    duration_hours,
    location,
    location_lat,
    location_lon,
    category_id,
    organizer_id,
    contacts,
    chat,
    max_volunteers,
    current_volunteers,
    status,
    cancelled_reason,
    completed_at,
    created_at,
    updated_at,
    category_name,
    slots_left::int4 AS slots_left,
    distance_km::double precision AS distance_km,
    application_status
FROM candidate
ORDER BY distance_km ASC, date ASC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;
