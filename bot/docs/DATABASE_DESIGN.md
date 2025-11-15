# База данных: Структура и запросы

## 1. Анализ необходимых запросов

### 1.1 Запросы для Волонтёра

#### Профиль и аутентификация
- `GetUserByID(id)` - получить пользователя при каждом апдейте
- `GetVolunteerProfile(id)` - получить профиль с рейтингом, рангом
- `UpdateVolunteerProfile(id, resume, interests, radius)` - обновить профиль
- `GetVolunteerRank(social_rating)` - вычислить ранг по рейтингу
- `UpdateSocialRating(id, points)` - начислить очки

#### Поиск мероприятий
- `GetEventsNearby(lat, lon, radius, categories, limit)` - поиск мероприятий рядом
- `GetEventsByCategories(user_id, limit)` - по интересам пользователя
- `GetEventDetails(event_id)` - детали мероприятия
- `GetEventOrganizer(event_id)` - инфо об организаторе

#### Заявки волонтёра
- `CreateApplication(volunteer_id, event_id)` - подать заявку
- `GetMyApplications(volunteer_id, status?)` - все заявки (с фильтром по статусу)
- `GetApplicationDetails(application_id)` - детали заявки
- `CheckExistingApplication(volunteer_id, event_id)` - проверить, не подана ли уже

#### Уведомления
- `GetUpcomingEvents(volunteer_id, days)` - мероприятия через N дней
- `GetApprovedEvents(volunteer_id)` - одобренные заявки для уведомлений

### 1.2 Запросы для Организатора

#### Профиль и верификация
- `CreateOrganizatorApplication(user_id, org_name, inn, contacts)` - подать заявку
- `GetOrganizatorVerificationStatus(user_id)` - статус проверки
- `GetOrganizatorProfile(id)` - профиль организации
- `UpdateOrganizatorProfile(id, org_name, contacts)` - обновить профиль

#### Мероприятия
- `CreateEvent(organizator_id, event_data)` - создать мероприятие
- `GetMyEvents(organizator_id, status?)` - мои мероприятия (фильтр: активные/завершённые)
- `UpdateEvent(event_id, data)` - редактировать мероприятие
- `CompleteEvent(event_id)` - завершить мероприятие
- `GetEventStats(event_id)` - статистика (кол-во заявок, одобренных, участников)

#### Заявки на мероприятие
- `GetEventApplications(event_id, status?)` - заявки на мероприятие
- `GetApplicationsForSwipe(event_id, organizator_id, limit)` - заявки для дайвинчика
- `ApproveApplication(application_id)` - одобрить заявку
- `RejectApplication(application_id)` - отклонить заявку
- `GetVolunteerForApplication(volunteer_id)` - инфо о волонтёре в заявке

#### Чат и участники
- `AddVolunteerToChat(event_id, volunteer_id, chat_id)` - добавить в чат после одобрения
- `GetEventParticipants(event_id)` - список участников
- `GetChatMembers(event_id)` - участники чата

#### Репорты
- `CreateReport(event_id, volunteer_id, organizator_id, reason)` - создать репорт
- `UpdateVolunteerTrustRating(volunteer_id, multiplier)` - снизить рейтинг

### 1.3 Запросы для Администратора

#### Верификация организаторов
- `GetPendingOrganizators(limit, offset)` - заявки на проверку
- `GetOrganizatorApplicationDetails(organizator_id)` - детали заявки
- `ApproveOrganizator(organizator_id, admin_id)` - одобрить
- `RejectOrganizator(organizator_id, admin_id, reason)` - отклонить

#### Управление пользователями
- `GetAllUsers(role?, limit, offset)` - список пользователей (volunteer, organizator, admin)
- `GetUserDetails(user_id)` - детали пользователя
- `BlockUser(user_id)` - заблокировать
- `UnblockUser(user_id)` - разблокировать
- `GetUserActivity(user_id)` - активность пользователя

#### Репорты
- `GetAllReports(limit, offset)` - все репорты
- `GetReportDetails(report_id)` - детали репорта
- `GetReportsByVolunteer(volunteer_id)` - репорты на конкретного волонтёра
- `MarkReportAsReviewed(report_id)` - отметить как рассмотренный

#### Категории
- `GetAllCategories()` - все категории
- `CreateCategory(name, description)` - создать категорию
- `UpdateCategory(id, name, description)` - обновить
- `DeleteCategory(id)` - удалить

