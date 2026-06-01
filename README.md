# SOA2

Короткий CRUD-сервис для товаров на Go. API описан в `api.yaml`, база данных - PostgreSQL, миграции лежат в `migrations`.

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
