# Backend 项目结构说明

## 目录结构

```
backend/
├── main.go              # 主入口文件
├── go.mod               # Go 模块依赖
├── go.sum               # 依赖版本锁定
├── database/            # 数据库相关
│   └── db.go           # 数据库初始化和操作
├── handlers/            # HTTP 处理器
│   ├── system.go       # 系统相关（health, status）
│   ├── file.go         # 文件处理相关
│   ├── ai.go           # AI 相关
│   └── template.go     # 模板管理相关
├── models/              # 数据模型
│   └── models.go       # 所有数据结构定义
├── routes/              # 路由配置
│   └── routes.go       # API 路由设置
└── services/            # 业务逻辑服务
    └── service.go      # 通用服务函数
```

## 技术栈

- **Web 框架**: [Gin](https://gin-gonic.com/) - 高性能 HTTP web 框架
- **数据库**: SQLite3 - 轻量级嵌入式数据库
- **Go 版本**: 1.21+

## 主要功能模块

### 1. 数据库模块 (`database/`)
- `Init()`: 初始化数据库连接和表结构
- `SaveHistory()`: 保存文件处理历史
- `GetHistory()`: 获取历史记录列表
- `ClearHistory()`: 清除所有历史记录

### 2. 处理器模块 (`handlers/`)

#### system.go - 系统接口
- `Health()`: 健康检查
- `Status()`: 获取服务状态

#### file.go - 文件处理
- `ProcessFile()`: 处理文件（重命名、移动、复制）
- `GetHistory()`: 获取文件处理历史
- `ClearHistory()`: 清除历史记录

#### ai.go - AI 功能
- `GetOllamaModels()`: 获取可用的 Ollama 模型列表
- `TestAIConnection()`: 测试 AI 服务连接
- `GetAIConfig()`: 获取当前 AI 配置
- `SaveAIConfig()`: 保存 AI 配置
- `AnalyzeFile()`: 使用 AI 分析文件

#### template.go - 模板管理
- `GetTemplates()`: 获取所有模板
- `ImportTemplate()`: 导入新模板
- `DeleteTemplate()`: 删除指定模板

### 3. 服务模块 (`services/`)
- `CopyFile()`: 文件复制工具
- `AnalyzeFileWithOllama()`: 调用 Ollama API 分析文件
- `GenerateTemplatePreview()`: 生成模板预览
- `TestAIConnection()`: 测试 AI 连接
- `GetOllamaModels()`: 获取 Ollama 模型列表

### 4. 模型模块 (`models/`)
定义所有数据结构：
- `Response`: 统一 API 响应格式
- `FileProcessRequest/Response`: 文件处理请求/响应
- `AIAnalysis`: AI 分析结果
- `HistoryRecord`: 历史记录
- `Template`: 模板结构
- `AIConfig`: AI 配置

### 5. 路由模块 (`routes/`)
- `SetupRoutes()`: 配置所有 API 路由

## API 端点

### 系统相关
- `GET /api/health` - 健康检查
- `GET /api/status` - 获取服务状态

### 文件处理
- `POST /api/files/process` - 处理文件

### 历史记录
- `GET /api/history` - 获取历史记录
- `POST /api/history/clear` - 清除历史记录

### 模板管理
- `GET /api/templates` - 获取模板列表
- `POST /api/templates/import` - 导入模板
- `DELETE /api/templates/:id` - 删除模板

### AI 功能
- `GET /api/ollama/models` - 获取 Ollama 模型列表
- `POST /api/ai/test-connection` - 测试 AI 连接
- `GET /api/ai/config` - 获取 AI 配置
- `POST /api/ai/config` - 保存 AI 配置
- `POST /api/ai/analyze` - AI 分析文件

## 运行说明

### 开发环境
```bash
# 安装依赖
go mod download

# 运行服务
go run main.go

# 或者编译后运行
go build -o blackhole-backend
./blackhole-backend
```

### 生产环境
```bash
# 编译
go build -o blackhole-backend

# 运行
./blackhole-backend
```

服务将在 `http://localhost:8080` 启动。

## 特性

### 1. 模块化设计
- 代码按功能分层，便于维护和扩展
- 清晰的职责分离（路由、处理、服务、数据）

### 2. Gin 框架优势
- 高性能路由
- 中间件支持
- 参数绑定和验证
- JSON 序列化/反序列化

### 3. CORS 支持
- 跨域请求支持，方便前端调用

### 4. 统一响应格式
所有 API 返回统一的响应结构：
```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

## 下一步优化建议

1. **添加配置文件**: 使用 viper 管理配置
2. **日志增强**: 使用 logrus 或 zap 替代标准日志
3. **错误处理**: 统一错误处理中间件
4. **数据持久化**: 模板存储移到数据库
5. **单元测试**: 为各个模块添加测试用例
6. **API 文档**: 集成 Swagger 自动生成 API 文档
7. **性能监控**: 添加性能指标收集
8. **安全增强**: API 鉴权、限流等