#### Статистика
- `GetTotalUsers(role?)` - общее кол-во пользователей
- `GetTotalEvents(status?)` - кол-во мероприятий
- `GetTotalApplications(status?)` - кол-во заявок
- `GetActiveVolunteers(days)` - активные волонтёры за период
- `GetTopVolunteers(limit)` - топ по рейтингу
- `GetTopOrganizators(limit)` - топ организаторов

---

## 2. Структура таблиц БД

### 2.1 Таблица `users` (основная)

```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY,                    -- Telegram ID
    username TEXT,                            -- Telegram username
    name TEXT NOT NULL,                       -- Имя пользователя
    role TEXT NOT NULL CHECK (role IN ('volunteer', 'organizator', 'admin')),
    state TEXT NOT NULL,                      -- Текущее состояние FSM
    is_blocked BOOLEAN DEFAULT FALSE,         -- Заблокирован ли пользователь
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_state ON users(state);
CREATE INDEX idx_users_blocked ON users(is_blocked);
```

**Запросы:**
```sql
-- Получить пользователя
SELECT * FROM users WHERE id = $1;

-- Обновить состояние FSM
UPDATE users SET state = $1, updated_at = NOW() WHERE id = $2;

-- Заблокировать пользователя
UPDATE users SET is_blocked = TRUE, updated_at = NOW() WHERE id = $1;
```

---

### 2.2 Таблица `volunteers` (профиль волонтёра)

```sql
CREATE TABLE volunteers (
    id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    resume TEXT,                              -- Резюме/описание навыков
    social_rating INT DEFAULT 0,              -- Социальный рейтинг (очки)
    trust_rating DECIMAL(3,2) DEFAULT 1.00,   -- Рейтинг добросовестности (0.00-1.00)
    rank TEXT DEFAULT 'newbie',               -- Ранг (newbie, experienced, pro, authority)
    search_radius INT DEFAULT 10,             -- Радиус поиска (км)
    interests TEXT[],                         -- Категории интересов (массив)
    dobro_token TEXT,                         -- Токен для интеграции с dobro.ru
    location_lat DECIMAL(10,8),               -- Последняя известная широта
    location_lon DECIMAL(11,8),               -- Последняя известная долгота
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_volunteers_rating ON volunteers(social_rating DESC);
CREATE INDEX idx_volunteers_location ON volunteers USING GIST (
    point(location_lon, location_lat)
);
```

**Запросы:**
```sql
-- Получить профиль волонтёра
SELECT u.*, v.* 
FROM users u 
JOIN volunteers v ON u.id = v.id 
WHERE u.id = $1;

-- Обновить профиль
UPDATE volunteers 
SET resume = $1, interests = $2, search_radius = $3, updated_at = NOW() 
WHERE id = $4;

-- Начислить очки и обновить ранг
UPDATE volunteers 
SET social_rating = social_rating + $1, 
    rank = CASE 
        WHEN social_rating + $1 >= 600 THEN 'authority'
        WHEN social_rating + $1 >= 300 THEN 'pro'
        WHEN social_rating + $1 >= 100 THEN 'experienced'
        ELSE 'newbie'
    END,
    updated_at = NOW()
WHERE id = $2;

-- Снизить рейтинг добросовестности
UPDATE volunteers 
SET trust_rating = trust_rating * $1, 
    updated_at = NOW() 
WHERE id = $2;

-- Топ волонтёров
SELECT u.name, u.username, v.social_rating, v.rank
FROM volunteers v
JOIN users u ON v.id = u.id
WHERE u.is_blocked = FALSE
ORDER BY v.social_rating DESC
LIMIT $1;
```

---

### 2.3 Таблица `organizators` (профиль организатора)

```sql
CREATE TABLE organizators (
    id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    organization_name TEXT NOT NULL,
    inn TEXT,                                 -- ИНН для верификации
    verification_status TEXT DEFAULT 'pending' 
        CHECK (verification_status IN ('pending', 'approved', 'rejected')),
    rejection_reason TEXT,                    -- Причина отклонения
    contacts TEXT,                            -- Контакты (JSON или текст)
    is_verified BOOLEAN DEFAULT FALSE,        -- Быстрая проверка верификации
    verified_at TIMESTAMP,                    -- Когда верифицирован
    verified_by BIGINT REFERENCES users(id),  -- Кто верифицировал (админ)
    events_count INT DEFAULT 0,               -- Кол-во созданных мероприятий
    completed_events_count INT DEFAULT 0,     -- Кол-во завершённых
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_organizators_status ON organizators(verification_status);
CREATE INDEX idx_organizators_verified ON organizators(is_verified);
```

