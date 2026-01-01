<template>
  <div
    class="floating-ball"
    @contextmenu.prevent="showContextMenu"
    @dragover.prevent="onDragOver"
    @dragleave.prevent="onDragLeave"
    @drop.prevent="onDrop"
  >
    <div 
      class="ball-inner"
      ref="ballRef"
      :class="{ 
        'is-hovered': isHovered,
        'is-absorbing': isAbsorbing,
        'is-drag-over': isDragOver,
        'is-success': isSuccess
      }"
      @mousedown="startDrag"
    >
      <!-- 常驻漩涡效果 -->
      <div class="inner-vortex">
        <div class="spiral spiral-1"></div>
        <div class="spiral spiral-2"></div>
        <div class="spiral spiral-3"></div>
      </div>
      <div v-if="!isSuccess" class="ball-core"></div>
      <div v-if="isSuccess" class="success-icon">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="20 6 9 17 4 12"></polyline>
        </svg>
      </div>
      <div v-if="isDragOver" class="vortex">
        <div class="vortex-ring ring-1"></div>
        <div class="vortex-ring ring-2"></div>
        <div class="vortex-ring ring-3"></div>
      </div>
    </div>
    
    <!-- 处理状态提示 -->
    <div v-if="processingFile" class="processing-overlay">
      <div class="processing-progress">{{ currentFileIndex }}/{{ totalFiles }}</div>
      <div class="processing-text">{{ processingMessage }}</div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const isDragging = ref(false)
const isHovered = ref(false)
const isDragOver = ref(false)
const isAbsorbing = ref(false)
const isSuccess = ref(false)
const processingFile = ref(false)
const processingMessage = ref('')
const currentFileIndex = ref(0)
const totalFiles = ref(0)
const ballRef = ref(null)
let lastScreenX = 0
let lastScreenY = 0

// 播放成功提示音
const playSuccessSound = () => {
  try {
    const audioContext = new (window.AudioContext || window.webkitAudioContext)()
    const oscillator = audioContext.createOscillator()
    const gainNode = audioContext.createGain()
    
    oscillator.connect(gainNode)
    gainNode.connect(audioContext.destination)
    
    // 播放两个音符形成提示音
    oscillator.frequency.setValueAtTime(880, audioContext.currentTime) // A5
    oscillator.frequency.setValueAtTime(1108.73, audioContext.currentTime + 0.1) // C#6
    
    gainNode.gain.setValueAtTime(0.3, audioContext.currentTime)
    gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + 0.3)
    
    oscillator.start(audioContext.currentTime)
    oscillator.stop(audioContext.currentTime + 0.3)
  } catch (e) {
    console.log('无法播放提示音:', e)
  }
}

// 显示成功状态
const showSuccess = () => {
  isSuccess.value = true
  playSuccessSound()
  setTimeout(() => {
    isSuccess.value = false
  }, 2000)
}

const onMouseMove = (e) => {
  if (!isDragging.value) return
  const deltaX = e.screenX - lastScreenX
  const deltaY = e.screenY - lastScreenY
  lastScreenX = e.screenX
  lastScreenY = e.screenY
  if (deltaX !== 0 || deltaY !== 0) {
    window.electronAPI?.moveWindow?.(deltaX, deltaY)
  }
}

const onMouseUp = () => {
  isDragging.value = false
  window.removeEventListener('mousemove', onMouseMove)
  window.removeEventListener('mouseup', onMouseUp)
}

const startDrag = (e) => {
  if (e.button !== 0) return
  e.preventDefault()
  e.stopPropagation()
  
  isDragging.value = true
  lastScreenX = e.screenX
  lastScreenY = e.screenY

  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseup', onMouseUp)
}

onUnmounted(() => {
  window.removeEventListener('mousemove', onMouseMove)
  window.removeEventListener('mouseup', onMouseUp)
})

const showContextMenu = () => {
  if (window.electronAPI?.showContextMenu) {
    window.electronAPI.showContextMenu()
  }
}

// 文件拖放处理
const onDragOver = (e) => {
  e.preventDefault()
  e.stopPropagation()
  // 确保窗口可以接收拖放
  if (e.dataTransfer) {
    e.dataTransfer.dropEffect = 'copy'
  }
  isDragOver.value = true
}

const onDragLeave = (e) => {
  e.preventDefault()
  e.stopPropagation()
  isDragOver.value = false
}

