const { contextBridge, ipcRenderer, webUtils } = require('electron')

// 暴露安全的 API 给渲染进程
contextBridge.exposeInMainWorld('electronAPI', {
  showContextMenu: () => ipcRenderer.send('show-context-menu'),
  moveWindow: (deltaX, deltaY) => ipcRenderer.send('move-window', deltaX, deltaY),
  onOpenSettings: (callback) => ipcRenderer.on('open-settings', callback),
  onDesktopModeChanged: (callback) => ipcRenderer.on('desktop-mode-changed', (event, value) => callback(value)),
  setIgnoreMouseEvents: (ignore) => ipcRenderer.send('set-ignore-mouse-events', ignore),
  // 打开文件夹选择对话框
  selectFolder: () => ipcRenderer.invoke('select-folder'),
  // 打开文件选择对话框
  selectFile: (options) => ipcRenderer.invoke('select-file', options),
  // 获取拖放文件的路径
  getPathForFile: (file) => {
    // Electron 22+ 使用 webUtils.getPathForFile
    if (webUtils && webUtils.getPathForFile) {
      try {
        return webUtils.getPathForFile(file)
      } catch (e) {
        console.error('getPathForFile error:', e)
        return file.path || null
      }
    }
    // 旧版本 Electron 直接使用 file.path
    return file.path || null
  },
  // 处理拖放的文件
  processDroppedFile: (filePath) => ipcRenderer.invoke('process-dropped-file', filePath)
})
