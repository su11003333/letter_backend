// backend/handlers/progress.go
package handlers

import (
	"backend/storage"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ProgressHandler 處理進度相關的請求
type ProgressHandler struct {
	store storage.Storage
}

// NewProgressHandler 創建一個新的進度處理器
func NewProgressHandler(store storage.Storage) *ProgressHandler {
	return &ProgressHandler{
		store: store,
	}
}

// GetUserProgress 獲取用戶的學習進度
func (h *ProgressHandler) GetUserProgress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// 獲取用戶進度
	progress := h.store.GetUserProgress(userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(progress)
}
