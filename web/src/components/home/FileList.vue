<template>
  <div class="file-list">
    <!-- 视图切换按钮 -->
    <div class="view-toggle">
      <n-button-group>
        <n-button
          :type="viewMode === 'list' ? 'primary' : 'default'"
          @click="viewMode = 'list'"
        >
          <template #icon>
            <n-icon><unordered-list-outlined /></n-icon>
          </template>
          列表
        </n-button>
        <n-button
          :type="viewMode === 'grid' ? 'primary' : 'default'"
          @click="viewMode = 'grid'"
        >
          <template #icon>
            <n-icon><appstore-outlined /></n-icon>
          </template>
          网格
        </n-button>
      </n-button-group>
    </div>

    <!-- 文件列表 -->
    <n-scrollbar>
      <transition name="mode-transition" mode="out-in">
        <div :class="[viewMode === 'list' ? 'list-mode' : 'grid-mode']">
          <template v-if="files.length > 0">
            <div
              v-for="file in files"
              :key="file.filePath"
              class="file-item"
              @dblclick="handleFileAction(file)"
            >
              <!-- 文件图标 -->
              <div class="file-icon">
                <n-icon :size="viewMode === 'grid' ? 32 : 24">
                  <template v-if="file.isDir">
                    <folder-outlined />
                  </template>
                  <template v-else>
                    <file-outlined />
                  </template>
                </n-icon>
              </div>

              <!-- 文件名称 -->
              <div class="file-name" :title="file.fileName">
                {{ file.fileName }}
              </div>

              <!-- 文件大小 -->
              <div v-if="!file.isDir" class="file-size">
                {{ formatFileSize(file.fileSize) }}
              </div>

              <!-- 修改时间 -->
              <div class="file-time">
                {{ formatTime(file.modTime) }}
              </div>

              <!-- 操作按钮 -->
              <div class="file-actions">
                <n-button
                  v-if="file.isDir"
                  size="small"
                  @click.stop="handleEnterDirectory(file)"
                >
                  <template #icon>
                    <n-icon><folder-open-outlined /></n-icon>
                  </template>
                  打开
                </n-button>
                <n-button
                  v-else
                  size="small"
                  @click.stop="handleDownload(file)"
                >
                  <template #icon>
                    <n-icon><download-outlined /></n-icon>
                  </template>
                  下载
                </n-button>
                <n-button
                  size="small"
                  @click.stop="handleDelete(file)"
                >
                  <template #icon>
                    <n-icon><delete-outlined /></n-icon>
                  </template>
                  删除
                </n-button>
              </div>
            </div>
          </template>
          <template v-else>
            <n-empty description="暂无文件" />
          </template>
        </div>
      </transition>
    </n-scrollbar>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useFileStore } from '@/stores/file'
import { formatFileSize, formatTime } from '@/utils/format'
import {
  FolderOutlined,
  FileOutlined,
  FolderOpenOutlined,
  DownloadOutlined,
  DeleteOutlined,
  UnorderedListOutlined,
  AppstoreOutlined
} from '@vicons/antd'
import { createDiscreteApi } from 'naive-ui'

const { dialog } = createDiscreteApi(['dialog'])
const fileStore = useFileStore()
const viewMode = ref('list')

const handleFileAction = (file: any) => {
  if (file.isDir) {
    handleEnterDirectory(file)
  } else {
    handleDownload(file)
  }
}

const handleEnterDirectory = (file: any) => {
  fileStore.enterDirectory(file.filePath)
}

const handleDownload = (file: any) => {
  fileStore.downloadFile(file.filePath)
}

const handleDelete = (file: any) => {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除 ${file.fileName} 吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => {
      fileStore.deleteFile(file.filePath)
    }
  })
}

defineProps<{
  files: any[]
}>()
</script>

<style scoped>
.file-list {
  margin-top: 16px;
}

.view-toggle {
  margin-bottom: 16px;
  display: flex;
  justify-content: flex-end;
}

/* 列表模式 */
.list-mode {
  .file-item {
    display: grid;
    grid-template-columns: auto minmax(200px, 1fr) auto auto auto;
    gap: 16px;
    align-items: center;
    padding: 12px 16px;
    border-radius: 8px;
    transition: background-color 0.3s;
    cursor: pointer;

    &:hover {
      background-color: var(--n-hover-color);
    }
  }

  .file-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
  }

  .file-name {
    font-weight: 500;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    min-width: 0;
  }

  .file-size {
    color: var(--n-text-color-3);
    font-size: 14px;
    white-space: nowrap;
    text-align: right;
    min-width: 80px;
  }

  .file-time {
    color: var(--n-text-color-3);
    font-size: 14px;
    white-space: nowrap;
    text-align: right;
    min-width: 150px; /* 确保时间列宽度一致 */
  }

  .file-actions {
    display: flex;
    gap: 8px;
    opacity: 0;
    transition: opacity 0.3s;
    white-space: nowrap;
  }

  .file-item:hover .file-actions {
    opacity: 1;
  }
}

/* 网格模式 */
.grid-mode {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
  padding: 16px;

  .file-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 16px;
    border-radius: 8px;
    transition: background-color 0.3s;
    cursor: pointer;
    text-align: center;

    &:hover {
      background-color: var(--n-hover-color);
    }
  }

  .file-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 48px;
    height: 48px;
    margin-bottom: 8px;
  }

  .file-name {
    font-weight: 500;
    margin-bottom: 4px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    width: 100%;
  }

  .file-size,
  .file-time {
    color: var(--n-text-color-3);
    font-size: 12px;
    margin-bottom: 4px;
  }

  .file-actions {
    display: flex;
    gap: 8px;
    margin-top: 8px;
    opacity: 0;
    transition: opacity 0.3s;
  }

  .file-item:hover .file-actions {
    opacity: 1;
  }
}

/* 确保模式切换时的平滑过渡 */
.mode-transition-enter-active,
.mode-transition-leave-active {
  transition: opacity 0.3s ease;
}

.mode-transition-enter-from,
.mode-transition-leave-to {
  opacity: 0;
}
</style> 