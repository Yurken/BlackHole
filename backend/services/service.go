package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"main/models"
)

// GlobalAIConfig 全局 AI 配置
var GlobalAIConfig = models.AIConfig{
	Provider: "ollama",
	BaseURL:  "http://localhost:11434",
	Model:    "qwen3-vl:4b",
}

// UserTemplates 用户模板存储
var UserTemplates = make(map[string]models.Template)

// CopyFile 复制文件
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// 复制文件权限
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, sourceInfo.Mode())
}

// MoveFile 移动文件（跨设备时回退为复制删除）
func MoveFile(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	if err := CopyFile(src, dst); err != nil {
		return err
	}

	return os.Remove(src)
}

// AnalyzeFile 使用配置的 AI 提供商分析文件
func AnalyzeFile(filePath string, model string) (*models.AIAnalysis, error) {
	if model == "" {
		model = GlobalAIConfig.Model
	}

	switch GlobalAIConfig.Provider {
	case "ollama":
		return analyzeWithOllama(filePath, model)
	case "openai", "deepseek", "qwen":
		return analyzeWithOpenAICompatible(filePath, model)
	default:
		return nil, fmt.Errorf("不支持的 AI 提供商: %s", GlobalAIConfig.Provider)
	}
}

// AnalyzeFileWithOllama 使用 Ollama API 分析文件（向后兼容）
func AnalyzeFileWithOllama(filePath string, model string) (*models.AIAnalysis, error) {
	return analyzeWithOllama(filePath, model)
}

// isImageFile 判断是否为图片文件
func isImageFile(ext string) bool {
	ext = strings.ToLower(ext)
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".tiff", ".tif"}
	for _, e := range imageExts {
		if ext == e {
			return true
		}
	}
	return false
}

// isPDFFile 判断是否为 PDF 文件
func isPDFFile(ext string) bool {
	return strings.ToLower(ext) == ".pdf"
}

// convertPDFToImage 将 PDF 第一页转为图片（使用 macOS qlmanage）
func convertPDFToImage(pdfPath string) (string, error) {
	// 创建临时目录
	tmpDir := os.TempDir()

	// 使用 macOS Quick Look 生成缩略图
	cmd := exec.Command("qlmanage", "-t", "-s", "1200", "-o", tmpDir, pdfPath)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("PDF转图片失败: %v", err)
	}

	// qlmanage 会生成 filename.pdf.png
	pdfName := filepath.Base(pdfPath)
	imagePath := filepath.Join(tmpDir, pdfName+".png")

	// 检查文件是否存在
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return "", fmt.Errorf("生成的图片不存在: %s", imagePath)
	}

	return imagePath, nil
}

