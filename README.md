markdown
# MaxBot Project

Проект включает в себя бота для волонтёров и организаторов, а также веб-интерфейс в виде карты с мероприятиями.

## Быстрый старт

### Предварительные требования

- Docker
- Docker Compose

### Установка и запуск

1. **Настройте переменные окружения**

   Создайте файл `.env` в корне проекта:
   ```bash
   cp env .env
   nano .env
Заполните файл следующими переменными:

env
TOKEN=your_bot_token_here
VITE_API_URL=http://localhost:8080
HTTP_ALLOWED_ORIGINS=http://localhost:5173
Запустите проект

```bash
docker-compose up -d
```
Запустите тестовые данные

Если скрипт не работает локально (требуется psql):

```bash
docker cp ./scripts/seed_mock_data.sh maxbot_postgres:/tmp/
docker exec maxbot_postgres chmod +x /tmp/seed_mock_data.sh
docker exec -e DATABASE_URL='postgres://postgres:postgres@localhost:5432/maxbot?sslmode=disable' -e ENV_FILE='' maxbot_postgres /tmp/seed_mock_data.sh
```
Доступ к сервисам
Frontend: http://localhost:5173

Backend API: http://localhost:8080

PostgreSQL: localhost:5432

Управление
Основные команды
bash
# Запуск всех сервисов
```bash
docker-compose up -d
```
# Остановка всех сервисов
```bash
docker-compose down
```
# Перезапуск сервисов
```bash
docker-compose restart
```
# Просмотр логов
```bash
docker-compose logs -f bot
docker-compose logs -f frontend
```
# Статус контейнеров
```bash
docker-compose ps
```


Моковые данные для фронтенда
Чтобы фронт быстро увидел карту, используйте bash-скрипт scripts/seed_mock_data.sh, который через psql создаёт категории, пользователей, организаторов, волонтёров и события вокруг Петербурга.

Локальный запуск
Запустите Postgres (например, docker compose up) и убедитесь, что DATABASE_URL в .env указывает на нужную БД.

Выполните:
```bash
./scripts/seed_mock_data.sh
```

Скрипт прочитает .env, очистит старые сиды через временную таблицу и зальёт актуальные данные. Повторный запуск безопасен — значения просто перезапишутся.

Запуск внутри контейнера Postgres
Если psql не установлен локально, можно прогнать скрипт прямо в maxbot_postgres:
```bash
docker cp ./scripts/seed_mock_data.sh maxbot_postgres:/tmp/seed_mock_data.sh
docker exec maxbot_postgres chmod +x /tmp/seed_mock_data.sh
docker exec \
  -e DATABASE_URL='postgres://postgres:postgres@localhost:5432/maxbot?sslmode=disable' \
  -e ENV_FILE='' \
  maxbot_postgres /tmp/seed_mock_data.sh
```
Важные параметры: DATABASE_URL обязателен даже внутри контейнера (используйте localhost:5432), а ENV_FILE='' отключает повторную загрузку .env.

Скрипт сам создаёт связи (медиа, участники, счётчики), поэтому его удобно запускать перед демонстрациями, чтобы откатить тестовую базу к известному состоянию.
