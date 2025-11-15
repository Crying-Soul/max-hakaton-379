# MaxBot Project

Проект включает в себя бота и веб-интерфейс для управления мероприятиями.

## Быстрый старт

### Предварительные требования

- Docker
- Docker Compose

### Установка и запуск

1. **Настройте переменные окружения**
   ```bash
Заполните .env файл:

env
TOKEN=your_bot_token_here
VITE_API_URL=http://localhost:8080
HTTP_ALLOWED_ORIGINS=http://localhost:5173
Запустите проект

bash
docker-compose up -d
Запустите тестовые данные

bash
chmod +x ./scripts/seed_mock_data.sh
./scripts/seed_mock_data.sh
Если скрипт не работает локально:

bash
docker cp ./scripts/seed_mock_data.sh maxbot_postgres:/tmp/
docker exec maxbot_postgres chmod +x /tmp/seed_mock_data.sh
docker exec -e DATABASE_URL='postgres://postgres:postgres@localhost:5432/maxbot?sslmode=disable' -e ENV_FILE='' maxbot_postgres /tmp/seed_mock_data.sh
Доступ к сервисам
Frontend: http://localhost:5173

Backend API: http://localhost:8080

PostgreSQL: localhost:5432

Управление
bash
# Запуск
docker-compose up -d

# Остановка
docker-compose down

# Перезапуск
docker-compose restart

# Логи
docker-compose logs -f bot
docker-compose logs -f frontend

# Статус
docker-compose ps
Тестовые данные
Скрипт scripts/seed_mock_data.sh заполняет базу:
Для перезаполнения данных просто запустите скрипт повторно.