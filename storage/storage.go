// backend/storage/storage.go
package storage

import (
	"backend/models"
)

// Storage 定義存儲介面
type Storage interface {
	// 用戶相關
	GetUsers() []models.User
	GetUserByID(id int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(user models.User) (*models.User, error)

	// 字元相關
	GetCharacters() []models.CharacterPreview
	GetCharacterByID(id int) (*models.Character, error)

	// 筆畫記錄相關
	CreateStrokeRecord(record models.StrokeRecord) (*models.StrokeRecord, error)
	GetStrokeRecordsByUserID(userID int) []models.StrokeRecord

	// 用戶進度相關
	GetUserProgress(userID int) models.UserProgress
	UpdateUserProgress(userID, characterID, strokeIndex int, score float64) error
}
