<script setup lang="ts">
import { RouterView } from 'vue-router'
import { zhCN, darkTheme } from 'naive-ui'
import { useOsTheme } from 'naive-ui'
import { computed } from 'vue'
import themeOverrides from '@/assets/theme/naive-ui-theme-overrides.json'

const osThemeRef = useOsTheme()
const theme = computed(() => (osThemeRef.value === 'dark' ? darkTheme : null))
</script>

<template>
  <n-config-provider
    class="height-100"
    :theme="theme"
    :theme-overrides="themeOverrides"
    :locale="zhCN"
  >
    <n-loading-bar-provider>
      <n-message-provider>
        <n-notification-provider>
          <n-dialog-provider>
            <n-modal-provider>
              <div class="app-container">
                <RouterView />
              </div>
            </n-modal-provider>
          </n-dialog-provider>
        </n-notification-provider>
      </n-message-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<style>
html,
body {
  margin: 0;
  padding: 0;
  width: 100%;
  height: 100%;
}

.height-100 {
  height: 100%;
}

.app-container {
  height: 100%;
  background-color: var(--body-color);
}

:root {
  --body-color: #f5f7fa;
}

.n-config-provider.n-theme-dark {
  --body-color: #101014;
}

/* 自定义滚动条样式 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 3px;
}

.n-theme-dark ::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.3);
}

.n-theme-dark ::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}

/* 全局过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-fade-enter-active {
  transition: all 0.3s ease-out;
}

.slide-fade-leave-active {
  transition: all 0.3s cubic-bezier(1, 0.5, 0.8, 1);
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateX(20px);
  opacity: 0;
}

/* 添加响应式布局支持 */
@media (max-width: 768px) {
  .file-manager {
    padding: 12px !important;
  }

  .home-card {
    border-radius: 0 !important;
  }
}

/* 优化打印样式 */
@media print {
  .app-container {
    background: white !important;
  }

  .n-watermark {
    display: none !important;
  }
}
</style>
