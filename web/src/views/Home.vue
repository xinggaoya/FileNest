<template>
  <div class="home-container">
    <!-- 顶部工具栏 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <n-breadcrumb>
          <n-breadcrumb-item
            v-for="(breadcrumb, index) in fileStore.breadcrumbs"
            :key="index"
            @click="navigateTo(breadcrumb.path)"
            :class="{ clickable: breadcrumb.path !== '' }"
          >
            <n-icon v-if="index === 0"><folder-outlined /></n-icon>
            {{ breadcrumb.label }}
          </n-breadcrumb-item>
        </n-breadcrumb>
      </div>

      <div class="toolbar-right">
        <n-space>
          <!-- 搜索框 -->
          <n-input
            v-model:value="searchKeyword"
            placeholder="搜索文件..."
            style="width: 200px"
            @keyup.enter="handleSearch"
            clearable
          >
            <template #suffix>
              <n-icon @click="handleSearch" style="cursor: pointer">
                <search-outlined />
              </n-icon>
            </template>
          </n-input>

          <!-- 视图切换 -->
          <n-button-group>
            <n-button
              :type="fileStore.viewMode === 'grid' ? 'primary' : 'default'"
              @click="fileStore.toggleViewMode()"
            >
              <n-icon><appstore-outlined /></n-icon>
            </n-button>
            <n-button
              :type="fileStore.viewMode === 'list' ? 'primary' : 'default'"
              @click="fileStore.toggleViewMode()"
            >
              <n-icon><bars-outlined /></n-icon>
            </n-button>
          </n-button-group>

          <!-- 粘贴按钮 -->
          <n-button
            v-if="fileStore.hasFilesToPaste"
            @click="handlePaste"
            :disabled="fileStore.isLoading"
          >
            <n-icon><snippets-outlined /></n-icon>
            粘贴
          </n-button>

          <!-- 新建文件夹按钮 -->
          <n-button @click="showCreateFolderModal = true">
            <n-icon><folder-add-outlined /></n-icon>
            新建文件夹
          </n-button>

          <!-- 上传文件按钮 -->
          <n-button type="primary" @click="showUploadModal = true">
            <n-icon><upload-outlined /></n-icon>
            上传文件
          </n-button>
        </n-space>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content">
      <!-- 侧边栏 -->
      <div class="sidebar">
        <div class="sidebar-section">
          <h4>快速导航</h4>
          <n-menu
            :options="quickNavOptions"
            :value="activeQuickNav"
            @update:value="handleQuickNavSelect"
          />
        </div>

        <div class="sidebar-section">
          <h4>统计信息</h4>
          <div class="stats-card">
            <div class="stat-item">
              <span class="stat-label">文件数:</span>
              <span class="stat-value">{{ fileStore.stats.totalFiles }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">文件夹数:</span>
              <span class="stat-value">{{ fileStore.stats.totalFolders }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">总大小:</span>
              <span class="stat-value">{{ formatFileSize(fileStore.stats.totalSize) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 文件列表区域 -->
      <div class="file-area">
        <n-spin :show="fileStore.isLoading">
          <!-- 返回按钮 -->
          <div v-if="fileStore.canGoBack" class="back-button">
            <n-button @click="fileStore.goBack()">
              <n-icon><arrow-left-outlined /></n-icon>
              返回上级
            </n-button>
          </div>

          <!-- 文件网格视图 -->
          <div v-if="fileStore.viewMode === 'grid'" class="file-grid">
            <div
              v-for="file in fileStore.displayFiles"
              :key="file.filePath"
              class="file-card"
              :class="{ selected: fileStore.isFileSelected(file.filePath) }"
              @click="handleFileClick(file)"
              @contextmenu.prevent="handleContextMenu(file, $event)"
            >
              <div class="file-icon">
                <n-icon size="48">
                  <folder-outlined v-if="file.isDir" />
                  <file-outlined v-else />
                </n-icon>
              </div>
              <div class="file-name">{{ file.fileName }}</div>
              <div class="file-info">
                <span class="file-size">{{ formatFileSize(file.fileSize) }}</span>
                <span class="file-time">{{ formatTime(file.modTime) }}</span>
              </div>

              <!-- 选择框 -->
              <div class="file-checkbox" @click.stop="fileStore.toggleFileSelection(file.filePath)">
                <n-checkbox :checked="fileStore.isFileSelected(file.filePath)" />
              </div>
            </div>
          </div>

          <!-- 文件列表视图 -->
          <div v-else class="file-list">
            <n-data-table
              :columns="tableColumns"
              :data="fileStore.displayFiles"
              :row-key="(row) => row.filePath"
              :checked-row-keys="Array.from(fileStore.selectedFiles)"
              @update:checked-row-keys="handleTableSelection"
              :row-props="tableRowProps"
              :pagination="{
                pageSize: 20,
                showSizePicker: true,
                pageSizes: [10, 20, 30, 40]
              }"
              :bordered="false"
              :single-line="false"
              size="large"
            />
          </div>

          <!-- 空状态 -->
          <n-empty
            v-if="fileStore.displayFiles.length === 0 && !fileStore.isLoading"
            description="暂无文件"
          />
        </n-spin>
      </div>
    </div>

    <!-- 创建文件夹对话框 -->
    <n-modal v-model:show="showCreateFolderModal" preset="dialog" title="创建文件夹">
      <template #default>
        <n-form>
          <n-form-item label="文件夹名称">
            <n-input
              v-model:value="newFolderName"
              placeholder="请输入文件夹名称"
              @keyup.enter="handleCreateFolder"
            />
          </n-form-item>
        </n-form>
      </template>
      <template #action>
        <n-space>
          <n-button @click="showCreateFolderModal = false">取消</n-button>
          <n-button type="primary" @click="handleCreateFolder">创建</n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 上传文件对话框 -->
    <n-modal v-model:show="showUploadModal" style="width: 80%; max-width: 800px">
      <n-card title="上传文件" :bordered="false" size="huge">
        <template #header-extra>
          <n-button text @click="showUploadModal = false">
            <n-icon><close-outlined /></n-icon>
          </n-button>
        </template>
        <FileUpload />
      </n-card>
    </n-modal>

    <!-- 文件预览对话框 -->
    <n-modal v-model:show="showPreviewModal" style="width: 80%; max-width: 1000px">
      <n-card :title="previewFileName" :bordered="false" size="huge">
        <template #header-extra>
          <n-button text @click="showPreviewModal = false">
            <n-icon><close-outlined /></n-icon>
          </n-button>
        </template>
        <div class="preview-content">
          <!-- 图片预览 -->
          <img v-if="isImage(previewFileName)" :src="previewUrl" class="preview-image" />
          <!-- 文本预览 -->
          <div v-else-if="isPreviewableText(previewFileName)" class="preview-text">
            <pre>{{ previewContent }}</pre>
          </div>
          <!-- 其他文件类型 -->
          <div v-else class="preview-unsupported">
            <n-empty description="该文件类型暂不支持预览">
              <template #extra>
                <n-button type="primary" @click="downloadFile"> 下载文件 </n-button>
              </template>
            </n-empty>
          </div>
        </div>
      </n-card>
    </n-modal>

    <!-- 重命名对话框 -->
    <n-modal v-model:show="showRenameModal" preset="dialog" title="重命名">
      <template #default>
        <n-form>
          <n-form-item label="新名称">
            <n-input
              v-model:value="newFileName"
              placeholder="请输入新名称"
              @keyup.enter="handleRename"
            />
          </n-form-item>
        </n-form>
      </template>
      <template #action>
        <n-space>
          <n-button @click="showRenameModal = false">取消</n-button>
          <n-button type="primary" @click="handleRename">确定</n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 右键菜单 -->
    <n-dropdown
      placement="bottom-start"
      trigger="manual"
      :x="contextMenuX"
      :y="contextMenuY"
      :options="contextMenuOptions"
      :show="showContextMenu"
      @clickoutside="showContextMenu = false"
      @select="handleContextMenuSelect"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, h, nextTick } from 'vue'

defineOptions({
  name: 'HomePage'
})
import { useFileStore } from '@/stores/file'
import { useMessageHandler } from '@/composables/useMessageHandler'
import FileUpload from '@/components/FileUpload.vue'
import type { FileInfo } from '@/types/file'
import type { DataTableColumns } from 'naive-ui'
import {
  FolderOutlined,
  FileOutlined,
  SearchOutlined,
  AppstoreOutlined,
  BarsOutlined,
  FolderAddOutlined,
  UploadOutlined,
  ArrowLeftOutlined,
  CloseOutlined,
  DeleteOutlined,
  EditOutlined,
  CopyOutlined,
  ScissorOutlined,
  DownloadOutlined,
  StarOutlined,
  SnippetsOutlined
} from '@vicons/antd'

const { handleResult } = useMessageHandler()
const fileStore = useFileStore()

// 响应式数据
const searchKeyword = ref('')
const showCreateFolderModal = ref(false)
const showUploadModal = ref(false)
const newFolderName = ref('')
const activeQuickNav = ref('all')

// 右键菜单相关
const showContextMenu = ref(false)
const contextMenuX = ref(0)
const contextMenuY = ref(0)
const contextMenuFile = ref<FileInfo | null>(null)

// 预览相关
const showPreviewModal = ref(false)
const previewFileName = ref('')
const previewUrl = ref('')
const previewContent = ref('')
const currentPreviewFile = ref<FileInfo | null>(null)

// 重命名相关
const showRenameModal = ref(false)
const newFileName = ref('')
const renameFile = ref<FileInfo | null>(null)

// 快速导航选项
const quickNavOptions = computed(() => [
  {
    label: '全部文件',
    key: 'all',
    icon: () => h(FolderOutlined)
  },
  {
    label: '收藏夹',
    key: 'favorites',
    icon: () => h(StarOutlined)
  }
])

// 表格列定义
const tableColumns: DataTableColumns<FileInfo> = [
  {
    type: 'selection',
    width: 50
  },
  {
    title: '名称',
    key: 'fileName',
    width: 300,
    render(row) {
      return h(
        'div',
        {
          style: 'display: flex; align-items: center; gap: 8px; cursor: pointer;'
        },
        [
          h(row.isDir ? FolderOutlined : FileOutlined, {
            style: {
              color: row.isDir ? '#1890ff' : '#666',
              width: '20px',
              height: '20px',
              minWidth: '20px',
              minHeight: '20px',
              flexShrink: 0
            },
            width: '20',
            height: '20'
          }),
          h(
            'span',
            {
              style: {
                color: '#333',
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap',
                fontSize: '15px'
              }
            },
            row.fileName
          )
        ]
      )
    }
  },
  {
    title: '大小',
    key: 'fileSize',
    width: 120,
    align: 'right',
    render(row) {
      return row.isDir ? '-' : formatFileSize(row.fileSize)
    }
  },
  {
    title: '类型',
    key: 'fileType',
    width: 100,
    render(row) {
      if (row.isDir) return '文件夹'
      const ext = row.fileType.toLowerCase()
      if (isImage(row.fileName)) return '图片'
      if (isPreviewableText(row.fileName)) return '文本'
      return ext ? ext.substring(1) : '未知'
    }
  },
  {
    title: '修改时间',
    key: 'modTime',
    width: 180,
    render(row) {
      return formatTime(row.modTime)
    }
  }
]

const menuIconStyle = { fontSize: '18px', marginRight: '4px', verticalAlign: 'middle' }

const contextMenuOptions = computed(() => {
  if (!contextMenuFile.value) return []

  const options = [
    {
      label: '重命名',
      key: 'rename',
      icon: () => h(EditOutlined, { style: menuIconStyle })
    },
    {
      label: '复制',
      key: 'copy',
      icon: () => h(CopyOutlined, { style: menuIconStyle })
    },
    {
      label: '剪切',
      key: 'cut',
      icon: () => h(ScissorOutlined, { style: menuIconStyle })
    },
    {
      type: 'divider' as const
    },
    {
      label: '删除',
      key: 'delete',
      icon: () => h(DeleteOutlined, { style: menuIconStyle })
    }
  ]

  if (!contextMenuFile.value.isDir) {
    options.splice(3, 0, {
      label: '下载',
      key: 'download',
      icon: () => h(DownloadOutlined, { style: menuIconStyle })
    })
  }

  return options
})

// 判断是否为可预览的文本文件
const isPreviewableText = (fileName: string): boolean => {
  const textExtensions = ['txt', 'md', 'json', 'js', 'ts', 'html', 'css', 'xml', 'yaml', 'yml']
  return textExtensions.includes(fileStore.getFileType(fileName))
}

// 判断是否为图片文件
const isImage = (fileName: string): boolean => {
  const imageExtensions = ['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg']
  return imageExtensions.includes(fileStore.getFileType(fileName))
}

// 预览文件
const previewFile = async (file: FileInfo) => {
  currentPreviewFile.value = file
  previewFileName.value = file.fileName
  showPreviewModal.value = true

  try {
    if (isImage(file.fileName)) {
      // 图片文件直接使用预览URL
      previewUrl.value = `/api/file/preview?path=${encodeURIComponent(file.filePath)}`
    } else if (isPreviewableText(file.fileName)) {
      // 文本文件获取内容
      const response = await fetch(`/api/file/preview?path=${encodeURIComponent(file.filePath)}`)
      if (!response.ok) throw new Error('获取文件内容失败')
      previewContent.value = await response.text()
    }
  } catch (error) {
    handleResult({
      success: false,
      error: error instanceof Error ? error.message : '预览文件失败'
    })
    showPreviewModal.value = false
  }
}

// 下载文件
const downloadFile = async (file?: FileInfo) => {
  let target = file || currentPreviewFile.value
  if (!target) return
  try {
    await fileStore.downloadFile(target.filePath, target.fileName)
  } catch (error) {
    handleResult({
      success: false,
      error: error instanceof Error ? error.message : '下载失败'
    })
  }
}

// 初始化
onMounted(async () => {
  const result = await fileStore.fetchFiles()
  if (result && !result.success) {
    handleResult(result)
  }

  // 获取收藏列表
  const favResult = await fileStore.fetchFavorites()
  if (favResult && !favResult.success) {
    handleResult(favResult)
  }
})

// 处理文件点击
const handleFileClick = async (file: FileInfo) => {
  if (file.isDir) {
    await fileStore.enterDirectory(file.filePath)
  } else {
    await previewFile(file)
  }
}

// 处理搜索
const handleSearch = async () => {
  if (searchKeyword.value.trim()) {
    const result = await fileStore.searchFile(searchKeyword.value.trim())
    if (result && !result.success) {
      handleResult(result)
    }
  } else {
    fileStore.exitSearchMode()
  }
}

// 导航到指定路径
const navigateTo = async (path: string) => {
  if (path !== '') {
    await fileStore.navigateTo(path)
  }
}

// 处理快速导航选择
const handleQuickNavSelect = (key: string) => {
  activeQuickNav.value = key
  // 根据选择执行相应操作
  if (key === 'favorites') {
    // 显示收藏夹
  }
}

// 处理表格选择
const handleTableSelection = (keys: string[]) => {
  fileStore.clearSelection()
  keys.forEach((key) => fileStore.selectFile(key))
}

// 创建文件夹
const handleCreateFolder = async () => {
  if (!newFolderName.value.trim()) return

  const result = await fileStore.createNewFolder(newFolderName.value.trim())
  handleResult(result)

  if (result.success) {
    showCreateFolderModal.value = false
    newFolderName.value = ''
  }
}

// 处理重命名
const handleRename = async () => {
  if (!renameFile.value || !newFileName.value.trim()) return

  const result = await fileStore.renameItem(renameFile.value.filePath, newFileName.value.trim())
  handleResult(result)

  if (result.success) {
    showRenameModal.value = false
    newFileName.value = ''
    renameFile.value = null
  }
}

// 处理右键菜单
const handleContextMenu = (file: FileInfo, event: MouseEvent) => {
  event.preventDefault()
  contextMenuFile.value = file
  contextMenuX.value = event.clientX
  contextMenuY.value = event.clientY
  showContextMenu.value = true
}

// 处理右键菜单选择
const handleContextMenuSelect = async (key: string) => {
  showContextMenu.value = false
  if (!contextMenuFile.value) return
  switch (key) {
    case 'delete': {
      const deleteResult = await fileStore.deleteFiles([contextMenuFile.value.filePath])
      handleResult(deleteResult)
      break
    }
    case 'rename': {
      renameFile.value = contextMenuFile.value
      newFileName.value = contextMenuFile.value.fileName
      showRenameModal.value = true
      break
    }
    case 'copy': {
      fileStore.selectFile(contextMenuFile.value.filePath)
      fileStore.copySelectedFiles()
      break
    }
    case 'cut': {
      fileStore.selectFile(contextMenuFile.value.filePath)
      fileStore.cutSelectedFiles()
      break
    }
    case 'download': {
      if (!contextMenuFile.value.isDir) {
        await downloadFile(contextMenuFile.value)
      }
      break
    }
  }
  contextMenuFile.value = null
}

// 处理粘贴
const handlePaste = async () => {
  const result = await fileStore.pasteFiles()
  handleResult(result)
}

// 格式化文件大小
const formatFileSize = (size: number): string => {
  if (size === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const index = Math.floor(Math.log(size) / Math.log(1024))
  return Math.round((size / Math.pow(1024, index)) * 100) / 100 + ' ' + units[index]
}

// 格式化时间
const formatTime = (timeStr: string): string => {
  const date = new Date(timeStr)
  return date.toLocaleString('zh-CN')
}

// 表格行右键菜单 rowProps 方案
const tableRowProps = (row: FileInfo) => ({
  onContextmenu: (e: MouseEvent) => {
    e.preventDefault()
    showContextMenu.value = false
    nextTick().then(() => {
      contextMenuFile.value = row
      contextMenuX.value = e.clientX
      contextMenuY.value = e.clientY
      showContextMenu.value = true
    })
  },
  style: { cursor: 'pointer' }
})
</script>

<style scoped>
.home-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #e8e8e8;
  background-color: white;
}

.toolbar-left .clickable {
  cursor: pointer;
}

.toolbar-left .clickable:hover {
  color: #1890ff;
}

.main-content {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.sidebar {
  width: 250px;
  background-color: #fafafa;
  border-right: 1px solid #e8e8e8;
  padding: 16px;
  overflow-y: auto;
}

.sidebar-section {
  margin-bottom: 24px;
}

.sidebar-section h4 {
  margin: 0 0 12px 0;
  color: #333;
  font-size: 14px;
}

.stats-card {
  background-color: white;
  border-radius: 6px;
  padding: 12px;
  border: 1px solid #e8e8e8;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.stat-item:last-child {
  margin-bottom: 0;
}

.stat-label {
  color: #666;
  font-size: 13px;
}

.stat-value {
  color: #333;
  font-weight: 500;
  font-size: 13px;
}

.file-area {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
}

.back-button {
  margin-bottom: 16px;
}

.file-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 16px;
}

.file-card {
  position: relative;
  background-color: white;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  padding: 16px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.file-card:hover {
  border-color: #1890ff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.file-card.selected {
  border-color: #1890ff;
  background-color: #f0f8ff;
}

.file-icon {
  color: #1890ff;
  margin-bottom: 8px;
}

.file-name {
  font-weight: 500;
  margin-bottom: 8px;
  word-break: break-all;
  font-size: 14px;
}

.file-info {
  font-size: 12px;
  color: #666;
}

.file-size {
  display: block;
  margin-bottom: 4px;
}

.file-time {
  display: block;
}

.file-checkbox {
  position: absolute;
  top: 8px;
  right: 8px;
}

.file-list {
  background-color: white;
  border-radius: 8px;
  padding: 16px;
  height: 100%;
  overflow: auto;
}

.file-list :deep(.n-data-table-td) {
  padding: 12px 16px;
}

.file-list :deep(.n-data-table-th) {
  padding: 12px 16px;
  font-weight: 500;
}

.file-list :deep(.n-data-table-tr:hover) {
  background-color: #f5f5f5;
}

.file-list :deep(.n-data-table-tr.selected) {
  background-color: #e6f7ff;
}

.preview-content {
  min-height: 400px;
  max-height: 70vh;
  overflow: auto;
}

.preview-image {
  max-width: 100%;
  max-height: 70vh;
  object-fit: contain;
}

.preview-text {
  background-color: #f5f5f5;
  padding: 16px;
  border-radius: 4px;
  font-family: monospace;
  white-space: pre-wrap;
  word-break: break-all;
}

.preview-unsupported {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 400px;
}

@media (max-width: 768px) {
  .main-content {
    flex-direction: column;
  }

  .sidebar {
    width: 100%;
    height: auto;
    border-right: none;
    border-bottom: 1px solid #e8e8e8;
  }

  .file-grid {
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 12px;
  }

  .toolbar {
    flex-direction: column;
    gap: 12px;
  }

  .toolbar-left,
  .toolbar-right {
    width: 100%;
  }
}

:deep(.n-dropdown-option) {
  min-height: 36px !important;
  padding: 6px 16px !important;
  font-size: 15px;
}
</style>
