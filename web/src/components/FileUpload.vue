<template>
  <div class="file-upload-container">
    <h3>文件上传</h3>

    <!-- 上传区域 -->
    <div
      class="upload-area"
      :class="{ 'drag-over': isDragOver }"
      @dragover.prevent="handleDragOver"
      @dragleave.prevent="handleDragLeave"
      @drop.prevent="handleDrop"
      @click="triggerFileInput"
    >
      <div class="upload-content">
        <n-icon size="48" class="upload-icon">
          <cloud-upload-outlined />
        </n-icon>
        <p class="upload-text">点击选择文件或拖拽文件到此处</p>
        <p class="upload-hint">支持多个文件同时上传</p>
      </div>

      <input
        ref="fileInputRef"
        type="file"
        multiple
        webkitdirectory
        style="display: none"
        @change="handleFileSelect"
      />
    </div>

    <!-- 上传选择按钮 -->
    <div class="upload-controls">
      <n-space>
        <n-button @click="triggerFileOnly">选择文件</n-button>
        <n-button @click="triggerFileInput">选择文件夹</n-button>
      </n-space>
    </div>

    <!-- 文件选择输入框（无文件夹支持） -->
    <input
      ref="fileOnlyInputRef"
      type="file"
      multiple
      style="display: none"
      @change="handleFileSelect"
    />

    <!-- 当前路径显示 -->
    <div class="current-path">
      <span>上传到：{{ currentPath || '根目录' }}</span>
    </div>

    <!-- 上传设置 -->
    <div class="upload-settings">
      <n-space>
        <n-checkbox v-model:checked="uploadConfig.override"> 覆盖同名文件 </n-checkbox>
        <n-button @click="clearCompleted" :disabled="Object.keys(uploadQueue).length === 0">
          清空已完成
        </n-button>
      </n-space>
    </div>

    <!-- 上传队列 -->
    <div class="upload-queue" v-if="Object.keys(uploadQueue).length > 0">
      <h4>上传队列</h4>
      <div class="upload-item" v-for="(item, key) in uploadQueue" :key="key">
        <div class="file-info">
          <n-icon class="file-icon">
            <file-outlined />
          </n-icon>
          <div class="file-details">
            <div class="file-name">{{ item.fileName }}</div>
            <div class="file-size">{{ formatFileSize(item.total) }}</div>
          </div>
        </div>

        <div class="upload-progress">
          <n-progress
            type="line"
            :percentage="item.percentage"
            :status="getProgressStatus(item.status)"
            :show-indicator="false"
          />
          <div class="progress-text">
            <span>{{ item.percentage }}%</span>
            <span class="status" :class="item.status">{{ getStatusText(item.status) }}</span>
          </div>
        </div>

        <div class="upload-actions">
          <n-button
            v-if="item.status === 'error'"
            size="small"
            type="primary"
            @click="retryUpload(key as string)"
          >
            重试
          </n-button>
          <n-button size="small" type="error" @click="removeFromQueue(key as string)">
            移除
          </n-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useFileStore } from '@/stores/file'
import { useMessageHandler } from '@/composables/useMessageHandler'
import { uploadFile } from '@/api/file'
import { CloudUploadOutlined, FileOutlined } from '@vicons/antd'
import type { UploadProgress, UploadConfig } from '@/types/file'

const { handleResult, message } = useMessageHandler()
const fileStore = useFileStore()

const fileInputRef = ref<HTMLInputElement>()
const fileOnlyInputRef = ref<HTMLInputElement>()
const isDragOver = ref(false)
const uploadQueue = ref<Record<string, UploadProgress>>({})

const uploadConfig = reactive({
  override: false
})

const currentPath = computed(() => fileStore.currentPath)

// 触发文件选择
const triggerFileInput = () => {
  fileInputRef.value?.click()
}

// 触发仅文件选择
const triggerFileOnly = () => {
  fileOnlyInputRef.value?.click()
}

// 处理文件选择
const handleFileSelect = (event: Event) => {
  const input = event.target as HTMLInputElement
  if (input.files) {
    handleFiles(Array.from(input.files))
  }
  // 清空 input 值，允许选择相同文件
  input.value = ''
}

// 处理拖拽事件
const handleDragOver = (event: DragEvent) => {
  isDragOver.value = true
}

const handleDragLeave = (event: DragEvent) => {
  isDragOver.value = false
}

const handleDrop = (event: DragEvent) => {
  isDragOver.value = false
  if (event.dataTransfer?.files) {
    handleFiles(Array.from(event.dataTransfer.files))
  }
}

// 处理文件列表
const handleFiles = (files: File[]) => {
  // 对文件进行排序，确保文件夹在上传前创建
  const sortedFiles = files.sort((a, b) => {
    const pathA = (a as any).webkitRelativePath || ''
    const pathB = (b as any).webkitRelativePath || ''
    return pathA.localeCompare(pathB)
  })

  sortedFiles.forEach((file) => {
    uploadFileToServer(file)
  })
}

