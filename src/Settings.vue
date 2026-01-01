<template>
  <div class="settings-container">
    <!-- 标签页 -->
    <div class="tabs">
      <button 
        v-for="tab in tabs" 
        :key="tab.id"
        class="tab"
        :class="{ active: currentTab === tab.id }"
        @click="currentTab = tab.id"
      >
        {{ tab.label }}
      </button>
    </div>
    
    <!-- 主体区域 -->
    <div class="main-body">
      <!-- 侧边栏 - 仅在规则页面显示 -->
      <div v-show="currentTab === 'rule'" class="sidebar">
        <div class="sidebar-header">
          <h2>规则列表</h2>
        </div>
        
        <div class="rules-list">
          <div 
            v-for="rule in rules" 
            :key="rule.id"
            class="rule-item"
            :class="{ active: currentRuleId === rule.id }"
            @click="selectRule(rule.id)"
          >
            <div class="rule-icon" :style="{ background: rule.color }">
              <img class="rule-icon-img" :src="getIcon(rule.icon)" alt="" />
            </div>
            <div class="rule-info">
              <div class="rule-name">{{ rule.name }}</div>
              <div class="rule-path">{{ rule.destination }}</div>
            </div>
          </div>
        </div>
        
        <button class="add-rule-btn" @click="addNewRule">
          <span>+</span> 添加规则
        </button>
      </div>

      <!-- 主内容区 -->
      <div class="main-content">
        <div class="content-wrapper">
          <!-- 规则配置 -->
          <div v-show="currentTab === 'rule'" class="tab-content">
            <div v-if="currentRule">
            <div class="section">
              <h3>规则基本信息</h3>
            <div class="form-group">
              <label>规则名称</label>
              <input v-model="currentRule.name" type="text" placeholder="例如：图片收集器" />
            </div>

            <div class="form-group">
              <label>规则图标</label>
              <button class="icon-picker-trigger" @click="openIconPicker">
                <span class="icon-preview">
                  <img :src="getIcon(currentRule.icon)" alt="" />
                </span>
                <span>选择图标</span>
              </button>
            </div>

            <div class="form-group">
              <label>图标颜色</label>
              <div class="color-picker">
                <div 
                  v-for="color in colors" 
                  :key="color"
                  class="color-option"
                  :style="{ background: color }"
                  :class="{ selected: currentRule.color === color }"
                  @click="currentRule.color = color"
                ></div>
              </div>
            </div>
          </div>

          <div class="section">
            <h3><img class="section-icon" :src="icons.folder" alt="" /> 目标路径</h3>
            <div class="form-group">
              <div class="path-input">
                <input 
                  v-model="currentRule.destination" 
                  type="text" 
                  placeholder="/Users/sen/Documents" 
                  readonly
                />
                <button @click="selectDestination">选择</button>
              </div>
            </div>
          </div>

          <div class="section">
            <h3><img class="section-icon" :src="icons.repeat" alt="" /> 操作模式</h3>
            <div class="button-group">
              <button 
                :class="{ active: currentRule.action === 'copy' }"
                @click="currentRule.action = 'copy'"
              >
                复制
              </button>
              <button 
                :class="{ active: currentRule.action === 'move' }"
                @click="currentRule.action = 'move'"
              >
                移动
              </button>
            </div>
          </div>

          <div class="section">
            <h3><img class="section-icon" :src="icons.filter" alt="" /> 文件类型过滤</h3>
            
            <div class="file-types">
              <button 
                v-for="type in fileTypes" 
                :key="type.id"
                class="file-type-btn"
                :class="{ selected: currentRule.fileTypes.includes(type.id) }"
                @click="toggleFileType(type.id)"
              >
                <img class="file-type-icon" :src="icons[type.icon]" alt="" />
                {{ type.label }}
              </button>
            </div>

            <div class="form-group">
              <label>自定义扩展名</label>
              <div class="extension-input">
                <input 
                  v-model="customExtension"
                  type="text" 
                  placeholder="如 psd, sketch"
                  @keyup.enter="addCustomExtension"
                />
                <button @click="addCustomExtension">添加</button>
              </div>
              <div class="extension-tags">
                <span 
                  v-for="ext in currentRule.customExtensions" 
                  :key="ext"
                  class="tag"
                >
                  {{ ext }}
                  <button @click="removeExtension(ext)">×</button>
                </span>
              </div>
            </div>

            <div class="form-group checkbox-group">
              <label>
                <input v-model="currentRule.allowAllFiles" type="checkbox" />
                允许所有文件类型
              </label>
            </div>
          </div>

          <div class="section">
            <h3><img class="section-icon" :src="icons.zap" alt="" /> 智能识别（文档 & 图片）</h3>
            <div class="ai-toggle">
              <label class="switch">
                <input v-model="currentRule.aiEnabled" type="checkbox" />
                <span class="slider"></span>
              </label>
            </div>
            <div class="ai-info">
              <p>开启后，AI 将分析文件内容并生成最佳文件名（最多20字），替换原文件名。</p>
            </div>
          </div>

          <div class="section">
            <h3><img class="section-icon" :src="icons.tag" alt="" /> 命名模板</h3>
            <div class="template-builder">
              <div class="template-components">
                <button 
                  v-for="component in nameComponents" 
                  :key="component.id"
                  @click="addNameComponent(component)"
                >
                  {{ component.label }}
                </button>
              </div>
              
              <div class="template-preview">
                <div class="preview-label">预览效果</div>
                <div class="preview-text">
                  example.pdf → {{ generatePreviewName() }}
                </div>
              </div>
            </div>

            <div class="template-parts">
              <div 
                v-for="(part, index) in currentRule.nameTemplate" 
                :key="index"
                class="template-part"
              >
                <span>{{ getComponentLabel(part) }}</span>
                <button @click="removeNameComponent(index)">×</button>
              </div>
            </div>
            
            <div v-if="currentRule.aiEnabled" class="ai-hint">
              <span class="hint-icon">💡</span> 已开启智能识别，模板中的"原文件名"将被 AI 生成的名称替换
            </div>
          </div>

          <div class="section">
            <h3><img class="section-icon" :src="icons.calendar" alt="" /> 日期来源</h3>
            <div class="date-sources">
              <label v-for="source in dateSources" :key="source.id">
                <input 
                  v-model="currentRule.dateSource" 
                  type="radio" 
                  :value="source.id"
                />
                {{ source.label }}
              </label>
            </div>
            <div v-if="currentRule.dateSource === 'content'" class="info-text">
              使用前请确认已开启智能识别功能
            </div>
          </div>

          <div class="section">
            <h3><img class="section-icon" :src="icons.star" alt="" /> 快捷访问</h3>
            <div class="form-group checkbox-group">
              <label>
                <input v-model="currentRule.quickAccess" type="checkbox" />
                显示在卫星上（最多 5 个）
              </label>
            </div>
          </div>

          <div class="section actions">
            <button class="btn-primary" @click="saveRule">保存规则</button>
            <button class="btn-toggle" @click="toggleRuleStatus">
              {{ currentRule.enabled ? '规则已启用' : '规则已禁用' }}
            </button>
            <button class="btn-danger" @click="deleteRule">删除规则</button>
          </div>
          </div>

          <div v-else class="empty-state">
            <p>请选择或创建一个规则</p>
          </div>
        </div>

        <!-- 常规标签页 -->
        <div v-show="currentTab === 'general'" class="tab-content">
          <!-- 启动 -->
          <div class="section">
            <h3><img class="section-icon" :src="icons.power" alt="" /> 启动</h3>
            <div class="setting-row">
              <div class="setting-info">
                <div class="setting-label">登录时启动 BlackHole</div>
              </div>
              <label class="switch">
                <input v-model="generalSettings.autoLaunch" type="checkbox" />
                <span class="slider"></span>
              </label>
            </div>
            <div class="setting-row no-border">
              <div class="setting-info">
                <div class="setting-label">状态</div>
              </div>
              <div class="setting-value text-muted">{{ generalSettings.autoLaunch ? '已启用' : '未找到' }}</div>
            </div>
          </div>

          <!-- 音效 -->
          <div class="section">
            <h3><img class="section-icon" :src="icons['volume-2']" alt="" /> 音效</h3>
            <div class="setting-row">
              <div class="setting-info">
                <div class="setting-label">静音模式</div>
                <div class="setting-description">关闭所有音效反馈</div>
              </div>
              <label class="switch">
                <input v-model="generalSettings.muteSound" type="checkbox" />
                <span class="slider"></span>
              </label>
            </div>
          </div>

          <!-- 外观 -->
          <div class="section">
            <h3><img class="section-icon" :src="icons.palette" alt="" /> 外观</h3>
            <div class="setting-row">
              <div class="setting-info">
                <div class="setting-label">闲置时自动变淡</div>
              </div>
              <label class="switch">
                <input v-model="generalSettings.autoFade" type="checkbox" />
                <span class="slider"></span>
              </label>
            </div>
            <div class="setting-row no-border">
              <div class="setting-info">
                <div class="setting-label">闲置透明度</div>
              </div>
              <div class="slider-container">
                <input 
                  v-model="generalSettings.idleOpacity" 
                  type="range" 
                  min="0" 
                  max="100" 
                  class="opacity-slider"
                />
                <span class="slider-value">{{ generalSettings.idleOpacity }}%</span>
              </div>
            </div>
          </div>

          <!-- 仅在桌面显示 -->
          <div class="section">
            <h3><img class="section-icon" :src="icons.monitor" alt="" /> 仅在桌面显示</h3>
            <div class="setting-row">
              <label class="switch">
                <input v-model="generalSettings.desktopOnly" type="checkbox" />
                <span class="slider"></span>
              </label>
            </div>
            <div class="info-banner">
              启用后，悬浮球将不再置顶，会被其他窗口遮挡。
            </div>
          </div>

          <!-- 语言 -->
          <div class="section">
            <h3><img class="section-icon" :src="icons.globe" alt="" /> 语言</h3>
            <div class="setting-row">
              <div class="setting-info">
                <div class="setting-label">语言</div>
                <div class="setting-description">切换语言后需重启应用后生效</div>
              </div>
              <select v-model="generalSettings.language" class="language-select">
                <option value="system">跟随系统</option>
                <option value="zh-CN">简体中文</option>
                <option value="en-US">English</option>
              </select>
            </div>
          </div>

          <!-- 快捷键 -->
          <div class="section">
            <h3><img class="section-icon" :src="icons.keyboard" alt="" /> 快捷键</h3>
            <div class="setting-row">
              <div class="setting-info">
                <div class="setting-label">显示/隐藏悬浮窗</div>
                <div class="setting-description">在任意位置按下此快捷键可显示/隐藏悬浮窗</div>
              </div>
              <div class="shortcut-display">
                <kbd>⌃</kbd>
                <kbd>⌘</kbd>
                <kbd>O</kbd>
              </div>
            </div>
            <div class="setting-row no-border">
              <div class="setting-info">
                <div class="setting-label">唤醒时跟随鼠标标</div>
                <div class="setting-description">关闭后，悬浮球将保持在屏幕上的固定位置</div>
              </div>
              <label class="switch">
                <input v-model="generalSettings.followMouse" type="checkbox" />
                <span class="slider"></span>
              </label>
            </div>
          </div>

          <!-- 关于 -->
          <div class="section">
            <h3><img class="section-icon" :src="icons.info" alt="" /> 关于</h3>
            <div class="setting-row">
              <div class="setting-info">
                <div class="setting-label">版本</div>
              </div>
              <div class="setting-value">1.3.1</div>
            </div>
            <div class="setting-row">
              <div class="setting-info">
                <div class="setting-label">构建</div>
              </div>
              <div class="setting-value">20251232</div>
            </div>
            <div class="setting-row no-border">
              <div class="setting-info">
                <div class="setting-label">联系我们</div>
              </div>
              <a href="mailto:jiangzhewen888@163.com" class="setting-value link">
                jiangzhewen888@163.com
              </a>
            </div>
          </div>
        </div>

        <!-- 模板标签页 -->
        <div v-show="currentTab === 'template'" class="tab-content template-page">
          <!-- 系统推荐 -->
          <div class="section">
            <h3><img class="section-icon" :src="icons.star" alt="" /> 系统推荐</h3>
            <div class="template-list">
              <div 
                v-for="template in systemTemplates" 
                :key="template.id"
                class="template-item"
              >
                <div class="template-icon">
                  <img :src="icons['file-text']" alt="" />
                </div>
                <div class="template-info">
                  <div class="template-name">
                    {{ template.name }}
                    <span class="system-badge">系统</span>
                  </div>
                  <div class="template-components">
                    <span 
                      v-for="(comp, idx) in template.components" 
                      :key="idx"
                      class="component-tag"
                      :class="comp.type"
                    >
                      {{ comp.label }}
                    </span>
                  </div>
                  <div class="template-preview">
                    {{ template.preview }}
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 我的模板 -->
          <div class="section">
            <div class="section-header">
              <h3><img class="section-icon" :src="icons.user" alt="" /> 我的模板</h3>
              <button class="import-btn" @click="importTemplate">
                <img class="button-icon" :src="icons.upload" alt="" /> 导入
              </button>
            </div>
            
            <div v-if="userTemplates.length === 0" class="empty-template">
              <div class="empty-icon">
                <img :src="icons['file-plus']" alt="" />
              </div>
              <div class="empty-text">暂无自定义模板</div>
            </div>

            <div v-else class="template-list">
              <div 
                v-for="template in userTemplates" 
                :key="template.id"
                class="template-item"
              >
                <div class="template-icon">
                  <img :src="icons['file-text']" alt="" />
                </div>
                <div class="template-info">
                  <div class="template-name">{{ template.name }}</div>
                  <div class="template-components">
                    <span 
                      v-for="(comp, idx) in template.components" 
                      :key="idx"
                      class="component-tag"
                      :class="comp.type"
                    >
                      {{ comp.label }}
                    </span>
                  </div>
                  <div class="template-preview">
                    {{ template.preview }}
                  </div>
                </div>
                <button class="delete-template-btn" @click="deleteTemplate(template.id)">
                  删除
                </button>
              </div>
            </div>
          </div>

          <!-- 底部信息 -->
          <div class="template-footer">
            <div class="template-count">
              共 {{ totalTemplateCount }} 个模板
            </div>
            <div class="template-hint">
              <img class="inline-icon" :src="icons.info" alt="" />
              在规则编辑器中保存模板
            </div>
          </div>
        </div>

        <!-- 本地AI标签页 -->
        <div v-show="currentTab === 'ai'" class="tab-content ai-config-page">
          <div class="ai-layout">
            <!-- 左侧提供商列表 -->
            <div class="ai-providers">
              <h3 class="providers-title">提供商</h3>
              <div class="provider-list">
                <div 
                  v-for="provider in aiProviders" 
                  :key="provider.id"
                  class="provider-item"
                  :class="{ active: selectedProvider === provider.id }"
                  @click="selectedProvider = provider.id"
                >
                  <img class="provider-icon" :src="icons[provider.icon]" alt="" />
                  <span class="provider-name">{{ provider.name }}</span>
                </div>
              </div>
            </div>

            <!-- 右侧配置区域 -->
            <div class="ai-config">
              <div v-if="currentProvider" class="config-content">
                <h2 class="config-title">{{ currentProvider.name }}</h2>
                
                <div class="config-status" :class="currentProvider.connectionStatus">
                  <span class="status-text">
                    <template v-if="currentProvider.connectionStatus === 'connected'">
                      ✓ 连接正常
                    </template>
                    <template v-else-if="currentProvider.connectionStatus === 'failed'">
                      ✗ 连接失败
                    </template>
                    <template v-else-if="currentProvider.connectionStatus === 'testing'">
                      测试中...
                    </template>
                    <template v-else>
                      未测试
                    </template>
                    <span v-if="currentProvider.lastTestTime" class="test-time">
                      ({{ currentProvider.lastTestTime }})
                    </span>
                  </span>
                  <button 
                    class="btn-auto-refresh" 
                    @click="testConnection(null, true)"
                    :disabled="currentProvider.connectionStatus === 'testing'"
                    title="刷新连接状态"
                  >
                    <img :src="icons['refresh-cw']" alt="" />
                  </button>
                </div>

                <div class="config-form">
                  <!-- Ollama 不需要 API Key -->
                  <div v-if="currentProvider.id !== 'ollama'" class="form-group">
                    <label class="form-label">
                      API Key
                      <img class="help-icon" :src="icons.info" alt="" />
                    </label>
                    <div class="password-input">
                      <input 
                        v-model="currentProvider.apiKey"
                        :type="showApiKey ? 'text' : 'password'"
                        placeholder="输入 API Key"
                        class="form-input"
                      />
                      <button 
                        class="toggle-visibility"
                        @click="showApiKey = !showApiKey"
                      >
                        <img :src="showApiKey ? icons['eye-off'] : icons.eye" alt="" />
                      </button>
                    </div>
                  </div>

                  <div class="form-group">
                    <label class="form-label">API Base URL</label>
                    <input 
                      v-model="currentProvider.baseUrl"
                      type="text"
                      :placeholder="currentProvider.defaultBaseUrl"
                      class="form-input"
                    />
                  </div>

                  <!-- Ollama 模型选择：从服务器加载 -->
                  <div v-if="currentProvider.id === 'ollama'" class="form-group">
                    <label class="form-label">
                      模型名称
                      <button class="btn-refresh" @click="loadOllamaModels" :disabled="ollamaModelsLoading">
                        <img
                          class="button-icon"
                          :class="{ 'icon-spin': ollamaModelsLoading }"
                          :src="icons['refresh-cw']"
                          alt=""
                        />
                        刷新
                      </button>
                    </label>
                    <select v-model="currentProvider.model" class="form-select" :disabled="ollamaModelsLoading">
                      <option value="">{{ ollamaModels.length === 0 ? '无可用模型' : '选择模型' }}</option>
                      <option v-for="model in ollamaModels" :key="model" :value="model">
                        {{ model }}
                      </option>
                    </select>
                    <small v-if="ollamaModels.length === 0" class="form-hint error">
                      未检测到模型，请先使用 <code>ollama pull</code> 下载模型
                    </small>
                  </div>

                  <!-- 其他服务商：手动输入模型名称 -->
                  <div v-else class="form-group">
                    <label class="form-label">模型名称</label>
                    <input 
                      v-model="currentProvider.model"
                      type="text"
                      placeholder="例如: gpt-4, claude-3-opus, gemini-pro"
                      class="form-input"
                    />
                  </div>

                  <div class="form-actions">
                    <button class="btn-test" @click="testConnection($event)">
                      测试连接
                    </button>
                    <button class="btn-save" @click="saveAIConfig">
                      保存配置
                    </button>
                  </div>
                </div>

                <div v-if="currentProvider.id === 'ollama'" class="ollama-info">
                  <h4>Ollama 本地部署说明</h4>
                  <ol>
                    <li>安装 Ollama: <code>brew install ollama</code></li>
                    <li>下载模型: <code>ollama pull qwen3-vl:4b</code></li>
                    <li>启动服务: <code>ollama serve</code></li>
                    <li>默认地址: <code>http://localhost:11434</code></li>
                  </ol>
                </div>
              </div>

              <div v-else class="empty-config">
                <p>请选择一个 AI 提供商</p>
              </div>
            </div>
          </div>
        </div>

        <!-- 历史记录标签页 -->
        <div v-show="currentTab === 'history'" class="tab-content history-page">
          <div class="section history-header-section">
            <div class="history-header">
              <h3>历史记录</h3>
              <div class="history-actions">
                <label class="switch-inline">
                  <span class="switch">
                    <input type="checkbox" v-model="privacyMode">
                    <span class="slider"></span>
                  </span>
                  <span class="label-text">隐私模式（关闭记录）</span>
                </label>
                <button class="btn-danger btn-small" @click="clearHistory" :disabled="historyLoading">
                  清除历史记录
                </button>
              </div>
            </div>
            
            <div class="history-stats">
              <div class="stat-item">
                <span class="stat-label">总操作次数</span>
                <span class="stat-value">{{ historyRecords.length }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">成功</span>
                <span class="stat-value success">{{ successCount }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">失败</span>
                <span class="stat-value error">{{ errorCount }}</span>
              </div>
            </div>
          </div>

          <div class="history-body">
            <div v-if="historyLoading" class="loading-state">
              <p>加载中...</p>
            </div>
            
            <div v-else-if="historyRecords.length === 0" class="empty-history">
              <svg width="64" height="64" viewBox="0 0 64 64" fill="none">
                <circle cx="32" cy="32" r="28" stroke="#CCC" stroke-width="2"/>
                <path d="M32 16V32L40 40" stroke="#CCC" stroke-width="2" stroke-linecap="round"/>
              </svg>
              <p>暂无操作历史</p>
              <small>拖动文件到悬浮球开始使用</small>
            </div>

            <div v-else class="history-list">
              <div v-for="record in historyRecords" :key="record.id" class="history-item">
                <div class="history-status" :class="normalizeHistoryStatus(record.status)">
                  <svg v-if="normalizeHistoryStatus(record.status) === 'success'" width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
                    <path d="M10 0C4.48 0 0 4.48 0 10s4.48 10 10 10 10-4.48 10-10S15.52 0 10 0zm-2 15l-5-5 1.41-1.41L8 12.17l7.59-7.59L17 6l-9 9z"/>
                  </svg>
                  <svg v-else width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
                    <path d="M10 0C4.48 0 0 4.48 0 10s4.48 10 10 10 10-4.48 10-10S15.52 0 10 0zm1 15H9v-2h2v2zm0-4H9V5h2v6z"/>
                  </svg>
                </div>
                
                <div class="history-content">
                  <div class="history-main">
                    <div class="file-change">
                      <div class="file-info">
                        <span class="file-name original">{{ record.original_name }}</span>
                        <span class="file-path">{{ record.original_path }}</span>
                      </div>
                      <svg class="arrow-icon" width="24" height="24" viewBox="0 0 24 24" fill="none">
                        <path d="M5 12h14m0 0l-7-7m7 7l-7 7" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                      </svg>
                      <div class="file-info">
                        <span class="file-name new">{{ record.new_name }}</span>
                        <span class="file-path">{{ record.new_path }}</span>
                      </div>
                    </div>
                  </div>
                  
                  <div class="history-meta">
                    <span class="meta-item">
                      <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
                        <path d="M7 0C3.13 0 0 3.13 0 7s3.13 7 7 7 7-3.13 7-7-3.13-7-7-7zm0 12.25c-2.89 0-5.25-2.36-5.25-5.25S4.11 1.75 7 1.75 12.25 4.11 12.25 7 9.89 12.25 7 12.25z" fill="currentColor"/>
                        <path d="M7.44 3.5H6.56v4.06l3.5 2.1.44-.73-3.06-1.82V3.5z" fill="currentColor"/>
                      </svg>
                      {{ formatTime(record.timestamp) }}
                    </span>
                    <span v-if="record.rule_name" class="meta-item rule-tag">
                      规则: {{ record.rule_name }}
                    </span>
                    <span class="meta-item action-tag" :class="record.action">
                      {{ record.action === 'copy' ? '已复制' : '已移动' }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        </div>
      </div>
    </div>
    <div v-if="iconPickerOpen" class="icon-picker-overlay" @click.self="closeIconPicker">
      <div class="icon-picker-modal">
        <div class="icon-picker-header">
          <div class="icon-picker-title">选择图标</div>
          <button class="icon-picker-done" @click="closeIconPicker">完成</button>
        </div>
        <div class="icon-picker-tabs">
          <button
            v-for="category in iconCategories"
            :key="category.id"
            class="icon-picker-tab"
            :class="{ active: selectedIconCategory === category.id }"
            @click="selectedIconCategory = category.id"
          >
            {{ category.label }}
          </button>
        </div>
        <div class="icon-picker-grid">
          <button
            v-for="icon in currentIconCategory.icons"
            :key="icon"
            class="icon-picker-item"
            :class="{ selected: currentRule?.icon === icon }"
            @click="selectRuleIcon(icon)"
          >
            <img :src="icons[icon]" alt="" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import { icons, getIcon, resolveIconKey } from './icons'

const currentTab = ref('rule')
const currentRuleId = ref(null)
const customExtension = ref('')
const iconPickerOpen = ref(false)
const selectedIconCategory = ref('file')

const iconCategories = [
  {
    id: 'file',
    label: '文件',
    icons: ['file', 'file-text', 'file-plus', 'file-check', 'file-search', 'folder', 'folder-open', 'folder-plus', 'archive']
  },
  {
    id: 'media',
    label: '媒体',
    icons: ['image', 'camera', 'video', 'film', 'music', 'headphones', 'mic', 'speaker', 'volume-2']
  },
  {
    id: 'symbol',
    label: '符号',
    icons: ['star', 'heart', 'check-circle', 'shield', 'lock', 'key', 'tag', 'puzzle', 'info']
  },
  {
    id: 'item',
    label: '物品',
    icons: ['package', 'gift', 'coffee', 'pen-tool', 'map-pin', 'paperclip', 'pencil', 'briefcase', 'link']
  },
  {
    id: 'business',
    label: '商务',
    icons: ['briefcase', 'bar-chart-3', 'pie-chart', 'credit-card', 'banknote', 'calculator', 'file-text', 'tag', 'folder']
  },
  {
    id: 'device',
    label: '设备',
    icons: ['laptop', 'monitor', 'smartphone', 'printer', 'mouse', 'keyboard', 'hard-drive', 'plug', 'cpu']
  },
  {
    id: 'arrow',
    label: '箭头',
    icons: ['arrow-up', 'arrow-down', 'arrow-left', 'arrow-right', 'repeat', 'refresh-cw', 'corner-down-left', 'corner-down-right', 'shuffle']
  }
]

// 常规设置
const generalSettings = ref({
  autoLaunch: false,
  muteSound: false,
  autoFade: true,
  idleOpacity: 40,
  desktopOnly: false,
  language: 'system',
  followMouse: true
})

const tabs = [
  { id: 'rule', label: '规则' },
  { id: 'general', label: '常规' },
  { id: 'template', label: '模板' },
  { id: 'ai', label: 'AI 配置' },
  { id: 'history', label: '历史记录' }
]

// AI 提供商配置
const selectedProvider = ref('ollama')
const showApiKey = ref(false)

const aiProviders = ref([
  { 
    id: 'anthropic', 
    name: 'Anthropic', 
    icon: 'anthropic',
    apiKey: '',
    baseUrl: '',
    defaultBaseUrl: 'https://api.anthropic.com/v1',
    connectionStatus: 'unknown',
    lastTestTime: null
  },
  { 
    id: 'azure', 
    name: 'Azure OpenAI', 
    icon: 'azure',
    apiKey: '',
    baseUrl: '',
    defaultBaseUrl: 'https://your-resource.openai.azure.com',
    connectionStatus: 'unknown',
    lastTestTime: null
  },
  { 
    id: 'deepseek', 
    name: 'DeepSeek', 
    icon: 'deepseek',
    apiKey: '',
    baseUrl: '',
    defaultBaseUrl: 'https://api.deepseek.com/v1',
    connectionStatus: 'unknown',
    lastTestTime: null
  },
  { 
    id: 'github', 
    name: 'GitHub Copilot', 
    icon: 'github',
    apiKey: '',
    baseUrl: '',
    defaultBaseUrl: 'https://api.githubcopilot.com',
    connectionStatus: 'unknown',
    lastTestTime: null
  },
  { 
    id: 'google', 
    name: 'Google AI', 
    icon: 'google',
    apiKey: '',
    baseUrl: '',
    defaultBaseUrl: 'https://generativelanguage.googleapis.com/v1beta',
    connectionStatus: 'unknown',
    lastTestTime: null
  },
  { 
    id: 'groq', 
    name: 'Groq', 
    icon: 'groq',
    apiKey: '',
    baseUrl: '',
    defaultBaseUrl: 'https://api.groq.com/openai/v1',
    connectionStatus: 'unknown',
    lastTestTime: null
  },
  { 
    id: 'mistral', 
    name: 'Mistral', 
    icon: 'mistral',
    apiKey: '',
    baseUrl: '',
    defaultBaseUrl: 'https://api.mistral.ai/v1',
    connectionStatus: 'unknown',
    lastTestTime: null
  },
  { 
    id: 'ollama', 
    name: 'Ollama', 
    icon: 'ollama',
    apiKey: '',
    baseUrl: '',
    defaultBaseUrl: 'http://localhost:11434',
    model: 'qwen3-vl:4b',
    connectionStatus: 'unknown',
    lastTestTime: null
  },
  { 
    id: 'openai', 
    name: 'OpenAI', 
    icon: 'openai',
    apiKey: '',
    baseUrl: '',
    defaultBaseUrl: 'https://api.openai.com/v1',
    connectionStatus: 'unknown',
    lastTestTime: null
  },
  { 
    id: 'openrouter', 
    name: 'OpenRouter', 
    icon: 'openrouter',
    apiKey: '',
    baseUrl: '',
    defaultBaseUrl: 'https://openrouter.ai/api/v1',
    connectionStatus: 'unknown',
    lastTestTime: null
  },
  { 
    id: 'perplexity', 
    name: 'Perplexity', 
    icon: 'perplexity',
    apiKey: '',
    baseUrl: '',
    defaultBaseUrl: 'https://api.perplexity.ai',
    connectionStatus: 'unknown',
    lastTestTime: null
  },
  { 
    id: 'together', 
    name: 'Together', 
    icon: 'together',
    apiKey: '',
    baseUrl: '',
    defaultBaseUrl: 'https://api.together.xyz/v1',
    connectionStatus: 'unknown',
    lastTestTime: null
  },
  { 
    id: 'xai', 
    name: 'xAI', 
    icon: 'x',
    apiKey: '',
    baseUrl: '',
    defaultBaseUrl: 'https://api.x.ai/v1',
    connectionStatus: 'unknown',
    lastTestTime: null
  }
])

const currentProvider = computed(() => {
  return aiProviders.value.find(p => p.id === selectedProvider.value)
})

// 历史记录
const historyRecords = ref([])
const historyLoading = ref(false)
const privacyMode = ref(false)

// Ollama 模型列表
const ollamaModels = ref([])
const ollamaModelsLoading = ref(false)

function normalizeHistoryStatus(status) {
  return status === 'failed' ? 'error' : status
}

const successCount = computed(() => {
  return historyRecords.value.filter(r => normalizeHistoryStatus(r.status) === 'success').length
})

const errorCount = computed(() => {
  return historyRecords.value.filter(r => normalizeHistoryStatus(r.status) === 'error').length
})

const colors = [
  '#FF4444', '#FF9500', '#FFCC00', '#34C759', 
  '#007AFF', '#AF52DE', '#8E8E93'
]

const fileTypes = [
  { id: 'image', label: '图片', icon: 'image' },
  { id: 'video', label: '视频', icon: 'video' },
  { id: 'document', label: '文档', icon: 'file-text' },
  { id: 'audio', label: '音频', icon: 'music' },
  { id: 'archive', label: '压缩包', icon: 'archive' },
  { id: 'code', label: '代码', icon: 'code' },
  { id: 'installer', label: '安装包', icon: 'package' },
  { id: 'folder', label: '文件夹', icon: 'folder' },
  { id: 'design', label: '设计', icon: 'pen-tool' },
  { id: 'ebook', label: '电子书', icon: 'book' }
]

const nameComponents = [
  { id: 'YYYY', label: 'YYYY' },
  { id: 'MM', label: 'MM' },
  { id: 'separator-', label: '格式 "-"' },
  { id: 'separator_', label: '格式 "_"' },
  { id: 'original', label: '原名' },
  { id: 'custom', label: '添加' }
]

const dateSources = [
  { id: 'current', label: '当前时间' },
  { id: 'created', label: '创建时间' },
  { id: 'modified', label: '修改时间' },
  { id: 'content', label: '内容日期 (AI)' }
]

// 系统推荐模板
const systemTemplates = ref([
  {
    id: 'date-archive',
    name: '日期归档',
    components: [
      { label: '年', type: 'date' },
      { label: '-', type: 'separator' },
      { label: '月', type: 'date' },
      { label: '-', type: 'separator' },
      { label: '日', type: 'date' },
      { label: '-', type: 'separator' },
      { label: '原名', type: 'original' }
    ],
    preview: '2025-12-21_文件名'
  },
  {
    id: 'timestamp',
    name: '时间戳命名',
    components: [
      { label: '年', type: 'date' },
      { label: '_', type: 'separator' },
      { label: '月', type: 'date' },
      { label: '_', type: 'separator' },
      { label: '日', type: 'date' },
      { label: '_', type: 'separator' },
      { label: '时', type: 'time' },
      { label: '_', type: 'separator' },
      { label: '原名', type: 'original' }
    ],
    preview: '20251221_14:30_文件名'
  },
  {
    id: 'sequence',
    name: '序号整理',
    components: [
      { label: '序号', type: 'number' },
      { label: '_', type: 'separator' },
      { label: '原名', type: 'original' }
    ],
    preview: '001_文件名'
  },
  {
    id: 'backup',
    name: '备份副本',
    components: [
      { label: '原名', type: 'original' },
      { label: '_backup', type: 'text' }
    ],
    preview: '文件名_backup'
  },
  {
    id: 'original-only',
    name: '仅原名',
    components: [
      { label: '原名', type: 'original' }
    ],
    preview: '文件名'
  }
])

// 用户自定义模板（从后端加载）
const userTemplates = ref([])

const totalTemplateCount = computed(() => {
  return systemTemplates.value.length + userTemplates.value.length
})

const defaultRules = [
  {
    id: 'rule_images',
    name: '图片收集器',
    icon: 'image',
    color: '#007AFF',
    destination: '/Users/sen/Documents',
    action: 'copy',
    keepOriginal: false,
    fileTypes: ['image'],
    customExtensions: [],
    allowAllFiles: false,
    nameTemplate: ['YYYY', 'separator-', 'MM', 'separator_', 'original'],
    dateSource: 'current',
    aiEnabled: false,
    quickAccess: true,
    enabled: true
  },
  {
    id: 'rule_documents',
    name: '文档归档',
    icon: 'file-text',
    color: '#FF9500',
    destination: '/Users/sen/Documents',
    action: 'move',
    keepOriginal: false,
    fileTypes: ['document'],
    customExtensions: ['psd'],
    allowAllFiles: false,
    nameTemplate: ['YYYY', 'separator-', 'original'],
    dateSource: 'modified',
    aiEnabled: true,
    quickAccess: false,
    enabled: true
  }
]

const rules = ref([])

const currentIconCategory = computed(() => {
  return iconCategories.find(category => category.id === selectedIconCategory.value) || iconCategories[0]
})

function toBackendRule(rule) {
  return {
    id: rule.id,
    name: rule.name,
    icon: resolveIconKey(rule.icon),
    color: rule.color,
    destination: rule.destination,
    action: rule.action,
    keep_original: rule.keepOriginal,
    file_types: rule.fileTypes,
    custom_extensions: rule.customExtensions,
    allow_all_files: rule.allowAllFiles,
    name_template: rule.nameTemplate,
    date_source: rule.dateSource,
    ai_enabled: rule.aiEnabled,
    quick_access: rule.quickAccess,
    enabled: rule.enabled
  }
}

function fromBackendRule(rule) {
  return {
    id: rule.id,
    name: rule.name || '新规则',
    icon: resolveIconKey(rule.icon),
    color: rule.color || colors[0],
    destination: rule.destination || '',
    action: rule.action || 'copy',
    keepOriginal: Boolean(rule.keep_original),
    fileTypes: Array.isArray(rule.file_types) ? rule.file_types : [],
    customExtensions: Array.isArray(rule.custom_extensions) ? rule.custom_extensions : [],
    allowAllFiles: Boolean(rule.allow_all_files),
    nameTemplate: Array.isArray(rule.name_template) && rule.name_template.length > 0 ? rule.name_template : ['original'],
    dateSource: rule.date_source || 'current',
    aiEnabled: Boolean(rule.ai_enabled),
    quickAccess: Boolean(rule.quick_access),
    enabled: rule.enabled !== false
  }
}

function setCurrentRuleId(ruleList) {
  if (ruleList.length === 0) {
    currentRuleId.value = null
    return
  }
  const exists = ruleList.some(rule => rule.id === currentRuleId.value)
  if (!exists) {
    currentRuleId.value = ruleList[0].id
  }
}

async function seedDefaultRules() {
  const created = []
  for (const rule of defaultRules) {
    try {
      const response = await fetch('http://localhost:18620/api/rules', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(toBackendRule(rule))
      })
      const data = await response.json()
      if (data.code === 0 && data.data) {
        created.push(fromBackendRule(data.data))
      }
    } catch (error) {
      console.error('创建默认规则失败:', error)
    }
  }

  rules.value = created.length > 0 ? created : defaultRules.map(rule => ({ ...rule }))
  setCurrentRuleId(rules.value)
}

async function loadRules() {
  try {
    const response = await fetch('http://localhost:18620/api/rules')
    const data = await response.json()
    if (data.code === 0) {
      const loadedRules = (data.data || []).map(fromBackendRule)
      if (loadedRules.length === 0) {
        await seedDefaultRules()
        return
      }
      rules.value = loadedRules
      setCurrentRuleId(loadedRules)
    } else {
      console.error('加载规则失败:', data.message)
    }
  } catch (error) {
    console.error('加载规则失败:', error)
  }
}

const currentRule = computed(() => {
  return rules.value.find(r => r.id === currentRuleId.value)
})

function openIconPicker() {
  if (!currentRule.value) return
  iconPickerOpen.value = true
}

function closeIconPicker() {
  iconPickerOpen.value = false
}

function selectRuleIcon(icon) {
  if (!currentRule.value) return
  currentRule.value.icon = icon
}

function selectRule(id) {
  currentRuleId.value = id
}

function addNewRule() {
  const newRule = {
    id: `rule_${Date.now()}`,
    name: '新规则',
    icon: 'folder',
    color: colors[0],
    destination: '',
    action: 'copy',
    keepOriginal: false,
    fileTypes: [],
    customExtensions: [],
    allowAllFiles: false,
    nameTemplate: ['original'],
    dateSource: 'current',
    aiEnabled: false,
    quickAccess: false,
    enabled: true,
    isNew: true
  }
  rules.value.push(newRule)
  currentRuleId.value = newRule.id
}

function toggleFileType(typeId) {
  const index = currentRule.value.fileTypes.indexOf(typeId)
  if (index > -1) {
    currentRule.value.fileTypes.splice(index, 1)
  } else {
    currentRule.value.fileTypes.push(typeId)
  }
}

function addCustomExtension() {
  if (customExtension.value.trim()) {
    const ext = customExtension.value.trim().toLowerCase()
    if (!currentRule.value.customExtensions.includes(ext)) {
      currentRule.value.customExtensions.push(ext)
    }
    customExtension.value = ''
  }
}

function removeExtension(ext) {
  const index = currentRule.value.customExtensions.indexOf(ext)
  if (index > -1) {
    currentRule.value.customExtensions.splice(index, 1)
  }
}

function addNameComponent(component) {
  currentRule.value.nameTemplate.push(component.id)
}

function removeNameComponent(index) {
  currentRule.value.nameTemplate.splice(index, 1)
}

function getComponentLabel(id) {
  const component = nameComponents.find(c => c.id === id)
  return component ? component.label : id
}

function generatePreviewName() {
  const template = currentRule.value.nameTemplate
  const parts = template.map(part => {
    if (part === 'YYYY') return '2026'
    if (part === 'MM') return '01'
    if (part.startsWith('separator')) return part.replace('separator', '')
    if (part === 'original') return 'example'
    return part
  })
  return parts.join('') + '.pdf'
}

function selectDestination() {
  // 通过 Electron API 打开文件夹选择对话框
  if (window.electronAPI && window.electronAPI.selectFolder) {
    window.electronAPI.selectFolder().then(folderPath => {
      if (folderPath) {
        currentRule.value.destination = folderPath
      }
    })
  } else {
    // 非 Electron 环境下的降级处理
    const path = prompt('请输入目标文件夹路径')
    if (path) {
      currentRule.value.destination = path
    }
  }
}

async function saveRule() {
  if (!currentRule.value) return

  const payload = toBackendRule(currentRule.value)
  try {
    const isNewRule = currentRule.value.isNew
    const url = isNewRule
      ? 'http://localhost:18620/api/rules'
      : `http://localhost:18620/api/rules/${currentRule.value.id}`
    const method = isNewRule ? 'POST' : 'PUT'
    const response = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    
    // 检查响应状态
    if (!response.ok) {
      const text = await response.text()
      console.error('保存失败，服务器响应:', response.status, text)
      alert(`保存失败: 服务器错误 (${response.status})`)
      return
    }
    
    // 获取响应文本
    const text = await response.text()
    if (!text || text.trim() === '') {
      console.error('保存失败：服务器返回空响应')
      alert('保存失败: 服务器返回空响应')
      return
    }
    
    // 解析JSON
    let data
    try {
      data = JSON.parse(text)
    } catch (parseError) {
      console.error('JSON解析失败:', parseError)
      console.error('响应内容:', text)
      alert('保存失败: 服务器响应格式错误')
      return
    }

    if (data.code === 0) {
      const updatedRule = fromBackendRule(data.data || payload)
      const index = rules.value.findIndex(r => r.id === currentRule.value.id)
      if (index > -1) {
        rules.value[index] = updatedRule
      } else {
        rules.value.push(updatedRule)
      }
      currentRuleId.value = updatedRule.id
      alert('规则已保存')
    } else {
      alert('保存失败: ' + data.message)
    }
  } catch (error) {
    console.error('保存规则失败:', error)
    alert('保存失败: ' + error.message)
  }
}

function toggleRuleStatus() {
  currentRule.value.enabled = !currentRule.value.enabled
}

async function deleteRule() {
  if (!confirm('确定要删除这条规则吗？')) return

  const ruleId = currentRuleId.value
  const rule = currentRule.value
  if (!ruleId || !rule) return

  if (rule.isNew) {
    const index = rules.value.findIndex(r => r.id === ruleId)
    if (index > -1) {
      rules.value.splice(index, 1)
      setCurrentRuleId(rules.value)
    }
    return
  }

  try {
    const response = await fetch(`http://localhost:18620/api/rules/${ruleId}`, {
      method: 'DELETE'
    })
    const data = await response.json()

    if (data.code === 0) {
      const index = rules.value.findIndex(r => r.id === ruleId)
      if (index > -1) {
        rules.value.splice(index, 1)
        setCurrentRuleId(rules.value)
      }
    } else {
      alert('删除失败: ' + data.message)
    }
  } catch (error) {
    console.error('删除规则失败:', error)
    alert('删除失败: ' + error.message)
  }
}

// 模板相关函数
async function importTemplate() {
  if (!currentRule.value || !currentRule.value.nameTemplate.length) {
    alert('请先配置命名模板')
    return
  }
  
  const templateName = prompt('请输入模板名称')
  if (!templateName) return
  
  try {
    // 将模板组件转换为后端格式
    const components = currentRule.value.nameTemplate.map(part => {
      if (part === 'YYYY') return { label: '年', type: 'year' }
      if (part === 'MM') return { label: '月', type: 'month' }
      if (part === 'DD') return { label: '日', type: 'day' }
      if (part === 'original') return { label: '原名', type: 'original' }
      if (part.startsWith('separator-')) return { label: part.split('-')[1], type: 'separator' }
      return { label: part, type: 'text' }
    })
    
    const response = await fetch('http://localhost:18620/api/templates/import', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name: templateName,
        components: components
      })
    })
    const data = await response.json()
    
    if (data.code === 0) {
      alert('模板导入成功')
      await loadUserTemplates()
    } else {
      alert('导入失败: ' + data.message)
    }
  } catch (error) {
    console.error('导入模板失败:', error)
    alert('导入失败: ' + error.message)
  }
}