**Запросы:**
```sql
-- Создать заявку организатора
INSERT INTO organizators (id, organization_name, inn, contacts)
VALUES ($1, $2, $3, $4);

-- Получить заявки на проверку
SELECT u.id, u.name, u.username, o.organization_name, o.inn, o.contacts, o.created_at
FROM organizators o
JOIN users u ON o.id = u.id
WHERE o.verification_status = 'pending'
ORDER BY o.created_at ASC
LIMIT $1 OFFSET $2;

-- Одобрить заявку
UPDATE organizators 
SET verification_status = 'approved', 
    is_verified = TRUE, 
    verified_at = NOW(), 
    verified_by = $2,
    updated_at = NOW()
WHERE id = $1;

-- Отклонить заявку
UPDATE organizators 
SET verification_status = 'rejected', 
    rejection_reason = $2,
    updated_at = NOW()
WHERE id = $1;

-- Получить статус верификации
SELECT verification_status, rejection_reason, is_verified
FROM organizators
WHERE id = $1;
```

---

### 2.4 Таблица `categories` (категории мероприятий)

```sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    icon TEXT,                                -- Эмодзи или иконка
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_categories_active ON categories(is_active);
```

**Запросы:**
```sql
-- Получить все активные категории
SELECT * FROM categories WHERE is_active = TRUE ORDER BY name;

-- Создать категорию
INSERT INTO categories (name, description, icon) VALUES ($1, $2, $3);

-- Обновить категорию
UPDATE categories SET name = $1, description = $2, icon = $3 WHERE id = $4;
```

---

### 2.5 Таблица `events` (мероприятия)

```sql
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    date TIMESTAMP NOT NULL,                  -- Дата и время начала
    duration INT,                             -- Длительность в часах
    location TEXT NOT NULL,                   -- Адрес
    location_lat DECIMAL(10,8) NOT NULL,      -- Широта
    location_lon DECIMAL(11,8) NOT NULL,      -- Долгота
    category_id INT REFERENCES categories(id),
    creator_id BIGINT REFERENCES organizators(id) ON DELETE CASCADE,
    contacts TEXT,                            -- Контакты для связи
    max_volunteers INT NOT NULL,              -- Макс. кол-во волонтёров
    current_volunteers INT DEFAULT 0,         -- Текущее кол-во одобренных
    reward_points INT NOT NULL,               -- Очки за участие
    status TEXT DEFAULT 'open' CHECK (status IN ('open', 'ongoing', 'completed', 'cancelled')),
    chat_id BIGINT,                           -- ID чата в Max
    chat_created_at TIMESTAMP,                -- Когда создан чат
    completed_at TIMESTAMP,                   -- Когда завершено
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_events_creator ON events(creator_id);
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_date ON events(date);
CREATE INDEX idx_events_category ON events(category_id);
CREATE INDEX idx_events_location ON events USING GIST (
    point(location_lon, location_lat)
);
```

**Запросы:**
```sql
-- Создать мероприятие
INSERT INTO events (
    title, description, date, duration, location, location_lat, location_lon,
    category_id, organizer_id, contacts, max_volunteers, reward_points
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING id;

-- Получить мероприятия рядом с волонтёром
SELECT e.*, c.name as category_name, o.organization_name, o.is_verified,
    earth_distance(
        ll_to_earth($1, $2),
        ll_to_earth(e.location_lat, e.location_lon)
    ) / 1000 AS distance_km
FROM events e
JOIN organizators o ON e.creator_id = o.id
LEFT JOIN categories c ON e.category_id = c.id
WHERE e.status = 'open'
  AND e.date > NOW()
  AND earth_distance(
        ll_to_earth($1, $2),
        ll_to_earth(e.location_lat, e.location_lon)
      ) <= $3 * 1000  -- радиус в метрах
ORDER BY 
    o.is_verified DESC,  -- Сначала верифицированные
    distance_km ASC,
    e.date ASC
LIMIT $4;

-- Получить мероприятия по категориям волонтёра
SELECT e.*, c.name as category_name, o.organization_name, o.is_verified
FROM events e
JOIN organizators o ON e.creator_id = o.id
LEFT JOIN categories c ON e.category_id = c.id
WHERE e.status = 'open'
  AND e.date > NOW()
  AND c.name = ANY($1)  -- массив интересов
ORDER BY 
    o.is_verified DESC,
    e.date ASC
LIMIT $2;

-- Мои мероприятия (организатор)
SELECT e.*, c.name as category_name,
    COUNT(a.id) FILTER (WHERE a.status = 'pending') as pending_count,
    COUNT(a.id) FILTER (WHERE a.status = 'approved') as approved_count
FROM events e
LEFT JOIN categories c ON e.category_id = c.id
LEFT JOIN applications a ON e.id = a.event_id
WHERE e.creator_id = $1
  AND ($2::TEXT IS NULL OR e.status = $2)
GROUP BY e.id, c.name
ORDER BY e.date DESC;

-- Завершить мероприятие
UPDATE events 
SET status = 'completed', completed_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- Увеличить счётчик одобренных волонтёров
UPDATE events 
SET current_volunteers = current_volunteers + 1, updated_at = NOW()
WHERE id = $1;
```

