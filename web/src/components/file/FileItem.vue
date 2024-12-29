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
            <n-button quaternary circle size="small" @click.stop="handleRename">
              <template #icon><edit-outlined /></template>
            </n-button>
            <n-button quaternary circle size="small" @click.stop="handleCopy">
              <template #icon><copy-outlined /></template>
            </n-button>
            <n-button quaternary circle size="small" @click.stop="handleMove">
              <template #icon><scissor-outlined /></template>
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
          <n-button quaternary circle size="small" @click.stop="handleRename">
            <template #icon><edit-outlined /></template>
          </n-button>
          <n-button quaternary circle size="small" @click.stop="handleCopy">
            <template #icon><copy-outlined /></template>
          </n-button>
          <n-button quaternary circle size="small" @click.stop="handleMove">
            <template #icon><scissor-outlined /></template>
          </n-button>
          <n-button quaternary circle size="small" @click.stop="handleDelete">
            <template #icon><delete-outlined /></template>
          </n-button>
        </n-button-group>
      </div>
    </template>
  </div>

  <!-- 预览弹窗 -->
  <n-modal v-model:show="showPreview" preset="card" style="max-width: 90vw; max-height: 90vh" :title="file.fileName">
    <div class="preview-content">
      <img v-if="isImage" :src="previewUrl" :alt="file.fileName"
        style="max-width: 100%; max-height: 80vh; object-fit: contain;" />
      <pre v-else-if="isText" v-text="previewContent" />
      <div v-else class="preview-unsupported">
        暂不支持预览该类型文件
      </div>
    </div>
  </n-modal>

  <!-- 重命名对话框 -->
  <n-modal v-model:show="showRenameModal" preset="card" title="重命名" style="width: 400px">
    <n-input v-model:value="newFileName" placeholder="请输入新名称" @keyup.enter="confirmRename" />
    <div style="display: flex; justify-content: flex-end; margin-top: 16px;">
      <n-button-group>
        <n-button @click="showRenameModal = false">取消</n-button>
        <n-button type="primary" @click="confirmRename" :disabled="!newFileName">确定</n-button>
      </n-button-group>
    </div>
  </n-modal>

  <!-- 复制对话框 -->
  <n-modal v-model:show="showCopyModal" preset="card" title="复制到" style="width: 500px">
    <div class="folder-select-modal">
      <div class="current-path">
        <n-breadcrumb>
          <n-breadcrumb-item @click="selectTargetPath('')">根目录</n-breadcrumb-item>
          <n-breadcrumb-item v-for="(folder, index) in targetPathArray" :key="index"
            @click="selectTargetPath(targetPathArray.slice(0, index + 1).join('/'))">
            {{ folder }}
          </n-breadcrumb-item>
        </n-breadcrumb>
      </div>
      <div class="folder-list">
        <n-scrollbar style="max-height: 300px">
          <div class="folder-item" @click="selectTargetPath('')">
            <n-icon><folder-outlined /></n-icon>
            <span>根目录</span>
          </div>
          <div v-for="folder in targetFolders" :key="folder.fileName" class="folder-item"
            @click="selectTargetPath(folder.filePath)">
            <n-icon><folder-outlined /></n-icon>
            <span>{{ folder.fileName }}</span>
          </div>
        </n-scrollbar>
      </div>
      <div class="selected-path">
        <span class="path-label">目标路径：</span>
        <span class="path-value">{{ targetPath || '根目录' }}</span>
      </div>
    </div>
    <div style="display: flex; justify-content: flex-end; margin-top: 16px;">
      <n-button-group>
        <n-button @click="showCopyModal = false">取消</n-button>
        <n-button type="primary" @click="confirmCopy">确定</n-button>
      </n-button-group>
    </div>
  </n-modal>

  <!-- 移动对话框 -->
  <n-modal v-model:show="showMoveModal" preset="card" title="移动到" style="width: 500px">
    <div class="folder-select-modal">
      <div class="current-path">
        <n-breadcrumb>
          <n-breadcrumb-item @click="selectTargetPath('')">根目录</n-breadcrumb-item>
          <n-breadcrumb-item v-for="(folder, index) in targetPathArray" :key="index"
            @click="selectTargetPath(targetPathArray.slice(0, index + 1).join('/'))">
            {{ folder }}
          </n-breadcrumb-item>
        </n-breadcrumb>
      </div>
      <div class="folder-list">
        <n-scrollbar style="max-height: 300px">
          <div class="folder-item" @click="selectTargetPath('')">
            <n-icon><folder-outlined /></n-icon>
            <span>根目录</span>
          </div>
          <div v-for="folder in targetFolders" :key="folder.fileName" class="folder-item"
            @click="selectTargetPath(folder.filePath)">
            <n-icon><folder-outlined /></n-icon>
            <span>{{ folder.fileName }}</span>
          </div>
        </n-scrollbar>
      </div>
      <div class="selected-path">
        <span class="path-label">目标路径：</span>
        <span class="path-value">{{ targetPath || '根目录' }}</span>
      </div>
    </div>
    <div style="display: flex; justify-content: flex-end; margin-top: 16px;">
      <n-button-group>
        <n-button @click="showMoveModal = false">取消</n-button>
        <n-button type="primary" @click="confirmMove">确定</n-button>
      </n-button-group>
    </div>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { FileInfo } from '@/types/file'
