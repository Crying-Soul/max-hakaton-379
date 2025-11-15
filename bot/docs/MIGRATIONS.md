# Миграции базы данных# Миграции БД



## Обзор## Обзор



Проект использует [golang-migrate](https://github.com/golang-migrate/migrate) для управления миграциями PostgreSQL.Проект использует [golang-migrate](https://github.com/golang-migrate/migrate) для управления миграциями базы данных. Миграции автоматически применяются при создании нового `UserRepository`.



Миграции находятся в `db/migrations/` и автоматически выполняются при старте приложения.## Расположение миграций



## Список миграцийВсе миграции находятся в папке `db/migrations/`.



### 20251110000001 - Базовые таблицы## Формат файлов миграций

**Файл:** `create_base_tables.up.sql`

Каждая миграция состоит из двух файлов:

Создаёт:- `{timestamp}_{description}.up.sql` - применение миграции

- Расширение `earthdistance` для геолокации- `{timestamp}_{description}.down.sql` - откат миграции

- Таблицу `users` - все пользователи бота (id, username, name, role, state, is_blocked)

- Таблицу `volunteers` - профили волонтёров (resume, ratings, rank, interests, location)Пример:

- Таблицу `organizators` - профили организаторов (organization_name, inn, verification_status)- `20251109210535_create_users_table.up.sql`

- Таблицу `categories` - категории мероприятий- `20251109210535_create_users_table.down.sql`

- Seed данные: 8 категорий по умолчанию

## Создание новой миграции

### 20251110000002 - Мероприятия и заявки

**Файл:** `create_events_and_applications.up.sql`### Вручную



Создаёт:1. Создайте два файла с одинаковым timestamp и описанием:

- Таблицу `events` - мероприятия```bash

- Таблицу `applications` - заявки волонтёров# Unix timestamp + описание

- Таблицу `event_participants` - участники чатаtouch db/migrations/$(date +%Y%m%d%H%M%S)_your_migration_name.up.sql

touch db/migrations/$(date +%Y%m%d%H%M%S)_your_migration_name.down.sql

### 20251110000003 - Репорты и уведомления```

**Файл:** `create_reports_and_notifications.up.sql`

2. Добавьте SQL код в `.up.sql` файл:

Создаёт:```sql

- Таблицу `reports` - репорты на волонтёров-- Создание новой таблицы

- Таблицу `notifications` - очередь уведомленийCREATE TABLE IF NOT EXISTS "public"."your_table" (

  "id" bigserial PRIMARY KEY,

### 20251110000004 - Триггеры  "name" text NOT NULL

**Файл:** `create_triggers.up.sql`);

```

Создаёт:

- Автообновление `updated_at`3. Добавьте откат в `.down.sql` файл:

- Счётчик мероприятий организатора```sql

-- Удаление таблицы

## Автоматический запускDROP TABLE IF EXISTS "public"."your_table";

```

Миграции выполняются при создании репозитория:

### Используя migrate CLI (опционально)

```go

repo, err := repository.NewUserRepository(dsn)Установите migrate CLI:

``````bash

go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

## SQLC интеграция```



После изменения миграций регенерируйте код:Создайте новую миграцию:

```bash

```bashmigrate create -ext sql -dir db/migrations -seq your_migration_name

go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate```

```

## Применение миграций

Миграции применяются автоматически при старте приложения в функции `NewUserRepository()`.

При создании репозитория:
```go
repo, err := repository.NewUserRepository(ctx, databaseURL)
// Миграции уже применены
```

## Откат миграций

Для отката миграций вручную можно использовать migrate CLI:

```bash
# Откат одной миграции
migrate -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" \
        -path db/migrations down 1

# Откат всех миграций
migrate -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" \
        -path db/migrations down -all
```

## Docker

При использовании Docker контейнеров, папка `db/migrations` автоматически копируется в образ (см. `Dockerfile`).

## Troubleshooting

### Ошибка "Dirty database version"

Если миграция прервалась с ошибкой, база может остаться в "dirty" состоянии. Исправьте проблему вручную:

```bash
migrate -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" \
        -path db/migrations force VERSION
```

где `VERSION` - это номер версии миграции.

### Миграции не применяются

Проверьте:
1. Правильность формата имен файлов (`.up.sql` и `.down.sql`)
2. Наличие папки `db/migrations` относительно рабочей директории приложения
3. Права доступа к файлам миграций
4. Логи приложения на наличие ошибок
