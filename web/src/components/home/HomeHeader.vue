<template>
  <div class="header-wrapper">
    <div class="left-section">
      <div class="nav-path">
        <n-breadcrumb>
          <n-breadcrumb-item @click="handleRootClick">
            <n-icon><folder-outlined /></n-icon>
            <n-spin size="small" v-if="fileStore.isLoading && !fileStore.currentPath.length" />
            <template v-else>根目录</template>
          </n-breadcrumb-item>
          <template v-for="(path, index) in fileStore.currentPath" :key="index">
            <n-icon><right-outlined /></n-icon>
            <n-breadcrumb-item @click="handleBreadcrumbClick(index)">
              <n-spin size="small" v-if="fileStore.isLoading && index === loadingIndex" />
              <template v-else>{{ path }}</template>
            </n-breadcrumb-item>
          </template>
        </n-breadcrumb>
      </div>
    </div>

    <!-- 右侧搜索框和操作按钮 -->
    <div class="right-section">
      <div class="search-and-actions">
        <!-- 搜索框 -->
        <div class="search-box">
          <n-input-group>
            <n-input
              v-model:value="searchKeyword"
              type="text"
              placeholder="搜索文件..."
              @keydown.enter="handleSearch"
            >
              <template #prefix>
                <n-icon><search-outlined /></n-icon>
              </template>
            </n-input>
            <n-button type="primary" ghost @click="handleSearch">
              搜索
            </n-button>
          </n-input-group>
        </div>

        <!-- 操作按钮 -->
        <div class="action-buttons">
          <n-button type="primary" @click="showUploadDrawer = true">
            <template #icon>
              <n-icon><upload-outlined /></n-icon>
            </template>
            上传
          </n-button>
          <n-button @click="handleCreateFolder">
            <template #icon>
              <n-icon><folder-add-outlined /></n-icon>
            </template>
            新建文件夹
          </n-button>
        </div>
      </div>
    </div>
  </div>

  <!-- 上传抽屉 -->
  <n-drawer
    v-model:show="showUploadDrawer"
    :width="500"
    placement="right"
    title="上传文件"
  >
    <n-drawer-content>
      <!-- 上传设置 -->
      <div class="upload-settings">
        <n-checkbox v-model:checked="uploadSettings.override">
          覆盖已存在的文件
        </n-checkbox>
      </div>

      <n-tabs type="segment" class="upload-tabs">
        <n-tab-pane name="file" tab="文件上传">
          <!-- 上传区域 -->
          <n-space vertical>
            <div class="upload-area">
              <n-upload
                :custom-request="handleUpload"
                :show-file-list="false"
                :directory="false"
                :max-size="1024 * 1024 * 500"
                @change="handleUploadChange"
                @before-upload="handleBeforeUpload"
              >
                <n-upload-dragger>
                  <div class="upload-dragger-content">
                    <n-icon size="48" depth="3">
                      <upload-outlined />
                    </n-icon>
                    <n-text>点击或拖拽文件到此处上传</n-text>
                    <div class="upload-hint">
                      <p>支持单个文件最大 500MB</p>
                      <p>支持多文件同时上传</p>
                    </div>
                  </div>
                </n-upload-dragger>
              </n-upload>
            </div>
          </n-space>
        </n-tab-pane>

        <n-tab-pane name="folder" tab="文件夹上传">
          <n-upload
            :custom-request="handleFolderUpload"
            :show-file-list="false"
            directory
            :max-size="1024 * 1024 * 500"
            @change="handleUploadChange"
          >
            <n-upload-dragger>
              <div class="upload-dragger-content">
                <n-icon size="48" depth="3">
                  <folder-add-outlined />
                </n-icon>
                <n-text>点击或拖拽文件夹到此处上传</n-text>
                <div class="upload-hint">
                  <p>支持文件夹及其子文件夹的上传</p>
                  <p>自动保持文件夹结构</p>
                </div>
              </div>
            </n-upload-dragger>
          </n-upload>
        </n-tab-pane>
      </n-tabs>

      <!-- 上传进度列表 -->
      <div class="upload-list" v-if="Object.keys(uploadStatus).length > 0">
        <div class="upload-list-header">
          <span>文件上传列表</span>
          <n-button text type="primary" @click="clearFinishedUploads">
            清除已完成
          </n-button>
        </div>
        <div class="upload-items">
          <div 
            class="upload-item" 
            v-for="(status, key) in uploadStatus" 
            :key="key"
            :class="{ 
              'is-finished': status.status === 'finished',
              'is-error': status.status === 'error'
            }"
          >
            <div class="upload-item-header">
              <n-icon size="20">
                <file-outlined />
              </n-icon>
              <span class="filename" :title="key">{{ key }}</span>
              <span class="filesize">{{ formatFileSize(status.size || 0) }}</span>
            </div>
            
            <div class="upload-item-body">
              <div class="progress-info">
                <span class="status">
                  {{ 
                    status.status === 'uploading' ? 
                      `上传中 - ${formatFileSize(status.speed || 0)}/s` :
                    status.status === 'finished' ? '上传完成' :
                    status.status === 'waiting' ? '等待上传' :
                    '上传失败'
                  }}
                </span>
                <span class="progress-text">{{ status.progress }}%</span>
              </div>
              <n-progress
                type="line"
                :percentage="status.progress"
                :status="status.status === 'error' ? 'error' :
                        status.status === 'finished' ? 'success' : 'info'"
                :show-indicator="false"
                :height="2"
              />
            </div>

            <div class="upload-item-footer">
              <div class="error-message" v-if="status.message">{{ status.message }}</div>
              <div class="actions">
                <n-button 
                  v-if="status.status === 'error'" 
                  text 
                  type="primary"
                  @click="retryUpload(key)"
                >
                  重试
                </n-button>
                <n-button 
                  text 
                  type="error"
                  @click="removeUpload(key)"
                >
                  移除
                </n-button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </n-drawer-content>
  </n-drawer>

  <!-- 创建文件夹对话框 -->
  <n-modal
    v-model:show="showCreateFolderDialog"
    preset="dialog"
    title="新建文件夹"
    positive-text="确定"
    negative-text="取消"
    @positive-click="handleConfirmCreateFolder"
  >
    <n-input
      v-model:value="newFolderName"
      placeholder="请输入文件夹名称"
      @keyup.enter="handleConfirmCreateFolder"
    />
  </n-modal>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import type { UploadFileInfo } from 'naive-ui'
