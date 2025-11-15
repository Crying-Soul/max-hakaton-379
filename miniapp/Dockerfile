FROM node:20-alpine

WORKDIR /app

# Копируем файлы зависимостей
COPY package*.json ./

# Устанавливаем зависимости
RUN npm install

# Копируем исходный код
COPY . .

# Сборка приложения
RUN npm run build

# Устанавливаем serve для раздачи статики
RUN npm install -g serve

# Экспонируем порт
EXPOSE 5173

# Запускаем сервер для раздачи статики
CMD ["serve", "-s", "dist", "-l", "5173"]