<template>
  <div class="home">
    <!-- 顶部统计和操作区域 -->
    <div class="top-section">
      <div class="stats-container">
        <file-stats />
      </div>
      
      <div class="action-group">
        <n-button-group>
          <n-button quaternary size="large" @click="fileStore.toggleViewMode">
            <template #icon>
              <n-icon>
                <appstore-outlined v-if="fileStore.viewMode === 'grid'" />
                <unordered-list-outlined v-else />
              </n-icon>
            </template>
            {{ fileStore.viewMode === 'grid' ? '列表视图' : '网格视图' }}
          </n-button>
          <n-button quaternary size="large" @click="fileStore.fetchFiles">
            <template #icon>
              <n-icon><reload-outlined /></n-icon>
            </template>
            刷新
          </n-button>
        </n-button-group>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="content-section">
      <n-card class="content-card">
        <!-- 头部导航和搜索 -->
        <div class="content-header">
          <home-header />
        </div>
        
        <!-- 文件列表区域 -->
        <div class="content-body">
          <n-empty v-if="fileStore.sortedFiles.length === 0" description="暂无文件" />
          <template v-else>
            <div class="file-list" :class="{ 'grid-view': fileStore.viewMode === 'grid' }">
              <file-item
                v-for="file in fileStore.sortedFiles"
                :key="file.filePath"
                :file="file"
              />
            </div>
          </template>
        </div>
      </n-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import HomeHeader from '@/components/home/HomeHeader.vue'
import FileItem from '@/components/file/FileItem.vue'
import FileStats from '@/components/home/FileStats.vue'
import { useFileStore } from '@/stores/file'
import {
  AppstoreOutlined,
  UnorderedListOutlined,
  ReloadOutlined
} from '@vicons/antd'

const fileStore = useFileStore()

onMounted(() => {
  fileStore.fetchFiles()
})
</script>

<style scoped>
.home {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: var(--n-color-background);
  padding: 16px;
  gap: 16px;
  box-sizing: border-box;
}

.top-section {
  flex: 0 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.stats-container {
  flex: 1;
  min-width: 300px;
}

.action-group {
  flex: 0 0 auto;
}

.content-section {
  flex: 1;
  min-height: 0;
  position: relative;
}

:deep(.content-card) {
  height: 100%;
}

:deep(.n-card__content) {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 0 !important;
}

.content-header {
  flex: 0 0 auto;
  border-bottom: 1px solid var(--n-border-color);
  padding: 16px;
}

.content-body {
  flex: 1;
  min-height: 0;
  position: relative;
  height: 100%;
  overflow: hidden;
}

.file-list {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow-y: scroll;
  padding: 16px;
}

.file-list.grid-view {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 16px;
}

/* 滚动条样式 */
.file-list::-webkit-scrollbar {
  width: 8px;
  background-color: transparent;
}

.file-list::-webkit-scrollbar-track {
  background: var(--n-color-background);
  border-radius: 4px;
}

.file-list::-webkit-scrollbar-thumb {
  background-color: rgba(0, 0, 0, 0.2);
  border-radius: 4px;
  border: 2px solid var(--n-color-background);
}

.file-list::-webkit-scrollbar-thumb:hover {
  background-color: rgba(0, 0, 0, 0.3);
}

/* 响应式布局 */
@media (max-width: 768px) {
  .home {
    padding: 8px;
  }

  .top-section {
    flex-direction: column;
  }

  .stats-container {
    width: 100%;
    min-width: 100%;
  }

  .action-group {
    width: 100%;
    display: flex;
    justify-content: flex-end;
  }

  .file-list.grid-view {
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 8px;
  }
}
</style>
