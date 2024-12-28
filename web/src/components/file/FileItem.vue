<template>
  <div class="file-item" :class="[
    { 'is-folder': file.isDir },
    { 'is-image': !file.isDir && isImage },
    { 'is-text': !file.isDir && isText },
    { 'is-other': !file.isDir && !isImage && !isText },
    { 'grid-mode': fileStore.viewMode === 'grid' }
  ]" @click="handleClick" @mouseenter="isHovered = true" @mouseleave="isHovered = false">
    <!-- 网格模式 -->
    <template v-if="fileStore.viewMode === 'grid'">
      <div class="grid-content">
        <div class="file-icon-wrapper">
          <div class="file-icon">
            <transition name="bounce">
              <n-icon size="32" v-if="file.isDir">
                <folder-open-outlined v-if="isHovered" />
                <folder-outlined v-else />
              </n-icon>
              <n-icon size="32" v-else>
                <component :is="getFileIcon(file.fileName)" />
              </n-icon>
            </transition>
          </div>
          <div class="file-type" v-if="!file.isDir">
            {{ file.fileName.split('.').pop()?.toUpperCase() }}
          </div>
        </div>
        <div class="file-info">
          <div class="file-name" :title="file.fileName">{{ file.fileName }}</div>
          <div class="file-meta">
            <span>{{ file.isDir ? '文件夹' : formatFileSize(file.fileSize) }}</span>
            <span class="dot">·</span>
            <span>{{ formatDate(file.modTime) }}</span>
          </div>
        </div>
        <div class="file-actions">
          <n-button-group>
            <n-button quaternary circle size="small" @click.stop="handlePreview" v-if="canPreview">
              <template #icon><eye-outlined /></template>
            </n-button>
            <n-button quaternary circle size="small" @click.stop="handleDownload" v-if="!file.isDir">
              <template #icon><download-outlined /></template>
            </n-button>
            <n-button quaternary circle size="small" @click.stop="handleDelete">
              <template #icon><delete-outlined /></template>
            </n-button>
          </n-button-group>
        </div>
      </div>
    </template>

    <!-- 列表模式 -->
    <template v-else>
      <div class="file-main">
        <div class="file-icon-wrapper">
          <div class="file-icon">
            <transition name="bounce">
              <n-icon size="24" v-if="file.isDir">
                <folder-open-outlined v-if="isHovered" />
                <folder-outlined v-else />
              </n-icon>
              <n-icon size="24" v-else>
                <component :is="getFileIcon(file.fileName)" />
              </n-icon>
            </transition>
          </div>
          <div class="file-type" v-if="!file.isDir">
            {{ file.fileName.split('.').pop()?.toUpperCase() }}
          </div>
        </div>
        <div class="file-name" :title="file.fileName">{{ file.fileName }}</div>
        <div class="file-size">{{ file.isDir ? '文件夹' : formatFileSize(file.fileSize) }}</div>
        <div class="file-date">{{ formatDate(file.modTime) }}</div>
      </div>
      <div class="file-actions">
        <n-button-group>
          <n-button quaternary circle size="small" @click.stop="handlePreview" v-if="canPreview">
            <template #icon><eye-outlined /></template>
          </n-button>
          <n-button quaternary circle size="small" @click.stop="handleDownload" v-if="!file.isDir">
            <template #icon><download-outlined /></template>
          </n-button>
          <n-button quaternary circle size="small" @click.stop="handleDelete">
            <template #icon><delete-outlined /></template>
          </n-button>
        </n-button-group>
      </div>
    </template>
  </div>

  <!-- 预览弹窗 -->
  <n-modal v-model:show="showPreview" preset="card" style="max-width: 90vw; max-height: 90vh">
    <template #header>
      <div class="preview-header">
        <span>{{ file.fileName }}</span>
      </div>
    </template>
    <div class="preview-content">
      <img v-if="isImage" :src="previewUrl" alt="预览图片" />
      <pre v-else-if="isText" v-text="previewContent" />
      <div v-else class="preview-unsupported">
        暂不支持预览该类型文件
      </div>
    </div>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { FileInfo } from '@/api/file/file'
import { getDownloadUrl } from '@/api/file/file'
import { useFileStore } from '@/stores/file'
import {
  FolderOutlined,
  FolderOpenOutlined,
  FileTextOutlined,
  FileImageOutlined,
  FileUnknownOutlined,
  EyeOutlined,
  DownloadOutlined,
  DeleteOutlined,
  CloseOutlined
} from '@vicons/antd'
import { createDiscreteApi } from 'naive-ui'

const { message } = createDiscreteApi(['message'])
const fileStore = useFileStore()

const props = defineProps<{
  file: FileInfo
}>()

const isHovered = ref(false)
const showPreview = ref(false)
const previewContent = ref('')

// 获取文件图标
const getFileIcon = (fileName: string) => {
  const ext = fileName.split('.').pop()?.toLowerCase()
  if (['jpg', 'jpeg', 'png', 'gif', 'webp'].includes(ext || '')) {
    return FileImageOutlined
  }
  if (['txt', 'md', 'json', 'js', 'ts', 'html', 'css'].includes(ext || '')) {
    return FileTextOutlined
  }
  return FileUnknownOutlined
}

