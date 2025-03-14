// backend/routes/routes.go
package routes

import (
	"backend/configs"
	"backend/handlers"
	"backend/middleware"
	"backend/storage/memory"

	"github.com/gorilla/mux"
)

// SetupRoutes 設置 API 路由
func SetupRoutes(config *configs.Config) *mux.Router {
	// 初始化儲存
	store := memory.NewMemoryStorage()

	// 初始化處理程序
	authHandler := handlers.NewAuthHandler(store, config)
	characterHandler := handlers.NewCharacterHandler(store)
	strokeHandler := handlers.NewStrokeHandler(store)
	progressHandler := handlers.NewProgressHandler(store)

	// 創建主路由器
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()

	// 公共路由
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")

	// 需要認證的路由
	authenticatedAPI := api.PathPrefix("").Subrouter()
	authenticatedAPI.Use(middleware.AuthMiddleware(config))

	// 字元相關路由
	authenticatedAPI.HandleFunc("/characters", characterHandler.GetCharacters).Methods("GET")
	authenticatedAPI.HandleFunc("/characters/{id}", characterHandler.GetCharacterByID).Methods("GET")

	// 筆畫記錄相關路由
	authenticatedAPI.HandleFunc("/strokes/record", strokeHandler.RecordStroke).Methods("POST")
	authenticatedAPI.HandleFunc("/users/{userId}/stroke-records", strokeHandler.GetUserStrokeRecords).Methods("GET")

	// 進度相關路由
	authenticatedAPI.HandleFunc("/users/{userId}/progress", progressHandler.GetUserProgress).Methods("GET")

	return router
}