import type { UploadCustomRequestOptions } from 'naive-ui'
import {
  FolderOutlined,
  FolderAddOutlined,
  UploadOutlined,
  RightOutlined,
  SearchOutlined,
  FileOutlined
} from '@vicons/antd'
import { uploadFile } from '@/api/file/file'
import { useFileStore } from '@/stores/file'
import { createDiscreteApi } from 'naive-ui'

const { message } = createDiscreteApi(['message'])
const fileStore = useFileStore()

// 搜索相关
const searchKeyword = ref('')
const handleSearch = () => {
  if (!searchKeyword.value.trim()) {
    return
  }
  fileStore.searchFile(searchKeyword.value)
}

// 上传相关
const showUploadDrawer = ref(false)
interface UploadStatus {
  progress: number
  status: 'waiting' | 'uploading' | 'finished' | 'error'
  message?: string
  size?: number
  speed?: number
  file?: File
  startTime?: number
}

const uploadStatus = reactive<Record<string, UploadStatus>>({})
const uploadQueue = ref<string[]>([])
const maxConcurrentUploads = 3
let activeUploads = 0

// 上传设置
const uploadSettings = reactive({
  override: false
})

// 处理上传前的检查
const handleBeforeUpload = async (data: { file: UploadFileInfo }) => {
  const file = data.file.file as File
  const fileName = file.name
  
  // 检查文件是否已存在
  if (uploadStatus[fileName]) {
    message.warning(`文件 ${fileName} 已在上传列表中`)
    return false
  }

  // 如果是文件夹上传，先创建文件夹结构
  const relativePath = (file as any).webkitRelativePath || ''
  if (relativePath) {
    const folders = relativePath.split('/').slice(0, -1)
    if (folders.length > 0) {
      try {
        let currentPath = fileStore.currentPathString || ''
        for (const folder of folders) {
          if (!folder) continue
          currentPath = currentPath ? `${currentPath}/${folder}` : folder
          try {
            await fileStore.createNewFolder(currentPath)
          } catch (error: any) {
            // 如果文件夹已存在，继续处理
            if (!error.message?.includes('已存在')) {
              throw error
            }
          }
        }
      } catch (error) {
        console.error('创建文件夹结构失败:', error)
        message.error('创建文件夹结构失败')
        return false
      }
    }
  }

  // 添加到上传状态和队列
  uploadStatus[fileName] = {
    progress: 0,
    status: 'waiting',
    size: file.size,
    file: file
  }
  uploadQueue.value.push(fileName)
  
  // 尝试开始上传
  processUploadQueue()
  
  // 返回 false 阻止默认上传行为
  return false
}

// 处理上传队列
const processUploadQueue = async () => {
  while (uploadQueue.value.length > 0 && activeUploads < maxConcurrentUploads) {
    const fileName = uploadQueue.value[0]
    const status = uploadStatus[fileName]
    
    if (status && status.status === 'waiting' && status.file) {
      activeUploads++
      uploadQueue.value.shift()
      await doUploadFile(fileName, status.file)
      activeUploads--
      processUploadQueue()
    }
  }
}

