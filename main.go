// backend/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// 定義類型
type Node struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Stroke struct {
	Nodes []Node `json:"nodes"`
}

type CharacterPreview struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Preview string `json:"preview"`
}

type Character struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	SVGUrl     string   `json:"svgUrl"`
	StrokeData []Stroke `json:"strokeData"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // 不在 JSON 中返回密碼
	Email    string `json:"email,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

type StrokeRecordRequest struct {
	UserID      int     `json:"userId"`
	CharacterID int     `json:"characterId"`
	StrokeIndex int     `json:"strokeIndex"`
	Path        []Node  `json:"path"`
	Score       float64 `json:"score"`
}

type StrokeRecord struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userId"`
	CharacterID int       `json:"characterId"`
	StrokeIndex int       `json:"strokeIndex"`
	Path        []Node    `json:"path"`
	Score       float64   `json:"score"`
	CreatedAt   time.Time `json:"createdAt"`
}

type StrokeRecordResponse struct {
	RecordID        int    `json:"recordId"`
	SimplifiedNodes []Node `json:"simplifiedNodes"`
}

// JWT 密鑰
var jwtKey = []byte("your_secret_key") // 在實際應用中應該存儲在環境變數中

// 模擬用戶數據
var users = []User{
	{ID: 1, Username: "demo", Password: "password", Email: "demo@example.com"},
}

// 模擬字元數據
var characters = []CharacterPreview{
	{ID: 1, Name: "一", Preview: "一"},
	{ID: 2, Name: "二", Preview: "二"},
	{ID: 3, Name: "三", Preview: "三"},
	{ID: 4, Name: "四", Preview: "四"},
	{ID: 5, Name: "五", Preview: "五"},
	{ID: 6, Name: "六", Preview: "六"},
	{ID: 7, Name: "七", Preview: "七"},
	{ID: 8, Name: "八", Preview: "八"},
	{ID: 9, Name: "九", Preview: "九"},
	{ID: 10, Name: "十", Preview: "十"},
}

// 模擬字元詳細數據
var characterDetails = map[int]Character{
	1: {
		ID:     1,
		Name:   "一",
		SVGUrl: "/assets/characters/yi.svg",
		StrokeData: []Stroke{
			{
				Nodes: []Node{
					{X: 150, Y: 300},
					{X: 300, Y: 300},
					{X: 450, Y: 300},
				},
			},
		},
	},
	2: {
		ID:     2,
		Name:   "二",
		SVGUrl: "/assets/characters/er.svg",
		StrokeData: []Stroke{
			{
				Nodes: []Node{
					{X: 150, Y: 250},
					{X: 300, Y: 250},
					{X: 450, Y: 250},
				},
			},
			{
				Nodes: []Node{
					{X: 150, Y: 350},
					{X: 300, Y: 350},
					{X: 450, Y: 350},
				},
			},
		},
	},
	3: {
		ID:     3,
		Name:   "三",
		SVGUrl: "/assets/characters/san.svg",
		StrokeData: []Stroke{
			{
				Nodes: []Node{
					{X: 150, Y: 200},
					{X: 300, Y: 200},
					{X: 450, Y: 200},
				},
			},
			{
				Nodes: []Node{
					{X: 150, Y: 300},
					{X: 300, Y: 300},
					{X: 450, Y: 300},
				},
			},
			{
				Nodes: []Node{
					{X: 150, Y: 400},
					{X: 300, Y: 400},
					{X: 450, Y: 400},
				},
			},
		},
	},
}

// 筆畫記錄
var strokeRecords = []StrokeRecord{}

// 處理函式

// 登入
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 找到用戶
	var user User
	userFound := false
	for _, u := range users {
		if u.Username == req.Username && u.Password == req.Password {
			user = u
			userFound = true
			break
		}
	}

	if !userFound {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 創建 JWT Token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// 隱藏密碼
	user.Password = ""

	// 返回用戶和令牌
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		User:  user,
		Token: tokenString,
	})
}

