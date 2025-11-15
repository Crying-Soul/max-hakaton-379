#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="${ENV_FILE:-${ROOT_DIR}/.env}"

if [[ -f "${ENV_FILE}" ]]; then
  echo "[seed] loading env from ${ENV_FILE}"
  set -a
  # shellcheck disable=SC1090
  source "${ENV_FILE}"
  set +a
fi

if [[ -z "${DATABASE_URL:-}" ]]; then
  echo "DATABASE_URL is not set. Provide it in the environment or via ENV_FILE" >&2
  exit 1
fi

PSQL_BIN="${PSQL_BIN:-psql}"

${PSQL_BIN} "${DATABASE_URL}" <<'SQL'
\set ON_ERROR_STOP on
BEGIN;

DROP TABLE IF EXISTS seed_event_ids;
CREATE TEMP TABLE seed_event_ids AS
SELECT id FROM events WHERE title IN (
    'Субботник в Парке 300-летия',
    'Восстанавливаем дворики Васильевского',
    'Уроки цифровой грамотности на Невском',
    'Помогаем приюту "Лучик"',
    'Экосплав "Нева без пластика"',
    'Кулинарный добробат в Невском районе',
    'Высаживаем деревья на Крестовском',
    'Наставники для подростков в Пушкине',
    'Тёплый вечер в доме престарелых',
    'Фестиваль добрых дел в Новой Голландии'
);

DELETE FROM event_media WHERE token LIKE 'seed:%' OR event_id IN (SELECT id FROM seed_event_ids);
DELETE FROM event_participants WHERE event_id IN (SELECT id FROM seed_event_ids);
DELETE FROM volunteer_applications WHERE event_id IN (SELECT id FROM seed_event_ids);
DELETE FROM events WHERE id IN (SELECT id FROM seed_event_ids);

INSERT INTO categories (name, description, is_active)
VALUES
    ('Экология', 'Субботники, сортировка отходов и благоустройство города', TRUE),
    ('Социальная помощь', 'Поддержка приютов, дом престарелых и точечная помощь людям', TRUE),
    ('Образование', 'Наставничество, обучение детей и взрослых цифровым навыкам', TRUE)
ON CONFLICT (name) DO UPDATE
SET description = EXCLUDED.description,
    is_active = EXCLUDED.is_active;

INSERT INTO users (id, username, name, role, state, is_blocked, location_lat, location_lon)
VALUES
    (1001, 'anna_helper', 'Анна Координатор', 'organizer', '0', FALSE, 55.751244, 37.618423),
    (1002, 'igor_curator', 'Игорь Куратор', 'organizer', '0', FALSE, 59.938955, 30.315644),
    (2001, 'maria_vol', 'Мария Волонтёр', 'volunteer', '0', FALSE, 55.744548, 37.605598),
    (2002, 'roman_help', 'Роман Доброволец', 'volunteer', '0', FALSE, 55.765057, 37.621389),
    (2003, 'svet_teacher', 'Светлана Наставник', 'volunteer', '0', FALSE, 59.93428, 30.335099)
ON CONFLICT (id) DO UPDATE
SET username = EXCLUDED.username,
    name = EXCLUDED.name,
    role = EXCLUDED.role,
    state = EXCLUDED.state,
    is_blocked = EXCLUDED.is_blocked,
    location_lat = EXCLUDED.location_lat,
    location_lon = EXCLUDED.location_lon,
    updated_at = NOW();

INSERT INTO organizers (id, organization_name, verification_status, rejection_reason, contacts, verified_at, verified_by)
VALUES
    (1001, 'Фонд Чистый Город', 'verified', NULL, '+7 (900) 123-45-67 — Анна', NOW(), NULL),
    (1002, 'Добровольцы Севера', 'verified', NULL, '+7 (921) 222-33-44 — Игорь', NOW(), NULL)
ON CONFLICT (id) DO UPDATE
SET organization_name = EXCLUDED.organization_name,
    verification_status = EXCLUDED.verification_status,
    rejection_reason = EXCLUDED.rejection_reason,
    contacts = EXCLUDED.contacts,
    verified_at = EXCLUDED.verified_at,
    updated_at = NOW();

INSERT INTO volunteers (id, cv, search_radius, category_ids)
VALUES
    (
        2001,
        'Организую эко-акции четвертый год',
        25,
        ARRAY[
            (SELECT id FROM categories WHERE name = 'Экология'),
            (SELECT id FROM categories WHERE name = 'Социальная помощь')
        ]
    ),
    (
        2002,
        'Люблю работать с людьми и помогать пожилым',
        30,
        ARRAY[
            (SELECT id FROM categories WHERE name = 'Социальная помощь')
        ]
    ),
    (
        2003,
        'Преподаю информатику и цифровую гигиену',
        15,
        ARRAY[
            (SELECT id FROM categories WHERE name = 'Образование')
        ]
    )
