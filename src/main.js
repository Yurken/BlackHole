import { createApp } from 'vue'
import App from './App.vue'

// 确保全局背景透明
document.documentElement.style.background = 'transparent'
document.body.style.background = 'transparent'

createApp(App).mount('#app')
