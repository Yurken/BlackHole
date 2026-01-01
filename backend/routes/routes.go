package routes

import (
	"main/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// API 路由组
	api := r.Group("/api")
	{
		// 系统相关
		api.GET("/health", handlers.Health)
		api.GET("/status", handlers.Status)

		// 文件处理
		api.POST("/files/process", handlers.ProcessFile)

		// 历史记录
		api.GET("/history", handlers.GetHistory)
		api.POST("/history/clear", handlers.ClearHistory)

		// 规则管理
		api.GET("/rules", handlers.GetRules)
		api.POST("/rules", handlers.CreateRule)
		api.PUT("/rules/:id", handlers.UpdateRule)
		api.DELETE("/rules/:id", handlers.DeleteRule)

		// 模板管理
		api.GET("/templates", handlers.GetTemplates)
		api.POST("/templates/import", handlers.ImportTemplate)
		api.DELETE("/templates/:id", handlers.DeleteTemplate)

		// AI 相关
		api.GET("/ollama/models", handlers.GetOllamaModels)
		api.POST("/ai/test-connection", handlers.TestAIConnection)
		api.GET("/ai/config", handlers.GetAIConfig)
		api.POST("/ai/config", handlers.SaveAIConfig)
		api.POST("/ai/analyze", handlers.AnalyzeFile)
	}
}