**Для геолокации нужно включить расширение:**
```sql
CREATE EXTENSION IF NOT EXISTS earthdistance CASCADE;
```

---

### 2.6 Таблица `applications` (заявки на участие)

```sql
CREATE TABLE applications (
    id SERIAL PRIMARY KEY,
    event_id INT REFERENCES events(id) ON DELETE CASCADE,
    volunteer_id BIGINT REFERENCES volunteers(id) ON DELETE CASCADE,
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    applied_at TIMESTAMP DEFAULT NOW(),
    reviewed_at TIMESTAMP,
    reviewed_by BIGINT REFERENCES organizators(id),
    
    UNIQUE(event_id, volunteer_id)  -- Нельзя подать заявку дважды
);

CREATE INDEX idx_applications_event ON applications(event_id);
CREATE INDEX idx_applications_volunteer ON applications(volunteer_id);
CREATE INDEX idx_applications_status ON applications(status);
```

**Запросы:**
```sql
-- Подать заявку
INSERT INTO applications (event_id, volunteer_id)
VALUES ($1, $2)
ON CONFLICT (event_id, volunteer_id) DO NOTHING
RETURNING id;

-- Мои заявки (волонтёр)
SELECT a.*, e.title, e.date, e.location, o.organization_name, c.name as category_name
FROM applications a
JOIN events e ON a.event_id = e.id
JOIN organizators o ON e.creator_id = o.id
LEFT JOIN categories c ON e.category_id = c.id
WHERE a.volunteer_id = $1
  AND ($2::TEXT IS NULL OR a.status = $2)
ORDER BY a.applied_at DESC;

-- Заявки на мероприятие (для организатора)
SELECT a.*, u.name, u.username, v.resume, v.social_rating, v.rank, v.trust_rating
FROM applications a
JOIN volunteers v ON a.volunteer_id = v.id
JOIN users u ON v.id = u.id
WHERE a.event_id = $1
  AND ($2::TEXT IS NULL OR a.status = $2)
ORDER BY 
    v.social_rating DESC,  -- Сначала с высоким рейтингом
    a.applied_at ASC;

-- Заявки для дайвинчика (случайный порядок, только pending)
SELECT a.*, u.name, u.username, v.resume, v.social_rating, v.rank, v.trust_rating
FROM applications a
JOIN volunteers v ON a.volunteer_id = v.id
JOIN users u ON v.id = u.id
WHERE a.event_id = $1
  AND a.status = 'pending'
ORDER BY RANDOM()
LIMIT $2;

-- Одобрить заявку
UPDATE applications 
SET status = 'approved', reviewed_at = NOW(), reviewed_by = $2
WHERE id = $1;

-- Отклонить заявку
UPDATE applications 
SET status = 'rejected', reviewed_at = NOW(), reviewed_by = $2
WHERE id = $1;

-- Проверить существование заявки
SELECT EXISTS(
    SELECT 1 FROM applications 
    WHERE event_id = $1 AND volunteer_id = $2
);
```

---

### 2.7 Таблица `event_participants` (участники мероприятия)

```sql
CREATE TABLE event_participants (
    id SERIAL PRIMARY KEY,
    event_id INT REFERENCES events(id) ON DELETE CASCADE,
    volunteer_id BIGINT REFERENCES volunteers(id) ON DELETE CASCADE,
    application_id INT REFERENCES applications(id),
    joined_chat_at TIMESTAMP DEFAULT NOW(),
    participated BOOLEAN DEFAULT FALSE,       -- Участвовал ли в итоге
    points_awarded INT,                       -- Начислено очков
    awarded_at TIMESTAMP,
    
    UNIQUE(event_id, volunteer_id)
);

CREATE INDEX idx_participants_event ON event_participants(event_id);
CREATE INDEX idx_participants_volunteer ON event_participants(volunteer_id);
```

