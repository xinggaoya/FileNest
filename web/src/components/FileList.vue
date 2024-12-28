<template>
  <div class="file-list-wrapper">
    <!-- 面包屑导航 -->
    <div class="breadcrumb-section">
      <n-breadcrumb>
        <n-breadcrumb-item @click="handleBack">
          <n-icon><folder-outlined /></n-icon>
          根目录
        </n-breadcrumb-item>
        <n-breadcrumb-item v-for="(path, index) in fileStore.currentPath" :key="index">
          {{ path }}
        </n-breadcrumb-item>
      </n-breadcrumb>
    </div>

    <!-- 文件列表 -->
    <div class="file-list">
      <n-spin :show="loading">
        <n-empty v-if="fileStore.sortedFiles.length === 0" description="暂无文件" />
        <n-list v-else>
          <n-list-item v-for="file in fileStore.sortedFiles" :key="file.filePath" @click="handleFileClick(file)">
            <n-thing>
              <template #icon>
                <n-icon size="20">
                  <folder-outlined v-if="file.isDir" />
                  <file-outlined v-else />
                </n-icon>
              </template>
              <template #header>{{ file.fileName }}</template>
              <template #description>
                <div class="file-info">
                  <span>{{ fileStore.formatFileSize(file.fileSize) }}</span>
                  <span class="time">{{ file.modTime }}</span>
                </div>
              </template>
            </n-thing>
          </n-list-item>
        </n-list>
      </n-spin>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useFileStore } from '@/stores/file'
import { useMessage } from 'naive-ui'
import { FolderOutlined, FileOutlined } from '@vicons/antd'
import type { FileInfo } from '@/api/file/file'

const message = useMessage()
const fileStore = useFileStore()
const loading = ref(false)

// 监听文件列表变化
watch(() => fileStore.sortedFiles, (newFiles) => {
  console.log('文件列表已更新:', newFiles)
}, { deep: true })

// 监听当前路径变化
watch(() => fileStore.currentPath, (newPath) => {
  console.log('当前路径已更新:', newPath)
})

// 初始化时获取文件列表
onMounted(async () => {
  try {
    loading.value = true
    console.log('开始获取文件列表')
    await fileStore.fetchFiles()
    console.log('成功获取文件列表')
  } catch (error) {
    console.error('获取文件列表失败:', error)
    message.error('获取文件列表失败，请稍后重试')
  } finally {
    loading.value = false
  }
})

// 处理���件点击事件
const handleFileClick = async (file: FileInfo) => {
  try {
    loading.value = true
    if (file.isDir) {
      console.log('进入目录:', file.filePath)
      await fileStore.enterDirectory(file.filePath)
    } else {
      console.log('点击文件:', file.filePath)
      // 处理文件点击
    }
  } catch (error) {
    console.error('处理文件点击失败:', error)
    message.error('操作失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

// 处理返回上级目录
const handleBack = async () => {
  try {
    loading.value = true
    console.log('返回上级目录')
    await fileStore.goBack()
  } catch (error) {
    console.error('返回上级目录失败:', error)
    message.error('返回上级目录失败，请稍后重试')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.file-list-wrapper {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 16px;
  gap: 16px;
}

.breadcrumb-section {
  flex: 0 0 auto;
}

.file-list {
  flex: 1 1 auto;
  min-height: 0;
  background-color: var(--n-card-color);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
}

:deep(.n-spin-container) {
  height: 100%;
  display: flex;
  flex-direction: column;
}

:deep(.n-empty) {
  margin: auto;
}

:deep(.n-list) {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

:deep(.n-list-item) {
  cursor: pointer;
  transition: background-color 0.2s;
  border-radius: 4px;
}

:deep(.n-list-item:hover) {
  background-color: rgba(0, 0, 0, 0.05);
}

.file-info {
  color: #6b7280;
  display: flex;
  gap: 16px;
}

/* 美化滚动条 */
:deep(.n-list::-webkit-scrollbar) {
  width: 6px;
}

:deep(.n-list::-webkit-scrollbar-thumb) {
  background-color: rgba(0, 0, 0, 0.2);
  border-radius: 3px;
}

:deep(.n-list::-webkit-scrollbar-track) {
  background-color: transparent;
}
</style> 