// 上传文件到服务器
const uploadFileToServer = async (file: File) => {
  const uploadId = `${file.name}-${Date.now()}`
  const relativePath = (file as any).webkitRelativePath || ''

  // 添加到上传队列
  const uploadProgress: UploadProgress = {
    fileName: file.name,
    loaded: 0,
    total: file.size,
    percentage: 0,
    status: 'pending',
    relativePath: relativePath
  }
  uploadQueue.value[uploadId] = uploadProgress

  try {
    // 更新状态为上传中
    uploadProgress.status = 'uploading'
    uploadQueue.value[uploadId] = { ...uploadProgress }

    // 准备上传配置
    const config: UploadConfig = {
      fileName: file.name,
      filePath: currentPath.value,
      fileSize: file.size,
      override: uploadConfig.override,
      relativePath: relativePath
    }

    // 开始上传
    await uploadFile(file, config, (progressEvent) => {
      if (progressEvent.total) {
        const percentage = Math.round((progressEvent.loaded / progressEvent.total) * 100)
        const updatedProgress = {
          ...uploadProgress,
          loaded: progressEvent.loaded,
          percentage: percentage
        }
        uploadQueue.value[uploadId] = updatedProgress
      }
    })

    // 上传成功
    uploadProgress.status = 'success'
    uploadProgress.percentage = 100
    uploadQueue.value[uploadId] = { ...uploadProgress }

    message.success(`文件 ${file.name} 上传成功`)

    // 刷新文件列表
    const refreshResult = await fileStore.fetchFiles()
    if (refreshResult && !refreshResult.success) {
      handleResult(refreshResult)
    }
  } catch (error) {
    // 上传失败
    uploadProgress.status = 'error'
    uploadProgress.error = error instanceof Error ? error.message : '上传失败'
    uploadQueue.value[uploadId] = { ...uploadProgress }

    message.error(`文件 ${file.name} 上传失败: ${uploadProgress.error}`)
  }
}

// 重试上传
const retryUpload = (uploadId: string) => {
  const item = uploadQueue.value[uploadId]
  if (item) {
    // 重置状态
    item.status = 'pending'
    item.percentage = 0
    item.loaded = 0
    item.error = undefined
    uploadQueue.value[uploadId] = { ...item }

    // 这里需要重新创建 File 对象，但由于安全限制，我们只能提示用户重新选择文件
    message.info('请重新选择文件进行上传')
    removeFromQueue(uploadId)
  }
}

// 从队列中移除
const removeFromQueue = (uploadId: string) => {
  delete uploadQueue.value[uploadId]
}

// 清空已完成的上传
const clearCompleted = () => {
  for (const key in uploadQueue.value) {
    const item = uploadQueue.value[key]
    if (item.status === 'success' || item.status === 'error') {
      delete uploadQueue.value[key]
    }
  }
}

// 获取进度条状态
const getProgressStatus = (status: UploadProgress['status']) => {
  switch (status) {
    case 'success':
      return 'success'
    case 'error':
      return 'error'
    default:
      return 'default'
  }
}

// 获取状态文本
const getStatusText = (status: UploadProgress['status']) => {
  switch (status) {
    case 'pending':
      return '等待中'
    case 'uploading':
      return '上传中'
    case 'success':
      return '成功'
    case 'error':
      return '失败'
    default:
      return ''
  }
}

// 格式化文件大小
const formatFileSize = (size: number): string => {
  if (size === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const index = Math.floor(Math.log(size) / Math.log(1024))
  return Math.round((size / Math.pow(1024, index)) * 100) / 100 + ' ' + units[index]
}
</script>

<style scoped>
.file-upload-container {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.upload-area {
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  padding: 40px 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background-color: #fafafa;
  margin-bottom: 16px;
}

.upload-area:hover,
.upload-area.drag-over {
  border-color: #1890ff;
  background-color: #f0f8ff;
}

.upload-content {
  color: #666;
}

.upload-icon {
  color: #999;
  margin-bottom: 16px;
}

.upload-text {
  font-size: 16px;
  margin: 8px 0;
  color: #333;
}

.upload-hint {
  font-size: 14px;
  color: #999;
  margin: 0;
}

.current-path {
  background-color: #f5f5f5;
  padding: 8px 12px;
  border-radius: 4px;
  margin-bottom: 16px;
  font-size: 14px;
  color: #666;
}

.upload-settings {
  margin-bottom: 20px;
}

.upload-queue {
  background-color: #f9f9f9;
  border-radius: 8px;
  padding: 16px;
}

.upload-queue h4 {
  margin: 0 0 16px 0;
  color: #333;
}

.upload-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px;
  background-color: white;
  border-radius: 6px;
  margin-bottom: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.file-info {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 200px;
}

.file-icon {
  color: #1890ff;
}

.file-details {
  flex: 1;
}

.file-name {
  font-weight: 500;
  color: #333;
  word-break: break-all;
}

.file-size {
  font-size: 12px;
  color: #999;
}

.upload-progress {
  flex: 1;
  min-width: 200px;
}

.progress-text {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  margin-top: 4px;
}

.status {
  font-weight: 500;
}

.status.pending {
  color: #faad14;
}

.status.uploading {
  color: #1890ff;
}

.status.success {
  color: #52c41a;
}

.status.error {
  color: #f5222d;
}

.upload-actions {
  display: flex;
  gap: 8px;
  min-width: 100px;
  justify-content: flex-end;
}

.upload-controls {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

@media (max-width: 768px) {
  .upload-item {
    flex-direction: column;
    align-items: stretch;
  }

  .upload-actions {
    justify-content: center;
  }
}
</style>