**Запросы:**
```sql
-- Добавить участника в чат
INSERT INTO event_participants (event_id, volunteer_id, application_id)
VALUES ($1, $2, $3)
ON CONFLICT (event_id, volunteer_id) DO NOTHING;

-- Получить участников мероприятия
SELECT ep.*, u.name, u.username, v.social_rating, v.rank
FROM event_participants ep
JOIN volunteers v ON ep.volunteer_id = v.id
JOIN users u ON v.id = u.id
WHERE ep.event_id = $1
ORDER BY ep.joined_chat_at;

-- Начислить очки участникам
UPDATE event_participants 
SET participated = TRUE, points_awarded = $2, awarded_at = NOW()
WHERE event_id = $1 AND volunteer_id = $3;
```

---

### 2.8 Таблица `reports` (репорты на волонтёров)

```sql
CREATE TABLE reports (
    id SERIAL PRIMARY KEY,
    event_id INT REFERENCES events(id) ON DELETE CASCADE,
    volunteer_id BIGINT REFERENCES volunteers(id) ON DELETE CASCADE,
    organizator_id BIGINT REFERENCES organizators(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'reviewed', 'dismissed')),
    admin_comment TEXT,
    reviewed_by BIGINT REFERENCES users(id),
    reviewed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_reports_volunteer ON reports(volunteer_id);
CREATE INDEX idx_reports_organizator ON reports(organizator_id);
CREATE INDEX idx_reports_status ON reports(status);
```

**Запросы:**
```sql
-- Создать репорт
INSERT INTO reports (event_id, volunteer_id, organizator_id, reason)
VALUES ($1, $2, $3, $4);

-- Получить все репорты
SELECT r.*, 
    u_vol.name as volunteer_name, u_vol.username as volunteer_username,
    o.organization_name,
    e.title as event_title
FROM reports r
JOIN volunteers v ON r.volunteer_id = v.id
JOIN users u_vol ON v.id = u_vol.id
JOIN organizators o ON r.organizator_id = o.id
JOIN events e ON r.event_id = e.id
WHERE ($1::TEXT IS NULL OR r.status = $1)
ORDER BY r.created_at DESC
LIMIT $2 OFFSET $3;

-- Репорты на конкретного волонтёра
SELECT r.*, o.organization_name, e.title
FROM reports r
JOIN organizators o ON r.organizator_id = o.id
JOIN events e ON r.event_id = e.id
WHERE r.volunteer_id = $1
ORDER BY r.created_at DESC;

-- Отметить как рассмотренный
UPDATE reports 
SET status = 'reviewed', admin_comment = $2, reviewed_by = $3, reviewed_at = NOW()
WHERE id = $1;
```

---

### 2.9 Таблица `notifications` (уведомления)

```sql
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    type TEXT NOT NULL,  -- 'event_reminder', 'application_approved', 'event_completed', etc.
    title TEXT,
    message TEXT NOT NULL,
    related_event_id INT REFERENCES events(id) ON DELETE CASCADE,
    related_application_id INT REFERENCES applications(id) ON DELETE CASCADE,
    is_sent BOOLEAN DEFAULT FALSE,
    sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_notifications_user ON notifications(user_id);
CREATE INDEX idx_notifications_sent ON notifications(is_sent);
CREATE INDEX idx_notifications_type ON notifications(type);
```

**Запросы:**
```sql
-- Создать уведомление
INSERT INTO notifications (user_id, type, title, message, related_event_id)
VALUES ($1, $2, $3, $4, $5);

-- Получить несент уведомления для отправки
SELECT * FROM notifications 
WHERE is_sent = FALSE 
ORDER BY created_at ASC
LIMIT 100;

-- Отметить как отправленное
UPDATE notifications SET is_sent = TRUE, sent_at = NOW() WHERE id = $1;
```

---

## 3. Дополнительные индексы и оптимизации

```sql
-- Composite индексы для частых запросов
CREATE INDEX idx_events_status_date ON events(status, date);
CREATE INDEX idx_applications_volunteer_status ON applications(volunteer_id, status);
CREATE INDEX idx_applications_event_status ON applications(event_id, status);

-- Partial индексы
CREATE INDEX idx_events_open ON events(date) WHERE status = 'open';
CREATE INDEX idx_applications_pending ON applications(applied_at) WHERE status = 'pending';
CREATE INDEX idx_issuers_pending ON issuers(created_at) WHERE verification_status = 'pending';
```