// encodeImageToBase64 将图片编码为 base64
func encodeImageToBase64(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// analyzeWithOllama 使用 Ollama API 分析文件
func analyzeWithOllama(filePath string, model string) (*models.AIAnalysis, error) {
	if model == "" {
		model = GlobalAIConfig.Model
	}

	fileName := filepath.Base(filePath)
	ext := filepath.Ext(fileName)
	nameWithoutExt := strings.TrimSuffix(fileName, ext)
	isImage := isImageFile(ext)
	isPDF := isPDFFile(ext)

	// 判断是否可以用视觉处理（图片或PDF）
	canUseVision := isImage || isPDF
	actualModel := model

	// 如果不是可视觉处理的文件，且用的是视觉模型，则切换到文本模型
	if !canUseVision && strings.Contains(model, "-vl") {
		actualModel = strings.Replace(model, "-vl", "", 1)
		log.Printf("[AI] 非视觉文件，切换模型: %s -> %s", model, actualModel)
	}

	// 获取图片数据（图片直接读取，PDF先转图片）
	var imageBase64 string
	var imagePath string
	if isImage {
		imagePath = filePath
	} else if isPDF {
		// PDF 转图片
		log.Printf("[AI] 正在将 PDF 转换为图片...")
		convertedPath, err := convertPDFToImage(filePath)
		if err != nil {
			log.Printf("[AI] PDF转图片失败: %v, 使用原名", err)
			return &models.AIAnalysis{
				SuggestedName: nameWithoutExt,
				Category:      "文档",
				Confidence:    0,
			}, nil
		}
		imagePath = convertedPath
		defer os.Remove(imagePath) // 处理完删除临时图片
		log.Printf("[AI] PDF 已转换为图片: %s", imagePath)
	}

	// 构建提示词 - 使用 /no_think 让 Qwen3 直接输出不思考
	var prompt string
	if canUseVision && imagePath != "" {
		// 对图片/PDF使用视觉模型
		if isPDF {
			prompt = `这是一份PDF文档的第一页。根据文档内容，返回一个简短的中文或者英文文件名（最多15个字词）。
只返回JSON: {"suggested_name": "文件名", "category": "分类", "confidence": 0.9}`
		} else {
			prompt = `根据图片内容，返回一个简短的中文或英文文件名（最多15个字词）。
只返回JSON: {"suggested_name": "文件名", "category": "分类", "confidence": 0.9}`
		}
		// 读取图片
		var err error
		imageBase64, err = encodeImageToBase64(imagePath)
		if err != nil {
			log.Printf("[AI] 读取图片失败: %v", err)
		}
	} else {
		prompt = fmt.Sprintf(`原文件名: %s
根据文件名含义，返回一个简短的中文或英文文件名（最多15个字词）。
只返回JSON: {"suggested_name": "文件名", "category": "分类", "confidence": 0.9}`, nameWithoutExt)
	}

	var reqBody map[string]interface{}
	var apiEndpoint string

	// 视觉模型使用 /api/chat，文本模型使用 /api/generate
	if canUseVision && imageBase64 != "" {
		// 使用 chat API（视觉模型）
		apiEndpoint = "/api/chat"
		message := map[string]interface{}{
			"role":    "user",
			"content": prompt,
		}
		if imageBase64 != "" {
			message["images"] = []interface{}{imageBase64}
		}
		reqBody = map[string]interface{}{
			"model":    actualModel,
			"messages": []interface{}{message},
			"stream":   false,
		}
	} else {
		// 使用 generate API（文本模型）
		apiEndpoint = "/api/generate"
		reqBody = map[string]interface{}{
			"model":  actualModel,
			"prompt": prompt,
			"stream": false,
			"options": map[string]interface{}{
				"num_predict": 100,
			},
		}
	}

	jsonData, _ := json.Marshal(reqBody)

	// 打印发送的请求（不打印完整图片数据）
	debugReqBody := make(map[string]interface{})
	for k, v := range reqBody {
		debugReqBody[k] = v
	}
	if canUseVision && imageBase64 != "" {
		if messages, ok := debugReqBody["messages"].([]interface{}); ok && len(messages) > 0 {
			if msg, ok := messages[0].(map[string]interface{}); ok {
				msg["images"] = []string{"<base64_data>"}
			}
		}
	}
	debugJSON, _ := json.MarshalIndent(debugReqBody, "", "  ")
	log.Printf("[AI] 发送请求到 %s: %s", apiEndpoint, string(debugJSON))

	baseURL := GlobalAIConfig.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	requestTimeout := 60 * time.Second
	if canUseVision {
		requestTimeout = 180 * time.Second
	}
	client := &http.Client{Timeout: requestTimeout}
	resp, err := client.Post(baseURL+apiEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("Ollama 响应超时: %v", err)
		}
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return nil, fmt.Errorf("Ollama 响应超时: %v", err)
		}
		return nil, fmt.Errorf("无法连接到 Ollama: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ollama API 返回错误 (状态码 %d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	// 根据不同的 API 解析不同的响应格式
	var content string
	if apiEndpoint == "/api/chat" {
		// chat API 返回格式
		var chatResult struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&chatResult); err != nil {
			return nil, fmt.Errorf("解析响应失败: %v", err)
		}
		content = chatResult.Message.Content
	} else {
		// generate API 返回格式
		var generateResult struct {
			Response string `json:"response"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&generateResult); err != nil {
			return nil, fmt.Errorf("解析响应失败: %v", err)
		}
		content = generateResult.Response
	}

	// 打印模型返回内容用于调试
	log.Printf("[AI] 模型: %s, 文件: %s", actualModel, fileName)
	log.Printf("[AI] 原始响应: %s", content)

	// 尝试从响应中提取 JSON
	var analysis models.AIAnalysis
	// 尝试直接解析
	if err := json.Unmarshal([]byte(content), &analysis); err != nil {
		// 如果失败，尝试查找 JSON 部分
		start := strings.Index(content, "{")
		end := strings.LastIndex(content, "}")
		if start != -1 && end != -1 && end > start {
			jsonStr := content[start : end+1]
			log.Printf("[AI] 提取的JSON: %s", jsonStr)
			if err := json.Unmarshal([]byte(jsonStr), &analysis); err != nil {
				// 解析失败，返回原文件名
				log.Printf("[AI] JSON解析失败: %v, 使用原名", err)
				return &models.AIAnalysis{
					SuggestedName: nameWithoutExt,
					Category:      "未知",
					Confidence:    0,
				}, nil
			}
		} else {
			// 响应为空或格式错误，返回原文件名
			log.Printf("[AI] 响应无JSON, 使用原名")
			return &models.AIAnalysis{
				SuggestedName: nameWithoutExt,
				Category:      "未知",
				Confidence:    0,
			}, nil
		}
	}

	log.Printf("[AI] 解析结果: suggested_name=%s, category=%s", analysis.SuggestedName, analysis.Category)

	// 如果返回的名称为空，使用原文件名
	if strings.TrimSpace(analysis.SuggestedName) == "" {
		log.Printf("[AI] 返回名称为空, 使用原名")
		analysis.SuggestedName = nameWithoutExt
		analysis.Confidence = 0
	}

	return &analysis, nil
}

// analyzeWithOpenAICompatible 使用 OpenAI 兼容 API 分析文件
func analyzeWithOpenAICompatible(filePath string, model string) (*models.AIAnalysis, error) {
	if model == "" {
		model = GlobalAIConfig.Model
	}

	fileName := filepath.Base(filePath)
	ext := filepath.Ext(fileName)

	// 构建提示词
	prompt := fmt.Sprintf(`你是一个文件命名专家。请根据文件名生成一个简洁、有意义的新文件名。

原文件名: %s
文件类型: %s

要求:
1. 新文件名必须简洁明了，最多20个字符
2. 不要包含扩展名
3. 不要包含特殊字符（如 / \ : * ? " < > |）
4. 使用中文或英文，保持简短

请只返回一个JSON对象，不要有任何其他文字:
{"suggested_name": "新文件名", "category": "文件分类", "confidence": 0.9}`, fileName, ext)

	// OpenAI API 请求格式
	reqBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
	}

	jsonData, _ := json.Marshal(reqBody)

	// 获取 base URL
	baseURL := GlobalAIConfig.BaseURL
	if baseURL == "" {
		switch GlobalAIConfig.Provider {
		case "openai":
			baseURL = "https://api.openai.com/v1"
		case "deepseek":
			baseURL = "https://api.deepseek.com/v1"
		case "qwen":
			baseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
		}
	}

	// 创建请求
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if GlobalAIConfig.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+GlobalAIConfig.APIKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("无法连接到 %s: %v", GlobalAIConfig.Provider, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%s API 返回错误 (状态码 %d): %s", GlobalAIConfig.Provider, resp.StatusCode, string(body))
	}

	// 解析响应
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("API 未返回有效响应")
	}

	content := result.Choices[0].Message.Content

	// 尝试从响应中提取 JSON
	var analysis models.AIAnalysis
	// 尝试直接解析
	if err := json.Unmarshal([]byte(content), &analysis); err != nil {
		// 如果失败，尝试查找 JSON 部分
		start := strings.Index(content, "{")
		end := strings.LastIndex(content, "}")
		if start != -1 && end != -1 && end > start {
			jsonStr := content[start : end+1]
			if err := json.Unmarshal([]byte(jsonStr), &analysis); err != nil {
				return nil, fmt.Errorf("解析 AI 响应失败: %v", err)
			}
		} else {
			return nil, fmt.Errorf("AI 响应格式错误")
		}
	}

	return &analysis, nil
}

// GenerateTemplatePreview 生成模板预览
func GenerateTemplatePreview(components []models.TemplateComponent) string {
	var parts []string
	now := time.Now()

	for _, c := range components {
		switch c.Type {
		case "year":
			parts = append(parts, now.Format("2006"))
		case "month":
			parts = append(parts, now.Format("01"))
		case "day":
			parts = append(parts, now.Format("02"))
		case "date":
			parts = append(parts, now.Format("2006-01-02"))
		case "original":
			parts = append(parts, "文件名")
		case "separator":
			parts = append(parts, c.Label)
		case "text":
			parts = append(parts, c.Label)
		default:
			parts = append(parts, c.Label)
		}
	}

	return strings.Join(parts, "")
}

// TestAIConnection 测试 AI 连接
func TestAIConnection(req models.AITestRequest) (map[string]interface{}, error) {
	switch req.Provider {
	case "ollama":
		baseURL := req.BaseURL
		if baseURL == "" {
			baseURL = "http://localhost:11434"
		}

		resp, err := http.Get(baseURL + "/api/tags")
		if err != nil {
			return nil, fmt.Errorf("无法连接到 Ollama 服务: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Ollama 服务返回错误")
		}

		return map[string]interface{}{
			"status":   "connected",
			"provider": "ollama",
		}, nil

	case "openai", "deepseek", "qwen":
		baseURL := req.BaseURL
		if baseURL == "" {
			switch req.Provider {
			case "openai":
				baseURL = "https://api.openai.com/v1"
			case "deepseek":
				baseURL = "https://api.deepseek.com/v1"
			case "qwen":
				baseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
			}
		}

		if req.APIKey == "" {
			return nil, fmt.Errorf("API Key 不能为空")
		}

		client := &http.Client{Timeout: 10 * time.Second}
		testReq, _ := http.NewRequest("GET", baseURL+"/models", nil)
		testReq.Header.Set("Authorization", "Bearer "+req.APIKey)

		resp, err := client.Do(testReq)
		if err != nil {
			return nil, fmt.Errorf("无法连接到 %s 服务: %v", req.Provider, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == 401 {
			return nil, fmt.Errorf("API Key 无效")
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%s 服务返回错误", req.Provider)
		}

		return map[string]interface{}{
			"status":   "connected",
			"provider": req.Provider,
		}, nil

	default:
		return nil, fmt.Errorf("不支持的 AI 提供商")
	}
}

// GetOllamaModels 获取 Ollama 模型列表
func GetOllamaModels() ([]string, error) {
	resp, err := http.Get("http://localhost:11434/api/tags")
	if err != nil {
		return nil, fmt.Errorf("无法连接到 Ollama 服务")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama 服务返回错误")
	}

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析模型列表失败")
	}

	models := make([]string, len(result.Models))
	for i, model := range result.Models {
		models[i] = model.Name
	}

	return models, nil
}