// 处理文件上传
const handleUpload = async ({ file, onFinish, onError, onProgress }: UploadCustomRequestOptions) => {
  if (!file) return
  await handleBeforeUpload({ file })
}

// 处理文件夹上传
const handleFolderUpload = async (options: UploadCustomRequestOptions) => {
  const { file } = options
  if (!file) return

  const fileName = file.name
  const relativePath = (file as any).webkitRelativePath || ''
  const folderPath = relativePath.split('/')
  
  // 检查文件是否已存在
  if (uploadStatus[fileName]) {
    message.warning(`文件 ${fileName} 已在上传列表中`)
    return
  }

  // 如果是文件夹内的文件
  if (folderPath.length > 1) {
    try {
      // 构建完整的文件夹路径
      let currentPath = fileStore.currentPathString || ''
      const folders = folderPath.slice(0, -1) // 去掉文件名，只保留文件夹路径
      
      // 逐级创建文件夹
      for (let i = 0; i < folders.length; i++) {
        const folder = folders[i]
        if (!folder) continue
        
        // 计算当前层级的完整路径
        if (currentPath) {
          currentPath = `${currentPath}/${folder}`
        } else {
          currentPath = folder
        }

        try {
          // 尝试创建文件夹
          await fileStore.createNewFolder(currentPath)
        } catch (error: any) {
          // 如果文件夹已存在，继续处理
          if (!error.message?.includes('已存在')) {
            throw error
          }
        }
      }
    } catch (error) {
      console.error('创建文件夹结构失败:', error)
      message.error('创建文件夹结构失败')
      return
    }
  }

  // 添加到上传状态和队列
  uploadStatus[fileName] = {
    progress: 0,
    status: 'waiting',
    size: file.file?.size,
    file: file.file as File,
    startTime: Date.now()
  }
  uploadQueue.value.push(fileName)
  
  // 开始上传
  processUploadQueue()

  // 调用上传完成回调
  options.onFinish?.()
}

// 处理上传变化
const handleUploadChange = (data: { fileList: UploadFileInfo[] }) => {
  console.log('Upload change:', data)
}

// 上传文件的具体实现
const doUploadFile = async (fileName: string, file: File) => {
  const status = uploadStatus[fileName]
  if (!status) return

  status.status = 'uploading'
  status.startTime = Date.now()
  status.speed = 0

  try {
    // 获取文件的相对路径
    const relativePath = (file as any).webkitRelativePath || ''
    let path = fileStore.currentPathString || ''
    
    // 如果是文件夹上传（有相对路径），则构建完整路径
    if (relativePath) {
      const folderPath = relativePath.split('/').slice(0, -1).join('/')
      path = path ? `${path}/${folderPath}` : folderPath
    }

    console.log('开始上传文件:', {
      fileName,
      relativePath,
      path,
      size: file.size,
      override: uploadSettings.override
    })

    await uploadFile({
      file: file,
      path: path,
      override: uploadSettings.override,
      onProgress: (progress: number) => {
        console.log(`文件 ${fileName} 上传进度:`, progress)
        status.progress = progress
        // 计算上传速度
        const elapsedTime = (Date.now() - (status.startTime || 0)) / 1000
        if (elapsedTime > 0) {
          status.speed = (file.size * (progress / 100)) / elapsedTime
        }
      },
      onSuccess: () => {
        console.log(`文件 ${fileName} 上传完成`)
        status.status = 'finished'
        status.progress = 100
        message.success(`${fileName} 上传成功`)
        fileStore.fetchFiles()
      },
      onError: (error: string) => {
        console.error(`文件 ${fileName} 上传失败:`, error)
        status.status = 'error'
        status.message = error
        message.error(`${fileName} 上传失败`)
      }
    })
  } catch (error) {
    console.error(`文件 ${fileName} 上传出错:`, error)
    status.status = 'error'
    status.message = error instanceof Error ? error.message : '上传失败，请重试'
    message.error(`${fileName} 上传失败`)
  }
}

// 重试上传
const retryUpload = (fileName: string) => {
  const status = uploadStatus[fileName]
  if (status && status.file) {
    status.status = 'waiting'
    status.progress = 0
    status.message = undefined
    uploadQueue.value.push(fileName)
    processUploadQueue()
  }
}

// 移除上传
const removeUpload = (fileName: string) => {
  delete uploadStatus[fileName]
  const queueIndex = uploadQueue.value.indexOf(fileName)
  if (queueIndex > -1) {
    uploadQueue.value.splice(queueIndex, 1)
  }
}

// 清除已完成的上传
const clearFinishedUploads = () => {
  Object.entries(uploadStatus).forEach(([fileName, status]) => {
    if (status.status === 'finished') {
      delete uploadStatus[fileName]
    }
  })
}

