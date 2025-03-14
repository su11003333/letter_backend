// backend/handlers/auth.go
package handlers

import (
	"backend/configs"
	"backend/models"
	"backend/storage"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// AuthHandler 處理認證相關的請求
type AuthHandler struct {
	store  storage.Storage
	config *configs.Config
}

// NewAuthHandler 創建一個新的認證處理器
func NewAuthHandler(store storage.Storage, config *configs.Config) *AuthHandler {
	return &AuthHandler{
		store:  store,
		config: config,
	}
}

// Login 處理登入請求
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// 驗證用戶名和密碼
	user, err := h.store.GetUserByUsername(req.Username)
	if err != nil || user.Password != req.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 創建 JWT Token
	expirationTime := time.Now().Add(h.config.JWTExpiryTime)
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.config.JWTSecret)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// 隱藏密碼
	user.Password = ""

	// 返回用戶和令牌
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.LoginResponse{
		User:  *user,
		Token: tokenString,
	})
}

// Register 處理註冊請求
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// 驗證用戶名和密碼
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// 創建新用戶
	newUser := models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    "", // 可以從請求中獲取或留空
	}

	user, err := h.store.CreateUser(newUser)
	if err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 創建 JWT Token
	expirationTime := time.Now().Add(h.config.JWTExpiryTime)
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.config.JWTSecret)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// 隱藏密碼
	user.Password = ""

	// 返回用戶和令牌
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.LoginResponse{
		User:  *user,
		Token: tokenString,
	})
}