const onDrop = async (e) => {
  e.preventDefault()
  e.stopPropagation()
  isDragOver.value = false
  
  console.log('Drop event triggered')
  console.log('dataTransfer:', e.dataTransfer)
  console.log('files count:', e.dataTransfer?.files?.length)
  
  const droppedFiles = Array.from(e.dataTransfer?.files || [])
  
  if (droppedFiles.length === 0) {
    console.warn('No files dropped')
    return
  }
  
  // 获取文件路径 - 使用 Electron API 或直接读取 path 属性
  const filesWithPaths = []
  for (const file of droppedFiles) {
    let filePath = null
    
    // 尝试使用 Electron 的 getPathForFile API
    if (window.electronAPI?.getPathForFile) {
      filePath = window.electronAPI.getPathForFile(file)
    }
    
    // 回退到 file.path（旧版 Electron）
    if (!filePath && file.path) {
      filePath = file.path
    }
    
    console.log(`File: ${file.name}, Path: ${filePath}`)
    
    if (filePath) {
      filesWithPaths.push({
        name: file.name,
        path: filePath,
        type: file.type,
        size: file.size
      })
    }
  }
  
  // 打印文件信息用于调试
  console.log('Files with paths:', filesWithPaths)
  
  if (filesWithPaths.length === 0) {
    console.error('No valid file paths found')
    processingFile.value = true
    processingMessage.value = '无法获取文件路径'
    setTimeout(() => {
      processingFile.value = false
    }, 2000)
    return
  }
  
  // 显示吸收动画
  isAbsorbing.value = true
  setTimeout(async () => {
    isAbsorbing.value = false
    await processFiles(filesWithPaths)
  }, 800)
}

// 处理文件
const processFiles = async (files) => {
  processingFile.value = true
  totalFiles.value = files.length
  currentFileIndex.value = 0
  processingMessage.value = '准备处理...'
  
  try {
    for (let i = 0; i < files.length; i++) {
      const file = files[i]
      currentFileIndex.value = i + 1
      processingMessage.value = file.name.length > 12 
        ? file.name.substring(0, 10) + '...' 
        : file.name
      
      // 调用后端 API 处理文件
      const response = await fetch('http://localhost:18620/api/files/process', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          file_path: file.path
        })
      })
      
      const result = await response.json()
      
      if (result.code === 0) {
        console.log('文件处理成功:', result.data)
        // 只在最后一个文件处理完成后显示成功状态
        if (i === files.length - 1) {
          showSuccess()
        }
      } else {
        console.error('文件处理失败:', result.message)
      }
      
      // 等待一小段时间再处理下一个文件
      if (i < files.length - 1) {
        await new Promise(resolve => setTimeout(resolve, 300))
      }
    }
    
    processingMessage.value = '完成！'
    setTimeout(() => {
      processingFile.value = false
    }, 1000)
    
  } catch (error) {
    console.error('文件处理错误:', error)
    processingMessage.value = '处理失败'
    setTimeout(() => {
      processingFile.value = false
    }, 2000)
  }
}

onMounted(() => {
  // 监听全局拖放事件用于调试
  const handleGlobalDragEnter = (e) => {
    if (e.dataTransfer?.types?.includes('Files')) {
      console.log('Global drag enter with files detected')
    }
  }
  
  const handleGlobalDrop = (e) => {
    console.log('Global drop event')
  }
  
  document.addEventListener('dragenter', handleGlobalDragEnter)
  document.addEventListener('drop', handleGlobalDrop)

  window.electronAPI?.onOpenSettings(() => {
    alert('打开设置页面（待实现）')
  })

  window.electronAPI?.onDesktopModeChanged(() => {})
})
</script>

<style scoped>
.floating-ball {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: default;
  user-select: none;
  touch-action: none;
  overflow: visible;
  background: transparent !important;
  pointer-events: auto;
}