import { getDownloadUrl, getFileList } from '@/api/file/file'
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
  EditOutlined,
  CopyOutlined,
  ScissorOutlined
} from '@vicons/antd'
import { createDiscreteApi } from 'naive-ui'
import { useDialog, NCheckbox } from 'naive-ui'
import { h } from 'vue'
import { NBreadcrumb, NBreadcrumbItem, NScrollbar } from 'naive-ui'

const { message, dialog } = createDiscreteApi(['message', 'dialog'])
const fileStore = useFileStore()

const props = defineProps<{
  file: FileInfo
}>()

const isHovered = ref(false)
const showPreview = ref(false)
const previewContent = ref('')
const forceDelete = ref(false)
const showRenameModal = ref(false)
const showCopyModal = ref(false)
const showMoveModal = ref(false)
const newFileName = ref('')
const targetPath = ref('')
const targetFolders = ref<FileInfo[]>([])
const targetPathArray = computed(() => targetPath.value ? targetPath.value.split('/') : [])

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
  fileStore.downloadFile(props.file.filePath)
}

// 处理删除
const handleDelete = async () => {
  if (props.file.isDir) {
    dialog.warning({
      title: '确认删除',
      content: () => {
        return h('div', [
          h('p', `确定要删除文件夹 ${props.file.fileName} 吗？`),
          h(NCheckbox, {
            checked: forceDelete.value,
            onUpdateChecked: (checked) => forceDelete.value = checked
          }, { default: () => '强制删除（删除文件夹及其所有内容）' })
        ])
      },
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: async () => {
        await fileStore.removeFile(props.file.filePath, forceDelete.value)
      }
    })
  } else {
    dialog.warning({
      title: '确认删除',
      content: `确定要删除文件 ${props.file.fileName} 吗？`,
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: async () => {
        await fileStore.removeFile(props.file.filePath)
      }
    })
  }
}

// 处理重命名
const handleRename = () => {
  newFileName.value = props.file.fileName
  showRenameModal.value = true
}

// 选择目标路径
const selectTargetPath = async (path: string) => {
  targetPath.value = path
  try {
    const { data } = await getFileList(path)
    targetFolders.value = data.filter((item: FileInfo) => item.isDir)
  } catch (error) {
    message.error('获取文件夹列表失败')
  }
}

// 处理复制
const handleCopy = async () => {
  targetPath.value = ''
  await selectTargetPath('')
  showCopyModal.value = true
}

// 处理移动
const handleMove = async () => {
  targetPath.value = ''
  await selectTargetPath('')
  showMoveModal.value = true
}

// 确认重命名
const confirmRename = async () => {
  if (!newFileName.value || newFileName.value === props.file.fileName) {
    showRenameModal.value = false
    return
  }
  try {
    await fileStore.renameFile(props.file.filePath, newFileName.value)
    showRenameModal.value = false
  } catch (error) {
    // 错误已在 store 中处理
  }
}

// 确认复制
const confirmCopy = async () => {
  if (!targetPath.value) {
    targetPath.value = '/'
  }
  try {
    await fileStore.copyFile(props.file.filePath, targetPath.value)
    showCopyModal.value = false
  } catch (error) {
    // 错误已在 store 中处理
  }
}

// 确认移动
const confirmMove = async () => {
  if (!targetPath.value) {
    targetPath.value = '/'
  }
  try {
    await fileStore.moveFile(props.file.filePath, targetPath.value)
    showMoveModal.value = false
  } catch (error) {
    // 错误已在 store 中处理
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

.file-size,
.file-date {
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

.modal-info {
  margin: 8px 0;
  color: var(--n-text-color-2);
  font-size: 14px;
}

.path-input {
  margin-top: 16px;
}

.input-label {
  margin-bottom: 8px;
  font-size: 14px;
  color: var(--n-text-color);
}

.modal-input {
  margin: 8px 0;
}

.path-tip {
  margin-top: 4px;
  font-size: 12px;
  color: var(--n-text-color-3);
}

.rename-modal,
.copy-modal,
.move-modal {
  padding: 8px 0;
}

.folder-select-modal {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.current-path {
  padding: 8px;
  background-color: var(--n-color-hover);
  border-radius: 4px;
}

.folder-list {
  border: 1px solid var(--n-border-color);
  border-radius: 4px;
  padding: 8px;
}

.folder-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.3s;
}

.folder-item:hover {
  background-color: var(--n-color-hover);
}

.selected-path {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  background-color: var(--n-color-hover);
  border-radius: 4px;
}

.path-label {
  color: var(--n-text-color-3);
}

.path-value {
  color: var(--n-text-color);
  font-weight: 500;
}
</style>