async function deleteTemplate(templateId) {
  if (!confirm('确定要删除这个模板吗？')) return
  
  try {
    const response = await fetch(`http://localhost:18620/api/templates/${templateId}`, {
      method: 'DELETE'
    })
    const data = await response.json()
    
    if (data.code === 0) {
      alert('模板删除成功')
      await loadUserTemplates()
    } else {
      alert('删除失败: ' + data.message)
    }
  } catch (error) {
    console.error('删除模板失败:', error)
    alert('删除失败: ' + error.message)
  }
}

// AI 配置相关函数
async function testConnection(event, silent = false) {
  if (!currentProvider.value) return
  
  const provider = currentProvider.value
  const testingBtn = event?.target
  if (testingBtn) {
    testingBtn.disabled = true
    testingBtn.textContent = '测试中...'
  }
  
  // 更新状态为测试中
  provider.connectionStatus = 'testing'
  
  try {
    const response = await fetch('http://localhost:18620/api/ai/test-connection', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        provider: provider.id,
        base_url: provider.baseUrl || provider.defaultBaseUrl,
        api_key: provider.apiKey,
        model: provider.model
      })
    })
    const data = await response.json()
    
    if (data.code === 0) {
      provider.connectionStatus = 'connected'
      provider.lastTestTime = new Date().toLocaleString('zh-CN')
      if (!silent) {
        alert(`${provider.name} 连接成功！`)
      }
    } else {
      provider.connectionStatus = 'failed'
      provider.lastTestTime = new Date().toLocaleString('zh-CN')
      if (!silent) {
        alert(`连接失败: ${data.message}`)
      }
    }
  } catch (error) {
    console.error('测试连接失败:', error)
    provider.connectionStatus = 'failed'
    provider.lastTestTime = new Date().toLocaleString('zh-CN')
    if (!silent) {
      alert(`连接失败: ${error.message}`)
    }
  } finally {
    if (testingBtn) {
      testingBtn.disabled = false
      testingBtn.textContent = '测试连接'
    }
  }
}

