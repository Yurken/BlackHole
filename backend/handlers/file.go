package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"main/database"
	"main/models"
	"main/services"

	"github.com/gin-gonic/gin"
)

// ProcessFile 处理文件
func ProcessFile(c *gin.Context) {
	var req models.FileProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    1000,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(req.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusOK, models.Response{
			Code:    1001,
			Message: "File not found",
		})
		return
	}

	// 获取文件名和扩展名
	originalName := filepath.Base(req.FilePath)
	ext := filepath.Ext(originalName)
	nameWithoutExt := strings.TrimSuffix(originalName, ext)

	// 查找规则
	var rule *models.Rule
	if req.RuleID != "" {
		foundRule, err := database.GetRule(req.RuleID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusOK, models.Response{
					Code:    3000,
					Message: "规则不存在",
				})
				return
			}
			c.JSON(http.StatusOK, models.Response{
				Code:    5000,
				Message: "获取规则失败: " + err.Error(),
			})
			return
		}
		rule = &foundRule
	} else {
		rules, err := database.GetRules()
		if err == nil {
			rule = services.MatchRuleForFile(req.FilePath, rules)
		}
	}

	ruleName := "默认规则"
	action := "copy"
	keepOriginal := false
	dateSource := "current"
	nameTemplate := []string{}
	useAI := req.UseAI

	// 目标目录
	homeDir, _ := os.UserHomeDir()
	destDir := filepath.Join(homeDir, "Documents", "BlackHole")

	if rule != nil {
		if rule.Name != "" {
			ruleName = rule.Name
		}
		if rule.Action != "" {
			action = rule.Action
		}
		if rule.Destination != "" {
			destDir = rule.Destination
		}
		keepOriginal = rule.KeepOriginal
		if rule.DateSource != "" {
			dateSource = rule.DateSource
		}
		if len(rule.NameTemplate) > 0 {
			nameTemplate = rule.NameTemplate
		}
		if rule.AIEnabled {
			useAI = true
		}
	}

	// 确保目标目录存在
	os.MkdirAll(destDir, 0755)

	// AI 分析结果
	var aiAnalysis *models.AIAnalysis
	aiName := ""
	if useAI {
		analysis, err := services.AnalyzeFile(req.FilePath, req.Model)
		if err != nil {
			log.Printf("AI 分析失败: %v", err)
			aiAnalysis = &models.AIAnalysis{
				SuggestedName: nameWithoutExt,
				Category:      "文档",
				Confidence:    0.5,
			}
			aiName = aiAnalysis.SuggestedName
		} else {
			aiAnalysis = analysis
			if aiAnalysis.SuggestedName != "" {
				aiName = aiAnalysis.SuggestedName
			}
		}
	}

	// 生成新文件名
	fileDate := services.SelectTimestamp(req.FilePath, dateSource)
	var newBase string
	if len(nameTemplate) > 0 {
		newBase = services.BuildNameFromTemplate(nameTemplate, originalName, aiName, fileDate)
	} else {
		base := nameWithoutExt
		if aiName != "" {
			base = aiName
		}
		newBase = fmt.Sprintf("%s_%s", fileDate.Format("2006-01-02"), base)
	}
	newName := newBase + ext
	destPath := filepath.Join(destDir, newName)

	operation := "copy"
	var processErr error
	if action == "move" && !keepOriginal {
		operation = "move"
		processErr = services.MoveFile(req.FilePath, destPath)
	} else {
		processErr = services.CopyFile(req.FilePath, destPath)
	}

	if processErr != nil {
		log.Printf("文件处理失败: %v", processErr)
		database.SaveHistory(req.FilePath, originalName, destPath, newName, ruleName, operation, "failed")
		c.JSON(http.StatusOK, models.Response{
			Code:    5000,
			Message: "文件处理失败: " + processErr.Error(),
		})
		return
	}

	response := models.FileProcessResponse{
		OriginalPath: req.FilePath,
		OriginalName: originalName,
		NewName:      newName,
		Destination:  destPath,
		RuleUsed:     ruleName,
		AIAnalysis:   aiAnalysis,
	}

	// 保存到历史记录
	database.SaveHistory(req.FilePath, originalName, destPath, newName, ruleName, operation, "success")

	log.Printf("处理文件: %s -> %s", req.FilePath, destPath)

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "处理成功",
		Data:    response,
	})
}

// GetHistory 获取历史记录
func GetHistory(c *gin.Context) {
	records, err := database.GetHistory()
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    5000,
			Message: "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    records,
	})
}

// ClearHistory 清除历史记录
func ClearHistory(c *gin.Context) {
	if err := database.ClearHistory(); err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    5000,
			Message: "Failed to clear history",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "历史记录已清除",
	})
}
