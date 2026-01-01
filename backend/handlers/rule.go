package handlers

import (
	"database/sql"
	"net/http"

	"main/database"
	"main/models"

	"github.com/gin-gonic/gin"
)

// GetRules 获取规则列表
func GetRules(c *gin.Context) {
	rules, err := database.GetRules()
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    5000,
			Message: "获取规则失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    rules,
	})
}

// CreateRule 创建规则
func CreateRule(c *gin.Context) {
	var rule models.Rule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    1000,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	created, err := database.CreateRule(rule)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    5000,
			Message: "创建规则失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "创建成功",
		Data:    created,
	})
}

// UpdateRule 更新规则
func UpdateRule(c *gin.Context) {
	ruleID := c.Param("id")
	if ruleID == "" {
		c.JSON(http.StatusOK, models.Response{
			Code:    1000,
			Message: "规则ID不能为空",
		})
		return
	}

	var rule models.Rule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    1000,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}
	rule.ID = ruleID

	updated, err := database.UpdateRule(rule)
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
			Message: "更新规则失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "更新成功",
		Data:    updated,
	})
}

// DeleteRule 删除规则
func DeleteRule(c *gin.Context) {
	ruleID := c.Param("id")
	if ruleID == "" {
		c.JSON(http.StatusOK, models.Response{
			Code:    1000,
			Message: "规则ID不能为空",
		})
		return
	}

	if err := database.DeleteRule(ruleID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, models.Response{
				Code:    3000,
				Message: "规则不存在",
			})
			return
		}
		c.JSON(http.StatusOK, models.Response{
			Code:    5000,
			Message: "删除规则失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "删除成功",
	})
}