.ball-inner {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background: radial-gradient(circle at 35% 35%, #1f1f1f 0%, #101010 55%, #060606 100%);
  border: 2px solid #3fd8ff;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: opacity 0.25s ease, border-color 0.3s ease;
  cursor: default;
  position: relative;
  -webkit-app-region: no-drag;
  opacity: 0.42;
  pointer-events: auto;
  overflow: hidden;
  animation: orb-breathe 2.5s ease-in-out infinite;
}

.ball-inner:hover,
.ball-inner.is-hovered {
  opacity: 1;
}

.ball-inner.is-drag-over {
  opacity: 1;
  border-color: #3fd8ff;
  box-shadow:
    0 0 10px rgba(63, 216, 255, 1),
    0 0 20px rgba(63, 216, 255, 0.6);
}

.ball-inner.is-absorbing {
  animation: absorb-pulse 0.8s ease-out;
}

.ball-inner.is-absorbing .spiral {
  animation-duration: 0.3s !important;
}

.ball-inner.is-absorbing .spiral-2 {
  animation-duration: 0.2s !important;
}

.ball-inner.is-absorbing .spiral-3 {
  animation-duration: 0.15s !important;
}

/* 成功状态样式 */
.ball-inner.is-success {
  opacity: 1;
  border-color: #22c55e;
  animation: success-glow 2s ease-in-out;
}

.success-icon {
  width: 24px;
  height: 24px;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
}

.success-icon svg {
  width: 100%;
  height: 100%;
}

@keyframes success-glow {
  0% {
    box-shadow:
      0 0 8px rgba(34, 197, 94, 0.8),
      0 0 15px rgba(34, 197, 94, 0.5),
      0 0 22px rgba(34, 197, 94, 0.3);
  }
  20% {
    box-shadow:
      0 0 12px rgba(34, 197, 94, 1),
      0 0 20px rgba(34, 197, 94, 0.7),
      0 0 28px rgba(34, 197, 94, 0.5);
  }
  100% {
    box-shadow:
      0 0 8px rgba(34, 197, 94, 0.6),
      0 0 15px rgba(34, 197, 94, 0.4),
      0 0 22px rgba(34, 197, 94, 0.2);
  }
}

.ball-core {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: radial-gradient(circle at 35% 35%, #2a2a2a 0%, #151515 70%, #0c0c0c 100%);
  box-shadow: inset 0 0 4px rgba(0, 0, 0, 0.55);
  pointer-events: none;
}

/* 内部漩涡效果 */
.inner-vortex {
  position: absolute;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  overflow: hidden;
}

.spiral {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 80%;
  height: 80%;
  border-radius: 50%;
  border: 1px solid transparent;
  border-top-color: rgba(63, 216, 255, 0.5);
  border-right-color: rgba(63, 216, 255, 0.3);
  transform: translate(-50%, -50%);
  animation: spiral-rotate 3s linear infinite;
  animation-play-state: paused;
}

.spiral-1 {
  width: 85%;
  height: 85%;
  border-top-color: rgba(63, 216, 255, 0.6);
  animation-duration: 3s;
}

.spiral-2 {
  width: 65%;
  height: 65%;
  border-top-color: rgba(63, 216, 255, 0.4);
  border-right-color: rgba(63, 216, 255, 0.2);
  animation-duration: 2s;
  animation-direction: reverse;
}

.spiral-3 {
  width: 45%;
  height: 45%;
  border-top-color: rgba(63, 216, 255, 0.3);
  border-right-color: rgba(63, 216, 255, 0.15);
  animation-duration: 1.5s;
}

@keyframes spiral-rotate {
  0% {
    transform: translate(-50%, -50%) rotate(0deg);
  }
  100% {
    transform: translate(-50%, -50%) rotate(360deg);
  }
}

.ball-inner:hover .spiral,
.ball-inner.is-hovered .spiral,
.ball-inner.is-drag-over .spiral {
  border-top-color: rgba(63, 216, 255, 0.8);
  animation-play-state: running;
}

.ball-inner:hover .spiral-2,
.ball-inner.is-hovered .spiral-2 {
  border-top-color: rgba(63, 216, 255, 0.6);
}

.ball-inner:hover .spiral-3,
.ball-inner.is-hovered .spiral-3 {
  border-top-color: rgba(63, 216, 255, 0.5);
}

/* 黑洞漩涡动画 */
.vortex {
  position: absolute;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.vortex-ring {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  border: 2px solid rgba(63, 216, 255, 0.6);
  border-radius: 50%;
  animation: vortex-spin 1.5s linear infinite;
}

.vortex-ring.ring-1 {
  width: 90%;
  height: 90%;
  animation-duration: 1s;
}

.vortex-ring.ring-2 {
  width: 100%;
  height: 100%;
  animation-duration: 1.5s;
  border-color: rgba(63, 216, 255, 0.4);
}

.vortex-ring.ring-3 {
  width: 110%;
  height: 110%;
  animation-duration: 2s;
  border-color: rgba(63, 216, 255, 0.3);
}

@keyframes vortex-spin {
  0% {
    transform: translate(-50%, -50%) rotate(0deg) scale(1);
    opacity: 1;
  }
  100% {
    transform: translate(-50%, -50%) rotate(360deg) scale(0.5);
    opacity: 0;
  }
}

@keyframes absorb-pulse {
  0% {
    border-color: #3fd8ff;
  }
  50% {
    border-color: #fff;
    box-shadow:
      0 0 12px rgba(63, 216, 255, 1),
      0 0 20px rgba(63, 216, 255, 0.6);
  }
  100% {
    border-color: #3fd8ff;
  }
}

@keyframes orb-breathe {
  0%, 100% {
    box-shadow:
      0 0 8px rgba(63, 216, 255, 0.4),
      0 0 15px rgba(63, 216, 255, 0.2);
  }
  50% {
    box-shadow:
      0 0 12px rgba(63, 216, 255, 0.7),
      0 0 20px rgba(63, 216, 255, 0.4);
  }
}

/* 处理状态提示 */
.processing-overlay {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  margin-top: 10px;
  padding: 8px 12px;
  background: rgba(0, 0, 0, 0.85);
  color: white;
  border-radius: 16px;
  font-size: 11px;
  white-space: nowrap;
  pointer-events: none;
  animation: fade-in 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(63, 216, 255, 0.3);
}

.processing-progress {
  background: linear-gradient(135deg, #3fd8ff 0%, #00b4d8 100%);
  color: #000;
  font-weight: 600;
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 8px;
  min-width: 28px;
  text-align: center;
}

.processing-text {
  color: rgba(255, 255, 255, 0.9);
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
}

@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateX(-50%) translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateX(-50%) translateY(0);
  }
}
</style>