---

## 4. Триггеры для автоматизации

```sql
-- Автоматическое обновление updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    
CREATE TRIGGER update_volunteers_updated_at BEFORE UPDATE ON volunteers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    
CREATE TRIGGER update_organizators_updated_at BEFORE UPDATE ON organizators
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    
CREATE TRIGGER update_events_updated_at BEFORE UPDATE ON events
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Счётчик мероприятий организатора
CREATE OR REPLACE FUNCTION update_organizator_events_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE organizators SET events_count = events_count + 1 WHERE id = NEW.creator_id;
    ELSIF TG_OP = 'UPDATE' AND OLD.status != 'completed' AND NEW.status = 'completed' THEN
        UPDATE organizators SET completed_events_count = completed_events_count + 1 WHERE id = NEW.creator_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER organizator_events_counter AFTER INSERT OR UPDATE ON events
    FOR EACH ROW EXECUTE FUNCTION update_organizator_events_count();
```

---

## 5. Миграции в правильном порядке

```sql
-- 1. Расширения
CREATE EXTENSION IF NOT EXISTS earthdistance CASCADE;

-- 2. Основные таблицы
CREATE TABLE users (...);
CREATE TABLE volunteers (...);
CREATE TABLE organizators (...);
CREATE TABLE categories (...);

-- 3. Таблицы с зависимостями
CREATE TABLE events (...);
CREATE TABLE applications (...);
CREATE TABLE event_participants (...);
CREATE TABLE reports (...);
CREATE TABLE notifications (...);

-- 4. Индексы
CREATE INDEX ...;

-- 5. Триггеры и функции
CREATE FUNCTION ...;
CREATE TRIGGER ...;
```

---

## 6. Seed данные для разработки

```sql
-- Категории по умолчанию
INSERT INTO categories (name, description, icon) VALUES
    ('Экология', 'Уборка территорий, посадка деревьев', '🌱'),
    ('Помощь пожилым', 'Поддержка пенсионеров и ветеранов', '👴'),
    ('Благоустройство', 'Ремонт, покраска, уборка', '🏗️'),
    ('Образование', 'Репетиторство, мастер-классы', '📚'),
    ('Животные', 'Помощь приютам, выгул собак', '🐕'),
    ('Спорт', 'Организация спортивных мероприятий', '⚽'),
    ('Культура', 'Помощь в организации мероприятий', '🎭'),
    ('Здоровье', 'Донорство, медицинская помощь', '💊');

-- Тестовый админ (ID нужно будет заменить реальным)
-- UPDATE users SET role = 'admin' WHERE id = YOUR_TELEGRAM_ID;
```

---

## 7. Запросы для статистики (Admin)

```sql
-- Общая статистика
SELECT 
    (SELECT COUNT(*) FROM users) as total_users,
    (SELECT COUNT(*) FROM users WHERE role = 'volunteer') as volunteers,
    (SELECT COUNT(*) FROM users WHERE role = 'organizator') as organizators,
    (SELECT COUNT(*) FROM events WHERE status = 'open') as open_events,
    (SELECT COUNT(*) FROM events WHERE status = 'completed') as completed_events,
    (SELECT COUNT(*) FROM applications WHERE status = 'pending') as pending_applications;

-- Активные волонтёры за последние N дней
SELECT COUNT(DISTINCT volunteer_id) 
FROM applications 
WHERE applied_at > NOW() - INTERVAL '$1 days';

-- Топ организаторов по завершённым мероприятиям
SELECT u.name, o.organization_name, o.completed_events_count, o.is_verified
FROM organizators o
JOIN users u ON o.id = u.id
WHERE o.is_verified = TRUE
ORDER BY o.completed_events_count DESC
LIMIT 10;
```

---

## Резюме

### Основные таблицы: 9
1. `users` - все пользователи
2. `volunteers` - профили волонтёров
3. `organizators` - профили организаторов
4. `categories` - категории мероприятий
5. `events` - мероприятия
6. `applications` - заявки на участие
7. `event_participants` - участники чата/мероприятия
8. `reports` - репорты на волонтёров
9. `notifications` - система уведомлений

### Ключевые фичи:
- ✅ Геолокация с earthdistance
- ✅ Полнотекстовый поиск (можно добавить)
- ✅ Автоматические триггеры для счётчиков
- ✅ Индексы для оптимизации запросов
- ✅ Unique constraints для предотвращения дубликатов
- ✅ Cascade delete для целостности данных
