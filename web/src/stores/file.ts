import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { FileInfo } from '@/types/file'
import { createDiscreteApi } from 'naive-ui'

const { message } = createDiscreteApi(['message'])

// 文件大小格式化
const formatFileSize = (size: number): string => {
  if (size === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const k = 1024
  const i = Math.floor(Math.log(size) / Math.log(k))
  return parseFloat((size / Math.pow(k, i)).toFixed(2)) + ' ' + units[i]
}

export const useFileStore = defineStore('file', () => {
  // 当前路径
  const currentPath = ref<string[]>([])
  const currentPathString = computed(() => currentPath.value.join('/'))

  // 文件列表
  const files = ref<FileInfo[]>([])
  const isLoading = ref(false)

  // 视图模式
  const viewMode = ref<'grid' | 'list'>('grid')
  const toggleViewMode = () => {
    viewMode.value = viewMode.value === 'grid' ? 'list' : 'grid'
  }

  // 排序后的文件列表
  const sortedFiles = computed(() => {
    const sorted = [...files.value]
    // 文件夹排在前面
    return sorted.sort((a, b) => {
      if (a.isDir && !b.isDir) return -1
      if (!a.isDir && b.isDir) return 1
      return a.fileName.localeCompare(b.fileName)
    })
  })

  // 获取文件列表
  const fetchFiles = async () => {
    try {
      isLoading.value = true
      const response = await fetch(`/api/file/list?path=${currentPathString.value}`)
      if (!response.ok) {
        throw new Error('获取文件列表失败')
      }
      const data = await response.json()
      if (!data || !data.data) {
        throw new Error('服务器返回数据格式错误')
      }
      files.value = Array.isArray(data.data) ? data.data : []
    } catch (error) {
      message.error(error instanceof Error ? error.message : '获取文件列表失败')
      files.value = [] // 出错时清空文件列表
    } finally {
      isLoading.value = false
    }
  }

  // 进入目录
  const enterDirectory = (path: string) => {
    if (path === '/') {
      currentPath.value = []
    } else {
      currentPath.value = path.split('/')
    }
    fetchFiles()
  }

  // 创建文件夹
  const createNewFolder = async (name: string) => {
    try {
      const path = currentPathString.value ? `${currentPathString.value}/${name}` : name
      const response = await fetch(`/api/file/create-folder?path=${path}`, {
        method: 'POST'
      })
      if (!response.ok) {
        throw new Error('创建文件夹失败')
      }
      message.success('创建文件夹成功')
      fetchFiles()
    } catch (error) {
      message.error(error instanceof Error ? error.message : '创建文件夹失败')
      throw error
    }
  }

  // 下载文件
  const downloadFile = async (path: string) => {
    try {
      const response = await fetch(`/api/file/download?path=${path}`)
      if (!response.ok) {
        throw new Error('下载文件失败')
      }
      const blob = await response.blob()
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = path.split('/').pop() || 'download'
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
      message.success('下载成功')
    } catch (error) {
      message.error(error instanceof Error ? error.message : '下载文件失败')
    }
  }

  // 删除文件
  const deleteFile = async (path: string) => {
    try {
      const response = await fetch(`/api/file/delete?path=${path}`, {
        method: 'DELETE'
      })
      if (!response.ok) {
        throw new Error('删除文件失败')
      }
      message.success('删除成功')
      fetchFiles()
    } catch (error) {
      message.error(error instanceof Error ? error.message : '删除文件失败')
    }
  }

  // 删除文件（别名，为了保持 API 一致性）
  const removeFile = deleteFile

  // 搜索文件
  const searchFile = async (keyword: string) => {
    try {
      isLoading.value = true
      const response = await fetch(`/api/file/search?keyword=${keyword}`)
      if (!response.ok) {
        throw new Error('搜索文件失败')
      }
      const data = await response.json()
      if (!data || !data.data) {
        throw new Error('服务器返回数据格式错误')
      }
      files.value = Array.isArray(data.data) ? data.data : []
    } catch (error) {
      message.error(error instanceof Error ? error.message : '搜索文件失败')
      files.value = [] // 出错时清空文件列表
    } finally {
      isLoading.value = false
    }
  }

  // 返回上级目录
  const goBack = async () => {
    if (currentPath.value.length > 0) {
      currentPath.value.pop()
      await fetchFiles()
    }
  }

  return {
    currentPath,
    currentPathString,
    files,
    isLoading,
    viewMode,
    sortedFiles,
    fetchFiles,
    enterDirectory,
    createNewFolder,
    downloadFile,
    deleteFile,
    removeFile,
    toggleViewMode,
    searchFile,
    goBack,
    formatFileSize
  }
})