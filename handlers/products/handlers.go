package products

import (
	"SOA2/handlers/errorhandler"
	"SOA2/pkg/api"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/oapi-codegen/runtime/types"
)

type ProductServer struct {
	Db *sql.DB
}

func (s *ProductServer) HandleCreateProducts(w http.ResponseWriter, r *http.Request) {
	var req api.ProductCreate
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorhandler.WriteError(w, http.StatusBadRequest, api.VALIDATIONERROR, "json decoding error", nil)
		return
	}
	var resp api.ProductResponse
	query := `
		INSERT INTO products (name, description, price, stock, category, status) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id, created_at, updated_at
	`
	err = s.Db.QueryRowContext(
		r.Context(),
		query,
		req.Name,
		req.Description,
		req.Price,
		req.Stock,
		req.Category,
		req.Status,
	).Scan(&resp.Id, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		log.Printf("Ошибка при вставке в БД: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	resp.Name = req.Name
	resp.Description = req.Description
	resp.Price = req.Price
	resp.Stock = req.Stock
	resp.Category = req.Category
	resp.Status = req.Status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Ошибка отправки ответа: %v", err)
	}
}

// GET /products
func (s *ProductServer) HandleGetProducts(w http.ResponseWriter, r *http.Request, params api.HandleGetProductsParams) {
	page := 0
	if params.Page != nil {
		page = *params.Page
	}
	size := 20
	if params.Size != nil {
		size = *params.Size
	}
	offset := page * size

	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argId := 1

	if params.Status != nil {
		whereClause += fmt.Sprintf(" AND status = $%d", argId)
		args = append(args, string(*params.Status))
		argId++
	}
	if params.Category != nil {
		whereClause += fmt.Sprintf(" AND category = $%d", argId)
		args = append(args, *params.Category)
		argId++
	}

	// 3. Узнаем общее количество подходящих товаров (для поля totalElements)
	countQuery := "SELECT COUNT(*) FROM products " + whereClause
	var totalElements int
	if err := s.Db.QueryRowContext(r.Context(), countQuery, args...).Scan(&totalElements); err != nil {
		log.Printf("Ошибка подсчета товаров: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	query := fmt.Sprintf(`
		SELECT id, name, description, price, stock, category, status, created_at, updated_at 
		FROM products 
		%s 
		ORDER BY created_at DESC 
		LIMIT $%d OFFSET $%d
	`, whereClause, argId, argId+1)

	args = append(args, size, offset)

	rows, err := s.Db.QueryContext(r.Context(), query, args...)
	if err != nil {
		log.Printf("Ошибка запроса списка товаров: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := make([]api.ProductResponse, 0)
	for rows.Next() {
		var p api.ProductResponse
		var desc sql.NullString

		if err := rows.Scan(&p.Id, &p.Name, &desc, &p.Price, &p.Stock, &p.Category, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			log.Printf("Ошибка чтения строки: %v", err)
			continue
		}
		if desc.Valid {
			p.Description = &desc.String
		}
		items = append(items, p)
	}

	resp := api.ProductListResponse{
		Items:         items,
		TotalElements: totalElements,
		Page:          page,
		Size:          size,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GET /products/{id}
func (s *ProductServer) HandleGetOneProduct(w http.ResponseWriter, r *http.Request, id types.UUID) {
	var resp api.ProductResponse
	var desc sql.NullString // Специальный тип для nullable полей

	query := `
		SELECT id, name, description, price, stock, category, status, created_at, updated_at 
		FROM products 
		WHERE id = $1
	`

	err := s.Db.QueryRowContext(r.Context(), query, id.String()).Scan(
		&resp.Id, &resp.Name, &desc, &resp.Price, &resp.Stock,
		&resp.Category, &resp.Status, &resp.CreatedAt, &resp.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		details := make(map[string]interface{})
		details["id"] = id.String()
		errorhandler.WriteError(w, http.StatusNotFound, api.PRODUCTNOTFOUND, "product not found", &details)
		return
	} else if err != nil {
		log.Printf("Ошибка при получении товара: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	if desc.Valid {
		resp.Description = &desc.String
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// PUT /products/{id}
func (s *ProductServer) HandleChangeOneProduct(w http.ResponseWriter, r *http.Request, id types.UUID) {
	var req api.ProductUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorhandler.WriteError(w, http.StatusBadRequest, api.VALIDATIONERROR, "json decoding error", nil)
		return
	}

	var resp api.ProductResponse
	var desc sql.NullString

	query := `
		UPDATE products 
		SET name = $1, description = $2, price = $3, stock = $4, category = $5, status = $6
		WHERE id = $7
		RETURNING id, name, description, price, stock, category, status, created_at, updated_at
	`

	err := s.Db.QueryRowContext(r.Context(), query,
		req.Name, req.Description, req.Price, req.Stock, req.Category, req.Status, id.String(),
	).Scan(
		&resp.Id, &resp.Name, &desc, &resp.Price, &resp.Stock,
		&resp.Category, &resp.Status, &resp.CreatedAt, &resp.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		details := make(map[string]interface{})
		details["id"] = id.String()
		errorhandler.WriteError(w, http.StatusNotFound, api.PRODUCTNOTFOUND, "product not found", &details)
		return
	} else if err != nil {
		log.Printf("Ошибка при обновлении товара: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	if desc.Valid {
		resp.Description = &desc.String
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DELETE /products/{id}
func (s *ProductServer) HandleDeleteOneProduct(w http.ResponseWriter, r *http.Request, id types.UUID) {
	query := `
		UPDATE products
		SET status = 'ARCHIVED'
		WHERE id = $1
	`

	result, err := s.Db.ExecContext(r.Context(), query, id.String())
	if err != nil {
		log.Printf("DELETE (soft) error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		details := make(map[string]interface{})
		details["id"] = id.String()
		errorhandler.WriteError(w, http.StatusNotFound, api.PRODUCTNOTFOUND, "product not found", &details)
		return
	}

	w.WriteHeader(http.StatusOK)
}
