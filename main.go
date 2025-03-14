// backend/main.go
package main

import (
	"backend/configs"
	"backend/routes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"
)

func main() {
	// 加載配置
	config := configs.LoadConfig()

	// 設置路由
	router := routes.SetupRoutes(config)

	// 設置 CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   config.AllowedOrigins,
		AllowedMethods:   config.AllowedMethods,
		AllowedHeaders:   config.AllowedHeaders,
		AllowCredentials: config.AllowCredentials,
		MaxAge:           300, // 5 分鐘的預檢快取
	})

	// 設置服務器
	server := &http.Server{
		Handler:      corsHandler.Handler(router),
		Addr:         ":" + config.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 啟動服務器
	fmt.Printf("Server is running on port %s...\n", config.Port)
	log.Fatal(server.ListenAndServe())
}