// 獲取所有字元
func getCharactersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characters)
}

// 獲取特定字元詳情
func getCharacterByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	character, exists := characterDetails[id]
	if !exists {
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(character)
}

// 記錄用戶筆畫數據
func recordStrokeHandler(w http.ResponseWriter, r *http.Request) {
	var req StrokeRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 創建新記錄
	newRecord := StrokeRecord{
		ID:          len(strokeRecords) + 1,
		UserID:      req.UserID,
		CharacterID: req.CharacterID,
		StrokeIndex: req.StrokeIndex,
		Path:        req.Path,
		Score:       req.Score,
		CreatedAt:   time.Now(),
	}

	// 添加到記錄
	strokeRecords = append(strokeRecords, newRecord)

	// 簡化筆畫路徑為三個關鍵節點
	simplifiedNodes := simplifyStroke(req.Path)

	// 返回響應
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StrokeRecordResponse{
		RecordID:        newRecord.ID,
		SimplifiedNodes: simplifiedNodes,
	})
}

// 簡化筆畫路徑為三個節點
func simplifyStroke(path []Node) []Node {
	if len(path) < 3 {
		return path // 如果路徑少於3個點，直接返回
	}

	result := make([]Node, 3)
	result[0] = path[0]           // 起點
	result[1] = path[len(path)/2] // 中間點
	result[2] = path[len(path)-1] // 終點

	return result
}

// 獲取用戶筆畫記錄
func getUserStrokeRecordsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// 過濾用戶記錄
	var userRecords []StrokeRecord
	for _, record := range strokeRecords {
		if record.UserID == userID {
			userRecords = append(userRecords, record)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userRecords)
}

// 分析用戶進度
func getUserProgressHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// 簡單的進度計算
	// 在實際應用中，這應該更複雜
	characterProgress := make(map[int]float64)
	recordCounts := make(map[int]int)

	for _, record := range strokeRecords {
		if record.UserID == userID {
			recordCounts[record.CharacterID]++
			characterProgress[record.CharacterID] += record.Score
		}
	}

	// 計算平均得分
	progressData := make(map[int]map[string]interface{})
	for charID, totalScore := range characterProgress {
		count := recordCounts[charID]
		avgScore := totalScore / float64(count)

		progressData[charID] = map[string]interface{}{
			"characterId": charID,
			"attempts":    count,
			"avgScore":    avgScore,
			"mastery":     avgScore * 100, // 簡化的熟練度計算
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(progressData)
}

// 認證中間件
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 從請求頭獲取令牌
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			// 某些路徑可能不需要認證
			if r.URL.Path == "/api/auth/login" || r.URL.Path == "/api/auth/register" {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 如果令牌使用 Bearer 前綴
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// 解析令牌
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 令牌有效，繼續處理請求
		next.ServeHTTP(w, r)
	})
}

func main() {
	// 加載環境變數
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default environment variables")
	}

	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()

	// 公共路由
	api.HandleFunc("/auth/login", loginHandler).Methods("POST")
	api.HandleFunc("/auth/register", loginHandler).Methods("POST") // 簡化，使用相同的處理函式

	// 需要認證的路由
	authenticatedAPI := api.PathPrefix("").Subrouter()
	authenticatedAPI.Use(authMiddleware)

	authenticatedAPI.HandleFunc("/characters", getCharactersHandler).Methods("GET")
	authenticatedAPI.HandleFunc("/characters/{id}", getCharacterByIDHandler).Methods("GET")
	authenticatedAPI.HandleFunc("/strokes/record", recordStrokeHandler).Methods("POST")
	authenticatedAPI.HandleFunc("/users/{userId}/stroke-records", getUserStrokeRecordsHandler).Methods("GET")
	authenticatedAPI.HandleFunc("/users/{userId}/progress", getUserProgressHandler).Methods("GET")

	// 設置 CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// 設置服務器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      corsHandler.Handler(router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	fmt.Printf("Server is running on port %s...\n", port)
	log.Fatal(server.ListenAndServe())
}
