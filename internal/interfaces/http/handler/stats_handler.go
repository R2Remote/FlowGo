package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"FLOWGO/internal/domain/entity"
	"FLOWGO/internal/infrastructure/database"
)

type StatsHandler struct{}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{}
}

// GetVisitStats 获取访问统计列表
func (h *StatsHandler) GetVisitStats(c *gin.Context) {
	var stats []entity.VisitStat

	// 按访问次数降序排列，取前 100 条
	// 这里可以根据需要加分页，但对于简单的功能先取 Top 100
	if err := database.DB.Order("count desc").Limit(100).Find(&stats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
