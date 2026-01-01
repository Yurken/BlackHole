package models

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// FileProcessRequest 文件处理请求
type FileProcessRequest struct {
	FilePath string `json:"file_path" binding:"required"`
	UseAI    bool   `json:"use_ai"`
	Model    string `json:"model"`
	RuleID   string `json:"rule_id,omitempty"`
}

// FileProcessResponse 文件处理响应
type FileProcessResponse struct {
	OriginalPath string      `json:"original_path"`
	OriginalName string      `json:"original_name"`
	NewName      string      `json:"new_name"`
	Destination  string      `json:"destination"`
	RuleUsed     string      `json:"rule_used"`
	AIAnalysis   *AIAnalysis `json:"ai_analysis,omitempty"`
}

// AIAnalysis AI 分析结果
type AIAnalysis struct {
	SuggestedName string  `json:"suggested_name"`
	Category      string  `json:"category"`
	Confidence    float64 `json:"confidence"`
}

// HistoryRecord 历史记录
type HistoryRecord struct {
	ID           int64  `json:"id"`
	OriginalPath string `json:"original_path"`
	OriginalName string `json:"original_name"`
	NewPath      string `json:"new_path"`
	NewName      string `json:"new_name"`
	RuleName     string `json:"rule_name"`
	Action       string `json:"action"` // copy or move
	Status       string `json:"status"` // success or failed
	Timestamp    string `json:"timestamp"`
}

// Template 模板结构
type Template struct {
	ID         string              `json:"id"`
	Name       string              `json:"name" binding:"required"`
	Components []TemplateComponent `json:"components" binding:"required"`
	Preview    string              `json:"preview"`
	CreatedAt  string              `json:"created_at"`
}

// TemplateComponent 模板组件
type TemplateComponent struct {
	Label string `json:"label"`
	Type  string `json:"type"`
}

// AIConfig AI 配置
type AIConfig struct {
	Provider string `json:"provider"`
	APIKey   string `json:"api_key"`
	BaseURL  string `json:"base_url"`
	Model    string `json:"model"`
}

// AITestRequest AI 测试连接请求
type AITestRequest struct {
	Provider string `json:"provider" binding:"required"`
	BaseURL  string `json:"base_url"`
	APIKey   string `json:"api_key"`
	Model    string `json:"model"`
}

// AIAnalyzeRequest AI 分析请求
type AIAnalyzeRequest struct {
	FilePath    string `json:"file_path" binding:"required"`
	AnalyzeType string `json:"analyze_type"`
}

// Rule 文件处理规则
type Rule struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Icon             string   `json:"icon"`
	Color            string   `json:"color"`
	Destination      string   `json:"destination"`
	Action           string   `json:"action"`
	KeepOriginal     bool     `json:"keep_original"`
	FileTypes        []string `json:"file_types"`
	CustomExtensions []string `json:"custom_extensions"`
	AllowAllFiles    bool     `json:"allow_all_files"`
	NameTemplate     []string `json:"name_template"`
	DateSource       string   `json:"date_source"`
	AIEnabled        bool     `json:"ai_enabled"`
	QuickAccess      bool     `json:"quick_access"`
	Enabled          bool     `json:"enabled"`
	CreatedAt        string   `json:"created_at,omitempty"`
	UpdatedAt        string   `json:"updated_at,omitempty"`
}
