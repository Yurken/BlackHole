package main

import (
	"fmt"
	"log"

	"main/database"
	"main/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	if err := database.Init(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 设置 Gin 模式
	// gin.SetMode(gin.ReleaseMode) // 生产环境使用

	// 创建 Gin 实例
	r := gin.Default()

	// 启用 CORS
	r.Use(corsMiddleware())

	// 设置路由
	routes.SetupRoutes(r)

	port := ":18620"

	fmt.Printf("🚀 Go backend server running on http://localhost%s\n", port)
	fmt.Println("📡 API endpoints:")
	fmt.Println("   - GET  /api/health            - 健康检查")
	fmt.Println("   - GET  /api/status            - 获取状态")
	fmt.Println("   - POST /api/files/process     - 处理文件")
	fmt.Println("   - GET  /api/history           - 获取历史记录")
	fmt.Println("   - POST /api/history/clear     - 清除历史记录")
	fmt.Println("   - GET  /api/ollama/models     - 获取Ollama模型列表")
	fmt.Println("   - GET  /api/templates         - 获取模板列表")
	fmt.Println("   - POST /api/templates/import  - 导入模板")
	fmt.Println("   - DELETE /api/templates/:id   - 删除模板")
	fmt.Println("   - POST /api/ai/test-connection- 测试AI连接")
	fmt.Println("   - GET/POST /api/ai/config     - AI配置")
	fmt.Println("   - POST /api/ai/analyze        - AI分析")
	fmt.Println("\n💡 使用说明:")
	fmt.Println("   1. 确保 Ollama 已启动: ollama serve")
	fmt.Println("   2. 下载模型: ollama pull qwen3-vl:4b")
	fmt.Println("   3. 拖动文件到悬浮球进行处理")

	// 启动服务器
	if err := r.Run(port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

// corsMiddleware CORS 中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
