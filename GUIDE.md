# BlackHole - 使用指南

## 项目概述

BlackHole 是一个基于 **Electron + Vue + Go** 的前后端分离悬浮球应用，专为 macOS 设计。

### 核心功能

✅ **桌面悬浮球** - 可自由拖拽的圆形悬浮窗
✅ **原生右键菜单** - 包含设置、显示模式切换、退出功能
✅ **前后端分离** - Vue 前端 + Go 后端 API
✅ **显示模式切换** - 支持"所有桌面显示"和"仅当前桌面显示"

## 快速开始

### 1. 安装依赖

```bash
npm install
```

### 2. 启动应用

```bash
npm run dev
```

这会同时启动：
- ✅ Vite 开发服务器（端口 5173）
- ✅ Electron 应用（悬浮球窗口）
- ✅ Go 后端服务（端口 8080）

### 3. 使用悬浮球

- **拖拽移动**：左键按住拖拽悬浮球
- **右键菜单**：右键点击悬浮球显示菜单
  - 设置 - 打开设置页面（待实现）
  - 显示模式切换 - 切换"所有桌面"或"仅当前桌面"
  - 退出应用 - 关闭应用

## 技术架构

### 前端 (Vue 3 + Vite)
- **框架**: Vue 3 Composition API
- **构建工具**: Vite 5
- **样式**: CSS3 渐变、阴影效果
- **功能**: 拖拽、右键菜单、事件监听

### 桌面 (Electron)
- **主进程**: [electron/main.js](electron/main.js)
  - 创建无边框透明窗口
  - 管理窗口层级和显示模式
  - 启动 Go 后端进程
  - 处理原生菜单
- **预加载脚本**: [electron/preload.js](electron/preload.js)
  - 安全的 IPC 通信桥接
  - 暴露受控的 API 给渲染进程

### 后端 (Go)
- **HTTP 服务器**: 监听 8080 端口
- **API 端点**:
  - `GET /api/health` - 健康检查
  - `GET /api/status` - 获取服务状态
- **CORS**: 支持跨域请求
- **JSON 响应**: 统一的响应格式

## 项目结构

```
BlackHole/
├── electron/              # Electron 主进程
│   ├── main.js           # 应用入口、窗口管理
│   └── preload.js        # IPC 通信桥接
├── src/                  # Vue 前端
│   ├── App.vue           # 悬浮球组件
│   └── main.js           # Vue 入口
├── backend/              # Go 后端
│   ├── main.go           # HTTP 服务器
│   └── go.mod            # Go 模块配置
├── index.html            # HTML 模板
├── vite.config.js        # Vite 配置
├── package.json          # 项目配置
└── README.md             # 项目说明
```

## 开发说明

### 窗口配置
- **尺寸**: 80x80 像素
- **特性**: 无边框、透明、置顶
- **位置**: 默认右下角
- **拖拽**: 支持全窗口拖拽

### 显示模式
1. **所有桌面显示** (默认)
   - `setVisibleOnAllWorkspaces(true)`
   - 在所有桌面和全屏应用中可见
   
2. **仅当前桌面显示**
   - `setVisibleOnAllWorkspaces(false)`
   - 只在当前桌面可见

### API 测试

```bash
# 测试健康检查
curl http://localhost:8080/api/health

# 测试状态接口
curl http://localhost:8080/api/status
```

或使用提供的测试脚本：
```bash
./test-api.sh
```

## 构建打包

### 构建前端
```bash
npm run build
```

### 打包 Mac 应用
```bash
npm run build:mac
```

生成的应用在 `release/` 目录。

## 常见问题

### 1. 端口占用
如果 8080 端口被占用：
```bash
lsof -ti:8080 | xargs kill -9
```

### 2. 悬浮球不显示
- 检查 Electron 应用是否启动
- 查看控制台是否有错误
- 确认 Vite 开发服务器正常运行

### 3. Go 后端启动失败
- 确认 Go 已安装（`go version`）
- 检查 backend/main.go 语法
- 查看终端错误信息

## 下一步开发

### 待实现功能
- [ ] 设置页面
- [ ] 自定义悬浮球样式
- [ ] 托盘图标
- [ ] 快捷键支持
- [ ] 数据持久化
- [ ] 更多 API 接口

### 代码改进
- [ ] 添加 TypeScript 支持
- [ ] 单元测试
- [ ] 错误处理优化
- [ ] 日志系统

## License

MIT

---

**提示**: 这是一个开发中的项目，欢迎贡献代码和提出建议！
