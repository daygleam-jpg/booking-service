# Hotel Booking Service

Учебный REST API-сервис для управления отелями, комнатами и бронированиями.

## Запуск проекта

### 1. Запуск PostgreSQL

В корневой директории проекта выполните:

```bash
sudo docker compose up -d
```

Параметры подключения:

- Database: `hotel_booking`
- User: `hotel_user`
- Password: `change_me_local`
- Host: `localhost`
- Port: `5432`

### 2. Применение миграций (UP)

Миграции выполняются вручную и строго по порядку.

**Создание таблицы отелей:**

```bash
sudo docker exec -i hotel-postgres psql -v ON_ERROR_STOP=1 -U hotel_user -d hotel_booking < migrations/001_create_hotels.up.sql
```

**Создание таблицы комнат:**

```bash
sudo docker exec -i hotel-postgres psql -v ON_ERROR_STOP=1 -U hotel_user -d hotel_booking < migrations/002_create_rooms.up.sql
```

**Создание таблицы бронирований:**

```bash
sudo docker exec -i hotel-postgres psql -v ON_ERROR_STOP=1 -U hotel_user -d hotel_booking < migrations/003_create_bookings.up.sql
```

Эти команды предназначены для первоначального создания таблиц в пустой базе данных. Не запускайте их повторно на уже созданной базе.

### 3. Откат миграций (DOWN)

Откат миграций выполняется строго в обратном порядке:

1. Сначала удаляется таблица бронирований.
2. Затем удаляется таблица комнат.
3. В последнюю очередь удаляется таблица отелей.

**Откат таблицы бронирований:**

```bash
sudo docker exec -i hotel-postgres psql -v ON_ERROR_STOP=1 -U hotel_user -d hotel_booking < migrations/003_create_bookings.down.sql
```

**Откат таблицы комнат:**

```bash
sudo docker exec -i hotel-postgres psql -v ON_ERROR_STOP=1 -U hotel_user -d hotel_booking < migrations/002_create_rooms.down.sql
```

**Откат таблицы отелей:**

```bash
sudo docker exec -i hotel-postgres psql -v ON_ERROR_STOP=1 -U hotel_user -d hotel_booking < migrations/001_create_hotels.down.sql
```

**Внимание:** DOWN-миграции удаляют таблицы вместе со всеми находящимися в них данными. Выполняйте их только при необходимости отката схемы базы данных.

После полного отката для повторного создания структуры необходимо выполнить все три UP-миграции.

### 4. Запуск приложения

В корневой директории проекта выполните:

```bash
go run ./cmd/api
```

Сервер будет доступен по адресу:

`http://localhost:8080`

## API Endpoints

### Hotels

| Method | Endpoint | Description |
|---|---|---|
| GET | `/hotels` | Получить все отели |
| GET | `/hotels/{id}` | Получить отель по ID |
| POST | `/hotels` | Создать отель |
| PUT | `/hotels/{id}` | Обновить отель |
| DELETE | `/hotels/{id}` | Удалить отель |

Пример создания отеля:

```bash
curl -i -X POST http://localhost:8080/hotels \
-H "Content-Type: application/json" \
-d '{
  "name": "Grand Hotel",
  "location": "Yerevan",
  "description": "Comfortable hotel"
}'
```

### Rooms

| Method | Endpoint | Description |
|---|---|---|
| GET | `/hotels/{id}/rooms` | Получить комнаты отеля |
| POST | `/hotels/{id}/rooms` | Создать комнату в отеле |

Пример создания комнаты:

```bash
curl -i -X POST http://localhost:8080/hotels/1/rooms \
-H "Content-Type: application/json" \
-d '{
  "number": "101"
}'
```

Возможные статусы комнаты:

- `available` — доступна.
- `booked` — забронирована.
- `maintenance` — находится на обслуживании.

Статус `booked` вычисляется на основании подтверждённых бронирований.

### Bookings

| Method | Endpoint | Description |
|---|---|---|
| POST | `/bookings` | Создать бронирование |
| DELETE | `/bookings/{id}` | Отменить бронирование |

Пример создания бронирования:

```bash
curl -i -X POST http://localhost:8080/bookings \
-H "Content-Type: application/json" \
-d '{
  "user_id": 1,
  "hotel_id": 1,
  "room_id": 1,
  "check_in": "2026-11-20",
  "check_out": "2026-11-22"
}'
```

При успешном создании сервер возвращает ID бронирования:

```json
{
  "id": 1,
  "message": "Booking created successfully"
}
```

Полученный ID используется для последующей отмены бронирования.

Пример отмены:

```bash
curl -i -X DELETE http://localhost:8080/bookings/1
```

Повторная отмена уже отменённого бронирования считается успешной операцией.

## Бизнес-правила бронирования

- Дата выезда должна быть позже даты заезда.
- Нельзя создавать бронирование на прошедшую дату.
- Комната не может иметь два подтверждённых бронирования на пересекающиеся даты.
- Последовательные бронирования разрешены: дата выезда одного гостя может совпадать с датой заезда следующего.
- Отменённые бронирования не блокируют доступность комнаты.
- Комнаты со статусом `maintenance` недоступны для бронирования.
- Бронирование может иметь статус `confirmed` или `cancelled`.
- Повторная отмена бронирования не считается ошибкой.
- Несуществующее бронирование при отмене возвращает `404 Not Found`.

## Проверка проекта

Запуск тестов:

```bash
go test ./...
```

Проверка сборки:

```bash
go build ./...
```

## Остановка PostgreSQL

Остановка контейнера:

```bash
sudo docker compose down
```

Для удаления контейнеров вместе с данными базы данных:

```bash
sudo docker compose down -v
```

**Внимание:** команда с `-v` удаляет Docker volume и все сохранённые данные PostgreSQL.