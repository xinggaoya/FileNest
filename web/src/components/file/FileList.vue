<template>
  <div class="file-list">
    <n-data-table
      :columns="columns"
      :data="fileStore.files"
      :pagination="false"
      :bordered="false"
    />
  </div>
</template>

<script setup lang="ts">
import { h } from 'vue'
import { NButton, NIcon, NSpace, NPopconfirm } from 'naive-ui'
import { useFileStore } from '@/stores/file'
import { DownloadOutlined, DeleteOutlined, StarOutlined, StarFilled } from '@vicons/antd'
import { getDownloadUrl, deleteFile } from '@/api/file/file'
import type { FileInfo } from '@/types/file'
import { isAxiosError } from 'axios'
import { createDiscreteApi } from 'naive-ui'

const { message } = createDiscreteApi(['message'])
const fileStore = useFileStore()

// 获取收藏列表
fileStore.fetchFavorites()

// 检查是否已收藏
const checkFavorite = (file: FileInfo): boolean => {
  return fileStore.favorites.some(f => f.path === file.filePath)
}

const columns = [
  {
    title: '文件名',
    key: 'fileName',
    render(row: FileInfo) {
      return h(
        'div',
        {
          class: 'file-name',
          onClick: () => handleFileClick(row)
        },
        row.fileName
      )
    }
  },
  {
    title: '大小',
    key: 'fileSize',
    render(row: FileInfo) {
      return row.isDir ? '-' : fileStore.formatFileSize(row.fileSize)
    }
  },
  {
    title: '修改时间',
    key: 'modTime'
  },
  {
    title: '操作',
    key: 'actions',
    render(row: FileInfo) {
      const isFavorite = checkFavorite(row)
      return h(NSpace, null, {
        default: () => [
          h(
            NButton,
            {
              text: true,
              onClick: () => handleFavoriteClick(row)
            },
            {
              default: () =>
                h(NIcon, null, {
                  default: () => h(isFavorite ? StarFilled : StarOutlined)
                })
            }
          ),
          h(
            NButton,
            {
              text: true,
              onClick: () => handleDownload(row)
            },
            {
              default: () =>
                h(NIcon, null, {
                  default: () => h(DownloadOutlined)
                })
            }
          ),
          h(
            NPopconfirm,
            {
              onPositiveClick: () => handleDelete(row)
            },
            {
              trigger: () =>
                h(
                  NButton,
                  {
                    text: true
                  },
                  {
                    default: () =>
                      h(NIcon, null, {
                        default: () => h(DeleteOutlined)
                      })
                  }
                ),
              default: () => '确认删除该文件？'
            }
          )
        ]
      })
    }
  }
]

const handleFileClick = (file: FileInfo) => {
  if (file.isDir) {
    fileStore.enterDirectory(file.fileName)
  }
}

const handleFavoriteClick = async (file: FileInfo) => {
  if (checkFavorite(file)) {
    await fileStore.removeFromFavorites(file.filePath)
  } else {
    await fileStore.addToFavorites(file.filePath)
  }
}

const handleDownload = (file: FileInfo) => {
  window.open(getDownloadUrl(file.filePath))
}

const handleDelete = async (file: FileInfo) => {
  try {
    await deleteFile(file.filePath)
    message.success('删除成功')
    fileStore.fetchFiles()
  } catch (error) {
    if (isAxiosError(error)) {
      message.error(error.response?.data?.msg || '删除失败：网络请求失败')
    } else {
      message.error('删除失败：未知错误')
    }
    console.error('删除失败:', error)
  }
}
</script>

<style scoped>
.file-list {
  width: 100%;
}

.file-name {
  cursor: pointer;
  color: var(--primary-color);
}

.file-name:hover {
  text-decoration: underline;
}
</style> 