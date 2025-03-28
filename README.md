# Тестовое задание

# User CRUD

Этот проект представляет собой REST API для управления пользователями, разработанный на Go с использованием фреймворка Echo и базы данных SQLite.

## Запуск

Для запуска проекта выполните следующие шаги:

1. Установите Go (если не установлен)
2. Клонируйте репозиторий:
   ```sh
   git clone https://github.com/Midichlorian007/portfolio.git
   cd portfolio
   ```
3. Установите зависимости:
   ```sh
   go mod tidy
   ```
4. Запустите сервер:
   ```sh
   go run ./cmd/app/main.go
   ```

По умолчанию сервер запустится на `http://localhost:9090`.

## Эндпоинты

| Метод  | Эндпоинт          | Описание                 |
|--------|-------------------|--------------------------|
| POST   | `/users`          | Создать пользователя     |
| GET    | `/users/:user_id` | Получить пользователя    |
| PUT    | `/users/:user_id` | Обновить пользователя    |
| DELETE | `/users/:user_id` | Удалить пользователя     |

## Вызов эндпоинтов в .http
В корне есть файл `.http` для выполнения запросов

