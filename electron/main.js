import { app, BrowserWindow, Menu, Tray, ipcMain, screen, dialog, nativeTheme } from 'electron'
import { fileURLToPath } from 'url'
import path from 'path'
import { spawn } from 'child_process'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

// 强制使用亮色主题
nativeTheme.themeSource = 'light'

let mainWindow = null
let settingsWindow = null
let goProcess = null
let isOnlyOnDesktop = false

// 启动 Go 后端服务
function startGoBackend() {
  const goPath = path.join(__dirname, '../backend/main')
  
  // 开发环境：使用 go run
  if (process.env.NODE_ENV === 'development') {
    goProcess = spawn('go', ['run', '.'], {
      cwd: path.join(__dirname, '../backend'),
      stdio: 'inherit'
    })
  } else {
    // 生产环境：运行编译后的二进制文件
    if (process.platform === 'darwin') {
      goProcess = spawn(goPath, [], { stdio: 'inherit' })
    }
  }

  goProcess.on('error', (err) => {
    console.error('Go backend error:', err)
  })

  goProcess.on('close', (code) => {
    console.log(`Go backend exited with code ${code}`)
  })
}

// 创建悬浮球窗口
function createWindow() {
  const { width, height } = screen.getPrimaryDisplay().workAreaSize
  
  mainWindow = new BrowserWindow({
    width: 120,  // 稍微大一点便于拖放
    height: 120,
    x: width - 100,
    y: height - 100,
    frame: false,
    transparent: true,
    alwaysOnTop: true,
    resizable: false,
    skipTaskbar: true,  // 不在任务栏/Dock显示
    hasShadow: false,
    focusable: true,
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      preload: path.join(__dirname, 'preload.js')
    }
  })
  
  // 隐藏 Dock 图标
  if (process.platform === 'darwin') {
    app.dock.hide()
  }

  console.log('Window created, position:', mainWindow.getPosition())

  // 设置窗口级别
  if (isOnlyOnDesktop) {
    mainWindow.setVisibleOnAllWorkspaces(false)
    mainWindow.setAlwaysOnTop(false, 'normal')
  } else {
    mainWindow.setVisibleOnAllWorkspaces(true, { visibleOnFullScreen: true })
    mainWindow.setAlwaysOnTop(true, 'floating', 1)
  }

  // 开发环境加载 Vite 开发服务器，生产环境加载打包后的文件
  if (process.env.NODE_ENV === 'development') {
    mainWindow.loadURL('http://localhost:5173')
  } else {
    mainWindow.loadFile(path.join(__dirname, '../dist/index.html'))
  }

  // 允许拖拽窗口
  mainWindow.setMovable(true)
  // 注意：默认不设置 ignoreMouseEvents，以支持拖放文件
  // 如果需要穿透，由渲染进程通过 IPC 控制
  // mainWindow.setIgnoreMouseEvents(true, { forward: true })
  
  // 防止页面导航
  mainWindow.webContents.on('will-navigate', (e) => {
    e.preventDefault()
  })

  mainWindow.webContents.on('context-menu', () => {
    const menu = createContextMenu()
    menu.popup({ window: mainWindow })
  })
}

// 创建右键菜单
function createContextMenu() {
  const template = [
    {
      label: '设置',
      click: () => {
        createSettingsWindow()
      }
    },
    {
      label: isOnlyOnDesktop ? '悬浮在所有窗口上' : '只在桌面显示',
      click: () => {
        isOnlyOnDesktop = !isOnlyOnDesktop
        
        if (isOnlyOnDesktop) {
          // 只在桌面显示：窗口级别设为普通，会被其他窗口遮挡
          mainWindow.setAlwaysOnTop(false)
          mainWindow.setVisibleOnAllWorkspaces(false)
          // 将窗口设置为桌面级别（在所有窗口后面）
          if (process.platform === 'darwin') {
            mainWindow.setWindowButtonVisibility(false)
          }
        } else {
          // 悬浮在所有窗口上
          mainWindow.setAlwaysOnTop(true, 'floating', 1)
          mainWindow.setVisibleOnAllWorkspaces(true, { visibleOnFullScreen: true })
        }
        
        mainWindow.webContents.send('desktop-mode-changed', isOnlyOnDesktop)
      }
    },
    { type: 'separator' },
    {
      label: '退出应用',
      click: () => {
        app.quit()
      }
    }
  ]

  return Menu.buildFromTemplate(template)
}

