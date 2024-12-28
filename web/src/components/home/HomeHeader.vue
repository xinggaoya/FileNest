<template>
  <div class="header-wrapper">
    <div class="left-section">
      <div class="nav-path">
        <n-breadcrumb>
          <n-breadcrumb-item @click="fileStore.enterDirectory('/')">
            <n-icon><folder-outlined /></n-icon>
            根目录
          </n-breadcrumb-item>
          <template v-for="(path, index) in fileStore.currentPath" :key="index">
            <n-icon><right-outlined /></n-icon>
            <n-breadcrumb-item>{{ path }}</n-breadcrumb-item>
          </template>
        </n-breadcrumb>
      </div>
    </div>

    <!-- 右侧搜索框和操作按钮 -->
    <div class="right-section">
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
        <n-space>
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
        </n-space>
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
      <n-space vertical>
        <n-upload
          :custom-request="handleUpload"
          :show-file-list="false"
          :directory="false"
          @change="handleUploadChange"
        >
          <n-upload-dragger>
            <div class="upload-dragger-content">
              <n-icon size="48" depth="3">
                <upload-outlined />
              </n-icon>
              <n-text>点击或拖拽文件到此处上传</n-text>
            </div>
          </n-upload-dragger>
        </n-upload>
        <n-divider>或者</n-divider>
        <n-upload
          :custom-request="handleFolderUpload"
          :show-file-list="false"
          directory
          @change="handleUploadChange"
        >
          <n-upload-dragger>
            <div class="upload-dragger-content">
              <n-icon size="48" depth="3">
                <folder-add-outlined />
              </n-icon>
              <n-text>点击或拖拽文件夹到此处上传</n-text>
            </div>
          </n-upload-dragger>
        </n-upload>

        <!-- 上传进度列表 -->
        <div class="upload-list" v-if="Object.keys(uploadStatus).length > 0">
          <div class="upload-item" v-for="(status, key) in uploadStatus" :key="key">
            <div class="upload-info">
              <span class="filename">{{ key }}</span>
              <span class="status">
                {{ status.status === 'uploading' ? '上传中...' :
                   status.status === 'finished' ? '上传完成' :
                   '上传失败' }}
              </span>
            </div>
            <n-progress
              type="line"
              :percentage="status.progress"
              :status="status.status === 'error' ? 'error' :
                      status.status === 'finished' ? 'success' : 'info'"
              :show-indicator="true"
            />
            <div class="error-message" v-if="status.message">{{ status.message }}</div>
          </div>
        </div>
      </n-space>
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
  SearchOutlined
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
  status: 'uploading' | 'finished' | 'error'
  message?: string
}

const uploadStatus = reactive<Record<string, UploadStatus>>({})

const handleUpload = async ({ file, onFinish, onError, onProgress }: UploadCustomRequestOptions) => {
  if (!file) return
  
  const fileName = file.name
  uploadStatus[fileName] = {
    progress: 0,
    status: 'uploading'
  }

  try {
    await uploadFile({
      file: file.file as File,
      path: fileStore.currentPathString,
      onProgress: (progress) => {
        uploadStatus[fileName].progress = progress
        onProgress?.({ percent: progress })
      },
      onSuccess: () => {
        uploadStatus[fileName].status = 'finished'
        message.success(`${fileName} 上传成功`)
        fileStore.fetchFiles()
        onFinish()
      },
      onError: (error: unknown) => {
        uploadStatus[fileName].status = 'error'
        uploadStatus[fileName].message = error instanceof Error ? error.message : '上传失败'
        message.error(`${fileName} 上传失败`)
        onError()
      }
    })
  } catch (error: unknown) {
    uploadStatus[fileName].status = 'error'
    uploadStatus[fileName].message = error instanceof Error ? error.message : '上传失败，请重试'
    message.error(`${fileName} 上传失败`)
    onError()
  }
}

const handleFolderUpload = async (options: UploadCustomRequestOptions) => {
  const { file, onFinish, onError, onProgress } = options
  if (!file) return

  const fileName = file.name
  const relativePath = (file as any).webkitRelativePath || ''
  const path = relativePath.split('/').slice(0, -1).join('/')
  const fullPath = fileStore.currentPathString
    ? `${fileStore.currentPathString}/${path}`
    : path
  
  uploadStatus[fileName] = {
    progress: 0,
    status: 'uploading'
  }

  try {
    await uploadFile({
      file: file.file as File,
      path: fullPath,
      onProgress: (progress) => {
        uploadStatus[fileName].progress = progress
        onProgress?.({ percent: progress })
      },
      onSuccess: () => {
        uploadStatus[fileName].status = 'finished'
        message.success(`${fileName} 上传成功`)
        fileStore.fetchFiles()
        onFinish()
      },
      onError: (error: unknown) => {
        uploadStatus[fileName].status = 'error'
        uploadStatus[fileName].message = error instanceof Error ? error.message : '上传失败'
        message.error(`${fileName} 上传失败`)
        onError()
      }
    })
  } catch (error: unknown) {
    uploadStatus[fileName].status = 'error'
    uploadStatus[fileName].message = error instanceof Error ? error.message : '上传失败，请重试'
    message.error(`${fileName} 上传失败`)
    onError()
  }
}

const handleUploadChange = (data: { fileList: UploadFileInfo[] }) => {
  console.log('Upload change:', data)
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
</script>

<style scoped>
.header-wrapper {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  gap: 20px;
  width: 100%;
  background-color: var(--n-card-color);
  border-radius: 8px;
  margin-bottom: 16px;
}

.left-section {
  flex: 0 0 auto;
  min-width: 200px;
  max-width: 30%;
}

.nav-path {
  width: 100%;
  overflow: hidden;
}

.right-section {
  flex: 1 1 auto;
  display: flex;
  align-items: center;
  gap: 16px;
  justify-content: flex-end;
  min-width: 0;
}

.search-box {
  flex: 1 1 auto;
  max-width: 400px;
  min-width: 200px;
  margin-right: 16px;
}

.action-buttons {
  flex: 0 0 auto;
  white-space: nowrap;
}

.upload-dragger-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 20px;
}

.upload-list {
  margin-top: 16px;
}

.upload-item {
  margin-bottom: 12px;
}

.upload-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
}

.error-message {
  color: var(--error-color);
  font-size: 12px;
  margin-top: 4px;
}
</style>

