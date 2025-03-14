// backend/handlers/stroke.go
package handlers

import (
	"backend/models"
	"backend/storage"
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// StrokeHandler 處理筆畫相關的請求
type StrokeHandler struct {
	store storage.Storage
}

// NewStrokeHandler 創建一個新的筆畫處理器
func NewStrokeHandler(store storage.Storage) *StrokeHandler {
	return &StrokeHandler{
		store: store,
	}
}

// RecordStroke 記錄筆畫
func (h *StrokeHandler) RecordStroke(w http.ResponseWriter, r *http.Request) {
	var req models.StrokeRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// 驗證請求
	if req.UserID <= 0 || req.CharacterID <= 0 || req.StrokeIndex < 0 || len(req.Path) < 2 {
		http.Error(w, "Invalid request parameters", http.StatusBadRequest)
		return
	}

	// 創建筆畫記錄
	newRecord := models.StrokeRecord{
		UserID:      req.UserID,
		CharacterID: req.CharacterID,
		StrokeIndex: req.StrokeIndex,
		Path:        req.Path,
		Score:       req.Score,
	}

	// 儲存記錄
	record, err := h.store.CreateStrokeRecord(newRecord)
	if err != nil {
		http.Error(w, "Error saving stroke record", http.StatusInternalServerError)
		return
	}

	// 簡化筆畫路徑為關鍵節點
	simplifiedNodes := h.simplifyStroke(req.Path)

	// 更新用戶進度
	err = h.store.UpdateUserProgress(req.UserID, req.CharacterID, req.StrokeIndex, req.Score)
	if err != nil {
		http.Error(w, "Error updating user progress", http.StatusInternalServerError)
		return
	}

	// 返回成功響應
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.StrokeRecordResponse{
		RecordID:        record.ID,
		SimplifiedNodes: simplifiedNodes,
	})
}

// GetUserStrokeRecords 獲取用戶筆畫記錄
func (h *StrokeHandler) GetUserStrokeRecords(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// 獲取用戶記錄
	records := h.store.GetStrokeRecordsByUserID(userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

// simplifyStroke 簡化筆畫路徑為關鍵節點
// 使用Douglas-Peucker演算法簡化路徑
func (h *StrokeHandler) simplifyStroke(path []models.Node) []models.Node {
	if len(path) <= 3 {
		return path
	}

	// 對於每個字的每個筆畫，我們都固定使用3個節點:
	// 起點、最佳中間點和終點
	result := make([]models.Node, 3)
	result[0] = path[0]           // 起點
	result[2] = path[len(path)-1] // 終點

	// 找出最佳中間點
	// 使用Douglas-Peucker演算法找到與直線偏差最大的點
	maxDist := 0.0
	maxIndex := 0

	for i := 1; i < len(path)-1; i++ {
		dist := h.perpendicularDistance(path[i], path[0], path[len(path)-1])
		if dist > maxDist {
			maxDist = dist
			maxIndex = i
		}
	}

	// 如果沒有顯著偏差，就選擇中間點
	if maxDist < 5.0 {
		result[1] = path[len(path)/2]
	} else {
		result[1] = path[maxIndex]
	}

	return result
}

// perpendicularDistance 計算點到線段的垂直距離
func (h *StrokeHandler) perpendicularDistance(point, lineStart, lineEnd models.Node) float64 {
	// 如果線段長度為0，返回點到起點的距離
	if lineStart.X == lineEnd.X && lineStart.Y == lineEnd.Y {
		return math.Sqrt(math.Pow(point.X-lineStart.X, 2) + math.Pow(point.Y-lineStart.Y, 2))
	}

	// 計算線段長度的平方
	lenSq := math.Pow(lineEnd.X-lineStart.X, 2) + math.Pow(lineEnd.Y-lineStart.Y, 2)

	// 計算點到線段的垂直距離
	area := math.Abs((lineEnd.X-lineStart.X)*(lineStart.Y-point.Y) - (lineStart.X-point.X)*(lineEnd.Y-lineStart.Y))
	return area / math.Sqrt(lenSq)
}
