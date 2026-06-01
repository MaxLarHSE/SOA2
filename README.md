# SOA2

REST API сервис для CRUD-операций с товарами на Go. API описан в `api.yaml`, база данных - PostgreSQL, миграции лежат в `migrations`.

## Endpoints

- `POST /products` - создать товар
- `GET /products` - получить список товаров
- `GET /products/{id}` - получить товар
- `PUT /products/{id}` - обновить товар
- `DELETE /products/{id}` - архивировать товар

## Запуск

```powershell
docker compose up --build
```

Сервис будет доступен на `http://localhost:8080`.

Для локальной сборки без Docker:

```powershell
make build
.\app
```

Перед локальным запуском нужен PostgreSQL и переменная `DATABASE_URL`, если база не запущена на дефолтном локальном адресе из `main.go`.
