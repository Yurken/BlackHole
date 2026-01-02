# BlackHole 🕳️

一个智能文件整理桌面悬浮球应用，基于 Electron + Vue 3 + Go 构建。支持 AI 视觉识别，自动根据文件内容生成有意义的文件名。

## ✨ 功能特性

### 核心功能
- 🎯 **桌面悬浮球** - 可自由拖拽，始终置顶
- 📁 **拖放文件处理** - 拖动文件到悬浮球自动重命名和整理
- 🤖 **AI 智能识别** - 支持图片和 PDF 文档内容识别
- 📊 **进度显示** - 实时显示多文件处理进度（如 `2/5 文件名...`）
- ✅ **成功反馈** - 处理完成后显示绿色对勾和提示音

### 视觉效果
- 🌀 **漩涡动画** - 悬停时触发青色螺旋旋转效果
- 🎨 **拖放动画** - 文件拖入时显示吸收漩涡特效
- 🔊 **音效反馈** - 成功处理播放提示音

### 智能特性
- 🏷️ **自定义规则** - 按文件类型、扩展名配置处理规则
- 📝 **命名模板** - 支持日期、原名、AI 名称等组件
- 🔄 **操作模式** - 复制或移动文件
- 📜 **历史记录** - 记录所有文件处理历史

## 🛠️ 技术栈

| 层级 | 技术 | 版本 |
|------|------|------|
| **前端框架** | Vue 3 | Composition API |
| **构建工具** | Vite | ^4.4.5 |
| **桌面框架** | Electron | ^28.0.0 |
| **后端语言** | Go | 1.21+ |
| **Web 框架** | Gin | - |
| **数据库** | SQLite | - |
| **AI 模型** | Ollama | qwen3-vl:4b |

## 🚀 快速开始

### 前置要求

1. **Node.js** >= 18.0.0
2. **Go** >= 1.21
3. **Ollama** (用于 AI 识别)
   ```bash
   # 安装 Ollama
   brew install ollama
   
   # 启动 Ollama 服务
   ollama serve
   
   # 下载视觉模型
   ollama pull qwen3-vl:4b
   ```

### 安装依赖

```bash
npm install
cd backend && go mod tidy && cd ..
```

### 启动开发环境

```bash
npm run start
```

这会同时启动：
- ✅ Vite 开发服务器 (端口 5173)
- ✅ Electron 应用窗口
- ✅ Go 后端服务 (端口 18620)

### 使用方法

1. **拖放文件** - 将文件拖到悬浮球上
2. **自动处理** - AI 识别内容并重命名
3. **查看进度** - 悬浮球下方显示处理进度
4. **配置规则** - 右键菜单 → 设置

## 📦 构建打包

```bash
# 构建 macOS 应用
npm run build:mac

# 构建后的应用在 dist/ 目录
```

## 🔧 API 接口

后端服务运行在 `http://localhost:18620`

### 核心接口

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/health` | GET | 健康检查 |
| `/api/files/process` | POST | 处理文件 |
| `/api/rules` | GET/POST | 规则管理 |
| `/api/ai/config` | GET/POST | AI 配置 |
| `/api/ai/test-connection` | POST | 测试 AI 连接 |
| `/api/history` | GET | 历史记录 |
| `/api/ollama/models` | GET | 获取可用模型 |

详细文档见 [BACKEND_API.md](BACKEND_API.md)

## 📂 项目结构

```
BlackHole/
├── electron/              # Electron 主进程
│   ├── main.js           # 窗口管理、菜单、Dock 控制
│   └── preload.js        # 安全桥接（IPC 通信）
├── src/                  # Vue 前端
│   ├── App.vue           # 悬浮球组件（拖放、动画）
│   ├── Settings.vue      # 设置页面
│   ├── main.js           # Vue 应用入口
│   └── assets/           # 静态资源
├── backend/              # Go 后端
│   ├── main.go           # 服务入口
│   ├── handlers/         # API 处理器
│   ├── services/         # 业务逻辑（AI 识别）
│   ├── database/         # SQLite 数据库
│   ├── models/           # 数据模型
│   └── routes/           # 路由配置
├── settings.html         # 设置页面独立窗口
├── vite.config.js        # Vite 构建配置
└── package.json          # 项目配置
```

## 🤖 AI 功能说明

### 支持的文件类型

| 类型 | 处理方式 |
|------|----------|
| **图片** (jpg, png, gif, etc.) | 使用 `qwen3-vl:4b` 视觉模型识别图片内容 |
| **PDF** | 自动提取第一页为图片，再用视觉模型识别 |
| **其他** | 根据文件名生成新名称 |

### PDF 处理机制

1. 使用 macOS `qlmanage` 将 PDF 第一页转为图片
2. 发送图片给视觉模型识别
3. 根据文档内容生成文件名
4. 自动删除临时图片

### 性能优化

- `/no_think` 指令 - 让模型直接输出，跳过思考过程
- 限制输出 token 数量（100 tokens）
- 超时设置 60 秒
- 失败回退机制 - AI 无响应时保留原文件名

## ⚙️ 配置说明

### 命名模板组件

- `{year}` - 四位年份 (2025)
- `{month}` - 两位月份 (01-12)
- `{day}` - 两位日期 (01-31)
- `{原文件名}` - 原始文件名（AI 开启时为 AI 生成名）
- 自定义文本 - 直接输入任意文字

### AI 提供商

支持以下 AI 服务：
- **Ollama** (本地) - 推荐用于隐私保护
- **OpenAI** - GPT-4V 等云端模型
- **DeepSeek** - 国内可用
- **Qwen** - 阿里云通义千问

## 🎨 UI 特性

### 悬浮球状态

| 状态 | 外观 |
|------|------|
| 空闲 | 半透明青色光晕 |
| 悬停 | 螺旋旋转动画 |
| 拖入 | 漩涡吸收效果 + 青色高亮 |
| 处理中 | 显示进度（2/5） + 文件名 |
| 成功 | 绿色对勾 + 音效 |

### Dock 行为

- **桌面显示模式** - 隐藏 Dock 图标
- **打开设置** - 显示 Dock 图标
- 右键菜单切换显示模式

## 🔍 调试日志

后端会输出详细的 AI 处理日志：

```
[AI] 正在将 PDF 转换为图片...
[AI] PDF 已转换为图片: /tmp/doc.pdf.png
[AI] 发送请求到 /api/chat: {...}
[AI] 模型: qwen3-vl:4b, 文件: example.pdf
[AI] 原始响应: {"suggested_name": "研究论文", ...}
[AI] 解析结果: suggested_name=研究论文, category=文档
```
