package handlers

import (
	"net/http"
	"time"

	"main/models"

	"github.com/gin-gonic/gin"
)

// Health 健康检查
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "OK",
		Data: map[string]interface{}{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		},
	})
}

// Status 获取状态
func Status(c *gin.Context) {
	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "OK",
		Data: map[string]interface{}{
			"status":  "running",
			"version": "1.0.0",
			"uptime":  time.Now().Unix(),
		},
	})
}
