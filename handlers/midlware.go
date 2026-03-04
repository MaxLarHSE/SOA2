package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type LogRequest struct {
	Request_id  uuid.UUID  `json:"request_id"`
	Method      string     `json:"method"`
	Endpoint    string     `json:"endpoint"`
	Status_code int        `json:"status_code"`
	DurationMs  int64      `json:"duration_ms"`
	Timestamp   time.Time  `json:"timestamp"`
	UserID      *uuid.UUID `json:"user_id"`

	Request_body json.RawMessage `json:"RequestBody"`
}

func RequestMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logReq := LogRequest{
			Request_id: uuid.New(),
			Method:     r.Method,
			Endpoint:   r.URL.Path,
			Timestamp:  start, // Время начала запроса
			UserID:     nil,   // Пока у нас нет авторизации, оставляем null (в будущем будешь доставать из токена)
		}

		if r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" {
			bodyBytes, _ := io.ReadAll(r.Body)
			logReq.Request_body = bodyBytes
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		w.Header().Set("X-Request-Id", logReq.Request_id.String())
		next.ServeHTTP(ww, r)
		logReq.Status_code = ww.Status()
		logReq.DurationMs = time.Since(start).Milliseconds()
		jsonLog, _ := json.Marshal(logReq)
		log.Println(string(jsonLog))
	})
}
