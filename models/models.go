// backend/models/models.go
package models

import "time"

// Node 代表筆畫中的一個點
type Node struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Stroke 代表一個完整的筆畫，由多個節點組成
type Stroke struct {
	Nodes []Node `json:"nodes"`
}

// CharacterPreview 用於字元選擇列表
type CharacterPreview struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Preview string `json:"preview"`
}

// Character 代表完整的字元資料
type Character struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	SVGUrl     string   `json:"svgUrl"`
	StrokeData []Stroke `json:"strokeData"`
}

// User 使用者資料
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // 不在 JSON 中返回密碼
	Email    string `json:"email,omitempty"`
}

// LoginRequest 登入請求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 登入回應
type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

// StrokeRecordRequest 筆畫記錄請求
type StrokeRecordRequest struct {
	UserID      int     `json:"userId"`
	CharacterID int     `json:"characterId"`
	StrokeIndex int     `json:"strokeIndex"`
	Path        []Node  `json:"path"`
	Score       float64 `json:"score"`
}

// StrokeRecord 筆畫記錄
type StrokeRecord struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userId"`
	CharacterID int       `json:"characterId"`
	StrokeIndex int       `json:"strokeIndex"`
	Path        []Node    `json:"path"`
	Score       float64   `json:"score"`
	CreatedAt   time.Time `json:"createdAt"`
}

// StrokeRecordResponse 筆畫記錄回應
type StrokeRecordResponse struct {
	RecordID        int    `json:"recordId"`
	SimplifiedNodes []Node `json:"simplifiedNodes"`
}

// CharacterProgress 字元進度
type CharacterProgress struct {
	CharacterID int     `json:"characterId"`
	Attempts    int     `json:"attempts"`
	AvgScore    float64 `json:"avgScore"`
	Mastery     float64 `json:"mastery"`
	LastStroke  int     `json:"lastStroke"`
}

// UserProgress 用戶進度映射 - 字元ID對應進度
type UserProgress map[int]CharacterProgress