// 判断是否可以预览
const canPreview = computed(() => {
  if (props.file.isDir) return false
  const ext = props.file.fileName.split('.').pop()?.toLowerCase()
  return ['jpg', 'jpeg', 'png', 'gif', 'webp', 'txt', 'md', 'json', 'js', 'ts', 'html', 'css'].includes(ext || '')
})

const isImage = computed(() => {
  const ext = props.file.fileName.split('.').pop()?.toLowerCase()
  return ['jpg', 'jpeg', 'png', 'gif', 'webp'].includes(ext || '')
})

const isText = computed(() => {
  const ext = props.file.fileName.split('.').pop()?.toLowerCase()
  return ['txt', 'md', 'json', 'js', 'ts', 'html', 'css'].includes(ext || '')
})

// 获取预览URL
const previewUrl = computed(() => {
  if (!props.file.isDir) {
    return getDownloadUrl(props.file.filePath)
  }
  return ''
})

// 格式化日期
const formatDate = (date: string) => {
  return new Date(date).toLocaleString()
}

// 格式化文件大小
const formatFileSize = (size: number) => {
  if (size < 1024) return size + ' B'
  if (size < 1024 * 1024) return (size / 1024).toFixed(2) + ' KB'
  if (size < 1024 * 1024 * 1024) return (size / 1024 / 1024).toFixed(2) + ' MB'
  return (size / 1024 / 1024 / 1024).toFixed(2) + ' GB'
}

// 处理点击事件
const handleClick = () => {
  if (props.file.isDir) {
    fileStore.enterDirectory(props.file.fileName)
  }
}

// 处理预览
const handlePreview = async () => {
  if (!canPreview.value) return
  
  if (isText.value) {
    try {
      const response = await fetch(previewUrl.value)
      previewContent.value = await response.text()
    } catch (error) {
      message.error('获取文件内容失败')
      return
    }
  }
  
  showPreview.value = true
}

// 处理下载
const handleDownload = () => {
  if (props.file.isDir) return
  window.open(previewUrl.value, '_blank')
}

// 处理删除
const handleDelete = async () => {
  if (confirm('确定要删除该文件吗？')) {
    await fileStore.removeFile(props.file.filePath)
  }
}
</script>

<style scoped>
.file-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  margin-bottom: 8px;
  transition: all 0.3s ease;
  cursor: pointer;
  border-radius: 8px;
  background-color: var(--n-card-color);
  border: 1px solid transparent;
}

.file-item:hover {
  background-color: var(--n-color-hover);
  border-color: var(--n-border-color-hover);
}

/* 网格模式样式 */
.file-item.grid-mode {
  flex-direction: column;
  padding: 12px;
  margin: 0;
  height: auto;
  min-height: 160px;
  max-height: 180px;
}

.grid-content {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.grid-mode .file-icon-wrapper {
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 4px;
}

.grid-mode .file-icon {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.03);
  border-radius: 8px;
}

.grid-mode .file-info {
  width: 100%;
  text-align: center;
  flex: 0 1 auto;
}

.grid-mode .file-name {
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 4px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 0;
  line-height: 1.4;
}

.grid-mode .file-meta {
  font-size: 12px;
  color: var(--n-text-color-3);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  line-height: 1.2;
}

.grid-mode .file-actions {
  opacity: 0;
  transition: opacity 0.3s ease;
  width: 100%;
  display: flex;
  justify-content: center;
  padding-top: 4px;
}

.grid-mode:hover .file-actions {
  opacity: 1;
}

/* 列表模式样式 */
.file-main {
  display: flex;
  align-items: center;
  gap: 16px;
  flex: 1;
  min-width: 0;
}

.file-icon-wrapper {
  position: relative;
  flex-shrink: 0;
}

.file-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.file-type {
  position: absolute;
  bottom: -2px;
  right: -2px;
  background-color: var(--n-primary-color);
  color: white;
  font-size: 9px;
  font-weight: bold;
  padding: 1px 3px;
  border-radius: 3px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.file-name {
  flex: 2;
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--n-text-color);
  min-width: 200px;
}

.file-size, .file-date {
  flex: 1;
  font-size: 13px;
  color: var(--n-text-color-3);
  white-space: nowrap;
}

.file-actions {
  flex-shrink: 0;
}

/* 动画 */
.bounce-enter-active {
  animation: bounce-in 0.3s;
}

@keyframes bounce-in {
  0% {
    transform: scale(0.3);
  }
  50% {
    transform: scale(1.1);
  }
  100% {
    transform: scale(1);
  }
}

/* 预览样式 */
.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.preview-content {
  max-height: 80vh;
  overflow: auto;
}

.preview-content img {
  max-width: 100%;
  height: auto;
}

.preview-content pre {
  white-space: pre-wrap;
  word-wrap: break-word;
  padding: 16px;
  background-color: var(--n-color-hover);
  border-radius: 8px;
}

.preview-unsupported {
  padding: 32px;
  text-align: center;
  color: var(--n-text-color-3);
}
</style> 