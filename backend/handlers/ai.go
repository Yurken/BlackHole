package handlers

import (
	"log"
	"net/http"

	"main/models"
	"main/services"

	"github.com/gin-gonic/gin"
)

// GetOllamaModels 获取 Ollama 模型列表
func GetOllamaModels(c *gin.Context) {
	modelList, err := services.GetOllamaModels()
	if err != nil {
		log.Println("Failed to fetch Ollama models:", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    5000,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    modelList,
	})
}

// TestAIConnection 测试 AI 连接
func TestAIConnection(c *gin.Context) {
	var req models.AITestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    1000,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	result, err := services.TestAIConnection(req)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    2000,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: req.Provider + " 连接成功",
		Data:    result,
	})
}

// GetAIConfig 获取 AI 配置
func GetAIConfig(c *gin.Context) {
	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    services.GlobalAIConfig,
	})
}

// SaveAIConfig 保存 AI 配置
func SaveAIConfig(c *gin.Context) {
	var req models.AIConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    1000,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// 更新全局配置
	services.GlobalAIConfig = req

	log.Printf("AI 配置已更新: Provider=%s, Model=%s", req.Provider, req.Model)

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "配置保存成功",
		Data:    services.GlobalAIConfig,
	})
}

// AnalyzeFile AI 分析文件
func AnalyzeFile(c *gin.Context) {
	var req models.AIAnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    1000,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	analysis, err := services.AnalyzeFile(req.FilePath, services.GlobalAIConfig.Model)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    2001,
			Message: "AI 分析失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "分析成功",
		Data:    analysis,
	})
}
