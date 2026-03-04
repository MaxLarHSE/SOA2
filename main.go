package main

//go:generate oapi-codegen -generate types,chi-server -package api -o api/api.gen.go api.yaml

import (
	"SOA2/handlers"
	"SOA2/handlers/errorhandler"
	"SOA2/handlers/products"
	"SOA2/pkg/api"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	_ "github.com/lib/pq"
)

func main() {
	log.SetFlags(0)
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://myuser:mypassword@localhost:5432/marketplace?sslmode=disable"
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка при создании пула соединений: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("База данных недоступна: %v", err)
	}
	log.Println("Успешное подключение к PostgreSQL!")

	server := &products.ProductServer{
		Db: db,
	}

	r := chi.NewRouter()
	r.Use(handlers.RequestMiddlware)
	r.Use(handlers.NewValidationMiddleware())

	api.HandlerWithOptions(server, api.ChiServerOptions{
		BaseRouter:       r,
		ErrorHandlerFunc: errorhandler.CustomErrorHandler,
	})
	log.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