// 自动测试所有提供商的连接
async function autoTestAllConnections() {
  for (const provider of aiProviders.value) {
    // 只测试已配置的提供商（有API Key或是Ollama）
    if (provider.id === 'ollama' || provider.apiKey) {
      const originalProvider = selectedProvider.value
      selectedProvider.value = provider.id
      await testConnection(null, true)
      selectedProvider.value = originalProvider
      // 避免请求过快
      await new Promise(resolve => setTimeout(resolve, 500))
    }
  }
}

// 定时自动测试当前提供商连接
let connectionTestInterval = null
function startAutoTest() {
  stopAutoTest()
  // 每60秒自动测试一次当前提供商
  connectionTestInterval = setInterval(() => {
    if (currentProvider.value && (currentProvider.value.id === 'ollama' || currentProvider.value.apiKey)) {
      testConnection(null, true)
    }
  }, 60000)
}

function stopAutoTest() {
  if (connectionTestInterval) {
    clearInterval(connectionTestInterval)
    connectionTestInterval = null
  }
}

async function saveAIConfig() {
  if (!currentProvider.value) return
  
  try {
    const config = {
      provider: currentProvider.value.id,
      api_key: currentProvider.value.apiKey || '',
      base_url: currentProvider.value.baseUrl || currentProvider.value.defaultBaseUrl,
      model: currentProvider.value.model
    }
    
    const response = await fetch('http://localhost:18620/api/ai/config', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(config)
    })
    const data = await response.json()
    
    if (data.code === 0) {
      alert('AI 配置保存成功！')
    } else {
      alert('保存失败: ' + data.message)
    }
  } catch (error) {
    console.error('保存配置失败:', error)
    alert('保存失败: ' + error.message)
  }
}

