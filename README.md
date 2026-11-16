# CloudBox

REST API файлового хранилища в стиле Dropbox, написанный на Go. Файлы хранятся в S3-совместимом хранилище, метаданные и квоты — в PostgreSQL. Проект построен по принципам Clean Architecture.

## Возможности

- **Аутентификация:** JWT access-токены и refresh-токены с ротацией. Refresh-токены хранятся только в виде хэшей, пароли — в виде bcrypt-хэшей.
- **Папки:** создание, список, переименование, перемещение, удаление.
- **Файлы:** загрузка с расчётом SHA-256, скачивание, переименование, перемещение, удаление.
- **Квота хранилища:** индивидуальный лимит места для каждого пользователя, который обновляется атомарно в транзакции.
- **Корзина:** восстановление, окончательное удаление и автоматическая очистка элементов старше 30 дней.

## Технологии

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-008ECF?style=for-the-badge&logo=gin&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![S3](https://img.shields.io/badge/S3_API-569A31?style=for-the-badge&logo=amazons3&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-000000?style=for-the-badge&logo=jsonwebtokens&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![Git](https://img.shields.io/badge/Git-F05032?style=for-the-badge&logo=git&logoColor=white)

| Область | Технологии |
|---|---|
| Язык и HTTP | Go, Gin |
| Данные | PostgreSQL, pgx, sqlc, golang-migrate |
| Хранилище файлов | SeaweedFS (S3 API), AWS SDK for Go v2 |
| Безопасность | JWT, bcrypt |
| Инфраструктура | Docker, Docker Compose, Make |
| Прочее | Uber zap, cleanenv, Git |

## Архитектура

```
HTTP-слой (Gin)  ->  usecase (бизнес-логика)  ->  repository (PostgreSQL, S3)
```

Бизнес-логика зависит только от интерфейсов, поэтому базу данных и хранилище файлов можно заменить, не меняя остальной код.

## Быстрый старт

```bash
cp .env.example .env        # задайте JWT_SECRET
docker compose up -d --build
curl localhost:8080/health
```

Локальная разработка:

```bash
docker compose up -d postgres seaweedfs
make migrate-up
go run ./cmd/api
```

## API

Все маршруты находятся под `/api/v1`. Для маршрутов с пометкой **auth** нужен заголовок `Authorization: Bearer <access_token>`.

| Метод | Путь | Auth | Описание |
|---|---|---|---|
| POST | `/auth/register` | | Регистрация |
| POST | `/auth/login` | | Вход |
| POST | `/auth/refresh` | | Обновление токенов |
| POST | `/auth/logout` | | Выход |
| GET | `/me` | auth | Профиль и использование квоты |
| POST, GET | `/folders` | auth | Создать папку, список папок |
| GET, PATCH, DELETE | `/folders/:id` | auth | Получить, переименовать, удалить |
| PATCH | `/folders/:id/move` | auth | Переместить папку |
| POST, GET | `/files` | auth | Загрузить файл, список файлов |
| GET, PATCH, DELETE | `/files/:id` | auth | Метаданные, переименование, удаление |
| GET | `/files/:id/download` | auth | Скачать файл |
| PATCH | `/files/:id/move` | auth | Переместить файл |
| GET | `/trash` | auth | Содержимое корзины |
| POST | `/trash/files/:id/restore`, `/trash/folders/:id/restore` | auth | Восстановить |
| DELETE | `/trash/files/:id`, `/trash/folders/:id` | auth | Удалить окончательно |

## Структура проекта

```
cmd/api                 точка входа
database                миграции и запросы sqlc
internal/entity         доменные типы
internal/usecase        бизнес-логика
internal/repository     PostgreSQL и S3
internal/delivery/http  обработчики, DTO, middleware
internal/worker         очистка корзины
pkg                     логгер, хэширование, JWT
```