ON CONFLICT (id) DO UPDATE
SET cv = EXCLUDED.cv,
    search_radius = EXCLUDED.search_radius,
    category_ids = EXCLUDED.category_ids;

-- Event inserts
INSERT INTO events (title, description, date, duration_hours, location, location_lat, location_lon, category_id, organizer_id, contacts, max_volunteers, current_volunteers, status)
VALUES
    (
        'Субботник в Парке 300-летия',
        'Собираем мусор с побережья Финского залива, красим лавочки и обновляем велопарковки.',
        NOW() + INTERVAL '36 hours',
        4,
        'Санкт-Петербург, Парк 300-летия',
        60.000421,
        30.199738,
        (SELECT id FROM categories WHERE name = 'Экология'),
        1001,
        '+7 (900) 123-45-67 — Анна',
        60,
        0,
        'open'
    ),
    (
        'Восстанавливаем дворики Васильевского',
        'Перекрашиваем ограждения, высаживаем цветы и ставим новые кормушки во дворах-колодцах.',
        NOW() + INTERVAL '60 hours',
        5,
        'Санкт-Петербург, 7-я линия В.О., 34',
        59.943104,
        30.280602,
        (SELECT id FROM categories WHERE name = 'Экология'),
        1001,
        '+7 (900) 123-45-67 — Анна',
        35,
        0,
        'open'
    ),
    (
        'Уроки цифровой грамотности на Невском',
        'Обучаем пенсионеров пользоваться Госуслугами и мессенджерами, помогаем настроить смартфоны.',
        NOW() + INTERVAL '72 hours',
        3,
        'Санкт-Петербург, Невский проспект, 46',
        59.932166,
        30.347060,
        (SELECT id FROM categories WHERE name = 'Образование'),
        1002,
        '+7 (921) 222-33-44 — Игорь',
        20,
        0,
        'open'
    ),
    (
        'Помогаем приюту "Лучик"',
        'Готовим корм, выгуливаем собак и фотографируем хвостатых для соцсетей.',
        NOW() + INTERVAL '96 hours',
        5,
        'Санкт-Петербург, пр. Обуховской Обороны, 120',
        59.862047,
        30.455221,
        (SELECT id FROM categories WHERE name = 'Социальная помощь'),
        1002,
        '+7 (921) 222-33-44 — Игорь',
        25,
        0,
        'open'
    ),
    (
        'Экосплав "Нева без пластика"',
        'Собираем мусор с воды и набережных, сортируем пластик и металлы, учим жителей раздельному сбору.',
        NOW() + INTERVAL '120 hours',
        6,
        'Санкт-Петербург, набережная Обводного канала, 74',
        59.915021,
        30.339277,
        (SELECT id FROM categories WHERE name = 'Экология'),
        1002,
        '+7 (921) 222-33-44 — Игорь',
        30,
        0,
        'open'
    ),
    (
        'Кулинарный добробат в Невском районе',
        'Готовим горячие обеды для малообеспеченных семей и развозим по заявкам соцслужб.',
        NOW() + INTERVAL '144 hours',
        4,
        'Санкт-Петербург, ул. Бабушкина, 36',
        59.891044,
        30.450781,
        (SELECT id FROM categories WHERE name = 'Социальная помощь'),
        1002,
        '+7 (921) 222-33-44 — Игорь',
        20,
        0,
        'open'
    ),
    (
        'Высаживаем деревья на Крестовском',
        'Сажаем молодые клёны и обновляем мульчу вдоль беговых дорожек.',
        NOW() + INTERVAL '168 hours',
        4,
        'Санкт-Петербург, Крестовский остров, Морской проспект',
        59.971512,
        30.259327,
        (SELECT id FROM categories WHERE name = 'Экология'),
        1001,
        '+7 (900) 123-45-67 — Анна',
        50,
        0,
        'open'
    ),
    (
        'Наставники для подростков в Пушкине',
        'Помогаем ребятам из центра постинтернатного сопровождения с учебой и профориентацией.',
        NOW() + INTERVAL '192 hours',
        3,
        'Санкт-Петербург, Пушкин, Московская улица, 4',
        59.724523,
        30.403621,
        (SELECT id FROM categories WHERE name = 'Образование'),
        1002,
        '+7 (921) 222-33-44 — Игорь',
        18,
        0,
        'open'
    ),
    (
        'Тёплый вечер в доме престарелых',
        'Готовим концерт, беседуем и приносим подарочные наборы жильцам пансионата.',
        NOW() + INTERVAL '216 hours',
        4,
        'Санкт-Петербург, проспект Металлистов, 56',
        59.968274,
        30.421975,
        (SELECT id FROM categories WHERE name = 'Социальная помощь'),
        1002,
        '+7 (921) 222-33-44 — Игорь',
        30,
        0,
        'open'
    ),
    (
        'Фестиваль добрых дел в Новой Голландии',
        'Собираем благотворительные наборы, проводим мастер-классы и рассказываем гостям о волонтёрстве.',
        NOW() + INTERVAL '240 hours',
        6,
        'Санкт-Петербург, Новая Голландия, наб. Адмиралтейского канала, 2',
        59.927864,
        30.290632,
        (SELECT id FROM categories WHERE name = 'Социальная помощь'),
        1001,
        '+7 (900) 123-45-67 — Анна',
        80,
        0,
        'open'
    );

