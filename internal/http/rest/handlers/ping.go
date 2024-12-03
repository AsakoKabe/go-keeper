package handlers

import (
	"encoding/json"
	"net/http"

	"go-keeper/internal/db/repository"
)

// PingHandler Структура для endpoints проверки состояния БД
type PingHandler struct {
	pingService repository.PingRepository
}

// NewPingHandler Конструктор для PingHandler
func NewPingHandler(pingService repository.PingRepository) *PingHandler {
	return &PingHandler{pingService: pingService}
}

// HealthDB Get запрос на проверку доступа к БД
func (h *PingHandler) HealthDB(w http.ResponseWriter, r *http.Request) {
	err := h.pingService.PingDB(r.Context())
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		message, _ := json.Marshal(map[string]string{"name": err.Error()})
		http.Error(w, string(message), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("ok"))
}
