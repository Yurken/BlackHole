package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"main/models"
	"main/services"

	"github.com/gin-gonic/gin"
)

// GetTemplates 获取模板列表
func GetTemplates(c *gin.Context) {
	templates := make([]models.Template, 0, len(services.UserTemplates))
	for _, t := range services.UserTemplates {
		templates = append(templates, t)
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    templates,
	})
}

// ImportTemplate 导入模板
func ImportTemplate(c *gin.Context) {
	var req struct {
		Name       string                     `json:"name" binding:"required"`
		Components []models.TemplateComponent `json:"components" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    1000,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// 生成模板ID
	templateID := fmt.Sprintf("template_%d", time.Now().UnixNano())

	// 生成预览
	preview := services.GenerateTemplatePreview(req.Components)

	template := models.Template{
		ID:         templateID,
		Name:       req.Name,
		Components: req.Components,
		Preview:    preview,
		CreatedAt:  time.Now().Format(time.RFC3339),
	}

	services.UserTemplates[templateID] = template

	log.Printf("导入模板: %s (%s)", req.Name, templateID)

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "导入成功",
		Data:    template,
	})
}

// DeleteTemplate 删除模板
func DeleteTemplate(c *gin.Context) {
	templateID := c.Param("id")

	if templateID == "" {
		c.JSON(http.StatusOK, models.Response{
			Code:    1000,
			Message: "模板ID不能为空",
		})
		return
	}

	if _, exists := services.UserTemplates[templateID]; !exists {
		c.JSON(http.StatusOK, models.Response{
			Code:    3000,
			Message: "模板不存在",
		})
		return
	}

	delete(services.UserTemplates, templateID)

	log.Printf("删除模板: %s", templateID)

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "删除成功",
	})
}