-- Event media
INSERT INTO event_media (event_id, token, uploaded_at)
SELECT e.id, token, NOW()
FROM events e
JOIN (VALUES
    ('Субботник в Парке 300-летия', 'seed:park300:1'),
    ('Субботник в Парке 300-летия', 'seed:park300:2'),
    ('Восстанавливаем дворики Васильевского', 'seed:vasileostrovsky:1'),
    ('Уроки цифровой грамотности на Невском', 'seed:nevsky-digital:1'),
    ('Помогаем приюту "Лучик"', 'seed:shelter-luchik:1'),
    ('Помогаем приюту "Лучик"', 'seed:shelter-luchik:2'),
    ('Экосплав "Нева без пластика"', 'seed:neva-cleanup:1'),
    ('Кулинарный добробат в Невском районе', 'seed:nevsky-kitchen:1'),
    ('Высаживаем деревья на Крестовском', 'seed:krestovsky-trees:1'),
    ('Наставники для подростков в Пушкине', 'seed:pushkin-mentors:1'),
    ('Тёплый вечер в доме престарелых', 'seed:metal-workers-home:1'),
    ('Фестиваль добрых дел в Новой Голландии', 'seed:new-holland-fest:1')
) AS media(title, token) ON e.title = media.title;

-- Event participants
INSERT INTO event_participants (event_id, volunteer_id, joined_chat_at)
SELECT e.id, v.volunteer_id, NOW()
FROM events e
JOIN (VALUES
    ('Субботник в Парке 300-летия', 2001),
    ('Восстанавливаем дворики Васильевского', 2001),
    ('Восстанавливаем дворики Васильевского', 2002),
    ('Уроки цифровой грамотности на Невском', 2003),
    ('Помогаем приюту "Лучик"', 2001),
    ('Помогаем приюту "Лучик"', 2002),
    ('Экосплав "Нева без пластика"', 2001),
    ('Экосплав "Нева без пластика"', 2002),
    ('Кулинарный добробат в Невском районе', 2002),
    ('Высаживаем деревья на Крестовском', 2001),
    ('Высаживаем деревья на Крестовском', 2002),
    ('Наставники для подростков в Пушкине', 2003),
    ('Тёплый вечер в доме престарелых', 2001),
    ('Тёплый вечер в доме престарелых', 2002),
    ('Фестиваль добрых дел в Новой Голландии', 2001),
    ('Фестиваль добрых дел в Новой Голландии', 2002),
    ('Фестиваль добрых дел в Новой Голландии', 2003)
) AS v(title, volunteer_id) ON e.title = v.title
ON CONFLICT (event_id, volunteer_id) DO NOTHING;

WITH counts AS (
    SELECT event_id, COUNT(*)::int AS total
    FROM event_participants
    WHERE event_id IN (SELECT id FROM events WHERE title IN (
        'Субботник в Парке 300-летия',
        'Восстанавливаем дворики Васильевского',
        'Уроки цифровой грамотности на Невском',
        'Помогаем приюту "Лучик"',
        'Экосплав "Нева без пластика"',
        'Кулинарный добробат в Невском районе',
        'Высаживаем деревья на Крестовском',
        'Наставники для подростков в Пушкине',
        'Тёплый вечер в доме престарелых',
        'Фестиваль добрых дел в Новой Голландии'
    ))
    GROUP BY event_id
)
UPDATE events e
SET current_volunteers = counts.total
FROM counts
WHERE e.id = counts.event_id;

COMMIT;
SQL

echo "[seed] mock data successfully inserted"