// 格式化文件大小
const formatFileSize = (size: number) => {
  if (size === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const k = 1024
  const i = Math.floor(Math.log(size) / Math.log(k))
  return parseFloat((size / Math.pow(k, i)).toFixed(2)) + ' ' + units[i]
}

// 文件夹操作相关
const showCreateFolderDialog = ref(false)
const newFolderName = ref('')

const handleCreateFolder = () => {
  showCreateFolderDialog.value = true
}

const handleConfirmCreateFolder = async () => {
  if (!newFolderName.value.trim()) {
    message.warning('请输入文件夹名称')
    return
  }

  try {
    await fileStore.createNewFolder(newFolderName.value.trim())
    message.success('文件夹创建成功')
    showCreateFolderDialog.value = false
    newFolderName.value = ''
  } catch (error: unknown) {
    message.error(error instanceof Error ? error.message : '创建文件夹失败')
  }
}

const loadingIndex = ref(-1)

// 处理根目录点击
const handleRootClick = async () => {
  try {
    loadingIndex.value = -1
    await fileStore.enterDirectory('/')
  } finally {
    loadingIndex.value = -1
  }
}

// 处理面包屑点击
const handleBreadcrumbClick = async (index: number) => {
  try {
    loadingIndex.value = index
    // 根据点击的索引构建路径
    const targetPath = fileStore.currentPath.slice(0, index + 1).join('/')
    await fileStore.enterDirectory('/' + targetPath)
  } finally {
    loadingIndex.value = -1
  }
}
</script>

<style scoped>
.header-wrapper {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  background: var(--n-card-color);
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  margin-bottom: 16px;
  gap: 24px;
  flex-wrap: wrap;
}

.left-section {
  flex: 1;
  min-width: 200px;
}

.nav-path {
  :deep(.n-breadcrumb) {
    font-size: 15px;
    
    .n-breadcrumb-item {
      cursor: pointer;
      transition: color 0.2s ease;
      
      &:hover {
        color: var(--n-primary-color);
      }
      
      .n-icon {
        margin-right: 4px;
      }
    }
  }
}

.right-section {
  flex: 2;
  max-width: 800px;
}

.search-and-actions {
  display: flex;
  gap: 16px;
  align-items: center;
  flex-wrap: wrap;
}

.search-box {
  flex: 1;
  min-width: 260px;
  
  :deep(.n-input) {
    .n-input__border,
    .n-input__state-border {
      box-shadow: none !important;
    }
    
    &:hover .n-input__border {
      border-color: var(--n-primary-color);
    }
  }
}

.action-buttons {
  display: flex;
  gap: 12px;
  
  .n-button {
    padding: 0 16px;
    height: 34px;
    
    .n-icon {
      margin-right: 4px;
    }
    
    &:hover {
      transform: translateY(-1px);
      transition: all 0.2s ease;
    }
  }
}

.upload-settings {
  margin-bottom: 16px;
}

.upload-tabs {
  margin-bottom: 24px;
}

.upload-dragger-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 32px;
  
  .n-icon {
    margin-bottom: 16px;
    color: var(--n-primary-color);
  }
  
  .n-text {
    font-size: 16px;
    margin-bottom: 8px;
  }
  
  .upload-hint {
    text-align: center;
    color: var(--n-text-color-3);
    font-size: 14px;
    
    p {
      margin: 4px 0;
    }
  }
}

.upload-list {
  margin-top: 24px;
  border-top: 1px solid var(--n-border-color);
  padding-top: 16px;
}

.upload-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  
  span {
    font-size: 15px;
    font-weight: 500;
  }
}

.upload-items {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.upload-item {
  background: var(--n-card-color);
  border-radius: 6px;
  padding: 12px;
  transition: all 0.3s ease;
  
  &.is-finished {
    background: rgba(var(--n-success-color), 0.1);
  }
  
  &.is-error {
    background: rgba(var(--n-error-color), 0.1);
  }
}

.upload-item-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  
  .filename {
    flex: 1;
    font-weight: 500;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  
  .filesize {
    color: var(--n-text-color-3);
    font-size: 13px;
  }
}

.upload-item-body {
  margin-bottom: 8px;
}

.progress-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
  font-size: 13px;
  color: var(--n-text-color-2);
}

.upload-item-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  
  .error-message {
    color: var(--n-error-color);
    font-size: 13px;
  }
  
  .actions {
    display: flex;
    gap: 8px;
  }
}

@media (max-width: 768px) {
  .header-wrapper {
    padding: 12px 16px;
    gap: 16px;
  }
  
  .right-section {
    flex: 1;
    max-width: none;
  }
  
  .search-and-actions {
    gap: 12px;
  }
  
  .search-box {
    min-width: 200px;
  }
  
  .action-buttons {
    .n-button {
      padding: 0 12px;
    }
  }
}
</style>

