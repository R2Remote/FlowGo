package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"FLOWGO/internal/domain/entity"
	"FLOWGO/internal/infrastructure/database"
)

// VisitLogger 访问统计中间件
// 记录请求来源 IP，非阻塞式写入数据库（使用 Goroutine）
func VisitLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		// 异步记录，避免阻塞请求
		go func(clientIP string) {
			// 如果数据库连接未初始化，直接跳过
			if database.DB == nil {
				return
			}

			// Upsert: 存在则 Count+1，不存在则 Insert
			database.DB.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "ip"}},
				DoUpdates: clause.Assignments(map[string]any{
					"count":      gorm.Expr("count + 1"),
					"last_seen":  time.Now(),
					"updated_at": time.Now(),
				}),
			}).Create(&entity.VisitStat{
				IP:       clientIP,
				Count:    1,
				LastSeen: time.Now(),
			})
		}(ip)

		c.Next()
	}
}
