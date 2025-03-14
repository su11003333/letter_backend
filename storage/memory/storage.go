// backend/storage/memory/storage.go
package memory

import (
	"backend/models"
	"errors"
	"fmt"
	"time"
)

// MemoryStorage 實現 Storage 接口的記憶體儲存
type MemoryStorage struct {
	users            []models.User
	characters       []models.CharacterPreview
	characterDetails map[int]models.Character
	strokeRecords    []models.StrokeRecord
	userProgress     map[int]models.UserProgress // userID -> characterID -> progress
	recordCounter    int                         // 用於生成唯一ID
}

// NewMemoryStorage 創建一個新的記憶體儲存
func NewMemoryStorage() *MemoryStorage {
	// 初始化模擬資料
	users := []models.User{
		{ID: 1, Username: "admin", Password: "password", Email: "admin@example.com"},
		// 這裡應該有您想要登入的用戶
	}

	characters := []models.CharacterPreview{
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

	characterDetails := map[int]models.Character{
		1: {
			ID:     1,
			Name:   "一",
			SVGUrl: "/assets/characters/yi.svg",
			StrokeData: []models.Stroke{
				{
					Nodes: []models.Node{
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
			StrokeData: []models.Stroke{
				{
					Nodes: []models.Node{
						{X: 150, Y: 250},
						{X: 300, Y: 250},
						{X: 450, Y: 250},
					},
				},
				{
					Nodes: []models.Node{
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
			StrokeData: []models.Stroke{
				{
					Nodes: []models.Node{
						{X: 150, Y: 200},
						{X: 300, Y: 200},
						{X: 450, Y: 200},
					},
				},
				{
					Nodes: []models.Node{
						{X: 150, Y: 300},
						{X: 300, Y: 300},
						{X: 450, Y: 300},
					},
				},
				{
					Nodes: []models.Node{
						{X: 150, Y: 400},
						{X: 300, Y: 400},
						{X: 450, Y: 400},
					},
				},
			},
		},
	}

	return &MemoryStorage{
		users:            users,
		characters:       characters,
		characterDetails: characterDetails,
		strokeRecords:    []models.StrokeRecord{},
		userProgress:     make(map[int]models.UserProgress),
		recordCounter:    1,
	}
}

// GetUsers 獲取所有用戶
func (s *MemoryStorage) GetUsers() []models.User {
	return s.users
}

// GetUserByID 根據ID獲取用戶
func (s *MemoryStorage) GetUserByID(id int) (*models.User, error) {
	for _, user := range s.users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}

// GetUserByUsername 根據用戶名獲取用戶
func (s *MemoryStorage) GetUserByUsername(username string) (*models.User, error) {
	for _, user := range s.users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}

// CreateUser 創建新用戶
func (s *MemoryStorage) CreateUser(user models.User) (*models.User, error) {
	// 檢查用戶名是否已存在
	for _, existingUser := range s.users {
		if existingUser.Username == user.Username {
			return nil, errors.New("username already exists")
		}
	}

	// 設置新用戶ID
	user.ID = len(s.users) + 1
	s.users = append(s.users, user)
	return &user, nil
}

// GetCharacters 獲取所有字元預覽
func (s *MemoryStorage) GetCharacters() []models.CharacterPreview {
	return s.characters
}

// GetCharacterByID 根據ID獲取字元詳情
func (s *MemoryStorage) GetCharacterByID(id int) (*models.Character, error) {
	character, exists := s.characterDetails[id]
	if !exists {
		return nil, fmt.Errorf("character with ID %d not found", id)
	}
	return &character, nil
}

// CreateStrokeRecord 創建筆畫記錄
func (s *MemoryStorage) CreateStrokeRecord(record models.StrokeRecord) (*models.StrokeRecord, error) {
	// 設置記錄ID和時間
	record.ID = s.recordCounter
	record.CreatedAt = time.Now()
	s.recordCounter++

	s.strokeRecords = append(s.strokeRecords, record)
	return &record, nil
}

// GetStrokeRecordsByUserID 獲取用戶的筆畫記錄
func (s *MemoryStorage) GetStrokeRecordsByUserID(userID int) []models.StrokeRecord {
	var userRecords []models.StrokeRecord
	for _, record := range s.strokeRecords {
		if record.UserID == userID {
			userRecords = append(userRecords, record)
		}
	}
	return userRecords
}

// GetUserProgress 獲取用戶進度
func (s *MemoryStorage) GetUserProgress(userID int) models.UserProgress {
	progress, exists := s.userProgress[userID]
	if !exists {
		return models.UserProgress{}
	}
	return progress
}

// UpdateUserProgress 更新用戶進度
func (s *MemoryStorage) UpdateUserProgress(userID, characterID, strokeIndex int, score float64) error {
	// 確保用戶進度映射存在
	progress, exists := s.userProgress[userID]
	if !exists {
		progress = models.UserProgress{}
		s.userProgress[userID] = progress
	}

	// 獲取字元進度或初始化
	charProgress, exists := progress[characterID]
	if !exists {
		charProgress = models.CharacterProgress{
			CharacterID: characterID,
			Attempts:    1,
			AvgScore:    score,
			Mastery:     score * 100,
			LastStroke:  strokeIndex,
		}
	} else {
		// 更新現有進度
		attempts := charProgress.Attempts
		avgScore := charProgress.AvgScore

		// 計算新的平均得分
		newAvgScore := (avgScore*float64(attempts) + score) / float64(attempts+1)

		charProgress.Attempts = attempts + 1
		charProgress.AvgScore = newAvgScore
		charProgress.Mastery = newAvgScore * 100

		// 更新最後筆畫索引（如果更大）
		if strokeIndex > charProgress.LastStroke {
			charProgress.LastStroke = strokeIndex
		}
	}

	// 儲存更新後的進度
	progress[characterID] = charProgress
	s.userProgress[userID] = progress

	return nil
}
