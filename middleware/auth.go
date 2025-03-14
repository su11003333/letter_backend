// backend/middleware/auth.go
package middleware

import (
	"backend/configs"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// 創建上下文鍵類型
type contextKey string

const userIDKey contextKey = "userID"

// GetUserIDFromContext 從上下文中獲取用戶ID
func GetUserIDFromContext(ctx context.Context) (interface{}, bool) {
	userID, ok := ctx.Value(userIDKey).(interface{})
	return userID, ok
}

// AuthMiddleware 認證中間件
func AuthMiddleware(config *configs.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 開放路徑 - 不需要認證
			openPaths := map[string]bool{
				"/api/auth/login":    true,
				"/api/auth/register": true,
			}

			// 檢查是否為開放路徑
			if openPaths[r.URL.Path] {
				next.ServeHTTP(w, r)
				return
			}

			// 從請求標頭獲取 JWT Token
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
				return
			}

			// 如果 Token 使用 Bearer 前綴，則去除前綴
			if strings.HasPrefix(tokenString, "Bearer ") {
				tokenString = strings.TrimPrefix(tokenString, "Bearer ")
			}

			// 驗證 Token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// 驗證簽名方法
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return config.JWTSecret, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
				return
			}

			// Token 有效，取得 claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Unauthorized: Invalid token claims", http.StatusUnauthorized)
				return
			}

			// 從 claims 獲取用戶 ID
			userID, ok := claims["user_id"]
			if !ok {
				http.Error(w, "Unauthorized: User ID not found in token", http.StatusUnauthorized)
				return
			}

			// 將用戶 ID 添加到上下文
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			r = r.WithContext(ctx)

			// 繼續執行下一個處理程序
			next.ServeHTTP(w, r)
		})
	}
}