// 显示右键菜单
ipcMain.on('show-context-menu', (event) => {
  const menu = createContextMenu()
  menu.popup(BrowserWindow.fromWebContents(event.sender))
})

// 移动窗口
ipcMain.on('move-window', (event, deltaX, deltaY) => {
  const win = BrowserWindow.fromWebContents(event.sender)
  if (win) {
    const [x, y] = win.getPosition()
    win.setPosition(x + deltaX, y + deltaY)
  }
})

ipcMain.on('set-ignore-mouse-events', (event, ignore) => {
  const win = BrowserWindow.fromWebContents(event.sender)
  if (!win) return

  if (ignore) {
    win.setIgnoreMouseEvents(true, { forward: true })
  } else {
    win.setIgnoreMouseEvents(false)
  }
})

// 应用常规设置
ipcMain.on('apply-settings', (event, settings) => {
  console.log('Applying settings:', settings)
  
  // 应用桌面显示模式
  if (settings.desktopOnly !== undefined && settings.desktopOnly !== isOnlyOnDesktop) {
    isOnlyOnDesktop = settings.desktopOnly
    if (mainWindow) {
      if (isOnlyOnDesktop) {
        mainWindow.setVisibleOnAllWorkspaces(false)
        mainWindow.setAlwaysOnTop(false, 'normal')
      } else {
        mainWindow.setVisibleOnAllWorkspaces(true, { visibleOnFullScreen: true })
        mainWindow.setAlwaysOnTop(true, 'floating', 1)
      }
    }
  }
  
  // 通知主窗口设置已更新
  if (mainWindow && !mainWindow.isDestroyed()) {
    mainWindow.webContents.send('settings-updated', settings)
  }
})

// 处理拖放的文件
ipcMain.handle('process-dropped-file', async (event, filePath) => {
  console.log('Processing dropped file:', filePath)
  // 这里可以添加额外的文件处理逻辑
  return { success: true, path: filePath }
})

// 选择文件夹
ipcMain.handle('select-folder', async () => {
  const result = await dialog.showOpenDialog({
    properties: ['openDirectory', 'createDirectory'],
    title: '选择目标文件夹'
  })
  
  if (result.canceled) {
    return null
  }
  
  return result.filePaths[0]
})

// 选择文件
ipcMain.handle('select-file', async (event, options = {}) => {
  const result = await dialog.showOpenDialog({
    properties: ['openFile'],
    title: options.title || '选择文件',
    filters: options.filters || []
  })
  
  if (result.canceled) {
    return null
  }
  
  return result.filePaths[0]
})

// 创建设置窗口
function createSettingsWindow() {
  if (settingsWindow) {
    settingsWindow.focus()
    return
  }
  
  // 打开设置窗口时临时显示 Dock 图标
  if (process.platform === 'darwin') {
    app.dock.show()
  }

  settingsWindow = new BrowserWindow({
    width: 1000,
    height: 700,
    minWidth: 800,
    minHeight: 600,
    title: 'BlackHole 设置',
    backgroundColor: '#ffffff',
    vibrancy: null,
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      preload: path.join(__dirname, 'preload.js')
    }
  })

  // 开发环境加载设置页面
  if (process.env.NODE_ENV === 'development') {
    settingsWindow.loadURL('http://localhost:5173/settings.html')
  } else {
    settingsWindow.loadFile(path.join(__dirname, '../dist/settings.html'))
  }

  settingsWindow.on('closed', () => {
    settingsWindow = null
    // 关闭设置窗口后隐藏 Dock 图标
    if (process.platform === 'darwin') {
      app.dock.hide()
    }
  })
}

app.whenReady().then(() => {
  // 启动 Go 后端
  startGoBackend()
  
  // 延迟启动窗口，等待后端启动
  setTimeout(() => {
    createWindow()
  }, 1000)

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow()
    }
  })
})

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit()
  }
})

app.on('before-quit', () => {
  // 关闭 Go 后端进程
  if (goProcess) {
    goProcess.kill()
  }
})
