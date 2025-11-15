markdown
# MaxBot Project

Проект включает в себя бота и веб-интерфейс для управления мероприятиями.

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

bash
docker-compose up -d
Запустите тестовые данные

Для локального запуска:

bash
chmod +x ./scripts/seed_mock_data.sh
./scripts/seed_mock_data.sh
Если скрипт не работает локально (требуется psql):

bash
docker cp ./scripts/seed_mock_data.sh maxbot_postgres:/tmp/
docker exec maxbot_postgres chmod +x /tmp/seed_mock_data.sh
docker exec -e DATABASE_URL='postgres://postgres:postgres@localhost:5432/maxbot?sslmode=disable' -e ENV_FILE='' maxbot_postgres /tmp/seed_mock_data.sh
Доступ к сервисам
Frontend: http://localhost:5173

Backend API: http://localhost:8080

PostgreSQL: localhost:5432

Управление
Основные команды
bash
# Запуск всех сервисов
docker-compose up -d

# Остановка всех сервисов
docker-compose down

# Перезапуск сервисов
docker-compose restart

# Просмотр логов
docker-compose logs -f bot
docker-compose logs -f frontend

# Статус контейнеров
docker-compose ps
Мониторинг и отладка
bash
# Просмотр всех логов в реальном времени
docker-compose logs -f

# Проверка использования ресурсов
docker-compose stats

# Пересборка образов
docker-compose up -d --build

# Остановка с удалением volumes
docker-compose down -v
Тестовые данные
Скрипт scripts/seed_mock_data.sh заполняет базу данных тестовыми данными:

Категории мероприятий

Пользователей и организаторов

Волонтеров

События вокруг Санкт-Петербурга

Связи между сущностями

Для перезаполнения данных просто запустите скрипт повторно:

bash
./scripts/seed_mock_data.sh
База данных