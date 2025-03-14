// backend/handlers/character.go
package handlers

import (
	"backend/storage"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CharacterHandler 處理字元相關的請求
type CharacterHandler struct {
	store storage.Storage
}

// NewCharacterHandler 創建一個新的字元處理器
func NewCharacterHandler(store storage.Storage) *CharacterHandler {
	return &CharacterHandler{
		store: store,
	}
}

// GetCharacters 獲取所有字元
func (h *CharacterHandler) GetCharacters(w http.ResponseWriter, r *http.Request) {
	characters := h.store.GetCharacters()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characters)
}

// GetCharacterByID 根據ID獲取字元詳情
func (h *CharacterHandler) GetCharacterByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	character, err := h.store.GetCharacterByID(id)
	if err != nil {
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(character)
}