// 历史记录相关函数
async function loadHistory() {
  if (privacyMode.value) return
  
  historyLoading.value = true
  try {
    const response = await fetch('http://localhost:18620/api/history')
    const data = await response.json()
    
    if (data.code === 0) {
      historyRecords.value = data.data || []
    } else {
      console.error('加载历史记录失败:', data.message)
    }
  } catch (error) {
    console.error('加载历史记录出错:', error)
  } finally {
    historyLoading.value = false
  }
}

async function clearHistory() {
  if (!confirm('确定要清除所有历史记录吗？此操作不可恢复。')) return
  
  historyLoading.value = true
  try {
    const response = await fetch('http://localhost:18620/api/history/clear', {
      method: 'POST'
    })
    const data = await response.json()
    
    if (data.code === 0) {
      historyRecords.value = []
      alert('历史记录已清除')
    } else {
      alert('清除失败: ' + data.message)
    }
  } catch (error) {
    console.error('清除历史记录出错:', error)
    alert('清除失败: ' + error.message)
  } finally {
    historyLoading.value = false
  }
}

function formatTime(timestamp) {
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now - date
  
  // 一分钟内
  if (diff < 60000) {
    return '刚刚'
  }
  
  // 一小时内
  if (diff < 3600000) {
    const minutes = Math.floor(diff / 60000)
    return `${minutes} 分钟前`
  }
  
  // 今天
  if (date.toDateString() === now.toDateString()) {
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
  
  // 昨天
  const yesterday = new Date(now)
  yesterday.setDate(yesterday.getDate() - 1)
  if (date.toDateString() === yesterday.toDateString()) {
    return '昨天 ' + date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
  
  // 其他日期
  return date.toLocaleString('zh-CN', { 
    month: '2-digit', 
    day: '2-digit', 
    hour: '2-digit', 
    minute: '2-digit' 
  })
}

// Ollama 模型加载函数
async function loadOllamaModels() {
  ollamaModelsLoading.value = true
  try {
    const response = await fetch('http://localhost:18620/api/ollama/models')
    const data = await response.json()
    
    if (data.code === 0) {
      ollamaModels.value = data.data || []
    } else {
      console.error('加载 Ollama 模型失败:', data.message)
      ollamaModels.value = []
    }
  } catch (error) {
    console.error('加载 Ollama 模型出错:', error)
    ollamaModels.value = []
  } finally {
    ollamaModelsLoading.value = false
  }
}

// 监听标签页切换
watch(currentTab, (newTab) => {
  if (newTab === 'history') {
    loadHistory()
  } else if (newTab === 'ai') {
    // 如果选中的是 Ollama，自动加载模型列表
    if (selectedProvider.value === 'ollama' && ollamaModels.value.length === 0) {
      loadOllamaModels()
    }
  }
})

// 监听提供商切换
watch(selectedProvider, (newProvider) => {
  if (newProvider === 'ollama' && ollamaModels.value.length === 0) {
    loadOllamaModels()
  }
  // 切换提供商时自动测试连接
  const provider = aiProviders.value.find(p => p.id === newProvider)
  if (provider && (provider.id === 'ollama' || provider.apiKey)) {
    testConnection(null, true)
  }
})

// 从后端加载用户模板
async function loadUserTemplates() {
  try {
    const response = await fetch('http://localhost:18620/api/templates')
    const data = await response.json()
    if (data.code === 0) {
      userTemplates.value = data.data || []
    } else {
      console.error('加载模板失败:', data.message)
    }
  } catch (error) {
    console.error('加载模板失败:', error)
  }
}

// 加载 AI 配置
async function loadAIConfig() {
  try {
    const response = await fetch('http://localhost:18620/api/ai/config')
    const data = await response.json()
    if (data.code === 0 && data.data) {
      const config = data.data
      // 更新对应的 provider 配置
      const provider = aiProviders.value.find(p => p.id === config.provider)
      if (provider) {
        selectedProvider.value = config.provider
        provider.apiKey = config.api_key || ''
        provider.baseUrl = config.base_url || provider.defaultBaseUrl
        provider.model = config.model || provider.defaultModel
      }
    }
  } catch (error) {
    console.error('加载 AI 配置失败:', error)
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadUserTemplates()
  loadRules()
  loadAIConfig()
  // 延迟启动自动测试，给配置加载时间
  setTimeout(() => {
    autoTestAllConnections()
    startAutoTest()
  }, 1000)
})

// 组件卸载时清理定时器
onBeforeUnmount(() => {
  stopAutoTest()
})
</script>

<style scoped>
/* ==================== 明暗主题基础样式 ==================== */
.settings-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #fff;
  transition: background-color 0.3s ease, color 0.3s ease;
}

/* 暗色主题 */
.settings-container.theme-dark {
  background: #1a1a1a;
  color: #e0e0e0;
}

.settings-container.theme-dark .sidebar {
  background: #252525;
  border-right-color: #333;
}

.settings-container.theme-dark .sidebar-header,
.settings-container.theme-dark .tabs {
  background: #1a1a1a;
  border-bottom-color: #333;
}

.settings-container.theme-dark .sidebar-header h2,
.settings-container.theme-dark .page-title {
  color: #e0e0e0;
}

.settings-container.theme-dark .tab {
  color: #999;
}

.settings-container.theme-dark .tab:hover {
  color: #007AFF;
}

.settings-container.theme-dark .tab.active {
  color: #007AFF;
}

.settings-container.theme-dark .rule-item:hover {
  background: #2a2a2a;
}

.settings-container.theme-dark .rule-item.active {
  background: #007AFF;
}

.settings-container.theme-dark .rule-name,
.settings-container.theme-dark .rule-path {
  color: inherit;
}

.settings-container.theme-dark .content-wrapper {
  background: #1a1a1a;
}

.settings-container.theme-dark .section {
  background: #252525;
  border-color: #333;
}

.settings-container.theme-dark .section h3 {
  color: #e0e0e0;
}

.settings-container.theme-dark .form-group label {
  color: #e0e0e0;
}

.settings-container.theme-dark .form-group input[type="text"],
.settings-container.theme-dark .form-select {
  background: #1a1a1a;
  border-color: #333;
  color: #e0e0e0;
}

.settings-container.theme-dark .button-group button {
  background: #1a1a1a;
  border-color: #333;
  color: #e0e0e0;
}

.settings-container.theme-dark .button-group button:hover {
  background: #2a2a2a;
}

.settings-container.theme-dark .button-group button.active {
  background: #007AFF;
  color: white;
  border-color: #007AFF;
}

.settings-container.theme-dark .file-type-btn {
  background: #1a1a1a;
  border-color: #333;
  color: #e0e0e0;
}

.settings-container.theme-dark .file-type-btn:hover {
  background: #2a2a2a;
}

.settings-container.theme-dark .file-type-btn.selected {
  background: #007AFF;
  border-color: #007AFF;
}

.settings-container.theme-dark .tag {
  background: #2a2a2a;
  color: #e0e0e0;
}

.settings-container.theme-dark .template-part {
  background: #2a2a2a;
  color: #e0e0e0;
}

.settings-container.theme-dark .info-text {
  background: #2a2a2a;
  color: #999;
}

.settings-container.theme-dark .ai-info p {
  color: #999;
}

.settings-container.theme-dark .ai-hint {
  background: linear-gradient(135deg, #1a2a35 0%, #1f2f3a 100%);
  border-color: #2a4a5a;
  color: #6ab0d0;
}

.settings-container.theme-dark .config-status {
  background: #2a2a2a;
}

.settings-container.theme-dark .config-form {
  background: transparent;
}

.settings-container.theme-dark .password-input input {
  background: #1a1a1a;
  border-color: #333;
  color: #e0e0e0;
}

.settings-container.theme-dark .ai-providers {
  background: #252525;
  border-right-color: #333;
}

.settings-container.theme-dark .providers-title {
  color: #999;
}

.settings-container.theme-dark .provider-item:hover {
  background: #2a2a2a;
}

.settings-container.theme-dark .provider-item.active {
  background: #333;
}

.settings-container.theme-dark .provider-name {
  color: #e0e0e0;
}

.settings-container.theme-dark .history-card {
  background: #252525;
  border-color: #333;
}

.settings-container.theme-dark .history-header,
.settings-container.theme-dark .file-name {
  color: #e0e0e0;
}

.settings-container.theme-dark .file-path,
.settings-container.theme-dark .meta-item {
  color: #999;
}

.settings-container.theme-dark .template-item {
  background: #252525;
  border-color: #333;
}

.settings-container.theme-dark .template-name {
  color: #e0e0e0;
}

.settings-container.theme-dark .template-preview {
  color: #999;
}

/* ==================== 原有样式 ==================== */

.header-bar {
  padding: 15px 20px;
  background: white;
  border-bottom: 1px solid #e0e0e0;
  flex-shrink: 0;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
  color: #333;
}

.main-body {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.sidebar {
  width: 280px;
  background: #f5f5f5;
  border-right: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #e0e0e0;
}

.sidebar-header h2 {
  font-size: 18px;
  font-weight: 600;
}

.rules-list {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
}

.rule-item {
  display: flex;
  align-items: center;
  padding: 12px;
  margin-bottom: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
}

.rule-item:hover {
  background: #e8e8e8;
}

.rule-item.active {
  background: #007AFF;
  color: white;
}

.rule-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
}

.rule-icon-img {
  width: 20px;
  height: 20px;
  filter: brightness(0) invert(1);
}

.rule-info {
  flex: 1;
  min-width: 0;
}

.rule-name {
  font-weight: 500;
  margin-bottom: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.rule-path {
  font-size: 12px;
  opacity: 0.7;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.add-rule-btn {
  margin: 10px;
  padding: 12px;
  background: #007AFF;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.add-rule-btn:hover {
  background: #0051D5;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.content-wrapper {
  flex: 1;
  overflow-y: auto;
  background: #f5f5f5;
}

.tabs {
  display: flex;
  border-bottom: 1px solid #e0e0e0;
  padding: 0 20px;
  background: white;
  flex-shrink: 0;
}

.tab {
  padding: 15px 20px;
  background: none;
  border: none;
  cursor: pointer;
  font-size: 14px;
  color: #666;
  border-bottom: 2px solid transparent;
  transition: all 0.2s;
}

.tab:hover {
  color: #007AFF;
}

.tab.active {
  color: #007AFF;
  border-bottom-color: #007AFF;
}

.tab-content {
  padding: 20px;
}

.section {
  margin-bottom: 30px;
  padding: 20px;
  background: white;
  border-radius: 12px;
  border: 1px solid #e0e0e0;
}

.section h3 {
  font-size: 16px;
  margin-bottom: 15px;
  font-weight: 600;
  display: flex;
  align-items: center;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
}

.form-group input[type="text"] {
  width: 100%;
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  font-size: 14px;
}

.color-picker {
  display: flex;
  gap: 10px;
}

.color-option {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.2s;
  border: 3px solid transparent;
}

.color-option:hover {
  transform: scale(1.1);
}

.color-option.selected {
  border-color: #333;
}

.path-input {
  display: flex;
  gap: 10px;
}

.path-input input {
  flex: 1;
}

.path-input button {
  padding: 10px 20px;
  background: #007AFF;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

.button-group {
  display: flex;
  gap: 10px;
  margin-bottom: 15px;
}

.button-group button {
  flex: 1;
  padding: 12px;
  background: #f5f5f5;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.button-group button:hover {
  background: #e8e8e8;
}

.button-group button.active {
  background: #007AFF;
  color: white;
  border-color: #007AFF;
}

.checkbox-group label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.file-types {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 10px;
  margin-bottom: 15px;
}

.file-type-btn {
  padding: 12px;
  background: #f5f5f5;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.2s;
}

.file-type-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

.file-type-btn:hover {
  background: #e8e8e8;
}

.file-type-btn.selected {
  background: #007AFF;
  color: white;
  border-color: #007AFF;
}

.file-type-btn.selected .file-type-icon {
  filter: brightness(0) invert(1);
}

.extension-input {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
}

.extension-input input {
  flex: 1;
}

.extension-input button {
  padding: 10px 20px;
  background: #007AFF;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

.extension-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag {
  padding: 6px 12px;
  background: #e0e0e0;
  border-radius: 16px;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.tag button {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 16px;
  padding: 0;
  line-height: 1;
}

.template-builder {
  margin-bottom: 15px;
}

.template-components {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 15px;
}

.template-components button {
  padding: 8px 16px;
  background: #f5f5f5;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
}

.template-components button:hover {
  background: #e8e8e8;
}

.template-preview {
  padding: 15px;
  background: #f9f9f9;
  border-radius: 8px;
}

.preview-label {
  font-size: 12px;
  color: #666;
  margin-bottom: 5px;
}

.preview-text {
  font-family: monospace;
  font-size: 14px;
}

.template-parts {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.template-part {
  padding: 6px 12px;
  background: #007AFF;
  color: white;
  border-radius: 6px;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.template-part button {
  background: none;
  border: none;
  color: white;
  cursor: pointer;
  font-size: 16px;
  padding: 0;
  line-height: 1;
}

.date-sources {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.date-sources label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.info-text {
  margin-top: 10px;
  padding: 10px;
  background: #fff3cd;
  border-radius: 6px;
  font-size: 13px;
  color: #856404;
}

.ai-toggle {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
}

.switch {
  position: relative;
  width: 50px;
  height: 28px;
}

.switch input {
  display: none;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  transition: 0.4s;
  border-radius: 28px;
}

.slider:before {
  position: absolute;
  content: "";
  height: 20px;
  width: 20px;
  left: 4px;
  bottom: 4px;
  background-color: white;
  transition: 0.4s;
  border-radius: 50%;
}

input:checked + .slider {
  background-color: #007AFF;
}

input:checked + .slider:before {
  transform: translateX(22px);
}

.ai-status {
  padding: 4px 12px;
  background: #FF9500;
  color: white;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.ai-info {
  font-size: 13px;
  color: #666;
}

.ai-info .usage {
  margin-top: 5px;
  color: #FF9500;
  display: flex;
  align-items: center;
  gap: 6px;
}

.ai-hint {
  margin-top: 12px;
  padding: 10px 14px;
  background: linear-gradient(135deg, #E8F5FF 0%, #F0F7FF 100%);
  border: 1px solid #B8E0FF;
  border-radius: 8px;
  font-size: 13px;
  color: #0066CC;
  display: flex;
  align-items: center;
  gap: 8px;
}

.ai-hint .hint-icon {
  font-size: 16px;
}

.actions {
  display: flex;
  gap: 10px;
}

.btn-primary {
  padding: 12px 24px;
  background: #007AFF;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
}

.btn-primary:hover {
  background: #0051D5;
}

.btn-toggle {
  padding: 12px 24px;
  background: #34C759;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
}

.btn-danger {
  padding: 12px 24px;
  background: #FF3B30;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  margin-left: auto;
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #999;
  font-size: 16px;
}

/* 常规页面样式 */
.setting-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 0;
  border-bottom: 1px solid #f0f0f0;
}

.setting-row.no-border {
  border-bottom: none;
}

.setting-info {
  flex: 1;
}

.setting-label {
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 4px;
}

.setting-description {
  font-size: 12px;
  color: #999;
}

.setting-value {
  font-size: 14px;
  color: #666;
}

.setting-value.link {
  color: #007AFF;
  text-decoration: none;
}

.setting-value.link:hover {
  text-decoration: underline;
}

.text-muted {
  color: #999;
}

.slider-container {
  display: flex;
  align-items: center;
  gap: 15px;
  min-width: 200px;
}

.opacity-slider {
  flex: 1;
  height: 4px;
  border-radius: 2px;
  background: #e0e0e0;
  outline: none;
  -webkit-appearance: none;
}

.opacity-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #007AFF;
  cursor: pointer;
}

.opacity-slider::-moz-range-thumb {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #007AFF;
  cursor: pointer;
  border: none;
}

.slider-value {
  font-size: 14px;
  color: #666;
  min-width: 45px;
  text-align: right;
}

.language-select {
  padding: 8px 12px;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  font-size: 14px;
  background: white;
  cursor: pointer;
  min-width: 150px;
}

.shortcut-display {
  display: flex;
  gap: 6px;
}

.shortcut-display kbd {
  padding: 6px 12px;
  background: #f5f5f5;
  border: 1px solid #d0d0d0;
  border-radius: 6px;
  font-size: 14px;
  font-family: -apple-system, BlinkMacSystemFont, sans-serif;
  box-shadow: 0 2px 0 #d0d0d0;
}

.info-banner {
  padding: 12px;
  background: #f9f9f9;
  border-radius: 8px;
  font-size: 13px;
  color: #666;
  margin-top: 10px;
}

.section-icon {
  width: 16px;
  height: 16px;
  margin-right: 8px;
  vertical-align: middle;
}

/* 模板页面样式 */
.template-page {
  padding-bottom: 80px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.import-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: #007AFF;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
}

.import-btn:hover {
  background: #0051D5;
}

.import-btn .button-icon {
  filter: brightness(0) invert(1);
}

.button-icon {
  width: 14px;
  height: 14px;
  display: inline-block;
}

.template-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.template-item {
  display: flex;
  align-items: flex-start;
  gap: 15px;
  padding: 15px;
  background: #f9f9f9;
  border-radius: 10px;
  transition: background 0.2s;
}

.template-item:hover {
  background: #f0f0f0;
}

.template-icon {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.template-icon img {
  width: 28px;
  height: 28px;
}

.template-info {
  flex: 1;
  min-width: 0;
}

.template-name {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.system-badge {
  padding: 2px 8px;
  background: #007AFF;
  color: white;
  border-radius: 10px;
  font-size: 11px;
  font-weight: 500;
}

.template-components {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 8px;
}

.component-tag {
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap;
}

.component-tag.date {
  background: #FFE5E5;
  color: #FF3B30;
}

.component-tag.time {
  background: #FFE5E5;
  color: #FF9500;
}

.component-tag.number {
  background: #E5F0FF;
  color: #007AFF;
}

.component-tag.separator {
  background: #F0F0F0;
  color: #666;
}

.component-tag.original {
  background: #E5F9E5;
  color: #34C759;
}

.component-tag.text {
  background: #F0F0F0;
  color: #333;
}

.template-preview {
  font-size: 13px;
  color: #999;
  font-family: monospace;
}

.delete-template-btn {
  padding: 6px 16px;
  background: transparent;
  color: #FF3B30;
  border: 1px solid #FF3B30;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  transition: all 0.2s;
}

.delete-template-btn:hover {
  background: #FF3B30;
  color: white;
}

.empty-template {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
}

.empty-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 15px;
  opacity: 0.5;
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-icon img {
  width: 56px;
  height: 56px;
}

.empty-text {
  font-size: 16px;
  color: #999;
}

.template-footer {
  position: fixed;
  bottom: 0;
  left: 280px;
  right: 0;
  padding: 15px 20px;
  background: white;
  border-top: 1px solid #e0e0e0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  z-index: 10;
}

.template-count {
  font-size: 13px;
  color: #666;
}

.template-hint {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #999;
}

.inline-icon {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
  display: inline-block;
}

/* AI 配置页面样式 */
.ai-config-page {
  padding: 0 !important;
}

.ai-layout {
  display: flex;
  height: calc(100vh - 100px);
}

.ai-providers {
  width: 280px;
  background: #f5f5f5;
  border-right: 1px solid #e0e0e0;
  overflow-y: auto;
  padding: 20px;
}

.providers-title {
  font-size: 16px;
  font-weight: 600;
  color: #666;
  margin-bottom: 15px;
}

.provider-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.provider-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 15px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
}

.provider-item:hover {
  background: #e8e8e8;
}

.provider-item.active {
  background: #e0e0e0;
}

.provider-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.provider-name {
  font-size: 14px;
  font-weight: 500;
}

.ai-config {
  flex: 1;
  overflow-y: auto;
  padding: 30px;
}

.config-content {
  max-width: 800px;
}

.config-title {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 20px;
}

.config-status {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: #f9f9f9;
  border-radius: 8px;
  margin-bottom: 30px;
  border: 2px solid transparent;
  transition: all 0.3s;
}

.config-status.connected {
  background: #E8F5E9;
  border-color: #4CAF50;
}

.config-status.failed {
  background: #FFEBEE;
  border-color: #F44336;
}

.config-status.testing {
  background: #FFF3E0;
  border-color: #FF9800;
}

.status-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

.status-text {
  font-size: 14px;
  color: #666;
  flex: 1;
}

.config-status.connected .status-text {
  color: #2E7D32;
  font-weight: 500;
}

.config-status.failed .status-text {
  color: #C62828;
  font-weight: 500;
}

.config-status.testing .status-text {
  color: #E65100;
  font-weight: 500;
}

.test-time {
  font-size: 12px;
  color: #999;
  margin-left: 8px;
}

.btn-auto-refresh {
  padding: 6px;
  background: transparent;
  border: 1px solid #ddd;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-auto-refresh:hover:not(:disabled) {
  background: #f0f0f0;
  border-color: #999;
}

.btn-auto-refresh:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.btn-auto-refresh img {
  width: 14px;
  height: 14px;
}

.icon-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.config-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-size: 14px;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 6px;
}

.btn-refresh {
  padding: 4px 10px;
  background: #f5f5f5;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  margin-left: auto;
  transition: all 0.2s;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.btn-refresh:hover:not(:disabled) {
  background: #e8e8e8;
}

.btn-refresh:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.form-select {
  padding: 10px 12px;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  font-size: 14px;
  background: white;
  cursor: pointer;
  transition: border-color 0.2s;
}

.form-select:focus {
  outline: none;
  border-color: #007AFF;
}

.form-select:disabled {
  background: #f5f5f5;
  cursor: not-allowed;
  opacity: 0.6;
}

.form-hint {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.form-hint.error {
  color: #FF3B30;
}

.form-hint code {
  padding: 2px 6px;
  background: #f5f5f5;
  border-radius: 3px;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 11px;
}

.help-icon {
  width: 14px;
  height: 14px;
  opacity: 0.6;
  cursor: help;
}

.form-input {
  padding: 10px 12px;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.2s;
}

.form-input:focus {
  outline: none;
  border-color: #007AFF;
}

.password-input {
  display: flex;
  gap: 8px;
}

.password-input .form-input {
  flex: 1;
}

.toggle-visibility {
  padding: 10px 15px;
  background: #f5f5f5;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;
}

.toggle-visibility img {
  width: 18px;
  height: 18px;
}

.toggle-visibility:hover {
  background: #e8e8e8;
}

.form-select {
  padding: 10px 12px;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  font-size: 14px;
  background: white;
  cursor: pointer;
}

.form-actions {
  display: flex;
  gap: 12px;
  margin-top: 10px;
}

.btn-test {
  padding: 10px 20px;
  background: white;
  color: #007AFF;
  border: 1px solid #007AFF;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.btn-test:hover {
  background: #007AFF;
  color: white;
}

.btn-save {
  padding: 10px 20px;
  background: #007AFF;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: background 0.2s;
}

.btn-save:hover {
  background: #0051D5;
}

.ollama-info {
  margin-top: 30px;
  padding: 20px;
  background: #f9f9f9;
  border-radius: 10px;
  border: 1px solid #e0e0e0;
}

.ollama-info h4 {
  font-size: 16px;
  margin-bottom: 15px;
}

.ollama-info ol {
  padding-left: 20px;
  line-height: 1.8;
}

.ollama-info code {
  padding: 2px 6px;
  background: #e0e0e0;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 13px;
}

.empty-config {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #999;
  font-size: 16px;
}

/* 历史记录页面样式 */
.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 20px;
}

.history-actions {
  display: flex;
  align-items: center;
  gap: 15px;
  flex-wrap: wrap;
}

.history-actions .btn-danger {
  margin-left: 0;
}

.switch-inline {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.switch-inline .switch {
  width: 40px;
  height: 22px;
}

.switch-inline .slider:before {
  height: 16px;
  width: 16px;
  left: 3px;
  bottom: 3px;
}

.switch-inline input:checked + .slider:before {
  transform: translateX(18px);
}

.label-text {
  font-size: 13px;
  color: #666;
  user-select: none;
}

.btn-danger {
  padding: 8px 16px;
  background: #FF3B30;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.3s;
}

.btn-danger.btn-small {
  padding: 6px 12px;
  font-size: 12px;
}

.btn-danger:hover:not(:disabled) {
  background: #D70015;
}

.btn-danger:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.history-page {
  background: #f5f5f5;
}

.history-header-section {
  margin-bottom: 20px;
}

.history-body {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.history-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 15px;
}

.stat-item {
  background: #f8f9fa;
  padding: 15px;
  border-radius: 8px;
  text-align: center;
}

.stat-label {
  display: block;
  font-size: 12px;
  color: #666;
  margin-bottom: 8px;
}

.stat-value {
  display: block;
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

.stat-value.success {
  color: #34C759;
}

.stat-value.error {
  color: #FF3B30;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  background: white;
  border-radius: 12px;
  border: 1px solid #e0e0e0;
  color: #666;
  font-size: 14px;
}

.empty-history {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  background: white;
  border-radius: 12px;
  border: 1px solid #e0e0e0;
}

.empty-history svg {
  margin-bottom: 15px;
}

.empty-history p {
  font-size: 16px;
  margin-bottom: 5px;
  color: #333;
}

.empty-history small {
  font-size: 13px;
  color: #666;
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.history-item {
  display: flex;
  gap: 15px;
  padding: 15px;
  background: white;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  transition: all 0.3s;
}

.history-item:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border-color: #d0d0d0;
}

.history-status {
  flex-shrink: 0;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

.history-status.success {
  background: #E8F5E9;
  color: #34C759;
}

.history-status.error {
  background: #FFEBEE;
  color: #FF3B30;
}

.history-content {
  flex: 1;
  min-width: 0;
}

.history-main {
  margin-bottom: 10px;
}

.file-change {
  display: flex;
  align-items: center;
  gap: 12px;
}

.file-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
  min-width: 0;
}

.file-name {
  font-size: 14px;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-name.original {
  color: #666;
}

.file-name.new {
  color: #007AFF;
}

.file-path {
  font-size: 12px;
  color: #999;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.arrow-icon {
  flex-shrink: 0;
  color: #999;
}

.history-meta {
  display: flex;
  align-items: center;
  gap: 15px;
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #999;
}

.meta-item svg {
  flex-shrink: 0;
  opacity: 0.7;
}

.rule-tag {
  padding: 2px 8px;
  background: #F0F0F0;
  border-radius: 4px;
  color: #666;
}

.action-tag {
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

.action-tag.copy {
  background: #E3F2FD;
  color: #007AFF;
}

.action-tag.move {
  background: #FFF3E0;
  color: #FF9500;
}

.icon-picker-trigger {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background: #f5f5f5;
  border: 1px solid #e0e0e0;
  border-radius: 10px;
  cursor: pointer;
  font-size: 14px;
  color: #333;
}

.icon-picker-trigger:hover {
  background: #ededed;
}

.icon-preview {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: #fff;
  border: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

.icon-picker-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.35);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.icon-picker-modal {
  width: 720px;
  max-width: 92vw;
  background: #f3f3f3;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.2);
}

.icon-picker-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  background: #f0f0f0;
  border-bottom: 1px solid #e0e0e0;
}

.icon-picker-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.icon-picker-done {
  background: #007AFF;
  color: white;
  border: none;
  padding: 6px 14px;
  border-radius: 10px;
  cursor: pointer;
  font-size: 14px;
}

.icon-picker-done:hover {
  background: #0051D5;
}

.icon-picker-tabs {
  display: flex;
  gap: 10px;
  padding: 12px 20px 4px;
  flex-wrap: wrap;
}

.icon-picker-tab {
  padding: 8px 14px;
  background: #e2e2e2;
  border: none;
  border-radius: 999px;
  cursor: pointer;
  font-size: 14px;
  color: #444;
}

.icon-picker-tab.active {
  background: #007AFF;
  color: #fff;
}

.icon-picker-grid {
  padding: 16px 20px 22px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(64px, 1fr));
  gap: 12px;
}

.icon-picker-item {
  height: 64px;
  border-radius: 14px;
  background: #e4e4e4;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: transform 0.12s ease, background 0.12s ease;
}

.icon-picker-item:hover {
  transform: translateY(-2px);
  background: #dcdcdc;
}

.icon-picker-item.selected {
  background: #007AFF;
  color: #fff;
}

.icon-picker-item img {
  width: 28px;
  height: 28px;
}

.icon-picker-item.selected img {
  filter: brightness(0) invert(1);
}

.icon-preview img {
  width: 20px;
  height: 20px;
  display: block;
}

.icon-spin {
  animation: icon-rotate 1s linear infinite;
}

@keyframes icon-rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

</style>
